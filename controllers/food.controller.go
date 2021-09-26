package foodController

import (
	"context"
	"fmt"
	"go-with-mongodb/database"
	myhelpers "go-with-mongodb/helpers"
	foods "go-with-mongodb/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var keyCollection string = "Foods"
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, keyCollection)
var validate = validator.New()

func PostFood(c *gin.Context) {
	//this is used to determine how long the API call should last
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var food foods.Food
	//bind the object that comes in with the declared varaible. thrrow an error if one occurs
	if err := c.BindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// use the validation packge to verify that all items coming in meet the requirements of the struct
	validationErr := validate.Struct(food)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}
	// assign the time stamps upon creation
	food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	//generate new ID for the object to be created
	food.ID = primitive.NewObjectID()

	// assign the the auto generated ID to the primary key attribute

	//replaced because we use a trigger for autoincrement value
	//food.Food_id = food.ID.Hex()
	var num = myhelpers.ToFixed(*food.Price, 2)
	food.Price = &num
	result, insertErr := foodCollection.InsertOne(ctx, food)
	if insertErr != nil {
		msg := fmt.Sprintf("Food item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	defer cancel()
	var newRecord bson.M
	err := foodCollection.FindOne(context.TODO(), bson.D{{"_id", result.InsertedID}}).Decode(&newRecord)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("found document %v", result)

	//return the id of the created object to the frontend
	c.JSON(http.StatusOK, newRecord)
}

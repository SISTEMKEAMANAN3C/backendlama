package controllers

import (
	"context"
	"golangsidang/database"
	"golangsidang/models"
	"golangsidang/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var faqAllCollection *mongo.Collection = database.OpenCollection(database.Client, "faq") // membuat collection baru

var validate_faq = validator.New()

func CreateFaq() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.FaqAll
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate_faq.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newUser := models.FaqAll{
			ID:         primitive.NewObjectID(),
			Question:   user.Question,
			Answer:     user.Answer,
			Created_at: time.Now(),
			Is_publish: user.Is_publish,
		}

		result, err := faqAllCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.GetallUser{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetFaq() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("id")
		var faq models.FaqAll
		defer cancel()

		objID, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid FAQ ID",
			})
			return
		}

		err = faqAllCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&faq)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "FAQ not found",
			})
			return
		}

		// Check if the FAQ is published (Is_publish is true) and Answer is not empty
		if faq.Is_publish && faq.Answer != nil && *faq.Answer != "" {
			c.JSON(http.StatusOK, faq)
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "FAQ not found",
			})
		}
	}
}

func EditFaq() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		absensiId := c.Param("id")
		var product models.FaqAll
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(absensiId)

		// Validate the request body
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, responses.FaqAll{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Set Is_publish to true
		product.Is_publish = true

		// Use the validator library to validate required fields
		if validationErr := validate_faq.Struct(&product); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.FaqAll{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"Question": product.Question, "Answer": product.Answer, "Is_publish": product.Is_publish}
		result, err := faqAllCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.FaqAll{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Get updated FAQ details
		var updatedAbsensi models.FaqAll
		if result.MatchedCount == 1 {
			err := faqAllCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedAbsensi)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.FaqAll{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.FaqAll{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedAbsensi}})
	}
}

func DeleteFaq() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userIDdel := c.Param("id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userIDdel)

		result, err := faqAllCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.FaqAll{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.FaqAll{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Absensi with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.FaqAll{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Absensi successfully deleted!"}},
		)
	}
}

func GetAllFaq() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var absensis []models.FaqAll
		defer cancel()

		results, err := faqAllCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.FaqAll{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.FaqAll
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, responses.FaqAll{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			absensis = append(absensis, singleUser)
		}

		c.JSON(http.StatusOK,
			responses.FaqAll{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": absensis}},
		)
	}
}

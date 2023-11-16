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

var userAllCollection *mongo.Collection = database.OpenCollection(database.Client, "user") // membuat collection baru

var validate_absensi = validator.New()

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate_absensi.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newUser := models.User{
			ID:           primitive.NewObjectID(),
			First_name:   user.First_name,
			Last_name:    user.Last_name,
			Password:     user.Password,
			Email:        user.Email,
			Phone:        user.Phone,
			Token:        user.Token,
			User_type:    user.User_type,
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
			Paseto_token: user.Paseto_token,
		}

		result, err := userAllCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.GetallUser{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetAuser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("usersGetId")
		var user models.User
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		err := userAllCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.GetallUser{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}})
	}
}

func EditAbsensi() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		absensiId := c.Param("absensiID")
		var absensi models.User
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(absensiId)

		//validate the request body
		if err := c.BindJSON(&absensi); err != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate_absensi.Struct(&absensi); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": absensi.First_name, "location": absensi.Last_name, "title": absensi.Email, "Phone": absensi.Phone, "User_type": absensi.User_type}
		result, err := userAllCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated user details
		var updatedAbsensi models.User
		if result.MatchedCount == 1 {
			err := userAllCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedAbsensi)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.GetallUser{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedAbsensi}})
	}
}

func DeleteAabsensi() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userIDdel := c.Param("usersID")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userIDdel)

		result, err := userAllCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.GetallUser{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Absensi with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.GetallUser{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Absensi successfully deleted!"}},
		)
	}
}

func GetAllAbssenis() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var absensis []models.User
		defer cancel()

		results, err := userAllCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.User
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			absensis = append(absensis, singleUser)
		}

		c.JSON(http.StatusOK,
			responses.GetallUser{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": absensis}},
		)
	}
}

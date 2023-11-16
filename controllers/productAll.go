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

var productAllCollection *mongo.Collection = database.OpenCollection(database.Client, "product") // membuat collection baru

var validate_product = validator.New()

func CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var product models.ProductAll
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, responses.GetAProduct{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate_product.Struct(&product); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.GetAProduct{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newUser := models.ProductAll{
			ID:          primitive.NewObjectID(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Size:        product.Size,
			Image:       product.Image,
		}

		result, err := productAllCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetAProduct{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.GetAProduct{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("productGetId")
		var user models.ProductAll
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		err := productAllCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetAProduct{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.GetAProduct{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}})
	}
}

func EditProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		absensiId := c.Param("productID")
		var product models.ProductAll
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(absensiId)

		//validate the request body
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, responses.GetAProduct{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate_product.Struct(&product); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.GetAProduct{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": product.Name, "Description": product.Description, "Price": product.Price, "Stock": product.Stock, "Size": product.Size, "Image": product.Image}
		result, err := productAllCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetAProduct{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated user details
		var updatedAbsensi models.ProductAll
		if result.MatchedCount == 1 {
			err := productAllCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedAbsensi)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.GetAProduct{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.GetAProduct{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedAbsensi}})
	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userIDdel := c.Param("productID")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userIDdel)

		result, err := productAllCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetAProduct{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.GetAProduct{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Product with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.GetAProduct{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Product successfully deleted!"}},
		)
	}
}

func GetAllProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var absensis []models.ProductAll
		defer cancel()

		results, err := productAllCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetAProduct{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var ProductUser models.ProductAll
			if err = results.Decode(&ProductUser); err != nil {
				c.JSON(http.StatusInternalServerError, responses.GetAProduct{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			absensis = append(absensis, ProductUser)
		}

		c.JSON(http.StatusOK,
			responses.GetAProduct{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": absensis}},
		)
	}
}

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

var ruanganAllCollection *mongo.Collection = database.OpenCollection(database.Client, "ruangan") // membuat collection baru

var validate_ruangan = validator.New()

func CreateRuang() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var ruangan models.Ruangan
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&ruangan); err != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate_ruangan.Struct(&ruangan); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newRuangan := models.Ruangan{
			ID:                   primitive.NewObjectID(),
			Nama_Ruangan:         ruangan.Nama_Ruangan,
			Deskripsi:            ruangan.Deskripsi,
			Foto:                 ruangan.Foto,
			User:                 ruangan.User,
			Status:               ruangan.Status,
			Tanggal_peminjaman:   ruangan.Tanggal_peminjaman,
			Tanggal_pengembalian: ruangan.Tanggal_pengembalian,
			Created_at:           time.Now(),
		}

		result, err := ruanganAllCollection.InsertOne(ctx, newRuangan)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusCreated, responses.GetallUser{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetRuangan() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		ruanganId := c.Param("id")
		var ruangan models.Ruangan
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(ruanganId)

		err := ruanganAllCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&ruangan)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetAProduct{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.GetAProduct{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": ruangan}})

	}
}

func EditRuangan() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		ruanganId := c.Param("id")
		var ruangan models.Ruangan
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(ruanganId)

		if err := c.BindJSON(&ruangan); err != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validatorErr := validate_ruangan.Struct(&ruangan); validatorErr != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validatorErr.Error()}})
			return
		}

		update := bson.M{
			"nama_ruangan":         ruangan.Nama_Ruangan,
			"deskripsi":            ruangan.Deskripsi,
			"foto":                 ruangan.Foto,
			"user":                 ruangan.User,
			"status":               ruangan.Status,
			"tanggal_peminjaman":   ruangan.Tanggal_peminjaman,
			"tanggal_pengembalian": ruangan.Tanggal_pengembalian,
			"created_at":           time.Now(),
		}
		result, err := ruanganAllCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		var ruanganUpdated models.Ruangan
		if result.MatchedCount == 1 {
			err := ruanganAllCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&ruanganUpdated)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}
		c.JSON(http.StatusOK, responses.GetallUser{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": ruanganUpdated}})
	}
}

func DeleteRuangan() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		ruanganId := c.Param("id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(ruanganId)
		result, err := ruanganAllCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.GetallUser{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "ruangan with specified ID not found!"}})
			return
		}
		c.JSON(http.StatusOK, responses.GetallUser{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "ruangan successfully deleted!"}})
	}
}

func GetAllRuangan() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var ruangans []models.Ruangan
		defer cancel()

		results, err := ruanganAllCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var RuanganS models.Ruangan
			if err = results.Decode(&RuanganS); err != nil {
				c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			ruangans = append(ruangans, RuanganS)
		}
		c.JSON(http.StatusOK, responses.GetallUser{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": ruangans}})
	}
}

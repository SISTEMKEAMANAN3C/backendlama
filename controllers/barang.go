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

var barangAllCollection *mongo.Collection = database.OpenCollection(database.Client, "barang") // membuat collection baru
var validate_barang = validator.New()

func CreateBarang() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var barang models.Barang
		defer cancel()

		if err := c.BindJSON(&barang); err != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate_barang.Struct(&barang); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newBarang := models.Barang{
			ID:                   primitive.NewObjectID(),
			Nama_Barang:          barang.Nama_Barang,
			Deskripsi:            barang.Deskripsi,
			Foto:                 barang.Foto,
			Stock:                barang.Stock,
			User:                 barang.User,
			Status:               barang.Status,
			Tanggal_peminjaman:   barang.Tanggal_peminjaman,
			Tanggal_pengembalian: barang.Tanggal_pengembalian,
			Created_at:           time.Now(),
		}
		result, err := barangAllCollection.InsertOne(ctx, newBarang)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusCreated, responses.GetallUser{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetBarang() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		barangId := c.Param("id")
		var barang models.Barang
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(barangId)
		err := barangAllCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&barang)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.GetallUser{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": barang}})
	}
}

func UpdateBarang() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		barangId := c.Param("id")
		var barang models.Barang
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(barangId)

		if err := c.BindJSON(&barang); err != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate_barang.Struct(&barang); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.GetallUser{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}
		update := bson.M{
			"nama_barang":          barang.Nama_Barang,
			"deskripsi":            barang.Deskripsi,
			"foto":                 barang.Foto,
			"user":                 barang.User,
			"stock":                barang.Stock,
			"status":               barang.Status,
			"tanggal_peminjaman":   barang.Tanggal_peminjaman,
			"tanggal_pengembalian": barang.Tanggal_pengembalian,
			"created_at":           barang.Created_at,
			"updated_at":           time.Now(),
		}
		result, err := barangAllCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.D{{"$set", update}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedBarang models.Barang
		if result.MatchedCount == 1 {
			err = barangAllCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedBarang)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			c.JSON(http.StatusOK, responses.GetallUser{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedBarang}})
		}
		c.JSON(http.StatusOK, responses.GetallUser{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedBarang}})

	}
}

func DeleteBarang() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		barangId := c.Param("id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(barangId)

		result, err := barangAllCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.GetallUser{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetAllBarang() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var baransg []models.Barang
		defer cancel()

		result, err := barangAllCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		defer result.Close(ctx)
		for result.Next(ctx) {
			var Barangs models.Barang
			if err = result.Decode(&Barangs); err != nil {
				c.JSON(http.StatusInternalServerError, responses.GetallUser{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			baransg = append(baransg, Barangs)
		}
		c.JSON(http.StatusOK, responses.GetallUser{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": baransg}})
	}
}

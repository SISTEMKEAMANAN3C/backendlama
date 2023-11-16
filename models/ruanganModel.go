package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ruangan struct {
	ID                   primitive.ObjectID `bson:"_id"`
	Nama_Ruangan         *string            `json:"nama_ruangan" bson:"nama_ruangan"`
	Deskripsi            *string            `json:"deskripsi" bson:"deskripsi"`
	Foto                 *string            `json:"foto" bson:"foto"`
	User                 User               `json:"user" bson:"user"`
	Status               *bool              `json:"status" bson:"status"`
	Tanggal_peminjaman   *string            `json:"tanggal_peminjaman" bson:"tanggal_peminjaman"`
	Tanggal_pengembalian *string            `json:"tanggal_pengembalian" bson:"tanggal_pengembalian"`
	Created_at           time.Time          `json:"created_at" bson:"created_at"`
}

type Barang struct {
	ID                   primitive.ObjectID `bson:"_id"`
	Nama_Barang          *string            `json:"nama_barang" bson:"nama_barang"`
	Deskripsi            *string            `json:"deskripsi" bson:"deskripsi"`
	Stock                *int               `json:"stock" bson:"stock"`
	Foto                 *string            `json:"foto" bson:"foto"`
	User                 User               `json:"user" bson:"user"`
	Status               bool               `json:"status" bson:"status"`
	Tanggal_peminjaman   *string            `json:"tanggal_peminjaman" bson:"tanggal_peminjaman"`
	Tanggal_pengembalian *string            `json:"tanggal_pengembalian" bson:"tanggal_pengembalian"`
	Created_at           time.Time          `json:"created_at" bson:"created_at"`
}

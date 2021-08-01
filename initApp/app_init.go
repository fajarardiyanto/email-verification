package initApp

import (
	"context"
	"log"
	"mail-service/structs"
	"mail-service/utils"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
)

var mainDB *qmgo.Database
var dbName = "DB"

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dbName = os.Getenv("DB_NAME")
	initMongoDB()
}

func initMongoDB() {
	if client, err := qmgo.NewClient(context.TODO(), &qmgo.Config{Uri: os.Getenv("MONGO_URI")}); err != nil {
		panic(err)
	} else {
		mainDB = client.Database(dbName)
	}

	mainDB.Collection("users").CreateOneIndex(context.TODO(), options.IndexModel{
		Key:    []string{"username"},
		Unique: true,
	})

	if err := mainDB.Collection("users").Find(context.TODO(), bson.M{"username": "admin"}).One(&structs.User{}); err == qmgo.ErrNoSuchDocuments {
		hashed, err := utils.HashedPassword("admin")
		if err != nil {
			log.Println(err)
		}
		mainDB.Collection("users").InsertOne(context.TODO(), structs.User{
			UserID:    uuid.New().String(),
			Name:      "admin",
			Username:  "admin",
			Email:     "admin@localhost",
			Password:  string(hashed),
			CreatedAt: time.Now().Unix(),
		})
	}
}

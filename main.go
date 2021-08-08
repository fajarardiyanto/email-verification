package main

import (
	"context"
	"mail-service/controllers"
	"mail-service/global"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/qiniu/qmgo"
)

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
		global.MainDB = client.Database(dbName)
	}
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.SignUp)

	r.Run(os.Getenv("BIND_ADDR"))
}

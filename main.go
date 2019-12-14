package main

import (
	"fmt"
	"github.com/adigunhammedolalekan/sms-forwarder/db"
	"github.com/adigunhammedolalekan/sms-forwarder/fn"
	"github.com/adigunhammedolalekan/sms-forwarder/http"
	"github.com/adigunhammedolalekan/sms-forwarder/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	conn, err := db.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	api := router.Group("/api")

	tokenGenerator := fn.NewJwtTokenGenerator([]byte(os.Getenv("JWT_SECRET")))
	userStore := store.NewUserStore(conn, tokenGenerator)

	handler := http.NewUserHttpHandler(userStore)
	api.POST("/user/new", handler.CreateUserHandler)
	api.POST("/user/authenticate", handler.AuthenticateUserHandler)

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
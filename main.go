package main

import (
	"eventify/config"
	"eventify/routes"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if err := config.ConnectDB(); err != nil {
		log.Fatal("DB connect error: ", err)
	}
	fmt.Println("DB CONNECTED SUCCESFULLY")

	r := routes.Setup()
	r.Run(":8080")
}
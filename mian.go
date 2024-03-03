package main

import (
	"github.com/AhmedSamy16/02-Subscribers-API/application"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := application.New()

	app.Start()
}

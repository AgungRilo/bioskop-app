package main

import (
	"os"

	"github.com/AgungRilo/bioskop-app/config"
	"github.com/AgungRilo/bioskop-app/routes"
	"github.com/joho/godotenv"
)

func main() {
	// .env opsional, variable dari environment tetap bisa dipakai
	godotenv.Load()

	config.ConnectDatabase()
	defer config.DB.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := routes.SetupRouter()
	r.Run(":" + port)
}

package main

import (
	"os"

	"github.com/AgungRilo/bioskop-app/config"
	"github.com/AgungRilo/bioskop-app/routes"
)

func main() {
	config.ConnectDatabase()
	defer config.DB.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	r := routes.SetupRouter()
	r.Run(":" + port)
}

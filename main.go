package main

import (
	"go-drive-project/config"
	"go-drive-project/logger"
	"go-drive-project/router"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.ViperConfig()
	PORT := viper.GetString("PORT")
	if PORT == "" {
		PORT = "8080" // Default to 8080 if not set
	}
	r := gin.Default()
	router.NewRouter(r)
	logger.Log.Println("Server is running on port:", PORT)
	if err := r.Run(":" + PORT); err != nil {
		logger.Log.Fatal("Failed to start server: ", err)
	}
}

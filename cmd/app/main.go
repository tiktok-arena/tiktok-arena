package main

import (
	"log"
	"tiktok-arena/configuration"
	"tiktok-arena/internal/api/app"
)

//	@title						TikTok arena API
//	@version					1.0
//	@description				API for TikTok arena application
//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@BasePath					/api/
func main() {
	err := configuration.LoadConfig(".env")
	if err != nil {
		log.Fatalln("Failed to load environment variables!", err.Error())
	}
	app.Run(&configuration.EnvConfig)
}

package main

import (
	"log"

	"gitlab.silkrode.com.tw/team_golang/kbc/km/storage/config"
	"gitlab.silkrode.com.tw/team_golang/kbc/km/storage/pkg/database"
	internal "gitlab.silkrode.com.tw/team_golang/kbc/km/storage/pkg/log"
	"gitlab.silkrode.com.tw/team_golang/kbc/km/storage/routes"
)

func main() {
	config, err := config.InitConfig("./config")
	if err != nil {
		log.Fatal("Failed to init config.")
	}

	logger, err := internal.InitLogger(config)
	if err != nil {
		log.Fatal("Failed to init logger.")
	}

	db, err := database.InitDB(config)
	if err != nil {
		log.Fatal("Failed to init db.")
	}

	//-----Restful-----
	r := routes.InitRouter(config, logger, db)
	r.Run(":" + config.Server.Port)

	//-----GraphQL-----
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = config.Server.Port
	// }
}

package main

import (
	"be-2/src/api/route"
	"be-2/src/config/database"
	"be-2/src/config/env"
	"be-2/src/config/storage"

	"github.com/joho/godotenv"
)

func main() {
	// load env file
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// init mysql
	database.InitMySQL()

	// migrate gorm
	database.MigrateMySQL()

	storage.InitGDrive()

	app := route.InitServer()
	app.Logger.Fatal(app.Start(":" + env.GetServerEnv()))
}

package models

import (
	"fmt"
	"log"

	"blue-admin.com/configs"
	"blue-admin.com/database"
)

func InitDatabase() {
	configs.NewEnvFile("./configs")
	database, err  := database.ReturnSession()
	fmt.Println("Connection Opened to Database")
	if err == nil {
		if err := database.AutoMigrate(
			&Role{},
			&App{},
			&User{},
			&Feature{},
			&Endpoint{},
			&Page{},
			&JWTSalt{},
		); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Database Migrated")
	} else {
		panic(err)
	}
}

func CleanDatabase() {
	configs.NewEnvFile("./configs")
	database, err := database.ReturnSession()
	if err == nil {
		fmt.Println("Connection Opened to Database")
		fmt.Println("Dropping Models if Exist")
		database.Migrator().DropTable(

			&Role{},

			&App{},

			&User{},

			&Feature{},

			&Endpoint{},

			&Page{},

			&JWTSalt{},

		)

		fmt.Println("Database Cleaned")
	} else {
		panic(err)
	}
}

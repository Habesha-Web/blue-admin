package models

import (
	"crypto/sha512"
	"encoding/hex"
	"log"

	"blue-admin.com/configs"
	"blue-admin.com/database"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

var Endpoints_JSON = make(map[string]string)

// Combine password and salt then hash them using the SHA-512
func hashfunc(password string) string {

	// var salt []byte
	// get salt from env variable
	salt := []byte(configs.AppConfig.Get("SECRETE_SALT"))

	// Convert password string to byte slice
	var pwdByte = []byte(password)

	// Create sha-512 hasher
	var sha512 = sha512.New()

	pwdByte = append(pwdByte, salt...)

	sha512.Write(pwdByte)

	// Get the SHA-512 hashed password
	var hashedPassword = sha512.Sum(nil)

	// Convert the hashed to hex string
	var hashedPasswordHex = hex.EncodeToString(hashedPassword)
	return hashedPasswordHex
}

func GetAppFeatures(app_uuid string) {
	db, _ := database.ReturnSession()
	var app App
	// var app_get = make(map[string]string)
	app_id, _ := uuid.Parse(app_uuid)
	if res := db.Model(&App{}).Preload(clause.Associations).Preload("Roles.Features").Preload("Roles.Features.Endpoints").Where("uuid = ?", app_id).First(&app); res.Error != nil {
		log.Fatal(res.Error.Error())
	}

	for _, value := range app.Roles {
		key := value.Name
		for _, value := range value.Features {
			for _, value := range value.Endpoints {
				Endpoints_JSON[value.Name] = key
			}
		}

	}

}

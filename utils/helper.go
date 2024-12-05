package utils

import (
	"context"
	"log"

	"blue-admin.com/configs"
	"blue-admin.com/database"
	"blue-admin.com/models"
	"gorm.io/gorm"
)

type ResourceMatrix struct {
	Name      string `gorm:"not null; unique;" json:"name,omitempty"`
	RoleName  string `gorm:"not null; unique;" json:"role_name,omitempty"`
	RoutePath string `gorm:"not null; unique;" json:"route_path,omitempty"`
}

var Endpoints_JSON = make(map[string]string)

func GetAppFeatures() {
	app_uuid := configs.AppConfig.Get("APP_ID")
	db, _ := database.ReturnSession()
	var role_matrix []ResourceMatrix

	query_string := `SELECT endpoints.name,roles.name as role_name FROM apps
		INNER JOIN roles ON apps.id = roles.app_id
		INNER JOIN features ON features.role_id = roles.id
		INNER JOIN endpoints ON features.id = endpoints.feature_id
		WHERE apps.uuid = ?
		  AND apps.active = true
		  AND roles.active = true
		  AND features.active = true
					ORDER BY apps.id`
	if res := db.Model(&models.App{}).Raw(query_string, app_uuid).Scan(&role_matrix); res.Error != nil {
		log.Fatal(res.Error.Error())
	}

	for _, value := range role_matrix {
		Endpoints_JSON[value.Name] = value.RoleName

	}
}

func GetAppFeaturesReturn(app_uuid string, db *gorm.DB, ctx context.Context) (map[string]string, error) {

	var role_matrix_list []ResourceMatrix

	query_string := `SELECT endpoints.name,roles.name as role_name FROM apps
		INNER JOIN roles ON apps.id = roles.app_id
		INNER JOIN features ON features.role_id = roles.id
		INNER JOIN endpoints ON features.id = endpoints.feature_id
		WHERE apps.uuid = ?
		  AND apps.active = true
		  AND roles.active = true
		  AND features.active = true
					ORDER BY apps.id`
	// app_id, _ := uuid.Parse(app_uuid)
	// if res := db.WithContext(ctx).Model(&models.App{}).Preload(clause.Associations).Preload("Roles.Features").Where("active = ?", true).Preload("Roles.Features.Endpoints").Where("uuid = ?", app_uuid).First(&app); res.Error != nil {
	if res := db.WithContext(ctx).Raw(query_string, app_uuid).Scan(&role_matrix_list); res.Error != nil {

		return nil, res.Error
	}
	role_matrix := make(map[string]string)

	for _, value := range role_matrix_list {
		role_matrix[value.Name] = value.RoleName
	}

	return role_matrix, nil
}
func GetAppFeaturesReturnPath(app_uuid string, db *gorm.DB, ctx context.Context) (map[string]string, error) {

	var role_matrix_list []ResourceMatrix

	query_string := `SELECT endpoints.route_path,roles.name as role_name FROM apps
		INNER JOIN roles ON apps.id = roles.app_id
		INNER JOIN features ON features.role_id = roles.id
		INNER JOIN endpoints ON features.id = endpoints.feature_id
		WHERE apps.uuid = ?
		  AND apps.active = true
		  AND roles.active = true
		  AND features.active = true
					ORDER BY apps.id`
	// app_id, _ := uuid.Parse(app_uuid)
	// if res := db.WithContext(ctx).Model(&models.App{}).Preload(clause.Associations).Preload("Roles.Features").Where("active = ?", true).Preload("Roles.Features.Endpoints").Where("uuid = ?", app_uuid).First(&app); res.Error != nil {
	if res := db.WithContext(ctx).Raw(query_string, app_uuid).Scan(&role_matrix_list); res.Error != nil {

		return nil, res.Error
	}
	role_matrix := make(map[string]string)

	for _, value := range role_matrix_list {
		role_matrix[value.RoutePath] = value.RoleName
	}

	return role_matrix, nil
}

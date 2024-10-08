package main

import (
	"blue-admin.com/manager"
)

//	@title			Swagger blue-admin API
//	@version		0.1
//	@description	This is blue-admin API OPENAPI Documentation.
//	@termsOfService	http://swagger.io/terms/
//  @BasePath  /api/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-APP-TOKEN
//	@description				Description for what is this security definition being used

//	@securityDefinitions.apikey Refresh
//	@in							header
//	@name						X-REFRESH-TOKEN
//	@description				Description for what is this security definition being used

// go build -tags netgo -ldflags '-s -w' -o app
func main() {
	manager.Execute()
}

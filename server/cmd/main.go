package main

import (
	"github.com/respondnow/respond/server"
	_ "github.com/respondnow/respond/server/docs"
)

//	@title			RespondNow Server
//	@version		1.0
//	@description	RespondNow Server APIs.
//	@termsOfService	https://www.harness.io/legal/subscription-terms
//
//	@contact.name	RespondNow

//	@host						localhost:8080
//	@BasePath					/
//
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization

func main() {
	server.Start()
}

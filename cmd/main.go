package main

//	@title			Payment Gateway API
//	@version		1.0.0
//	@description	This is a production-ready payment gateway service that demonstrates enterprise-grade architecture with idempotent operations, asynchronous processing, and reliable message delivery.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

import (
	"context"
	"os"
	"os/signal"
	"payment-gateway/initiator"
	"syscall"

	_ "payment-gateway/docs"
)

// This will be generated

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	initiator.Initialize(ctx)
}

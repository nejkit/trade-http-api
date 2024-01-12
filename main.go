package main

import (
	"context"
	"trade-http-api/app"
)

func main() {
	ctx := context.Background()
	app.InitApi(ctx)
}

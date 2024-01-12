package app

import (
	"context"
	"net/http"
	"trade-http-api/handlers"

	"github.com/gin-gonic/gin"
)

type ServerApi struct {
	handler handlers.HttpHandler
	host    string
}

func NewServer(handler handlers.HttpHandler, host string) ServerApi {
	return ServerApi{handler: handler, host: host}
}

func (s *ServerApi) StartServe(ctx context.Context) {
	router := gin.Default()

	router.POST("/create-asset", func(ctx *gin.Context) {
		s.handler.HandleCreateAsset(ctx)
	})

	router.POST("/emmit-asset", func(ctx *gin.Context) {
		s.handler.HandleEmmitAsset(ctx)
	})

	router.GET("/asset/:assetid", func(ctx *gin.Context) {
		s.handler.HandleGetAssets(ctx)
	})

	server := &http.Server{
		Addr:    s.host,
		Handler: router,
	}

	go server.ListenAndServe()

	select {
	case <-ctx.Done():
		server.Close()
		break
	}

}

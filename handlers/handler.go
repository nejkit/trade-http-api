package handlers

import (
	"context"
	"net/http"
	"trade-http-api/models"
	"trade-http-api/validators"

	"github.com/gin-gonic/gin"
)

type IBpsProvider interface {
	CreateAsset(ctx context.Context, req models.BpsCreateAssetRequest) (*models.BpsCreateAssetResponse, error)
	EmmitAsset(ctx context.Context, req models.BpsEmmitAssetRequest) error
	GetAssets(ctx context.Context, id string) (*models.BpsGetAssetsResponse, error)
}

type HttpHandler struct {
	bpsProvider IBpsProvider
}

func NewHttpHandler(bpsProvider IBpsProvider) HttpHandler {
	return HttpHandler{bpsProvider: bpsProvider}
}

func (h *HttpHandler) HandleCreateAsset(ctx *gin.Context) {
	var request models.BpsCreateAssetRequest

	err := ctx.BindJSON(&request)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Bad JSON type")
		return
	}

	err = validators.ValidateCreateAssetRequest(request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Validation error": err.Error()})
		return
	}

	resp, err := h.bpsProvider.CreateAsset(ctx, request)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) HandleEmmitAsset(ctx *gin.Context) {
	var request models.BpsEmmitAssetRequest

	err := ctx.BindJSON(&request)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Bad JSON type")
		return
	}

	err = validators.ValidateEmmitAssetRequest(request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Validation error": err.Error()})
		return
	}

	err = h.bpsProvider.EmmitAsset(ctx, request)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *HttpHandler) HandleGetAssets(ctx *gin.Context) {

	assetId := ctx.Param("assetid")

	if assetId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"Validation error": "InvalidAssetId"})
		return
	}

	resp, err := h.bpsProvider.GetAssets(ctx, assetId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

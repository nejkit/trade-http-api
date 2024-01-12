package validators

import (
	"errors"
	"trade-http-api/models"

	"github.com/go-playground/validator/v10"
)

var (
	validateErrorInvalidCurrencyCode      = errors.New("InvalidCurrencyCode")
	validateErrorInvalidAmount            = errors.New("InvalidAmount")
	validateErrorInvalidEmmitBalancesInfo = errors.New("InvalidEmmitBalancesInfo")
	validateErrorInvalidAssetId           = errors.New("InvalidAssetId")
)

func ValidateCreateAssetRequest(request models.BpsCreateAssetRequest) error {
	if request.EmmitBalanceInfos == nil {
		return nil
	}

	v := validator.New()

	for _, emmitData := range request.EmmitBalanceInfos {
		err := v.Struct(emmitData)

		if err != nil {
			return mapError(err)
		}
	}
	return nil
}

func ValidateEmmitAssetRequest(request models.BpsEmmitAssetRequest) error {
	v := validator.New()

	err := v.Struct(request)

	if err != nil {
		return mapError(err)
	}

	for _, emmitData := range request.EmmitBalanceInfos {
		err = v.Struct(emmitData)

		if err != nil {
			return mapError(err)
		}
	}
	return nil
}

func mapError(err error) error {
	validErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		return err
	}

	for _, ve := range validErrors {
		if ve.Field() == "AssetId" {
			return validateErrorInvalidAssetId
		}
		if ve.Field() == "EmmitBalanceInfos" {
			return validateErrorInvalidEmmitBalancesInfo
		}
		if ve.Field() == "CurrencyCode" {
			return validateErrorInvalidCurrencyCode
		}
		if ve.Field() == "Amount" {
			return validateErrorInvalidAmount
		}
	}

	return nil
}

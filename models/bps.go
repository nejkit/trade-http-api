package models

import "time"

type EmmitBalanceInfo struct {
	CurrencyCode string  `json:"currency_code,omitempty" validate:"required,len=3"`
	Amount       float64 `json:"amount,omitempty" validate:"required,gt=0.00"`
}

type BalanceInfo struct {
	CurrencyCode string  `json:"currency_code,omitempty"`
	Amount       float64 `json:"amount,omitempty"`
	LockedAmount float64 `json:"locked_amount,omitempty"`
}

type BpsCreateAssetRequest struct {
	EmmitBalanceInfos []EmmitBalanceInfo `json:"emmit_balance_infos,omitempty"`
}

type BpsCreateAssetResponse struct {
	AssetId string `json:"asset_id,omitempty"`
}

type BpsEmmitAssetRequest struct {
	AssetId           string             `json:"asset_id,omitempty" validate:"required"`
	EmmitBalanceInfos []EmmitBalanceInfo `json:"emmit_balance_infos,omitempty" validate:"required"`
}

type BpsGetAssetsResponse struct {
	AssetId      string        `json:"asset_id,omitempty"`
	CreatedDate  time.Time     `json:"created_date,omitempty"`
	BalancesInfo []BalanceInfo `json:"balances_info,omitempty"`
}

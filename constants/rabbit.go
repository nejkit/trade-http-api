package constants

const (
	BpsExchange = "bps.forward"
)

const (
	BpsCreateAssetResponseQueueName = "q.trade-api.response.CreateAsset"
	BpsGetAssetsResponseQueueName   = "q.trade-api.response.GetAsset"
)

const (
	RkCreateAssetRequest  = "r.trade-api.CreateAssetRequest.#"
	RkCreateAssetResponse = "r.#.CreateAssetResponse.#"

	RkEmmitAssetRequest = "r.trade-api.EmmitAssetRequest.#"

	RkGetAssetsRequest  = "r.trade-api.GetAssetsRequest.#"
	RkGetAssetsResponse = "r.#.GetAssetsResponse.#"
)

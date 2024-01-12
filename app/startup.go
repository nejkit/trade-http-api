package app

import (
	"context"
	"fmt"
	"os"
	"trade-http-api/config"
	"trade-http-api/constants"
	"trade-http-api/external/balances"
	"trade-http-api/handlers"
	"trade-http-api/providers"
	"trade-http-api/rabbit"
)

func InitApi(ctx context.Context) {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	var cfg config.Config
	// if err := env.Parse(&cfg); err != nil {
	// 	return
	// }

	cfg.RabbitUser = "guest"
	cfg.RabbitPassword = "guest"
	cfg.RabbitHost = "localhost"
	cfg.RabbitPort = "5672"
	cfg.ApiHost = "localhost:8080"

	rabbitConnection, err := rabbit.GetRabbitConnection(buildRabbitUrl(cfg))

	commonChannel, err := rabbitConnection.Channel()

	if err != nil {
		cancel()
		return
	}

	err = rabbit.InitRabbitInfrastructure(commonChannel)

	if err != nil {
		cancel()
		return
	}

	senderChannel, err := rabbitConnection.Channel()

	if err != nil {
		cancel()
		return
	}

	createAssetResponseChannel, err := rabbitConnection.Channel()

	if err != nil {
		cancel()
		return
	}

	getAssetResponseChannel, err := rabbitConnection.Channel()

	if err != nil {
		cancel()
		return
	}

	sender := rabbit.NewSender(senderChannel)

	createAssetStorage := rabbit.NewRabbitStorage[balances.BpsCreateAssetResponse]()
	getAssetStorage := rabbit.NewRabbitStorage[balances.BpsGetAssetInfoResponse]()

	createAssetProcessor := rabbit.NewProcessor[balances.BpsCreateAssetResponse](rabbit.GetIdFuncForCreateAsset(), rabbit.GetParserFuncForCreateAsset(), &createAssetStorage)
	getAssetsProcessor := rabbit.NewProcessor[balances.BpsGetAssetInfoResponse](rabbit.GetIdFuncForGetAsset(), rabbit.GetParserFuncForGetAssets(), &getAssetStorage)

	createAssetListener, err := rabbit.NewListener[balances.BpsCreateAssetResponse](ctxWithCancel, createAssetResponseChannel, createAssetProcessor, constants.BpsCreateAssetResponseQueueName)
	getAssetsListener, err := rabbit.NewListener[balances.BpsGetAssetInfoResponse](ctxWithCancel, getAssetResponseChannel, getAssetsProcessor, constants.BpsGetAssetsResponseQueueName)

	bpsProvider := providers.NewBpsProvider(&createAssetStorage, &getAssetStorage, sender)

	handler := handlers.NewHttpHandler(bpsProvider)

	server := NewServer(handler, cfg.ApiHost)

	go createAssetListener.Run(ctxWithCancel)
	go getAssetsListener.Run(ctxWithCancel)

	go server.StartServe(ctxWithCancel)

	exit := make(chan os.Signal, 1)
	for {
		select {
		case <-exit:
			{
				cancel()
				return
			}
		}
	}
}

func buildRabbitUrl(cfg config.Config) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.RabbitUser, cfg.RabbitPassword, cfg.RabbitHost, cfg.RabbitPort)
}

package rabbit

import (
	"time"
	"trade-http-api/constants"
	"trade-http-api/external/balances"
	"trade-http-api/staticerr"

	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func GetRabbitConnection(connectionString string) (*amqp091.Connection, error) {
	timeout := time.After(time.Minute * 5)
	for {
		select {
		case <-timeout:
			return nil, staticerr.ErrorRabbitConnectionFail
		default:
			connect, err := amqp091.Dial(connectionString)

			if err != nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}

			return connect, nil
		}
	}
}

func InitRabbitInfrastructure(channel *amqp091.Channel) error {
	defer channel.Close()

	_, err := channel.QueueDeclare(constants.BpsCreateAssetResponseQueueName, true, false, false, false, nil)

	if err != nil {
		return err
	}

	_, err = channel.QueueDeclare(constants.BpsGetAssetsResponseQueueName, true, false, false, false, nil)

	if err != nil {
		return err
	}

	err = channel.QueueBind(constants.BpsCreateAssetResponseQueueName, constants.RkCreateAssetResponse, constants.BpsExchange, false, nil)

	if err != nil {
		return err
	}

	err = channel.QueueBind(constants.BpsGetAssetsResponseQueueName, constants.RkGetAssetsResponse, constants.BpsExchange, false, nil)

	if err != nil {
		return err
	}

	return nil
}

func GetIdFuncForCreateAsset() IdFunc[balances.BpsCreateAssetResponse] {
	return func(bcar *balances.BpsCreateAssetResponse) string {
		return bcar.Id
	}
}

func GetIdFuncForGetAsset() IdFunc[balances.BpsGetAssetInfoResponse] {
	return func(bgair *balances.BpsGetAssetInfoResponse) string {
		return bgair.Id
	}
}

func GetParserFuncForCreateAsset() ParserFunc[balances.BpsCreateAssetResponse] {
	return func(d amqp091.Delivery) (*balances.BpsCreateAssetResponse, error) {
		var protoMsg balances.BpsCreateAssetResponse
		err := proto.Unmarshal(d.Body, &protoMsg)
		if err != nil {
			return nil, err
		}
		return &protoMsg, nil
	}
}

func GetParserFuncForGetAssets() ParserFunc[balances.BpsGetAssetInfoResponse] {
	return func(d amqp091.Delivery) (*balances.BpsGetAssetInfoResponse, error) {
		var protoMsg balances.BpsGetAssetInfoResponse
		err := proto.Unmarshal(d.Body, &protoMsg)
		if err != nil {
			return nil, err
		}
		return &protoMsg, nil
	}
}

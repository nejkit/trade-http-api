package providers

import (
	"context"
	"errors"
	"trade-http-api/external/balances"
	"trade-http-api/models"

	"github.com/google/uuid"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IStorage[T any] interface {
	GetMessageById(id string) (*T, error)
}

type ISender interface {
	SendMessage(ctx context.Context, msg protoreflect.ProtoMessage, exchange, rk string) error
}

type BpsProvider struct {
	creationAssetStorage IStorage[balances.BpsCreateAssetResponse]
	getAssetsStorage     IStorage[balances.BpsGetAssetInfoResponse]
	sender               ISender
}

func NewBpsProvider(cas IStorage[balances.BpsCreateAssetResponse], gas IStorage[balances.BpsGetAssetInfoResponse], sender ISender) *BpsProvider {
	return &BpsProvider{creationAssetStorage: cas, getAssetsStorage: gas, sender: sender}
}

func (p *BpsProvider) CreateAsset(ctx context.Context, req models.BpsCreateAssetRequest) (*models.BpsCreateAssetResponse, error) {
	protoReq := mapCreateAssetToProto(req)
	err := p.sender.SendMessage(ctx, protoReq, "", "")
	if err != nil {
		return nil, err
	}

	resp, err := p.creationAssetStorage.GetMessageById(protoReq.Id)

	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.ErrorCode.String())
	}

	return &models.BpsCreateAssetResponse{AssetId: resp.AssetId}, nil
}

func (p *BpsProvider) EmmitAsset(ctx context.Context, req models.BpsEmmitAssetRequest) error {
	protoReq := mapEmmitAssetToProto(req)

	err := p.sender.SendMessage(ctx, protoReq, "", "")

	if err != nil {
		return err
	}

	return nil
}

func (p *BpsProvider) GetAssets(ctx context.Context, id string) (*models.BpsGetAssetsResponse, error) {
	protoReq := balances.BbsGetAssetInfoRequest{Id: uuid.NewString(), AssetId: id}

	err := p.sender.SendMessage(ctx, &protoReq, "", "")

	if err != nil {
		return nil, err
	}

	resp, err := p.getAssetsStorage.GetMessageById(id)

	if err != nil {
		return nil, err
	}

	return mapGetAssetsResponse(resp), nil
}

func mapEmmitDataToProto(data []models.EmmitBalanceInfo) []*balances.EmmitBalanceInfo {
	var protoData []*balances.EmmitBalanceInfo
	for _, emmitData := range data {
		protoData = append(protoData, &balances.EmmitBalanceInfo{
			CurrencyName: emmitData.CurrencyCode,
			Amount:       emmitData.Amount,
		})
	}
	return protoData
}

func mapCreateAssetToProto(req models.BpsCreateAssetRequest) *balances.BpsCreateAssetRequest {
	protoReq := balances.BpsCreateAssetRequest{
		Id: uuid.NewString(),
	}
	if req.EmmitBalanceInfos == nil {
		return &protoReq
	}

	protoReq.EmmitInfo = mapEmmitDataToProto(req.EmmitBalanceInfos)

	return &protoReq
}

func mapEmmitAssetToProto(req models.BpsEmmitAssetRequest) *balances.BpsEmmitAssetRequest {
	return &balances.BpsEmmitAssetRequest{
		Id:               uuid.NewString(),
		AssetId:          req.AssetId,
		EmitBalancesInfo: mapEmmitDataToProto(req.EmmitBalanceInfos),
	}
}

func mapGetAssetsResponse(protoResp *balances.BpsGetAssetInfoResponse) *models.BpsGetAssetsResponse {
	resp := models.BpsGetAssetsResponse{AssetId: protoResp.AssetId, CreatedDate: protoResp.CreatedDate.AsTime()}

	for _, balInfo := range protoResp.BalancesInfo {
		resp.BalancesInfo = append(resp.BalancesInfo, models.BalanceInfo{
			CurrencyCode: balInfo.CurrencyName,
			Amount:       balInfo.Amount,
			LockedAmount: balInfo.LockedAmount,
		})
	}

	return &resp
}

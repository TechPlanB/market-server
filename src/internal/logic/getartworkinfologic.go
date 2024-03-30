package logic

import (
	"context"
	"market/internal/common/errorx"
	"market/internal/svc"
	"market/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetArtworkInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArtworkInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetArtworkInfoLogic {
	return GetArtworkInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArtworkInfoLogic) GetArtworkInfo(req types.ArtworkInfoReq) (*types.ArtworkInfoResponse, error) {
	artworks, err := l.svcCtx.NftTokenModel.FindInfo(req.ID, req.IdInContract)
	if err != nil {
		return nil, errorx.NewDefaultError("can not find artwork records")
	}

	properties, err := l.svcCtx.NftTokenModel.FindProperties(req.ID)
	if err != nil {
		return nil, errorx.NewDefaultError("can not find token properties")
	}

	list := make([]types.Property, 0)

	for _, property := range properties {
		list = append(list, types.Property{
			Key:   property.Key,
			Value: property.Value.String,
		})
	}

	status := "not_on_sale"
	count := artworks.Count

	var fixed_price float64
	if artworks.Status.String == "active" {
		status = "on_sale"
		count = artworks.SaleCount.Int64
		fixed_price = artworks.TotalPrice.Float64 / float64(count)
	}

	ArtworksInfo := &types.ArtworkInfoBody{
		ID:             artworks.Id,
		Name:           artworks.Name.String,
		Description:    artworks.Description.String,
		CreatedAt:      artworks.CreatedAt.String(),
		CategoryName:   artworks.CategoryName,
		CollectionName: artworks.CollectionName,
		Count:          count,
		TokenId:        artworks.TokenId,
		TokenStandard:  artworks.TokenStandard,
		TokenType:      0,
		Owner:          artworks.Owner,
		BlindBox:       artworks.BlindBox,
		NftAddress:     artworks.NftAddress,
		BlockChain:     "Binance Smart Chain",
		ProtocolFee:    0,
		FixedPrice:     fixed_price,
		TotalPrice:     artworks.TotalPrice.Float64,
		Properties:     list,
		Status:         status,
		IdInContract:   artworks.AuctionId.String,
		ImageUrl:       artworks.ImageUrl.String,
		SaleAddress:    artworks.SaleAddress.String,
	}

	return &types.ArtworkInfoResponse{
		Code: 200,
		Msg:  "success",
		Data: *ArtworksInfo,
	}, nil
}

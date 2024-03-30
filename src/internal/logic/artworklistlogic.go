package logic

import (
	"context"
	"market/internal/common/errorx"
	"market/internal/svc"
	"market/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ArtworkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArtworkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) ArtworkListLogic {
	return ArtworkListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArtworkListLogic) ArtworkList(req types.ArtworkListReq) (*types.ArtworkListResponse, error) {
	artworks, count, err := l.svcCtx.NftTokenModel.Find(req)
	if err != nil {
		return nil, errorx.NewDefaultError("failed to fetch list")
	}
	list := make([]types.ArtworkListBody, 0)
	for _, artwork := range artworks {
		list = append(list, types.ArtworkListBody{
			ID:            artwork.Id,
			Name:          artwork.Name.String,
			Count:         artwork.Count,
			FixedPrice:    artwork.FixedPrice.Float64,
			HighestPrice:  artwork.HighestPrice.Float64,
			TokenID:       artwork.TokenId,
			ImageUrl:      artwork.ImageUrl.String,
			TokenStandard: artwork.TokenStandard.String,
			NftAddress:    artwork.NftAddress,
			IdInContract:  artwork.AuctionId.String,
			SaleAddress:   artwork.SaleAddress.String,
		})
	}

	return &types.ArtworkListResponse{
		Code:  200,
		Msg:   "success",
		Total: count,
		List:  list,
	}, nil
}

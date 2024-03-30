package logic

import (
	"context"
	"fmt"
	"market/internal/svc"
	"market/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type HistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) HistoryLogic {
	return HistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HistoryLogic) History(req types.HistoryReq) (*types.HistoryResponse, error) {
	list := make([]types.HistoryBody, 0)
	fmt.Println("req:", req)
	histories, total, err := l.svcCtx.NftSalesModel.Histories(req.Address, req.Type, req.PageNo, req.PageSize)
	if err != nil {
		return nil, err
	}
	if total > 0 {
		for _, history := range histories {
			Type := ""
			if req.Address == history.Owner {
				Type = "Sell"
			} else if req.Address == history.BuyerAddress.String {
				Type = "Buy"
			}
			Amount, _ := history.FixedPrice.Float64()
			list = append(list, types.HistoryBody{
				Id:         history.Id,
				NftTokenId: history.NftTokenId,
				Amount:     Amount,
				Type:       Type,
				Time:       history.UpdatedAt.Unix(),
				Count:      history.Count,
				Hash:       history.OperateTxHash.String,
				TokenName:  history.TokenName,
			})
		}
	}

	return &types.HistoryResponse{
		Code:  200,
		Msg:   "success",
		Total: total,
		List:  list,
	}, nil
	//artworks, count, err := l.svcCtx.NftTokenModel.Find(req)
	//if err != nil {
	//	return nil, errorx.NewDefaultError("failed to fetch list")
	//}
	//list := make([]types.ArtworkListBody, 0)
	//for _, artwork := range artworks {
	//	list = append(list, types.ArtworkListBody{
	//		ID:            artwork.Id,
	//		Name:          artwork.Name.String,
	//		Count:         artwork.Count,
	//		FixedPrice:    artwork.FixedPrice.Float64,
	//		HighestPrice:  artwork.HighestPrice.Float64,
	//		TokenID:       artwork.TokenId,
	//		ImageUrl:      artwork.ImageUrl.String,
	//		TokenStandard: artwork.TokenStandard.String,
	//		NftAddress:    artwork.NftAddress,
	//		IdInContract:  artwork.AuctionId.String,
	//		SaleAddress:   artwork.SaleAddress.String,
	//	})
	//}
	//
	//return &types.ArtworkListResponse{
	//	Code:  200,
	//	Msg:   "success",
	//	Total: count,
	//	List:  list,
	//}, nil
}

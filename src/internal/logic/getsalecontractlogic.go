package logic

import (
	"context"
	"market/internal/common/errorx"

	"market/internal/svc"
	"market/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetSaleContractLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSaleContractLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetSaleContractLogic {
	return GetSaleContractLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSaleContractLogic) GetSaleContract(req types.SaleContractReq) (*types.SaleContractInfoResponse, error) {
	saleContract, err := l.svcCtx.NftSaleContractsModel.FindActiveByCode(req.Code)
	if err != nil {
		return nil, errorx.NewDefaultError("can not find record")
	}

	saleContractInfoBody := types.SaleContractInfoBody{
		Code:            saleContract.Code,
		ContractAddress: saleContract.ContractAddress,
	}

	return &types.SaleContractInfoResponse{
		Data: saleContractInfoBody,
	}, nil
}

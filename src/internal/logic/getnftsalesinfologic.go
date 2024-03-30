package logic

import (
	"context"
	"encoding/json"
	"log"
	"market/internal/common/errorx"
	"market/internal/svc"
	"market/internal/types"

	"github.com/shopspring/decimal"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetNFTSalesInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNFTSalesInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetNFTSalesInfoLogic {
	return GetNFTSalesInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type Metadata struct {
	Image       string       `json:"image"`
	ExternalURL string       `json:"external_url"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Attributes  []Attributes `json:"attributes"`
}
type Attributes struct {
	Rarity  string `json:"Rarity"`
	Wisdom  string `json:"Wisdom"`
	Luck    string `json:"Luck"`
	Size    string `json:"Size"`
	Race    string `json:"Race"`
	Courage string `json:"Courage"`
	Stealth string `json:"Stealth"`
	Level   string `json:"Level"`
	Healthy string `json:"Healthy"`
}

func (l *GetNFTSalesInfoLogic) GetNFTSalesInfo(req types.NFTSalesInfoReq) (*types.NFTSalesInfoResponse, error) {
	sales, err := l.svcCtx.NftSalesModel.FindInfo(req.ID)
	if err != nil {
		return nil, errorx.NewDefaultError("can not find nft sales records")
	}

	nftToken,err  :=l.svcCtx.NftTokenModel.FindOne(sales.NftTokenId)
	if err != nil {
		return nil, errorx.NewDefaultError("can not find  nft token records")
		log.Println(err)
	}
	tb, err := l.svcCtx.TbNftTokenModel.FindOneByContractAddressTokenId(sales.NftAddress, nftToken.TokenId)
	attributesMap := make(map[string]string)
	if err != nil {
		log.Println(err)
	}
	if tb != nil {
		//解析json
		itemInfoBytes := []byte(tb.Metadata.String)
		var ItemInfo Metadata
		err = json.Unmarshal(itemInfoBytes, &ItemInfo)
		data, _ := json.Marshal(&ItemInfo.Attributes[0])
		json.Unmarshal(data, &attributesMap)
	}else {
		//没数据找tokens表
		itemInfoBytes := []byte(nftToken.Metadata.String)
		var ItemInfo Metadata
		err = json.Unmarshal(itemInfoBytes, &ItemInfo)
		if len(ItemInfo.Attributes) > 0{
			data, _ := json.Marshal(&ItemInfo.Attributes[0])
			json.Unmarshal(data, &attributesMap)
		}
	}

	list := make([]types.Property, 0)
	for key, Value := range attributesMap {
		list = append(list, types.Property{
			Key:   key,
			Value: Value,
		})
	}

	NFTSalesInfo := &types.NFTSalesInfoBody{
		ID:             sales.Id,
		Name:           sales.Name,
		Description:    sales.Description.String,
		CreatedAt:      sales.CreatedAt.String(),
		CategoryName:   sales.CategoryName,
		CollectionName: sales.CollectionName,
		ImageUrl:       sales.ImageUrl,
		Count:          sales.Count,
		FixedPrice:     sales.FixedPrice.DivRound(decimal.NewFromInt(sales.Count), 2),
		TotalPrice:     sales.FixedPrice,
		HighestPrice:   sales.HighestPrice,
		SaleAddress:    sales.SaleAddress,
		SaleType:       sales.SaleType,
		IdInContract:   sales.IdInContract,
		TokenId:        sales.TokenId,
		TokenStandard:  sales.TokenStandard,
		TokenType:      sales.TokenType,
		Owner:          sales.Owner,
		NftAddress:     sales.NftAddress,
		BlockChain:     sales.BlockChain,
		ProtocolFee:    sales.ProtocolFee,
		NeedCheck:      sales.NeedCheck,
		StartTime:      sales.StartTime.Int64,
		EndTime:        sales.EndTime.Int64,
		Status:         sales.Status,
		Properties:     list,
	}

	return &types.NFTSalesInfoResponse{
		Code: 200,
		Msg:  "success",
		Data: *NFTSalesInfo,
	}, nil
}

func (l *GetNFTSalesInfoLogic) NFTOnSalesToday(req types.NFTSalesStatisticsInfoReq) (*types.NFTOnSalesStatisticsResponse, error) {

	NTF_Cntract_Address := make(map[string]string, 0)
	NTF_Cntract_Address["MPB"] = "0x061c6eeca7b14cf4ec1b190dd879008dd7d7e470"
	NTF_Cntract_Address["BMM"] = "0x6efdd0380c9dde9c50ae99669d8ffd9439efcdbd"
	NTF_Cntract_Address["MML"] = "0xb5665e1038c4e17c58ab55b5c591fab52ce83fc4"
	NTF_Cntract_Address["MDB_YELLOW"] = "0x872028d15bd53d11561856909d4d3eaa4e846b64"
	NTF_Cntract_Address["MDB_PURPLE"] = "0x067c327cf85aa843cc52e89664f9fcb26bc68944"
	NTF_Cntract_Address["MDB_PINKPURPLE"] = "0x1f70bb2Ac9b51D5C2Ee32E76188d619d6EAF882a"
	NTF_Cntract_Address["SNDOG"] = "0x406f9a5779571e2d8abefb367cdc90d848b88471"
	NTF_Cntract_Address["MMA"] = "0x63872646b05f9094ec6e6de03042d31ce24457e1"
	NTF_Cntract_Address["MMP"] = "0xb9e879719c1271d044a5a3007797613478552d36"
	NTF_Cntract_Address["TESLA"] = "0xd4d63f37a13cbc99094afcf8187142af783e29b8"
	NTF_Cntract_Address["DING"] = "0x061c6eeca7b14cf4ec1b190dd879008dd7d7e470"
	NTF_Cntract_Address["MRM"] = "0x982B5345D0f213ecb2a8e6e24336909f59B1d6E3"
	NTF_Cntract_Address["NEW_YELLOW_DIAMOND"] = "0x5dc3FeD851e07715965E5727592CE33d14b7828D"
	NTF_Cntract_Address["NEW_POTION"] = "0x51353799F8550c9010a8b0CbFE6C02cA96E026E2"
	NTF_Cntract_Address["MME"] = "0x0cf6ec310531a65bc198452961b975db30eaf4ca"

	//单日挂单NFT量
	pendingOrdersNum, err := l.svcCtx.NftSalesModel.FindPendingOrdersNumToday(req.StartTime, req.EndTime)
	if err != nil {
		return nil, errorx.NewDefaultError("failed to fetch pending-NFT num")
	}
	//单日成交NFT量
	onSalesNFTNum, err := l.svcCtx.NftSalesModel.FindNFTDealNumToday(req.StartTime, req.EndTime)
	if err != nil {
		return nil, errorx.NewDefaultError("failed to fetch on-sales-NFT num")
	}
	//单日成交raca量
	onSalesRaacNum, err := l.svcCtx.NftSalesModel.FindRacaDealNumToday(req.StartTime, req.EndTime)
	if err != nil {
		return nil, errorx.NewDefaultError("failed to fetch on-sales-Raac num")
	}
	//细分单日成交raca量/NFT量

	list := make([]types.NTFProportion, 0)
	//查询所有合约列表
	for key, adress := range NTF_Cntract_Address {
		num, err := l.svcCtx.NftSalesModel.FindNFTTradeNumToday(adress, req.StartTime, req.EndTime)
		if err != nil {
			return nil, errorx.NewDefaultError("failed to fetch on-sales num")
		}
		proportion := types.NTFProportion{
			Name:       key,
			CntractNum: num,
			NtfNum:     onSalesNFTNum,
		}
		list = append(list, proportion)
	}
	return &types.NFTOnSalesStatisticsResponse{
		PendingOrdersNum: pendingOrdersNum,
		OnSalesNFTNum:    onSalesNFTNum,
		OnSalesRaacNum:   onSalesRaacNum,
		ProportionList:   list,
	}, nil
}

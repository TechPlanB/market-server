package svc

import (
	"log"
	"market/internal/config"
	"market/internal/model"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                  config.Config
	NftDropsModel           model.NftDropsModel
	NftSalesModel           model.NftSalesModel
	NftSaleContractsModel   model.NftSaleContractsModel
	NftChainsModel          model.NftChainsModel
	NftCategoriesModel      model.NftCategoriesModel
	NftCollectionsModel     model.NftCollectionsModel
	NftTransactionsModel    model.NftTransactionsModel
	NftListenersModel       model.NftListenersModel
	NftTokenModel           model.NftTokenModel
	NftTokenAttributesModel model.NftTokenAttributesModel
	TbNftTokenModel         model.TbNftTokenModel
	MinioClient             minio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	minioClient, err := minio.New(c.Minio.Url, &minio.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKey, c.Minio.AccessSecret, ""),
		Secure: c.Minio.UseSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	marketConn := sqlx.NewMysql(c.MarketDataSource)
	gameConn := sqlx.NewMysql(c.GameDataSource)
	return &ServiceContext{
		Config:                  c,
		NftDropsModel:           model.NewNftDropsModel(marketConn, c.CacheRedis),
		NftSalesModel:           model.NewNftSalesModel(marketConn, c.CacheRedis),
		NftSaleContractsModel:   model.NewNftSaleContractsModel(marketConn, c.CacheRedis),
		NftChainsModel:          model.NewNftChainsModel(marketConn, c.CacheRedis),
		NftCategoriesModel:      model.NewNftCategoriesModel(marketConn, c.CacheRedis),
		NftCollectionsModel:     model.NewNftCollectionsModel(marketConn, c.CacheRedis),
		NftTransactionsModel:    model.NewNftTransactionsModel(marketConn, c.CacheRedis),
		NftListenersModel:       model.NewNftListenersModel(marketConn, c.CacheRedis),
		NftTokenModel:           model.NewNftTokenModel(marketConn, c.CacheRedis),
		NftTokenAttributesModel: model.NewNftTokenAttributesModel(marketConn, c.CacheRedis),
		TbNftTokenModel:         model.NewTbNftTokenModel(gameConn, c.CacheRedis),
		MinioClient:             *minioClient,
	}
}

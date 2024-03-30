package config

import (
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/rest"
)

type Minio struct {
	Url          string
	AccessKey    string
	AccessSecret string
	UseSSL       bool
	BucketName   string
	OutUrl       string
}

type Chain struct {
	NetworkUrl                    string
	GetNFTTokenIdsContractAddress string
}

type Config struct {
	rest.RestConf
	MarketDataSource string
	GameDataSource string
	Minio      Minio
	CacheRedis cache.CacheConf
	Chain      Chain
}

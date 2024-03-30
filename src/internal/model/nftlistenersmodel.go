package model

import (
	"database/sql"
	"fmt"
	"market/internal/config"
	"market/internal/contracts/erc1155"
	"market/internal/contracts/erc721"
	"market/internal/contracts/utils"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/common/log"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	nftListenersFieldNames          = builderx.RawFieldNames(&NftListeners{})
	nftListenersRows                = strings.Join(nftListenersFieldNames, ",")
	nftListenersRowsExpectAutoSet   = strings.Join(stringx.Remove(nftListenersFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	nftListenersRowsWithPlaceHolder = strings.Join(stringx.Remove(nftListenersFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheNftListenersIdPrefix              = "cache:nftListeners:id:"
	cacheNftListenersContractAddressPrefix = "cache:nftListeners:contractAddress:"
)

type (
	NftListenersModel interface {
		Insert(data NftListeners) (sql.Result, error)
		FindOne(id int64) (*NftListeners, error)
		FindOneByContractAddress(contractAddress string) (*NftListeners, error)
		FindForSync() ([]NftListeners, error)
		Update(data NftListeners) error
		Delete(id int64) error
		SyncNFT(contractAddress string, tokenId int64, userAddress string, nftTokenModel NftTokenModel, chainConfig config.Chain) error
		Sync721WithoutTokenId(contractAddress string, userAddress string, nftTokenModel NftTokenModel, chainConfig config.Chain) error
	}

	defaultNftListenersModel struct {
		sqlc.CachedConn
		table string
	}

	NftListeners struct {
		UpdatedAt       time.Time `db:"updated_at"` // 更新时间
		Id              int64     `db:"id"`
		ContractAddress string    `db:"contract_address"` // 合约地址
		//ChainId         int64     `db:"chain_id"`         // NFT区块链表ID
		//NftType         string    `db:"nft_type"`         // NFT类别
		StartBlock     int64          `db:"start_block"`      // 开始同步区块高度
		StartBlockHash string         `db:"start_block_hash"` // 开始同步区块hash
		CreatedAt      time.Time      `db:"created_at"`       // 创建时间
		Name           sql.NullString `db:"name"`
		Symbol         sql.NullString `db:"symbol"`
		CategoryId     int64          `db:"category_id"`
		CollectionId   int64          `db:"collection_id"`
		NftType        sql.NullString `db:"nft_type"`
	}
)

func NewNftListenersModel(conn sqlx.SqlConn, c cache.CacheConf) NftListenersModel {
	return &defaultNftListenersModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`nft_listeners`",
	}
}

func (m *defaultNftListenersModel) Insert(data NftListeners) (sql.Result, error) {
	nftListenersContractAddressKey := fmt.Sprintf("%s%v", cacheNftListenersContractAddressPrefix, data.ContractAddress)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, nftListenersRowsExpectAutoSet)
		return conn.Exec(query, data.UpdatedAt, data.ContractAddress /*data.ChainId, data.NftType,*/, data.StartBlock, data.StartBlockHash, data.CreatedAt, data.Name, data.CategoryId, data.CollectionId, data.NftType)
	}, nftListenersContractAddressKey)
	return ret, err
}

func (m *defaultNftListenersModel) FindOne(id int64) (*NftListeners, error) {
	nftListenersIdKey := fmt.Sprintf("%s%v", cacheNftListenersIdPrefix, id)
	var resp NftListeners
	err := m.QueryRow(&resp, nftListenersIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftListenersRows, m.table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultNftListenersModel) FindForSync() ([]NftListeners, error) {
	query := fmt.Sprintf("select nl.* from %s nl where symbol != null or symbol != '' order by nl.name", m.table)

	var resp []NftListeners
	err := m.CachedConn.QueryRowsNoCache(&resp, query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return resp, err
}

func (m *defaultNftListenersModel) FindOneByContractAddress(contractAddress string) (*NftListeners, error) {
	var resp NftListeners
	query := fmt.Sprintf("select %s from %s where `contract_address` = ? limit 1", nftListenersRows, m.table)
	err := m.QueryRowNoCache(&resp, query, contractAddress)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultNftListenersModel) Update(data NftListeners) error {
	nftListenersIdKey := fmt.Sprintf("%s%v", cacheNftListenersIdPrefix, data.Id)
	nftListenersContractAddressKey := fmt.Sprintf("%s%v", cacheNftListenersContractAddressPrefix, data.ContractAddress)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, nftListenersRowsWithPlaceHolder)
		return conn.Exec(query, data.UpdatedAt, data.ContractAddress /*data.ChainId, data.NftType,*/, data.StartBlock, data.StartBlockHash, data.CreatedAt, data.Name, data.CategoryId, data.CollectionId, data.NftType, data.Id)
	}, nftListenersIdKey, nftListenersContractAddressKey)
	return err
}

func (m *defaultNftListenersModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	nftListenersIdKey := fmt.Sprintf("%s%v", cacheNftListenersIdPrefix, id)
	nftListenersContractAddressKey := fmt.Sprintf("%s%v", cacheNftListenersContractAddressPrefix, data.ContractAddress)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, nftListenersContractAddressKey, nftListenersIdKey)
	return err
}

func (m *defaultNftListenersModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheNftListenersIdPrefix, primary)
}

func (m *defaultNftListenersModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftListenersRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultNftListenersModel) SyncNFT(contractAddress string, tokenId int64, userAddress string, nftTokenModel NftTokenModel, chainConfig config.Chain) error {
	nftListeners, err := m.FindOneByContractAddress(contractAddress)
	if err != nil {
		log.Error("sync NFT: get NFT listeners", err)
		return err
	}

	if nftListeners.NftType.String == "BEP1155" {
		m.sync1155(contractAddress, tokenId, userAddress, nftListeners, nftTokenModel, chainConfig.NetworkUrl)
	} else {
		m.sync721WithTokenId(contractAddress, tokenId, nftListeners, nftTokenModel, chainConfig.NetworkUrl)
	}

	return nil
}

func (m *defaultNftListenersModel) sync721WithTokenId(contractAddress string, tokenId int64, nftListeners *NftListeners, nftTokenModel NftTokenModel, networkUrl string) error {
	client, err := ethclient.Dial(networkUrl)
	if err != nil {
		log.Error("sync NFT: dial to chain", err)
		return err
	}

	address := common.HexToAddress(contractAddress)
	instance, err := erc721.NewErc721(address, client)
	if err != nil {
		log.Error("sync NFT: get contract instance", err)
		return err
	}

	owner, err := instance.OwnerOf(&bind.CallOpts{}, big.NewInt(tokenId))
	if err != nil {
		log.Error("sync NFT: get NFT owner from chain", err)
		return err
	}

	nftToken, err := nftTokenModel.FindOneByNftListenerTokenId(nftListeners.Id, tokenId)
	switch err {
	case nil:
		if strings.ToLower(nftToken.Owner) != strings.ToLower(owner.String()) {
			nftToken.Owner = strings.ToLower(owner.String())
			err = nftTokenModel.Update(*nftToken)
			if err != nil {
				log.Error("sync NFT: update NFT owner", err)
				return err
			}
		}
	case sqlc.ErrNotFound:
		nftToken = &NftToken{
			Symbol:       nftListeners.Symbol,
			TokenId:      tokenId,
			Owner:        strings.ToLower(owner.String()),
			Count:        1,
			CategoryId:   nftListeners.CategoryId,
			CollectionId: nftListeners.CollectionId,
			NftListener:  nftListeners.Id,
		}
		_, err := nftTokenModel.Insert(*nftToken)
		if err != nil {
			log.Error("sync NFT: insert NFT token", err)
			return err
		}
	default:
		log.Error("sync NFT: find NFT token", err)
		return err
	}

	return nil
}

func (m *defaultNftListenersModel) sync1155(contractAddress string, tokenId int64, userAddress string, nftListeners *NftListeners, nftTokenModel NftTokenModel, networkUrl string) error {
	log.Infof(fmt.Sprintf("chain network: %s\n", networkUrl))
	client, err := ethclient.Dial(networkUrl)
	if err != nil {
		log.Error("sync NFT: dial to chain", err)
		return err
	}

	address := common.HexToAddress(contractAddress)
	instance, err := erc1155.NewErc1155(address, client)
	if err != nil {
		log.Error("sync NFT: get contract instance", err)
		return err
	}

	balanceOf, err := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(userAddress), big.NewInt(tokenId))
	if err != nil {
		log.Error("sync NFT: get NFT owner from chain", err)
		return err
	}

	nftToken, err := nftTokenModel.FindOneByNftListenerTokenIdOwner(nftListeners.Id, tokenId, strings.ToLower(userAddress))
	switch err {
	case nil:
		if nftToken.Count != balanceOf.Int64() {
			nftToken.Count = balanceOf.Int64()
			err = nftTokenModel.Update(*nftToken)
			if err != nil {
				log.Error("sync NFT: update NFT owner", err)
				return err
			}
		}
	case sqlc.ErrNotFound:
		nftToken = &NftToken{
			Symbol:       nftListeners.Symbol,
			TokenId:      tokenId,
			Owner:        strings.ToLower(userAddress),
			Count:        balanceOf.Int64(),
			CategoryId:   nftListeners.CategoryId,
			CollectionId: nftListeners.CollectionId,
			NftListener:  nftListeners.Id,
		}
		_, err := nftTokenModel.Insert(*nftToken)
		if err != nil {
			log.Error("sync NFT: insert NFT token", err)
			return err
		}
	default:
		log.Error("sync NFT: find NFT token", err)
		return err
	}

	return nil
}

func (m *defaultNftListenersModel) Sync721WithoutTokenId(contractAddress string, userAddress string, nftTokenModel NftTokenModel, chainConfig config.Chain) error {
	nftListeners, err := m.FindOneByContractAddress(contractAddress)
	if err != nil {
		log.Error("sync NFT: get NFT listeners", err)
		return err
	}

	go m.sync721WithoutTokenId(contractAddress, userAddress, nftListeners, nftTokenModel, chainConfig)

	return nil
}

func (m *defaultNftListenersModel) sync721WithoutTokenId(contractAddress string, userAddress string, nftListeners *NftListeners, nftTokenModel NftTokenModel, chainConfig config.Chain) error {
	client, err := ethclient.Dial(chainConfig.NetworkUrl)
	if err != nil {
		log.Error("sync NFT: dial to chain", err)
		return err
	}

	address := common.HexToAddress(chainConfig.GetNFTTokenIdsContractAddress)
	instance, err := utils.NewUtils(address, client)
	if err != nil {
		log.Error("sync NFT: get utils contract instance", err)
		return err
	}

	var cnt int64 = 0
	var pageNum int64 = 1
	var pageSize int64 = 500
	for cnt >= (pageNum-1)*pageSize {
		balanceBigInt, tokenIdsBigInt, err := instance.GetTokenIds(&bind.CallOpts{}, common.HexToAddress(contractAddress), common.HexToAddress(userAddress), big.NewInt(pageNum), big.NewInt(pageSize))
		if err != nil {
			log.Error("sync NFT: get token ids", err)
			return err
		}
		cnt = balanceBigInt.Int64()
		for _, tokenIdBigInt := range tokenIdsBigInt {
			m.syncNFTRec(userAddress, nftListeners, nftTokenModel, tokenIdBigInt.Int64())
		}
		pageNum++
	}

	return nil
}

func (m *defaultNftListenersModel) syncNFTRec(userAddress string, nftListeners *NftListeners, nftTokenModel NftTokenModel, tokenId int64) error {
	nftToken, err := nftTokenModel.FindOneByNftListenerTokenId(nftListeners.Id, tokenId)
	switch err {
	case nil:
		if strings.ToLower(nftToken.Owner) != strings.ToLower(userAddress) {
			nftToken.Owner = strings.ToLower(userAddress)
			err = nftTokenModel.Update(*nftToken)
			if err != nil {
				log.Error("sync NFT: update NFT owner", err)
				return err
			}
		}
	case sqlc.ErrNotFound:
		nftToken = &NftToken{
			Symbol:       nftListeners.Symbol,
			TokenId:      tokenId,
			Owner:        strings.ToLower(userAddress),
			Count:        1,
			CategoryId:   nftListeners.CategoryId,
			CollectionId: nftListeners.CollectionId,
			NftListener:  nftListeners.Id,
		}
		_, err := nftTokenModel.Insert(*nftToken)
		if err != nil {
			log.Error("sync NFT: insert NFT token", err)
			return err
		}
	default:
		log.Error("sync NFT: find NFT token", err)
		return err
	}
	return nil
}

package model

import (
	"database/sql"
	"fmt"
	"market/internal/types"
	"strings"
	"time"

	"github.com/prometheus/common/log"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	nftDropsFieldNames          = builderx.RawFieldNames(&NftDrops{})
	nftDropsRows                = strings.Join(nftDropsFieldNames, ",")
	nftDropsRowsExpectAutoSet   = strings.Join(stringx.Remove(nftDropsFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	nftDropsRowsWithPlaceHolder = strings.Join(stringx.Remove(nftDropsFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheNftDropsIdPrefix                     = "cache:nftDrops:id:"
	cacheNftDropsBuyAddressIdInContractPrefix = "cache:nftDrops:buyAddress:idInContract:"
)

type (
	NftDropsModel interface {
		Insert(data NftDrops) (sql.Result, error)
		FindOne(id int64) (*NftDrops, error)
		FindOneByBuyAddressIdInContract(buyAddress string, idInContract int64) (*NftDrops, error)
		Update(data NftDrops) error
		Delete(id int64) error
		FindInfo(id int64) (*NftDropsInfo, error)
		Find(filter types.NFTDropsReq) ([]NftDrops, int64, error)
	}

	NftDropsInfo struct {
		CollectionName string        `db:"collection_name"` // NFT对应收藏类名称
		ImageUrl       string        `db:"image_url"`       // NFT图片url
		StartTime      int64         `db:"start_time"`      // 开始时间-时间戳
		TokenId        sql.NullInt64 `db:"token_id"`        // BEP1155，对应NFT的tokenId
		Description    string        `db:"description"`     // 描述，支持富文本
		Price          float64       `db:"price"`           // 价格
		Deleted        int64         `db:"deleted"`         // 是否删除
		DeletedAt      sql.NullTime  `db:"deleted_at"`      // 删除日期
		UpdatedAt      time.Time     `db:"updated_at"`      // 更新时间
		Id             int64         `db:"id"`
		NftAddress     string        `db:"nft_address"`    // NFT对应合约地址
		NeedCheck      int64         `db:"need_check"`     // 是否需要检查条件（条件检查调用合约）
		Name           string        `db:"name"`           // NFT名称
		Count          int64         `db:"count"`          // 当前数量
		CategoryName   string        `db:"category_name"`  // NFT类别名称
		ChainId        int64         `db:"chain_id"`       // NFT区块链表ID
		EndTime        sql.NullInt64 `db:"end_time"`       // 结束时间-时间戳
		BuyAddress     string        `db:"buy_address"`    // NFT对应销售合约地址
		IdInContract   int64         `db:"id_in_contract"` // NFT商品在销售合约中的id
		ShowLeft       int64         `db:"show_left"`      // 是否显示剩余数量
		TokenType      int64         `db:"token_type"`     // 价格类型，1代表RACA，2代表BNB
		CreatedAt      time.Time     `db:"created_at"`     // 创建时间
	}

	defaultNftDropsModel struct {
		sqlc.CachedConn
		table string
	}

	NftDropsTotal struct {
		Total int64 `json:"total"`
	}

	NftDrops struct {
		CollectionId int64         `db:"collection_id"` // NFT对应收藏类ID
		ImageUrl     string        `db:"image_url"`     // NFT图片url
		StartTime    int64         `db:"start_time"`    // 开始时间-时间戳
		TokenId      sql.NullInt64 `db:"token_id"`      // BEP1155，对应NFT的tokenId
		Description  string        `db:"description"`   // 描述，支持富文本
		Price        float64       `db:"price"`         // 价格
		Deleted      int64         `db:"deleted"`       // 是否删除
		DeletedAt    sql.NullTime  `db:"deleted_at"`    // 删除日期
		UpdatedAt    time.Time     `db:"updated_at"`    // 更新时间
		Id           int64         `db:"id"`
		NftAddress   string        `db:"nft_address"`    // NFT对应合约地址
		NeedCheck    int64         `db:"need_check"`     // 是否需要检查条件（条件检查调用合约）
		Name         string        `db:"name"`           // NFT名称
		Count        int64         `db:"count"`          // 当前数量
		CategoryId   int64         `db:"category_id"`    // NFT类别ID
		ChainId      int64         `db:"chain_id"`       // NFT区块链表ID
		EndTime      sql.NullInt64 `db:"end_time"`       // 结束时间-时间戳
		BuyAddress   string        `db:"buy_address"`    // NFT对应销售合约地址
		IdInContract int64         `db:"id_in_contract"` // NFT商品在销售合约中的id
		ShowLeft     int64         `db:"show_left"`      // 是否显示剩余数量
		TokenType    int64         `db:"token_type"`     // 价格类型，1代表RACA，2代表BNB
		CreatedAt    time.Time     `db:"created_at"`     // 创建时间
	}
)

func NewNftDropsModel(conn sqlx.SqlConn, c cache.CacheConf) NftDropsModel {
	return &defaultNftDropsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`nft_drops`",
	}
}

func (m *defaultNftDropsModel) Insert(data NftDrops) (sql.Result, error) {
	nftDropsBuyAddressIdInContractKey := fmt.Sprintf("%s%v:%v", cacheNftDropsBuyAddressIdInContractPrefix, data.BuyAddress, data.IdInContract)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, nftDropsRowsExpectAutoSet)
		return conn.Exec(query, data.CollectionId, data.ImageUrl, data.StartTime, data.TokenId, data.Description, data.Price, data.Deleted, data.DeletedAt, data.UpdatedAt, data.NftAddress, data.NeedCheck, data.Name, data.Count, data.CategoryId, data.ChainId, data.EndTime, data.BuyAddress, data.IdInContract, data.ShowLeft, data.TokenType, data.CreatedAt)
	}, nftDropsBuyAddressIdInContractKey)
	return ret, err
}

func (m *defaultNftDropsModel) FindOne(id int64) (*NftDrops, error) {
	nftDropsIdKey := fmt.Sprintf("%s%v", cacheNftDropsIdPrefix, id)
	var resp NftDrops
	err := m.QueryRow(&resp, nftDropsIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftDropsRows, m.table)
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

func (m *defaultNftDropsModel) FindOneByBuyAddressIdInContract(buyAddress string, idInContract int64) (*NftDrops, error) {
	nftDropsBuyAddressIdInContractKey := fmt.Sprintf("%s%v:%v", cacheNftDropsBuyAddressIdInContractPrefix, buyAddress, idInContract)
	var resp NftDrops
	err := m.QueryRowIndex(&resp, nftDropsBuyAddressIdInContractKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `buy_address` = ? and `id_in_contract` = ? limit 1", nftDropsRows, m.table)
		if err := conn.QueryRow(&resp, query, buyAddress, idInContract); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultNftDropsModel) Update(data NftDrops) error {
	nftDropsBuyAddressIdInContractKey := fmt.Sprintf("%s%v:%v", cacheNftDropsBuyAddressIdInContractPrefix, data.BuyAddress, data.IdInContract)
	nftDropsIdKey := fmt.Sprintf("%s%v", cacheNftDropsIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, nftDropsRowsWithPlaceHolder)
		return conn.Exec(query, data.CollectionId, data.ImageUrl, data.StartTime, data.TokenId, data.Description, data.Price, data.Deleted, data.DeletedAt, data.UpdatedAt, data.NftAddress, data.NeedCheck, data.Name, data.Count, data.CategoryId, data.ChainId, data.EndTime, data.BuyAddress, data.IdInContract, data.ShowLeft, data.TokenType, data.CreatedAt, data.Id)
	}, nftDropsIdKey, nftDropsBuyAddressIdInContractKey)
	return err
}

func (m *defaultNftDropsModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	nftDropsIdKey := fmt.Sprintf("%s%v", cacheNftDropsIdPrefix, id)
	nftDropsBuyAddressIdInContractKey := fmt.Sprintf("%s%v:%v", cacheNftDropsBuyAddressIdInContractPrefix, data.BuyAddress, data.IdInContract)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, nftDropsIdKey, nftDropsBuyAddressIdInContractKey)
	return err
}

func (m *defaultNftDropsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheNftDropsIdPrefix, primary)
}

func (m *defaultNftDropsModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftDropsRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultNftDropsModel) Find(filter types.NFTDropsReq) ([]NftDrops, int64, error) {

	query := fmt.Sprintf("select %s from %s where 1=1 ", nftDropsRows, m.table)
	//query := fmt.Sprintf("select %s from %s where `name` like CONCAT('%%', ?, '%%')", nftDropsRows, m.table)
	countQuery := fmt.Sprintf("select count(1) from %s where 1=1 ", m.table)
	//countQuery := fmt.Sprintf("select count(1) from %s where `name` like CONCAT('%%', ?, '%%')", m.table)
	if filter.Name != "" {
		query = fmt.Sprintf("%s and `name` = '%s' ", query, filter.Name)
		countQuery = fmt.Sprintf("%s and `name` = '%s' ", query, filter.Name)
	}
	if filter.Category != 0 {
		query = fmt.Sprintf("%s and `category_id` = %d ", query, filter.Category)
		countQuery = fmt.Sprintf("%s and `category_id` = %d ", countQuery, filter.Category)
	}

	offset := (filter.PageNo - 1) * filter.PageSize
	size := filter.PageSize
	query = fmt.Sprintf("%s order by %s %s limit %d offset %d", query, filter.SortBy, filter.Order, size, offset)

	var totalResp NftDropsTotal
	err := m.CachedConn.QueryRowNoCache(&totalResp, countQuery)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]NftDrops, 0)
	if totalResp.Total > 0 {
		err = m.CachedConn.QueryRowsNoCache(&resp, query)
		if err != nil {
			return nil, 0, err
		}
	}

	return resp, totalResp.Total, err
}

func (m *defaultNftDropsModel) FindInfo(id int64) (*NftDropsInfo, error) {
	nftDropsIdKey := fmt.Sprintf("%s%v", cacheNftDropsIdPrefix, id)
	var resp NftDropsInfo
	err := m.QueryRow(&resp, nftDropsIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		leftJoin := "left join nft_categories nca on nca.id = nd.category_id left join nft_collections nco on nco.id = nd.collection_id"
		query := fmt.Sprintf("select nd.*,nca.name as category_name,nco.name as collection_name from %s nd %s", m.table, leftJoin)
		query = fmt.Sprintf("%s where nd.id = ? limit 1", query)
		return conn.QueryRow(v, query, id)
	})

	log.Info(err)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	nftSaleContractsFieldNames          = builderx.RawFieldNames(&NftSaleContracts{})
	nftSaleContractsRows                = strings.Join(nftSaleContractsFieldNames, ",")
	nftSaleContractsRowsExpectAutoSet   = strings.Join(stringx.Remove(nftSaleContractsFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	nftSaleContractsRowsWithPlaceHolder = strings.Join(stringx.Remove(nftSaleContractsFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheNftSaleContractsIdPrefix   = "cache:nftSaleContracts:id:"
	cacheNftSaleContractsCodePrefix = "cache:nftSaleContracts:code:"
)

type (
	NftSaleContractsModel interface {
		Insert(data NftSaleContracts) (sql.Result, error)
		FindOne(id int64) (*NftSaleContracts, error)
		Update(data NftSaleContracts) error
		Delete(id int64) error
		FindActiveByCode(code string) (*NftSaleContracts, error)
	}

	defaultNftSaleContractsModel struct {
		sqlc.CachedConn
		table string
	}

	NftSaleContracts struct {
		Id              int64     `db:"id"`
		Code            string    `db:"code"`             // 编码
		ContractAddress string    `db:"contract_address"` // 合约地址
		IsActive        int64     `db:"is_active"`        // 是否有效
		CreatedAt       time.Time `db:"created_at"`       // 创建时间
		UpdatedAt       time.Time `db:"updated_at"`       // 更新时间
	}
)

func NewNftSaleContractsModel(conn sqlx.SqlConn, c cache.CacheConf) NftSaleContractsModel {
	return &defaultNftSaleContractsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`nft_sale_contracts`",
	}
}

func (m *defaultNftSaleContractsModel) Insert(data NftSaleContracts) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, nftSaleContractsRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.Code, data.ContractAddress, data.IsActive, data.CreatedAt, data.UpdatedAt)

	return ret, err
}

func (m *defaultNftSaleContractsModel) FindOne(id int64) (*NftSaleContracts, error) {
	nftSaleContractsIdKey := fmt.Sprintf("%s%v", cacheNftSaleContractsIdPrefix, id)
	var resp NftSaleContracts
	err := m.QueryRow(&resp, nftSaleContractsIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftSaleContractsRows, m.table)
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

func (m *defaultNftSaleContractsModel) Update(data NftSaleContracts) error {
	nftSaleContractsIdKey := fmt.Sprintf("%s%v", cacheNftSaleContractsIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, nftSaleContractsRowsWithPlaceHolder)
		return conn.Exec(query, data.Code, data.ContractAddress, data.IsActive, data.CreatedAt, data.UpdatedAt, data.Id)
	}, nftSaleContractsIdKey)
	return err
}

func (m *defaultNftSaleContractsModel) Delete(id int64) error {

	nftSaleContractsIdKey := fmt.Sprintf("%s%v", cacheNftSaleContractsIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, nftSaleContractsIdKey)
	return err
}

func (m *defaultNftSaleContractsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheNftSaleContractsIdPrefix, primary)
}

func (m *defaultNftSaleContractsModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftSaleContractsRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultNftSaleContractsModel) FindActiveByCode(code string) (*NftSaleContracts, error) {
	nftSaleContractsCodeKey := fmt.Sprintf("%s%v", cacheNftSaleContractsCodePrefix, code)
	var resp NftSaleContracts
	err := m.QueryRow(&resp, nftSaleContractsCodeKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `code` = ? and `is_active` = 1 limit 1", nftSaleContractsRows, m.table)
		return conn.QueryRow(v, query, code)
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

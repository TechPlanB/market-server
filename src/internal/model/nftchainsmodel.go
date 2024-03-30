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
	nftChainsFieldNames          = builderx.RawFieldNames(&NftChains{})
	nftChainsRows                = strings.Join(nftChainsFieldNames, ",")
	nftChainsRowsExpectAutoSet   = strings.Join(stringx.Remove(nftChainsFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	nftChainsRowsWithPlaceHolder = strings.Join(stringx.Remove(nftChainsFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheNftChainsIdPrefix   = "cache:nftChains:id:"
	cacheNftChainsCodePrefix = "cache:nftChains:code:"
)

type (
	NftChainsModel interface {
		Insert(data NftChains) (sql.Result, error)
		FindOne(id int64) (*NftChains, error)
		FindOneByCode(code string) (*NftChains, error)
		Update(data NftChains) error
		Delete(id int64) error
	}

	defaultNftChainsModel struct {
		sqlc.CachedConn
		table string
	}

	NftChains struct {
		Name          string         `db:"name"`             // 区块链名称
		Code          string         `db:"code"`             // 区块链编码
		ChainId       int64          `db:"chain_id"`         // 区块链chain_id
		BackupRpcUrl0 sql.NullString `db:"backup_rpc_url_0"` // 备用地址0
		Deleted       int64          `db:"deleted"`          // 是否删除
		Id            int64          `db:"id"`
		RpcUrl        string         `db:"rpc_url"`          // 连接rpc地址
		BackupRpcUrl1 sql.NullString `db:"backup_rpc_url_1"` // 备用地址1
		BackupRpcUrl2 sql.NullString `db:"backup_rpc_url_2"` // 备用地址2
		BackupRpcUrl3 sql.NullString `db:"backup_rpc_url_3"` // 备用地址3
		DeletedAt     sql.NullTime   `db:"deleted_at"`       // 删除时间
		CreatedAt     time.Time      `db:"created_at"`       // 创建时间
		UpdatedAt     time.Time      `db:"updated_at"`       // 更新时间
	}
)

func NewNftChainsModel(conn sqlx.SqlConn, c cache.CacheConf) NftChainsModel {
	return &defaultNftChainsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`nft_chains`",
	}
}

func (m *defaultNftChainsModel) Insert(data NftChains) (sql.Result, error) {
	nftChainsCodeKey := fmt.Sprintf("%s%v", cacheNftChainsCodePrefix, data.Code)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, nftChainsRowsExpectAutoSet)
		return conn.Exec(query, data.Name, data.Code, data.ChainId, data.BackupRpcUrl0, data.Deleted, data.RpcUrl, data.BackupRpcUrl1, data.BackupRpcUrl2, data.BackupRpcUrl3, data.DeletedAt, data.CreatedAt, data.UpdatedAt)
	}, nftChainsCodeKey)
	return ret, err
}

func (m *defaultNftChainsModel) FindOne(id int64) (*NftChains, error) {
	nftChainsIdKey := fmt.Sprintf("%s%v", cacheNftChainsIdPrefix, id)
	var resp NftChains
	err := m.QueryRow(&resp, nftChainsIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftChainsRows, m.table)
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

func (m *defaultNftChainsModel) FindOneByCode(code string) (*NftChains, error) {
	nftChainsCodeKey := fmt.Sprintf("%s%v", cacheNftChainsCodePrefix, code)
	var resp NftChains
	err := m.QueryRowIndex(&resp, nftChainsCodeKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `code` = ? limit 1", nftChainsRows, m.table)
		if err := conn.QueryRow(&resp, query, code); err != nil {
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

func (m *defaultNftChainsModel) Update(data NftChains) error {
	nftChainsIdKey := fmt.Sprintf("%s%v", cacheNftChainsIdPrefix, data.Id)
	nftChainsCodeKey := fmt.Sprintf("%s%v", cacheNftChainsCodePrefix, data.Code)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, nftChainsRowsWithPlaceHolder)
		return conn.Exec(query, data.Name, data.Code, data.ChainId, data.BackupRpcUrl0, data.Deleted, data.RpcUrl, data.BackupRpcUrl1, data.BackupRpcUrl2, data.BackupRpcUrl3, data.DeletedAt, data.CreatedAt, data.UpdatedAt, data.Id)
	}, nftChainsIdKey, nftChainsCodeKey)
	return err
}

func (m *defaultNftChainsModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	nftChainsIdKey := fmt.Sprintf("%s%v", cacheNftChainsIdPrefix, id)
	nftChainsCodeKey := fmt.Sprintf("%s%v", cacheNftChainsCodePrefix, data.Code)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, nftChainsCodeKey, nftChainsIdKey)
	return err
}

func (m *defaultNftChainsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheNftChainsIdPrefix, primary)
}

func (m *defaultNftChainsModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftChainsRows, m.table)
	return conn.QueryRow(v, query, primary)
}

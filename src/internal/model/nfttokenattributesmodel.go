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
	nftTokenAttributesFieldNames          = builderx.RawFieldNames(&NftTokenAttributes{})
	nftTokenAttributesRows                = strings.Join(nftTokenAttributesFieldNames, ",")
	nftTokenAttributesRowsExpectAutoSet   = strings.Join(stringx.Remove(nftTokenAttributesFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	nftTokenAttributesRowsWithPlaceHolder = strings.Join(stringx.Remove(nftTokenAttributesFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheNftTokenAttributesIdPrefix = "cache:nftTokenAttributes:id:"
)

type (
	NftTokenAttributesModel interface {
		Insert(data NftTokenAttributes) (sql.Result, error)
		FindOne(id int64) (*NftTokenAttributes, error)
		Update(data NftTokenAttributes) error
		Delete(id int64) error
	}

	defaultNftTokenAttributesModel struct {
		sqlc.CachedConn
		table string
	}

	NftTokenAttributes struct {
		CreatedAt  time.Time      `db:"created_at"` // 创建时间
		UpdatedAt  time.Time      `db:"updated_at"` // 更新时间
		Id         int64          `db:"id"`
		NftTokenId int64          `db:"nft_token_id"` // nft_token表id
		Key        string         `db:"key"`          // 属性key
		Value      sql.NullString `db:"value"`        // 属性value
	}
)

func NewNftTokenAttributesModel(conn sqlx.SqlConn, c cache.CacheConf) NftTokenAttributesModel {
	return &defaultNftTokenAttributesModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`nft_token_attributes`",
	}
}

func (m *defaultNftTokenAttributesModel) Insert(data NftTokenAttributes) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, nftTokenAttributesRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.CreatedAt, data.UpdatedAt, data.NftTokenId, data.Key, data.Value)

	return ret, err
}

func (m *defaultNftTokenAttributesModel) FindOne(id int64) (*NftTokenAttributes, error) {
	nftTokenAttributesIdKey := fmt.Sprintf("%s%v", cacheNftTokenAttributesIdPrefix, id)
	var resp NftTokenAttributes
	err := m.QueryRow(&resp, nftTokenAttributesIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftTokenAttributesRows, m.table)
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

func (m *defaultNftTokenAttributesModel) Update(data NftTokenAttributes) error {
	nftTokenAttributesIdKey := fmt.Sprintf("%s%v", cacheNftTokenAttributesIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, nftTokenAttributesRowsWithPlaceHolder)
		return conn.Exec(query, data.CreatedAt, data.UpdatedAt, data.NftTokenId, data.Key, data.Value, data.Id)
	}, nftTokenAttributesIdKey)
	return err
}

func (m *defaultNftTokenAttributesModel) Delete(id int64) error {

	nftTokenAttributesIdKey := fmt.Sprintf("%s%v", cacheNftTokenAttributesIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, nftTokenAttributesIdKey)
	return err
}

func (m *defaultNftTokenAttributesModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheNftTokenAttributesIdPrefix, primary)
}

func (m *defaultNftTokenAttributesModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftTokenAttributesRows, m.table)
	return conn.QueryRow(v, query, primary)
}

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
	nftCollectionsFieldNames          = builderx.RawFieldNames(&NftCollections{})
	nftCollectionsRows                = strings.Join(nftCollectionsFieldNames, ",")
	nftCollectionsRowsExpectAutoSet   = strings.Join(stringx.Remove(nftCollectionsFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	nftCollectionsRowsWithPlaceHolder = strings.Join(stringx.Remove(nftCollectionsFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheNftCollectionsIdPrefix   = "cache:nftCollections:id:"
	cacheNftCollectionsCodePrefix = "cache:nftCollections:code:"
)

type (
	NftCollectionsModel interface {
		Insert(data NftCollections) (sql.Result, error)
		FindOne(id int64) (*NftCollections, error)
		FindOneByCode(code string) (*NftCollections, error)
		Update(data NftCollections) error
		Delete(id int64) error
	}

	defaultNftCollectionsModel struct {
		sqlc.CachedConn
		table string
	}

	NftCollections struct {
		Description    sql.NullString `db:"description"` // 描述
		CreatedAt      time.Time      `db:"created_at"`  // 创建时间
		UpdatedAt      time.Time      `db:"updated_at"`  // 更新时间
		Deleted        int64          `db:"deleted"`     // 是否删除
		DeletedAt      sql.NullTime   `db:"deleted_at"`  // 删除时间
		Id             int64          `db:"id"`
		Name           string         `db:"name"`            // 收藏家名称
		Code           string         `db:"code"`            // 收藏家编码
		AccountAddress sql.NullString `db:"account_address"` // 对应账户地址
	}
)

func NewNftCollectionsModel(conn sqlx.SqlConn, c cache.CacheConf) NftCollectionsModel {
	return &defaultNftCollectionsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`nft_collections`",
	}
}

func (m *defaultNftCollectionsModel) Insert(data NftCollections) (sql.Result, error) {
	nftCollectionsCodeKey := fmt.Sprintf("%s%v", cacheNftCollectionsCodePrefix, data.Code)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, nftCollectionsRowsExpectAutoSet)
		return conn.Exec(query, data.Description, data.CreatedAt, data.UpdatedAt, data.Deleted, data.DeletedAt, data.Name, data.Code, data.AccountAddress)
	}, nftCollectionsCodeKey)
	return ret, err
}

func (m *defaultNftCollectionsModel) FindOne(id int64) (*NftCollections, error) {
	nftCollectionsIdKey := fmt.Sprintf("%s%v", cacheNftCollectionsIdPrefix, id)
	var resp NftCollections
	err := m.QueryRow(&resp, nftCollectionsIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftCollectionsRows, m.table)
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

func (m *defaultNftCollectionsModel) FindOneByCode(code string) (*NftCollections, error) {
	nftCollectionsCodeKey := fmt.Sprintf("%s%v", cacheNftCollectionsCodePrefix, code)
	var resp NftCollections
	err := m.QueryRowIndex(&resp, nftCollectionsCodeKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `code` = ? limit 1", nftCollectionsRows, m.table)
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

func (m *defaultNftCollectionsModel) Update(data NftCollections) error {
	nftCollectionsCodeKey := fmt.Sprintf("%s%v", cacheNftCollectionsCodePrefix, data.Code)
	nftCollectionsIdKey := fmt.Sprintf("%s%v", cacheNftCollectionsIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, nftCollectionsRowsWithPlaceHolder)
		return conn.Exec(query, data.Description, data.CreatedAt, data.UpdatedAt, data.Deleted, data.DeletedAt, data.Name, data.Code, data.AccountAddress, data.Id)
	}, nftCollectionsIdKey, nftCollectionsCodeKey)
	return err
}

func (m *defaultNftCollectionsModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	nftCollectionsIdKey := fmt.Sprintf("%s%v", cacheNftCollectionsIdPrefix, id)
	nftCollectionsCodeKey := fmt.Sprintf("%s%v", cacheNftCollectionsCodePrefix, data.Code)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, nftCollectionsIdKey, nftCollectionsCodeKey)
	return err
}

func (m *defaultNftCollectionsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheNftCollectionsIdPrefix, primary)
}

func (m *defaultNftCollectionsModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftCollectionsRows, m.table)
	return conn.QueryRow(v, query, primary)
}

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
	nftCategoriesFieldNames          = builderx.RawFieldNames(&NftCategories{})
	nftCategoriesRows                = strings.Join(nftCategoriesFieldNames, ",")
	nftCategoriesRowsExpectAutoSet   = strings.Join(stringx.Remove(nftCategoriesFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	nftCategoriesRowsWithPlaceHolder = strings.Join(stringx.Remove(nftCategoriesFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheNftCategoriesIdPrefix   = "cache:nftCategories:id:"
	cacheNftCategoriesCodePrefix = "cache:nftCategories:code:"
)

type (
	NftCategoriesModel interface {
		Insert(data NftCategories) (sql.Result, error)
		FindOne(id int64) (*NftCategories, error)
		FindOneByCode(code string) (*NftCategories, error)
		Update(data NftCategories) error
		Delete(id int64) error
		Find() ([]NftCategories, error)
	}

	defaultNftCategoriesModel struct {
		sqlc.CachedConn
		table string
	}

	NftCategories struct {
		DeletedAt sql.NullTime `db:"deleted_at"` // 删除时间
		CreatedAt time.Time    `db:"created_at"` // 创建时间
		UpdatedAt time.Time    `db:"updated_at"` // 更新时间
		Id        int64        `db:"id"`
		Name      string       `db:"name"`    // 类别名称
		Code      string       `db:"code"`    // 类别编码
		Deleted   int64        `db:"deleted"` // 是否删除
	}
)

func NewNftCategoriesModel(conn sqlx.SqlConn, c cache.CacheConf) NftCategoriesModel {
	return &defaultNftCategoriesModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`nft_categories`",
	}
}

func (m *defaultNftCategoriesModel) Insert(data NftCategories) (sql.Result, error) {
	nftCategoriesCodeKey := fmt.Sprintf("%s%v", cacheNftCategoriesCodePrefix, data.Code)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, nftCategoriesRowsExpectAutoSet)
		return conn.Exec(query, data.DeletedAt, data.CreatedAt, data.UpdatedAt, data.Name, data.Code, data.Deleted)
	}, nftCategoriesCodeKey)
	return ret, err
}

func (m *defaultNftCategoriesModel) FindOne(id int64) (*NftCategories, error) {
	nftCategoriesIdKey := fmt.Sprintf("%s%v", cacheNftCategoriesIdPrefix, id)
	var resp NftCategories
	err := m.QueryRow(&resp, nftCategoriesIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftCategoriesRows, m.table)
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

func (m *defaultNftCategoriesModel) FindOneByCode(code string) (*NftCategories, error) {
	nftCategoriesCodeKey := fmt.Sprintf("%s%v", cacheNftCategoriesCodePrefix, code)
	var resp NftCategories
	err := m.QueryRowIndex(&resp, nftCategoriesCodeKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `code` = ? limit 1", nftCategoriesRows, m.table)
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

func (m *defaultNftCategoriesModel) Update(data NftCategories) error {
	nftCategoriesIdKey := fmt.Sprintf("%s%v", cacheNftCategoriesIdPrefix, data.Id)
	nftCategoriesCodeKey := fmt.Sprintf("%s%v", cacheNftCategoriesCodePrefix, data.Code)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, nftCategoriesRowsWithPlaceHolder)
		return conn.Exec(query, data.DeletedAt, data.CreatedAt, data.UpdatedAt, data.Name, data.Code, data.Deleted, data.Id)
	}, nftCategoriesIdKey, nftCategoriesCodeKey)
	return err
}

func (m *defaultNftCategoriesModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	nftCategoriesIdKey := fmt.Sprintf("%s%v", cacheNftCategoriesIdPrefix, id)
	nftCategoriesCodeKey := fmt.Sprintf("%s%v", cacheNftCategoriesCodePrefix, data.Code)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, nftCategoriesIdKey, nftCategoriesCodeKey)
	return err
}

func (m *defaultNftCategoriesModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheNftCategoriesIdPrefix, primary)
}

func (m *defaultNftCategoriesModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftCategoriesRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultNftCategoriesModel) Find() ([]NftCategories, error) {
	query := fmt.Sprintf("select %s from %s", nftCategoriesRows, m.table)
	var resp []NftCategories
	err := m.CachedConn.QueryRowsNoCache(&resp, query)
	return resp, err
}

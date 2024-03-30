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
	nftTransactionsFieldNames          = builderx.RawFieldNames(&NftTransactions{})
	nftTransactionsRows                = strings.Join(nftTransactionsFieldNames, ",")
	nftTransactionsRowsExpectAutoSet   = strings.Join(stringx.Remove(nftTransactionsFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	nftTransactionsRowsWithPlaceHolder = strings.Join(stringx.Remove(nftTransactionsFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheNftTransactionsIdPrefix = "cache:nftTransactions:id:"
)

type (
	NftTransactionsModel interface {
		Insert(data NftTransactions) (sql.Result, error)
		FindOne(id int64) (*NftTransactions, error)
		Update(data NftTransactions) error
		Delete(id int64) error
	}

	defaultNftTransactionsModel struct {
		sqlc.CachedConn
		table string
	}

	NftTransactions struct {
		Id          int64     `db:"id"`
		NftListener int64     `db:"nft_listener"` // nft_token表id
		To          string    `db:"to"`           // token转移to
		BlockNumber int64     `db:"block_number"` // 区块高度
		CreatedAt   time.Time `db:"created_at"`   // 创建时间
		TokenId     int64     `db:"token_id"`     // NFT对应的tokenId
		From        string    `db:"from"`         // token转移from
		TxTime      time.Time `db:"tx_time"`      // 交易时间
		UpdatedAt   time.Time `db:"updated_at"`   // 更新时间
	}
)

func NewNftTransactionsModel(conn sqlx.SqlConn, c cache.CacheConf) NftTransactionsModel {
	return &defaultNftTransactionsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`nft_transactions`",
	}
}

func (m *defaultNftTransactionsModel) Insert(data NftTransactions) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, nftTransactionsRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.NftListener, data.To, data.BlockNumber, data.CreatedAt, data.TokenId, data.From, data.TxTime, data.UpdatedAt)

	return ret, err
}

func (m *defaultNftTransactionsModel) FindOne(id int64) (*NftTransactions, error) {
	nftTransactionsIdKey := fmt.Sprintf("%s%v", cacheNftTransactionsIdPrefix, id)
	var resp NftTransactions
	err := m.QueryRow(&resp, nftTransactionsIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftTransactionsRows, m.table)
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

func (m *defaultNftTransactionsModel) Update(data NftTransactions) error {
	nftTransactionsIdKey := fmt.Sprintf("%s%v", cacheNftTransactionsIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, nftTransactionsRowsWithPlaceHolder)
		return conn.Exec(query, data.NftListener, data.To, data.BlockNumber, data.CreatedAt, data.TokenId, data.From, data.TxTime, data.UpdatedAt, data.Id)
	}, nftTransactionsIdKey)
	return err
}

func (m *defaultNftTransactionsModel) Delete(id int64) error {

	nftTransactionsIdKey := fmt.Sprintf("%s%v", cacheNftTransactionsIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, nftTransactionsIdKey)
	return err
}

func (m *defaultNftTransactionsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheNftTransactionsIdPrefix, primary)
}

func (m *defaultNftTransactionsModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftTransactionsRows, m.table)
	return conn.QueryRow(v, query, primary)
}

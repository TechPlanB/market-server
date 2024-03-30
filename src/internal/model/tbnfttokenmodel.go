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
	tbNftTokenFieldNames          = builderx.RawFieldNames(&TbNftToken{})
	tbNftTokenRows                = strings.Join(tbNftTokenFieldNames, ",")
	tbNftTokenRowsExpectAutoSet   = strings.Join(stringx.Remove(tbNftTokenFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	tbNftTokenRowsWithPlaceHolder = strings.Join(stringx.Remove(tbNftTokenFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheTbNftTokenIdPrefix                     = "cache:tbNftToken:id:"
	cacheTbNftTokenContractAddressTokenIdPrefix = "cache:tbNftToken:contractAddress:tokenId:"
)

type (
	TbNftTokenModel interface {
		Insert(data TbNftToken) (sql.Result, error)
		FindOne(id int64) (*TbNftToken, error)
		FindOneByContractAddressTokenId(contractAddress string, tokenId int64) (*TbNftToken, error)
		Update(data TbNftToken) error
		Delete(id int64) error
	}

	defaultTbNftTokenModel struct {
		sqlc.CachedConn
		table string
	}

	TbNftToken struct {
		Id              int64          `db:"id"`
		Symbol          string         `db:"symbol"`           // NFT symbol
		ContractAddress string         `db:"contract_address"` // NFT contract address
		TokenId         int64          `db:"token_id"`         // NFT tokenId
		Metadata        sql.NullString `db:"metadata"`         // NFT Metadata
		Name            sql.NullString `db:"name"`             // NFT Name
		Description     sql.NullString `db:"description"`      // NFT Description
		ImageUrl        sql.NullString `db:"image_url"`        // IMAGE URL
		Owner           string         `db:"owner"`            // NFT Owner
		CreatedAt       time.Time      `db:"created_at"`       // created time
		UpdatedAt       time.Time      `db:"updated_at"`       // updated time
		Status          int64          `db:"status"`           // 冻结状态
	}
)

func NewTbNftTokenModel(conn sqlx.SqlConn, c cache.CacheConf) TbNftTokenModel {
	return &defaultTbNftTokenModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`tb_nft_token`",
	}
}

func (m *defaultTbNftTokenModel) Insert(data TbNftToken) (sql.Result, error) {
	tbNftTokenContractAddressTokenIdKey := fmt.Sprintf("%s%v:%v", cacheTbNftTokenContractAddressTokenIdPrefix, data.ContractAddress, data.TokenId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, tbNftTokenRowsExpectAutoSet)
		return conn.Exec(query, data.Symbol, data.ContractAddress, data.TokenId, data.Metadata, data.Name, data.Description, data.ImageUrl, data.Owner, data.CreatedAt, data.UpdatedAt, data.Status)
	}, tbNftTokenContractAddressTokenIdKey)
	return ret, err
}

func (m *defaultTbNftTokenModel) FindOne(id int64) (*TbNftToken, error) {
	tbNftTokenIdKey := fmt.Sprintf("%s%v", cacheTbNftTokenIdPrefix, id)
	var resp TbNftToken
	err := m.QueryRow(&resp, tbNftTokenIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tbNftTokenRows, m.table)
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

func (m *defaultTbNftTokenModel) FindOneByContractAddressTokenId(contractAddress string, tokenId int64) (*TbNftToken, error) {
	tbNftTokenContractAddressTokenIdKey := fmt.Sprintf("%s%v:%v", cacheTbNftTokenContractAddressTokenIdPrefix, contractAddress, tokenId)
	var resp TbNftToken
	err := m.QueryRowIndex(&resp, tbNftTokenContractAddressTokenIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `contract_address` = ? and `token_id` = ? limit 1", tbNftTokenRows, m.table)
		if err := conn.QueryRow(&resp, query, contractAddress, tokenId); err != nil {
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

func (m *defaultTbNftTokenModel) Update(data TbNftToken) error {
	tbNftTokenIdKey := fmt.Sprintf("%s%v", cacheTbNftTokenIdPrefix, data.Id)
	tbNftTokenContractAddressTokenIdKey := fmt.Sprintf("%s%v:%v", cacheTbNftTokenContractAddressTokenIdPrefix, data.ContractAddress, data.TokenId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tbNftTokenRowsWithPlaceHolder)
		return conn.Exec(query, data.Symbol, data.ContractAddress, data.TokenId, data.Metadata, data.Name, data.Description, data.ImageUrl, data.Owner, data.CreatedAt, data.UpdatedAt, data.Status, data.Id)
	}, tbNftTokenIdKey, tbNftTokenContractAddressTokenIdKey)
	return err
}

func (m *defaultTbNftTokenModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	tbNftTokenIdKey := fmt.Sprintf("%s%v", cacheTbNftTokenIdPrefix, id)
	tbNftTokenContractAddressTokenIdKey := fmt.Sprintf("%s%v:%v", cacheTbNftTokenContractAddressTokenIdPrefix, data.ContractAddress, data.TokenId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, tbNftTokenIdKey, tbNftTokenContractAddressTokenIdKey)
	return err
}

func (m *defaultTbNftTokenModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheTbNftTokenIdPrefix, primary)
}

func (m *defaultTbNftTokenModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tbNftTokenRows, m.table)
	return conn.QueryRow(v, query, primary)
}


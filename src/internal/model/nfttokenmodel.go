package model

import (
	"database/sql"
	"errors"
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
	nftTokenFieldNames          = builderx.RawFieldNames(&NftToken{})
	nftTokenRows                = strings.Join(nftTokenFieldNames, ",")
	nftTokenRowsExpectAutoSet   = strings.Join(stringx.Remove(nftTokenFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	nftTokenRowsWithPlaceHolder = strings.Join(stringx.Remove(nftTokenFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheNftTokenIdPrefix                 = "cache:nftToken:id:"
	cacheNftTokenNftListenerTokenIdPrefix = "cache:nftToken:nftListener:tokenId:"
)

type (
	NftTokenModel interface {
		Insert(data NftToken) (sql.Result, error)
		FindOne(id int64) (*NftToken, error)
		FindOneByNftListenerTokenId(nftListener int64, tokenId int64) (*NftToken, error)
		FindOneByNftListenerTokenIdOwner(nftListener int64, tokenId int64, owner string) (*NftToken, error)
		Update(data NftToken) error
		Delete(id int64) error
		Find(filter types.ArtworkListReq) ([]NftTokenList, int64, error)
		FindInfo(id int64, id_in_contract string) (*NftTokenInfo, error)
		FindInfoByContractAndToken(contractAddress string, tokenId int64) (*MyNftToken, error)
		FindProperties(id int64) ([]NftTokenProperties, error)
		FindTokenIds(owner string, contractAddress string) ([]int, error)
		FindTokenIdsByName(owner string, contractAddress string, name string) ([]int, error)
	}

	ArtworkListTotal struct {
		Total int64 `json:"total"`
	}

	defaultNftTokenModel struct {
		sqlc.CachedConn
		table string
	}

	NftTokenProperties struct {
		Key   string         `db:"key"`   // 属性key
		Value sql.NullString `db:"value"` // 属性value
	}

	NftTokenList struct {
		Id            int64           `db:"id"`
		Symbol        sql.NullString  `db:"symbol"`        // 符号
		CreatedAt     time.Time       `db:"created_at"`    // 创建时间
		Description   sql.NullString  `db:"description"`   // 元数据中的description
		TokenId       int64           `db:"token_id"`      // NFT tokenId
		AuctionId     sql.NullString  `db:"auction_id"`    // Auction Id
		NftAddress    string          `db:"nft_address"`   // 合约地址
		Owner         string          `db:"owner"`         // NFT拥有者address
		Count         int64           `db:"count"`         // 数量,BEP721对应数量默认为1
		FixedPrice    sql.NullFloat64 `db:"fixed_price"`   // 出售价格
		HighestPrice  sql.NullFloat64 `db:"highest_price"` // 拍卖最高出价
		Metadata      sql.NullString  `db:"metadata"`      // 元数据
		Name          sql.NullString  `db:"name"`          // 元数据中的name
		CategoryId    int64           `db:"category_id"`   // 类别id
		CollectionId  int64           `db:"collection_id"` // 收藏家id
		NftListener   int64           `db:"nft_listener"`  // 监听id
		ImageUrl      sql.NullString  `db:"image_url"`     // 元数据中的图片url
		ExternalUrl   sql.NullString  `db:"external_url"`  // 元数据中的external_url
		UpdatedAt     time.Time       `db:"updated_at"`    // 更新时间
		SaleAddress   sql.NullString  `db:"sale_address"`
		TokenStandard sql.NullString  `db:"token_standard"`
	}

	MyNftToken struct {
		Name        sql.NullString `db:"name"`        // 元数据中的name
		ImageUrl    sql.NullString `db:"image_url"`   // 元数据中的图片url
		Description sql.NullString `db:"description"` // 元数据中的description
	}

	NftTokenInfo struct {
		Id             int64           `db:"id"`
		Symbol         sql.NullString  `db:"symbol"`          // 符号
		CreatedAt      time.Time       `db:"created_at"`      // 创建时间
		Description    sql.NullString  `db:"description"`     // 元数据中的description
		CategoryName   string          `db:"category_name"`   // NFT类别名称
		CollectionName string          `db:"collection_name"` // NFT对应收藏类名称
		BlockChain     string          `db:"block_chain"`     // 链信息
		TokenType      int64           `db:"token_type"`      // Token类型
		TokenStandard  string          `db:"token_standard"`  // Token标准
		TokenId        int64           `db:"token_id"`        // NFT tokenId
		Owner          string          `db:"owner"`           // NFT拥有者address
		BlindBox       bool            `db:"blind_box"`       // 是否是盲盒 0：否 1：是
		Count          int64           `db:"count"`           // 数量,BEP721对应数量默认为1
		Metadata       sql.NullString  `db:"metadata"`        // 元数据
		Name           sql.NullString  `db:"name"`            // 元数据中的name
		CategoryId     int64           `db:"category_id"`     // 类别id
		CollectionId   int64           `db:"collection_id"`   // 收藏家id
		NftAddress     string          `db:"nft_address"`     // 合约地址
		ProtocolFee    float64         `db:"protocol_fee"`
		NftListener    int64           `db:"nft_listener"` // 监听id
		ImageUrl       sql.NullString  `db:"image_url"`    // 元数据中的图片url
		ExternalUrl    sql.NullString  `db:"external_url"` // 元数据中的external_url
		UpdatedAt      time.Time       `db:"updated_at"`   // 更新时间
		FixedPrice     sql.NullFloat64 `db:"fixed_price"`  // 出售价格
		TotalPrice     sql.NullFloat64 `db:"total_price"`
		Status         sql.NullString  `db:"status"`
		AuctionId      sql.NullString  `db:"auction_id"` // Auction Id
		SaleCount      sql.NullInt64   `db:"sale_count"`
		SaleAddress    sql.NullString  `db:"sale_address"`
	}

	NftToken struct {
		Id           int64          `db:"id"`
		Symbol       sql.NullString `db:"symbol"` // 符号
		CreatedAt    time.Time      `db:"created_at" gorm:"autoCreateTime"`
		UpdatedAt    time.Time      `db:"updated_at" gorm:"autoUpdateTime"`
		Description  sql.NullString `db:"description"`   // 元数据中的description
		TokenId      int64          `db:"token_id"`      // NFT tokenId
		Owner        string         `db:"owner"`         // NFT拥有者address
		Count        int64          `db:"count"`         // 数量,BEP721对应数量默认为1
		Metadata     sql.NullString `db:"metadata"`      // 元数据
		Name         sql.NullString `db:"name"`          // 元数据中的name
		CategoryId   int64          `db:"category_id"`   // 类别id
		CollectionId int64          `db:"collection_id"` // 收藏家id
		NftListener  int64          `db:"nft_listener"`  // 监听id
		ImageUrl     sql.NullString `db:"image_url"`     // 元数据中的图片url
		ExternalUrl  sql.NullString `db:"external_url"`  // 元数据中的external_url
	}
)

func NewNftTokenModel(conn sqlx.SqlConn, c cache.CacheConf) NftTokenModel {
	return &defaultNftTokenModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`nft_tokens`",
	}
}

func (m *defaultNftTokenModel) FindInfoByContractAndToken(contractAddress string, tokenId int64) (*MyNftToken, error) {
	var resp MyNftToken
	sqlStr := `
		select nt.name, nt.description, nt.image_url
		from %s nt
         inner join %s nl on nt.nft_listener = nl.id
		where nt.token_id=? and nl.contract_address=?
	`

	query := fmt.Sprintf(sqlStr, TABLE_NAME_NFT_TOKENS, TABLE_NAME_NFT_LISTENERS)
	var err error
	err = m.QueryRowNoCache(&resp, query, tokenId, contractAddress)

	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultNftTokenModel) Insert(data NftToken) (sql.Result, error) {
	nftTokenNftListenerTokenIdKey := fmt.Sprintf("%s%v:%v", cacheNftTokenNftListenerTokenIdPrefix, data.NftListener, data.TokenId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, nftTokenRowsExpectAutoSet)
		return conn.Exec(query, data.Symbol, data.CreatedAt, data.Description, data.TokenId, data.Owner, data.Count, data.Metadata, data.Name, data.CategoryId, data.CollectionId, data.NftListener, data.ImageUrl, data.ExternalUrl, data.UpdatedAt)
	}, nftTokenNftListenerTokenIdKey)
	return ret, err
}

func (m *defaultNftTokenModel) FindOne(id int64) (*NftToken, error) {
	//nftTokenIdKey := fmt.Sprintf("%s%v", cacheNftTokenIdPrefix, id)
	var resp NftToken
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftTokenRows, m.table)
	err := m.QueryRowNoCache(&resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultNftTokenModel) FindOneByNftListenerTokenId(nftlistener int64, tokenId int64) (*NftToken, error) {
	var resp NftToken
	query := fmt.Sprintf("select %s from %s where `nft_listener` = ? and `token_id` = ? limit 1", nftTokenRows, m.table)
	err := m.QueryRowNoCache(&resp, query, nftlistener, tokenId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultNftTokenModel) FindOneByNftListenerTokenIdOwner(nftlistener int64, tokenId int64, owner string) (*NftToken, error) {
	var resp NftToken
	query := fmt.Sprintf("select %s from %s where `nft_listener` = ? and `token_id` = ? and `owner` = ? limit 1", nftTokenRows, m.table)
	err := m.QueryRowNoCache(&resp, query, nftlistener, tokenId, owner)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultNftTokenModel) Update(data NftToken) error {
	nftTokenNftListenerTokenIdKey := fmt.Sprintf("%s%v:%v", cacheNftTokenNftListenerTokenIdPrefix, data.NftListener, data.TokenId)
	nftTokenIdKey := fmt.Sprintf("%s%v", cacheNftTokenIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, nftTokenRowsWithPlaceHolder)
		return conn.Exec(query, data.Symbol, data.CreatedAt, data.Description, data.TokenId, data.Owner, data.Count, data.Metadata, data.Name, data.CategoryId, data.CollectionId, data.NftListener, data.ImageUrl, data.ExternalUrl, data.UpdatedAt, data.Id)
	}, nftTokenIdKey, nftTokenNftListenerTokenIdKey)
	return err
}

func (m *defaultNftTokenModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	nftTokenIdKey := fmt.Sprintf("%s%v", cacheNftTokenIdPrefix, id)
	nftTokenNftListenerTokenIdKey := fmt.Sprintf("%s%v:%v", cacheNftTokenNftListenerTokenIdPrefix, data.NftListener, data.TokenId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, nftTokenIdKey, nftTokenNftListenerTokenIdKey)
	return err
}

func (m *defaultNftTokenModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheNftTokenIdPrefix, primary)
}

func (m *defaultNftTokenModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftTokenRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultNftTokenModel) Find(filter types.ArtworkListReq) ([]NftTokenList, int64, error) {
	query := ""
	countQuery := ""
	var resp []NftTokenList
	var totalResp ArtworkListTotal
	err := errors.New("")
	if filter.Status == "not_on_sale" {
		where := "where nt.count > 0 and nt.owner = ?"

		if filter.Category > 0 {
			where = fmt.Sprintf("%s and nt.category_id = %d", where, filter.Category)
		}

		leftJoin := "left join nft_listeners nl on nl.id = nt.nft_listener"
		query = fmt.Sprintf("select nt.*, null as auction_id, null as highest_price, null as fixed_price, nl.contract_address as nft_address, null as sale_address, nl.nft_type as token_standard from %s nt %s  %s ", m.table, leftJoin, where)
		countQuery = fmt.Sprintf("select count(1) from %s nt %s %s", m.table, leftJoin, where)
		if filter.Category != 0 {
			query = fmt.Sprintf("%s and nt.category_id = %d", query, filter.Category)
			countQuery = fmt.Sprintf("%s and nt.category_id = %d", countQuery, filter.Category)
		}
		offset := (filter.PageNo - 1) * filter.PageSize
		size := filter.PageSize
		query = fmt.Sprintf("%s limit %d offset %d", query, size, offset)
		fmt.Printf("%s \n", query)
		err = m.CachedConn.QueryRowsNoCache(&resp, query, filter.Address)
		if err != nil {
			log.Error(err)
			return nil, 0, err
		}
		err = m.CachedConn.QueryRowNoCache(&totalResp, countQuery, filter.Address)
		if err != nil {
			return nil, 0, err
		}

	} else if filter.Status == "on_sale" {
		where := "where ns.owner = ? and ns.status = 'active'"

		if filter.Category > 0 {
			where = fmt.Sprintf("%s and nt.category_id = %d", where, filter.Category)
		}

		leftJoin := "left join nft_sales ns on ns.nft_token_id = nt.id left join nft_listeners nl on nl.id = nt.nft_listener"
		query = fmt.Sprintf("select nt.id, nl.symbol, nt.created_at, nt.description, nt.token_id, nt.owner, ns.count, nt.metadata, nt.name, nt.category_id, nt.collection_id, nt.nft_listener, nt.image_url, nt.external_url, nt.updated_at, ns.id_in_contract as auction_id, ns.highest_price as highest_price, ns.fixed_price as fixed_price, nl.contract_address as nft_address, ns.sale_address as sale_address, nl.nft_type as token_standard from %s nt %s %s", m.table, leftJoin, where)
		countQuery = fmt.Sprintf("select count(1) from %s nt %s %s", m.table, leftJoin, where)
		if filter.Category != 0 {
			query = fmt.Sprintf("%s and nt.category_id = %d", query, filter.Category)
			countQuery = fmt.Sprintf("%s and nt.category_id = %d", countQuery, filter.Category)
		}
		offset := (filter.PageNo - 1) * filter.PageSize
		size := filter.PageSize
		query = fmt.Sprintf("%s limit %d offset %d", query, size, offset)
		fmt.Printf("%s \n", query)
		err = m.CachedConn.QueryRowsNoCache(&resp, query, filter.Address)
		if err != nil {
			log.Error(err)
			return nil, 0, err
		}
		err = m.CachedConn.QueryRowNoCache(&totalResp, countQuery, filter.Address)
		if err != nil {
			return nil, 0, err
		}
	}

	return resp, totalResp.Total, err
}

func (m *defaultNftTokenModel) FindInfo(id int64, id_in_contract string) (*NftTokenInfo, error) {
	//nftTokenIdKey := fmt.Sprintf("%s%v", cacheNftTokenIdPrefix, id)
	var resp NftTokenInfo
	sqlStr := `
		select nt.id,
			   nt.symbol,
			   nt.category_id,
			   nt.nft_listener,
			   nt.token_id,
			   nt.owner,
               nt.blind_box,
			   nt.collection_id,
			   nt.count,
			   nt.metadata,
			   nt.name,
			   nt.description,
			   nt.image_url,
			   nt.external_url,
			   nt.created_at,
			   nt.updated_at,
			   nt.tx_time,
			   nt.block_number,
			   nca.name                  as category_name,
			   nco.name                  as collection_name,
			   nl.contract_address       as nft_address,
			   0                         as token_type,
			   0                         as protocol_fee,
			   nl.nft_type               as token_standard,
			   'Binance Smart Chain'     as block_chain,
			   ns.fixed_price            as total_price,
			   ns.status,
			   ns.id_in_contract         as auction_id,
			   ns.count                  as sale_count,
			   ns.sale_address           as sale_address
		from %s nt
				 left join %s ns on nt.id = ns.nft_token_id and ns.status = 'active'
				 left join %s nl on nl.id = nt.nft_listener
				 left join %s nca on nt.category_id = nca.id
				 left join %s nco on nt.collection_id = nco.id
	`

	query := fmt.Sprintf(sqlStr, TABLE_NAME_NFT_TOKENS, TABLE_NAME_NFT_SALES, TABLE_NAME_NFT_LISTENERS, TABLE_NAME_NFT_CATEGORIES, TABLE_NAME_NFT_COLLECTIONS)
	var err error
	if id_in_contract != "" {
		query = fmt.Sprintf("%s where nt.id = ? and ns.id_in_contract = ? limit 1", query)
		err = m.QueryRowNoCache(&resp, query, id, id_in_contract)
	} else {
		query = fmt.Sprintf("%s where nt.id = ? limit 1", query)
		err = m.QueryRowNoCache(&resp, query, id)
	}
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultNftTokenModel) FindProperties(id int64) ([]NftTokenProperties, error) {

	leftJoin := "left join nft_token_attributes nta on nt.id = nta.nft_token_id"
	query := fmt.Sprintf("select nta.key, nta.value from %s nt %s", m.table, leftJoin)
	query = fmt.Sprintf("%s where nt.id = %d", query, id)

	var resp []NftTokenProperties

	err := m.CachedConn.QueryRowsNoCache(&resp, query)
	if err != nil {
		log.Error(err)
		return resp, nil
	}
	return resp, err
}

func (m *defaultNftTokenModel) FindTokenIds(owner string, contractAddress string) ([]int, error) {
	where := " and nt.nft_listener in (select id from nft_listeners where contract_address = ?)"
	query := fmt.Sprintf("select nt.token_id from %s nt  where nt.owner = ? %s ", m.table, where)
	var resp []int
	err := m.CachedConn.QueryRowsNoCache(&resp, query, owner, contractAddress)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return resp, err
}

func (m *defaultNftTokenModel) FindTokenIdsByName(owner string, contractAddress string, name string) ([]int, error) {
	where := "inner join nft_listeners nl on nt.nft_listener = nl.id and nl.contract_address = ?"
	query := fmt.Sprintf("select nt.token_id from %s nt %s where nt.owner = ? and nt.name = ? ", m.table, where)
	var resp []int
	err := m.CachedConn.QueryRowsNoCache(&resp, query, contractAddress, owner, name)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return resp, err
}

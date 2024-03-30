package model

import (
	"database/sql"
	"fmt"
	"log"
	"market/internal/types"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	nftSalesFieldNames          = builderx.RawFieldNames(&NftSales{})
	nftSalesRows                = strings.Join(nftSalesFieldNames, ",")
	nftSalesRowsExpectAutoSet   = strings.Join(stringx.Remove(nftSalesFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	nftSalesRowsWithPlaceHolder = strings.Join(stringx.Remove(nftSalesFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheNftSalesIdPrefix                      = "cache:nftSales:id:"
	cacheNftSalesSaleAddressIdInContractPrefix = "cache:nftSales:saleAddress:idInContract:"
)

type (
	NftSalesModel interface {
		Insert(data NftSales) (sql.Result, error)
		FindOne(id int64) (*NftSales, error)
		Histories(address, _type string, page, pageSize int) ([]NftSalesHistory, int64, error)
		FindOneBySaleAddressIdInContract(saleAddress string, idInContract string) (*NftSales, error)
		Update(data NftSales) error
		Delete(id int64) error
		Find(filter types.NFTSalesReq) ([]NftSalesList, int64, error)
		FindInfo(id int64) (*NftSalesInfo, error)
		FindProperties(id int64) ([]NftSalesProperties, error)
		FindPendingOrdersNumToday(startTime, endTime string) (int64, error)
		FindNFTDealNumToday(startTime, endTime string) (int64, error)
		FindRacaDealNumToday(startTime, endTime string) (int64, error)
		FindNFTTradeNumToday(contractAddress, startTime, endTime string) (int64, error)
	}

	NftSalesProperties struct {
		Key   string         `db:"key"`   // 属性key
		Value sql.NullString `db:"value"` // 属性value
	}

	OrderNumTodayResp struct {
		Total int64 `json:"total"`
	}
	RacaDealNumResp struct {
		Total decimal.Decimal `json:"total"`
	}
	NftSalesList struct {
		ProtocolFeeReceiptAddress sql.NullString  `db:"protocol_fee_receipt_address"` // 手续费收款address
		NftTokenId                int64           `db:"nft_token_id"`                 // NFT_TOKEN表id
		SaleAddress               string          `db:"sale_address"`                 // 出售合约address
		EndingPrice               decimal.Decimal `db:"ending_price"`                 // 拍卖最终价格
		Status                    string          `db:"status"`                       // 状态
		UpdatedAt                 time.Time       `db:"updated_at"`                   // 更新时间
		FixedPrice                decimal.Decimal `db:"fixed_price"`                  // 出售价格
		StartingPrice             decimal.Decimal `db:"starting_price"`               // 拍卖起始价格
		TokenType                 int64           `db:"token_type"`                   // 价格类型，1代表RACA，2代表BNB
		StartTime                 sql.NullInt64   `db:"start_time"`                   // 开始时间戳
		EndTime                   sql.NullInt64   `db:"end_time"`                     // 结束时间戳
		ProtocolFee               float64         `db:"protocol_fee"`                 // 手续费
		CreatedAt                 time.Time       `db:"created_at"`                   // 创建时间
		Id                        int64           `db:"id"`
		Name                      string          `db:"name"`           // NFT名称
		TokenId                   string          `db:"token_id"`       // NFT token id
		ImageUrl                  string          `db:"image_url"`      // NFT图片url
		IdInContract              string          `db:"id_in_contract"` // 合约内id
		HighestPrice              decimal.Decimal `db:"highest_price"`  // 拍卖最高出价
		BuyerAddress              sql.NullString  `db:"buyer_address"`  // 买家address
		SaleType                  string          `db:"sale_type"`      // 出售方式
		Count                     int64           `db:"count"`          // 销售数量
		Owner                     string          `db:"owner"`          // NFT所有者address
	}

	NftSalesInfo struct {
		ProtocolFeeReceiptAddress sql.NullString  `db:"protocol_fee_receipt_address"` // 手续费收款address
		NftTokenId                int64           `db:"nft_token_id"`                 // NFT_TOKEN表id
		SaleAddress               string          `db:"sale_address"`                 // 出售合约address
		EndingPrice               decimal.Decimal `db:"ending_price"`                 // 拍卖最终价格
		Status                    string          `db:"status"`                       // 状态
		UpdatedAt                 time.Time       `db:"updated_at"`                   // 更新时间
		FixedPrice                decimal.Decimal `db:"fixed_price"`                  // 出售价格
		StartingPrice             decimal.Decimal `db:"starting_price"`               // 拍卖起始价格
		TokenStandard             string          `db:"token_standard"`               // Token标准
		BlockChain                string          `db:"block_chain"`                  // 链信息
		TokenType                 int64           `db:"token_type"`                   // 价格类型，1代表RACA，2代表BNB
		StartTime                 sql.NullInt64   `db:"start_time"`                   // 开始时间戳
		EndTime                   sql.NullInt64   `db:"end_time"`                     // 结束时间戳
		ProtocolFee               decimal.Decimal `db:"protocol_fee"`                 // 手续费
		CreatedAt                 time.Time       `db:"created_at"`                   // 创建时间
		Id                        int64           `db:"id"`
		Name                      string          `db:"name"`            // NFT名称
		Description               sql.NullString  `db:"description"`     // 描述，支持富文本
		CategoryName              string          `db:"category_name"`   // NFT类别名称
		CollectionName            string          `db:"collection_name"` // NFT对应收藏类名称
		ImageUrl                  string          `db:"image_url"`       // NFT图片url
		IdInContract              string          `db:"id_in_contract"`  // 合约内id
		NftAddress                string          `db:"nft_address"`     // 合约地址
		HighestPrice              decimal.Decimal `db:"highest_price"`   // 拍卖最高出价
		BuyerAddress              sql.NullString  `db:"buyer_address"`   // 买家address
		SaleType                  string          `db:"sale_type"`       // 出售方式
		Count                     int64           `db:"count"`           // 销售数量
		Owner                     string          `db:"owner"`           // NFT所有者address
		NeedCheck                 int64           `db:"need_check"`
		TokenId                   int64           `db:"token_id"`
	}

	defaultNftSalesModel struct {
		sqlc.CachedConn
		table string
	}

	NftSalesTotal struct {
		Total int64 `json:"total"`
	}

	NftSales struct {
		ProtocolFeeReceiptAddress sql.NullString  `db:"protocol_fee_receipt_address"` // 手续费收款address
		NftTokenId                int64           `db:"nft_token_id"`                 // NFT_TOKEN表id
		SaleAddress               string          `db:"sale_address"`                 // 出售合约address
		EndingPrice               decimal.Decimal `db:"ending_price"`                 // 拍卖最终价格
		Status                    string          `db:"status"`                       // 状态
		UpdatedAt                 time.Time       `db:"updated_at"`                   // 更新时间
		FixedPrice                decimal.Decimal `db:"fixed_price"`                  // 出售价格
		StartingPrice             decimal.Decimal `db:"starting_price"`               // 拍卖起始价格
		TokenType                 int64           `db:"token_type"`                   // 价格类型，1代表RACA，2代表BNB
		StartTime                 sql.NullInt64   `db:"start_time"`                   // 开始时间戳
		EndTime                   sql.NullInt64   `db:"end_time"`                     // 结束时间戳
		ProtocolFee               float64         `db:"protocol_fee"`                 // 手续费
		CreatedAt                 time.Time       `db:"created_at"`                   // 创建时间
		Id                        int64           `db:"id"`
		IdInContract              string          `db:"id_in_contract"` // 合约内id
		HighestPrice              decimal.Decimal `db:"highest_price"`  // 拍卖最高出价
		BuyerAddress              sql.NullString  `db:"buyer_address"`  // 买家address
		SaleType                  string          `db:"sale_type"`      // 出售方式
		Count                     int64           `db:"count"`          // 销售数量
		Owner                     string          `db:"owner"`          // NFT所有者address
		CreateTxHash              sql.NullString  `db:"create_tx_hash"`
		OperateTxHash             sql.NullString  `db:"operate_tx_hash"`
	}

	NftSalesHistory struct {
		NftTokenId    int64           `db:"nft_token_id"` // NFT_TOKEN表id
		UpdatedAt     time.Time       `db:"updated_at"`   // 更新时间
		FixedPrice    decimal.Decimal `db:"fixed_price"`  // 出售价格
		Id            int64           `db:"id"`
		BuyerAddress  sql.NullString  `db:"buyer_address"` // 买家address
		Count         int64           `db:"count"`         // 销售数量
		Owner         string          `db:"owner"`         // NFT所有者address
		OperateTxHash sql.NullString  `db:"operate_tx_hash"`
		TokenName     string          `db:"token_name"`
	}
)

func NewNftSalesModel(conn sqlx.SqlConn, c cache.CacheConf) NftSalesModel {
	return &defaultNftSalesModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`nft_sales`",
	}
}

func (m *defaultNftSalesModel) Insert(data NftSales) (sql.Result, error) {
	nftSalesSaleAddressIdInContractKey := fmt.Sprintf("%s%v:%v", cacheNftSalesSaleAddressIdInContractPrefix, data.SaleAddress, data.IdInContract)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, nftSalesRowsExpectAutoSet)
		return conn.Exec(query, data.ProtocolFeeReceiptAddress, data.NftTokenId, data.SaleAddress, data.EndingPrice, data.Status, data.UpdatedAt, data.FixedPrice, data.StartingPrice, data.TokenType, data.StartTime, data.EndTime, data.ProtocolFee, data.CreatedAt, data.IdInContract, data.HighestPrice, data.BuyerAddress, data.SaleType, data.Count, data.Owner)
	}, nftSalesSaleAddressIdInContractKey)
	return ret, err
}

func (m *defaultNftSalesModel) FindOne(id int64) (*NftSales, error) {
	//nftSalesIdKey := fmt.Sprintf("%s%v", cacheNftSalesIdPrefix, id)
	var resp NftSales
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftSalesRows, m.table)
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

func (m *defaultNftSalesModel) FindOneBySaleAddressIdInContract(saleAddress string, idInContract string) (*NftSales, error) {
	nftSalesSaleAddressIdInContractKey := fmt.Sprintf("%s%v:%v", cacheNftSalesSaleAddressIdInContractPrefix, saleAddress, idInContract)
	var resp NftSales
	err := m.QueryRowIndex(&resp, nftSalesSaleAddressIdInContractKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `sale_address` = ? and `id_in_contract` = ? limit 1", nftSalesRows, m.table)
		if err := conn.QueryRow(&resp, query, saleAddress, idInContract); err != nil {
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

func (m *defaultNftSalesModel) Update(data NftSales) error {
	nftSalesIdKey := fmt.Sprintf("%s%v", cacheNftSalesIdPrefix, data.Id)
	nftSalesSaleAddressIdInContractKey := fmt.Sprintf("%s%v:%v", cacheNftSalesSaleAddressIdInContractPrefix, data.SaleAddress, data.IdInContract)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, nftSalesRowsWithPlaceHolder)
		return conn.Exec(query, data.ProtocolFeeReceiptAddress, data.NftTokenId, data.SaleAddress, data.EndingPrice, data.Status, data.UpdatedAt, data.FixedPrice, data.StartingPrice, data.TokenType, data.StartTime, data.EndTime, data.ProtocolFee, data.CreatedAt, data.IdInContract, data.HighestPrice, data.BuyerAddress, data.SaleType, data.Count, data.Owner, data.Id)
	}, nftSalesIdKey, nftSalesSaleAddressIdInContractKey)
	return err
}

func (m *defaultNftSalesModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	nftSalesIdKey := fmt.Sprintf("%s%v", cacheNftSalesIdPrefix, id)
	nftSalesSaleAddressIdInContractKey := fmt.Sprintf("%s%v:%v", cacheNftSalesSaleAddressIdInContractPrefix, data.SaleAddress, data.IdInContract)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, nftSalesIdKey, nftSalesSaleAddressIdInContractKey)
	return err
}

func (m *defaultNftSalesModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheNftSalesIdPrefix, primary)
}

func (m *defaultNftSalesModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", nftSalesRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultNftSalesModel) Histories(address, _type string, page, pageSize int) ([]NftSalesHistory, int64, error) {
	cmd := " `status` = 'executed' "
	if _type == "buy" {
		cmd += " and `nft_sales`.`buyer_address` = '" + address + "' "
	} else if _type == "sell" {
		cmd += " and `nft_sales`.`owner` = '" + address + "' "
	} else {
		cmd += " and (`nft_sales`.`owner` = '" + address + "' or `nft_sales`.`buyer_address` = '" + address + "') "
	}
	var total int64
	totalSql := fmt.Sprintf(" select count(*) from %s where %s ", m.table, cmd)
	err := m.QueryRowNoCache(&total, totalSql)
	if err != nil {
		return nil, 0, err
	}
	fmt.Printf("total:%d\n", total)
	list := make([]NftSalesHistory, 0)
	if total == 0 {
		return list, 0, nil
	}
	offset := (page - 1) * pageSize
	fmt.Printf("offset: %d\n", offset)
	if int64(offset) >= total {
		return list, total, nil
	}

	listSql := fmt.Sprintf(" select nft_sales.id, nft_sales.nft_token_id, nft_sales.count, nft_sales.fixed_price, nft_sales.owner, nft_sales.buyer_address, nft_sales.updated_at, nft_sales.operate_tx_hash ,CONCAT(IFNULL(nft_tokens.name,'UNKNOWN'),'#', nft_tokens.token_id) token_name from nft_sales left join nft_tokens on nft_sales.nft_token_id = nft_tokens.id where %s order by `nft_sales`.`updated_at` desc limit ? offset ? ", cmd)
	err = m.QueryRowsNoCache(&list, listSql, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (m *defaultNftSalesModel) Find(filter types.NFTSalesReq) ([]NftSalesList, int64, error) {
	leftJoin := "left join nft_tokens nt on nt.id = ns.nft_token_id left join nft_listeners nl on nt.nft_listener = nl.id"
	//query := fmt.Sprintf("select ns.*, nt.name as name, nt.token_id as token_id, nt.image_url as image_url from %s ns %s where ns.status = 'active' and nt.name like CONCAT('%%', ?, '%%') ", m.table, leftJoin)
	//countQuery := fmt.Sprintf("select count(1) from %s ns %s where ns.status = 'active' and nt.name like CONCAT('%%', ?, '%%')", m.table, leftJoin)
	query := fmt.Sprintf("select ns.*, nt.name as name, nt.token_id as token_id, nt.image_url as image_url from %s ns %s where ns.status = 'active' ", m.table, leftJoin)
	countQuery := fmt.Sprintf("select count(1) from %s ns %s where ns.status = 'active' ", m.table, leftJoin)

	if filter.Name != "" {
		query = fmt.Sprintf("%s and `nt`.`name` = '%s' ", query, filter.Name)
		countQuery = fmt.Sprintf("%s and `nt`.`name` = '%s' ", countQuery, filter.Name)
	}

	if filter.Category > 0 {
		query = fmt.Sprintf("%s and nt.category_id = %d", query, filter.Category)
		countQuery = fmt.Sprintf("%s and nt.category_id = %d", countQuery, filter.Category)
	}

	offset := (filter.PageNo - 1) * filter.PageSize
	size := filter.PageSize
	sort := filter.SortBy
	if sort == "single_price" {
		sort = "fixed_price/ns.count"
	}
	query = fmt.Sprintf("%s order by %s %s limit %d offset %d", query, sort, filter.Order, size, offset)

	var totalResp NftSalesTotal
	err := m.CachedConn.QueryRowNoCache(&totalResp, countQuery)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]NftSalesList, 0)
	if totalResp.Total > 0 {
		err = m.CachedConn.QueryRowsNoCache(&resp, query)
		if err != nil {
			log.Fatal(err)
			return nil, 0, err
		}
	}

	fmt.Printf("************ %d ************\n", len(resp))

	return resp, totalResp.Total, err
}

func (m *defaultNftSalesModel) FindInfo(id int64) (*NftSalesInfo, error) {
	//nftSalesIdKey := fmt.Sprintf("%s%v", cacheNftSalesIdPrefix, id)
	var resp NftSalesInfo
	leftJoin := " left join nft_tokens nt on nt.id = ns.nft_token_id left join nft_listeners nl on nl.id = nt.nft_listener left join nft_categories nca on nca.id = nt.category_id left join nft_collections nco on nco.id = nt.collection_id"
	query := fmt.Sprintf("select ns.*, nt.token_id as token_id, nt.name as name, nt.description as description, nt.image_url as image_url, nl.contract_address as nft_address, nca.name as category_name, nco.name as collection_name, 0 as need_check, 'Binance Smart Chain' as block_chain, nl.nft_type as token_standard from %s ns %s", m.table, leftJoin)
	query = fmt.Sprintf("%s where ns.id = ? limit 1", query)
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

func (m *defaultNftSalesModel) FindProperties(id int64) ([]NftSalesProperties, error) {

	leftJoin := "left join nft_token_attributes nta on ns.nft_token_id = nta.nft_token_id"
	query := fmt.Sprintf("select nta.key, nta.value from %s ns %s", m.table, leftJoin)
	query = fmt.Sprintf("%s where ns.id = %d", query, id)

	var resp []NftSalesProperties

	err := m.CachedConn.QueryRowsNoCache(&resp, query)
	if err != nil {
		log.Printf("Error get attributes error %v\n", err)
		return resp, nil
	}
	return resp, err
}

func (m *defaultNftSalesModel) FindPendingOrdersNumToday(startTime, endTime string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s ns WHERE ns.created_at > ? AND ns.created_at < ?;", m.table)
	var resp OrderNumTodayResp
	err := m.CachedConn.QueryRowNoCache(&resp, query, startTime, endTime)
	if err != nil {
		return 0, err
	}
	return resp.Total, err
}

func (m *defaultNftSalesModel) FindNFTDealNumToday(startTime, endTime string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s ns WHERE ns.buyer_address is not null and ns.created_at > ?  AND ns.created_at < ? ;", m.table)
	var resp OrderNumTodayResp
	err := m.CachedConn.QueryRowNoCache(&resp, query, startTime, endTime)
	if err != nil {
		return 0, err
	}
	return resp.Total, err
}
func (m *defaultNftSalesModel) FindRacaDealNumToday(startTime, endTime string) (int64, error) {
	query := fmt.Sprintf("SELECT COALESCE(sum(ns.fixed_price),0) FROM %s ns WHERE  ns.status= 'executed'  and ns.created_at >?  AND ns.created_at < ? ;", m.table)
	var resp RacaDealNumResp
	err := m.CachedConn.QueryRowNoCache(&resp, query, startTime, endTime)
	if err != nil {
		return 0, err
	}
	return resp.Total.IntPart(), err
}
func (m *defaultNftSalesModel) FindNFTTradeNumToday(contractAddress, startTime, endTime string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s ns WHERE ns.nft_token_id IN (SELECT nt.id FROM nft_tokens nt "+
		"WHERE nt.token_id =(SELECT nl.id FROM nft_listeners nl WHERE nl.contract_address =? )) "+
		"and ns.created_at >?  AND ns.created_at < ? ;", m.table)
	var resp OrderNumTodayResp
	err := m.CachedConn.QueryRowNoCache(&resp, query, contractAddress, startTime, endTime)
	if err != nil {
		return 0, err
	}
	return resp.Total, err
}

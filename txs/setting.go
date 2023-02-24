package txs

import (
	"database/sql"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strconv"
	"strings"
)

// 一些数据库连接的配置
var (
	database = "eth" // 数据库名
	// 连接相关
	driver      = "mysql" // 数据库引擎
	user        = "morph"
	passwd      = "morphdag"
	protocol    = "tcp" //连接协议
	port        = "3306"
	useDatabase = "USE " + database

	dataSource string

	tables = []string{"txs1", "txs2", "txs3", "txs4", "txs5", "txs6", "txs7", "txs8", "txs9", "txs10",
		"txs11", "txs12", "txs13", "txs14", "txs15", "txs16", "txs17", "txs18", "txs19", "txs20",
		"txs21", "txs22", "txs23", "txs24", "txs25", "txs26", "txs27", "txs28", "txs29", "txs30",
		"txs31", "txs32", "txs33", "txs34", "txs35", "txs36", "txs37", "txs38", "txs39", "txs40",
		"txs41", "txs42", "txs43", "txs44", "txs45", "txs46", "txs47", "txs48", "txs49", "txs50",
		"txs51", "txs52", "txs53", "txs54", "txs55", "txs56", "txs57", "txs58", "txs59", "txs60"} // 表名

	sqlServers []*sql.DB
)

func SetHost(host string) {
	dataSource = user + ":" + passwd + "@" + protocol + "(" + host + ":" + port + ")/" // 用户名:密码@tcp(ip:端口)/
}

type MorphTransfer struct {
	From  string
	To    string
	Type  uint8
	Index uint16
}

type MorphTransaction struct {
	Hash        string
	BlockNumber *big.Int
	BlockHash   string
	From        string
	To          string

	Transfer []*MorphTransfer
}

func NewMorphTransfer() *MorphTransfer {
	return &MorphTransfer{
		From: "",
		To:   "",
	}
}

func NewMorphTransaction() *MorphTransaction {
	return &MorphTransaction{
		Hash:        "",
		BlockNumber: big.NewInt(0),
		BlockHash:   "",
		From:        "",
		To:          "",
	}
}

// 以下是导出交易的相关字段需求 (注: 结构体里面的数据要给json包访问 -> 需要首字母大写)

type Transfer struct {
	From *Balance `json:"from"`
	To   *Balance `json:"to"`
	// type 类型(1: 合约的调用者转账给某一接收方, 可能是嵌套合约的调用; 2: 合约创建者发送到合约账户; 3: 将手续费添加给矿工, 只有To字段)
	// type 类型(4: 多扣除的手续费退还, 只有To字段; 5: 从交易发起者账户预扣除交易费, 只有From字段; 6: 合约销毁)
	// type 类型(7: 给挖出叔父区块的矿工奖励, 只有To字段, Nonce值为0; 8: 给挖出该区块的矿工奖励, 只有To字段, Nonce值为0; 9: The DAO 硬分叉相关)
	// type 类型(9: The DAO 硬分叉相关, Nonce值为0)
	// 类型 7, 8 放在了该区块最后一个事务的交易中, 最后处理
	// 类型 9 放在了该区块第一个事务的交易中, 优先处理
	Type  uint8  `json:"type"` // 一个交易中, 3 4 5 类型最多只有一个, 且需要结合起来看(5-4: 交易发起者需要扣除的手续费 != 3: 给矿工的手续费)
	Nonce uint16 `json:"nonce"`
}

// BigInt 为了反序列化*big.Int
type BigInt struct {
	big.Int
}

type Balance struct {
	Address  string  `json:"address"`
	Value    *BigInt `json:"value"`
	BeforeTx *BigInt `json:"beforeTx"`
	AfterTx  *BigInt `json:"afterTx"`
}

type Transaction struct {
	Count            uint16          `json:"count"`
	Transfer         []*Transfer     `json:"balance"`
	BlockHash        *common.Hash    `json:"blockHash"`
	BlockNumber      *hexutil.Big    `json:"blockNumber"`
	From             common.Address  `json:"from"`
	Hash             *common.Hash    `json:"hash"`
	Contract         bool            `json:"contract"` // 是否为合约
	Nonce            hexutil.Uint64  `json:"nonce"`
	To               *common.Address `json:"to"`
	TransactionIndex *hexutil.Uint64 `json:"transactionIndex"`
	Value            *hexutil.Big    `json:"value"`
}

func (i *BigInt) UnmarshalJSON(b []byte) error {
	var a = string(b)
	sp := strings.Split(a, "e")
	var val string
	if len(sp) == 1 {
		val = sp[0]
	} else {
		num := strings.Split(sp[0], ".")
		leafLen := 0
		oLen, err := strconv.Atoi(sp[1])
		if err != nil {
			return err
		}
		if len(num) == 1 {
			leafLen = oLen
		} else if len(num[1]) < oLen {
			leafLen = oLen - len(num[1])
		}
		if len(num) == 1 {
			val = num[0]
		} else {
			val = num[0] + num[1]
		}
		for k := 0; k < leafLen; k++ {
			val += "0"
		}
	}
	i.SetString(val, 10) // 基于10进制

	return nil
}

package storagermw

import (
	"database/sql"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strconv"
	"strings"
)

// 一些数据库连接的配置
var (
	database = "test" // 数据库名
	// 连接相关
	driver      = "mysql" // 数据库引擎
	user        = "morph"
	passwd      = "morphdag"
	protocol    = "tcp" //连接协议
	port        = "3306"
	useDatabase = "USE " + database

	dataSource string

	tableOfTxs = []string{"txs1", "txs2", "txs3", "txs4", "txs5", "txs6", "txs7", "txs8", "txs9", "txs10",
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

type transfer interface { // 转账
	GetLabel() uint8 // 0: 普通转账(state), 1: ERC20类转账(storage), 2: KECCAK256, 3: push20, 4: ContractCall
}

type Transaction struct {

	Type         uint64
	ChainID      *big.Int
	InputData    []byte
	Gas          uint64 // Gas Limit
	MaxFeePerGas *big.Int   // MaxFeePerGas(GasFeeCap) which is the maximum you are willing to pay per unit of gas
	// MaxPriorityFeePerGas(GasTipCap), which is optional, determined by the user, and is paid directly to miners.
	MaxPriorityFeePerGas *big.Int
	GasPrice             *big.Int // The price of each unit of gas, in wei.
	Value                *big.Int
	Nonce                uint64
	To                   *common.Address

	// Signature values
	R *big.Int
	S *big.Int
	V *big.Int

	BlockNumber *big.Int
	Hash        *common.Hash
	From        *common.Address
	Index       uint64

	Transfers *[]transfer
}

func newTransaction(number string, hash string, info *string) *Transaction {
	tx := &Transaction{}
	infos := strings.Split(*info, "|")
	var tmp  = big.NewInt(0)
	//type_, chainID, inputData, gas, maxFeePerGas, maxPriorityFeePerGas, gasPrice, value, nonce, to,
	//	r, s, v, trs, from, index
	tx.Type, _ = strconv.ParseUint(infos[0], 0, 64)
	if infos[1] == ""{
		tx.ChainID = nil
	}else {
		tx.ChainID, _ = tmp.SetString(infos[1], 16)
	}
	tx.InputData = common.Hex2Bytes(infos[2])

	tx.Gas, _ = strconv.ParseUint(infos[3], 0, 64)
	if infos[4] == ""{
		tx.MaxFeePerGas = nil
	}else {
		tx.MaxFeePerGas, _ = tmp.SetString(infos[4], 16)
	}
	if infos[5] == ""{
		tx.MaxPriorityFeePerGas = nil
	}else {
		tx.MaxPriorityFeePerGas, _ = tmp.SetString(infos[5], 16)
	}
	if infos[6] == ""{
		tx.GasPrice = nil
	}else {
		tx.GasPrice, _ = tmp.SetString(infos[6], 16)
	}
	if infos[7] == ""{
		tx.Value = nil
	}else {
		tx.Value, _ = tmp.SetString(infos[7], 16)
	}

	tx.Nonce, _ = strconv.ParseUint(infos[8], 0, 64)

	to := common.HexToAddress(infos[9])
	tx.To = &to

	//// Signature values
	if infos[10] == ""{
		tx.R = nil
	}else {
		tx.R, _ = tmp.SetString(infos[10], 16)
	}
	if infos[11] == ""{
		tx.S = nil
	}else {
		tx.S, _ = tmp.SetString(infos[11], 16)
	}
	if infos[12] == ""{
		tx.V = nil
	}else {
		tx.V, _ = tmp.SetString(infos[12], 16)
	}

	tx.BlockNumber, _ = tmp.SetString(number, 16)

	h := common.HexToHash(hash)
	tx.Hash = &h

	tx.Transfers = newTransfers(infos[13])

	from := common.HexToAddress(infos[14])
	tx.From = &from
	tx.Index, _ = strconv.ParseUint(infos[15], 0, 64)

	return tx
}

func newTransfers(trs string) *[]transfer {
	var allTransfers []transfer
	transfers := strings.Split(trs, " ")
	for _, tr := range transfers{
		var tr_ transfer
		if tr[0] == '0'{
			tr_ = transfer(newStateTransition(tr))
		}else if tr[0] == '1' {
			tr_ = transfer(newStorageTransition(tr))
		}else if tr[0] == '4'{
			tr_ = transfer(newContractCall(tr))
		}
		allTransfers = append(allTransfers, tr_)
	}
	return &allTransfers
}
package txs

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	_ "github.com/go-sql-driver/mysql" // 导入包但不使用, init()
	"math/big"
	"sync"
)

var transactionChan chan []Transaction

func GetTxsByBlockNumber(blockNumber *big.Int) []*MorphTransaction {
	//start := time.Now()
	var wg sync.WaitGroup

	transactionChan = make(chan []Transaction, 60)
	wg.Add(len(tables))
	var txs []Transaction

	for index, table := range tables {
		go getTxsByBlockNumber(table, blockNumber, index, &wg)
	}
	var k = 0
	for {
		if k == 60 {
			close(transactionChan)
		}
		k += 1
		val, ok := <-transactionChan
		if !ok {
			break
		}
		if len(val) == 0 {
			continue
		}
		txs = append(txs, val...)
	}

	wg.Wait()
	transactionChan = nil
	mtxs := dealTxs(&txs) // 处理一个区块的事务
	//fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Get the transactions of block number("+blockNumber.String()+") spend "+time.Since(start).String())
	return mtxs
}

func GetTxByHash(hash string) []*MorphTransaction {
	//start := time.Now()
	var wg sync.WaitGroup

	transactionChan = make(chan []Transaction, 60)
	wg.Add(len(tables))
	var txs []Transaction

	for index, table := range tables {
		go getTxByHash(table, hash, index, &wg)
	}
	var k = 0
	for {
		if k == 60 {
			close(transactionChan)
		}
		k += 1
		val, ok := <-transactionChan
		if !ok {
			break
		}
		if len(val) == 0 {
			continue
		}
		txs = append(txs, val...)
	}

	wg.Wait()
	transactionChan = nil
	mtxs := dealTxs(&txs) // 处理一个区块的事务
	//fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Get the transactions by hash(\""+hash+"\") spend "+time.Since(start).String())
	return mtxs
}

// dealTxs deal transactions of a block
func dealTxs(txs *[]Transaction) []*MorphTransaction {
	txSet := make(map[common.Hash]bool, 500)
	var morphTxs []*MorphTransaction
	dealNil := false
	for _, tx := range *txs {
		var morphTx *MorphTransaction = nil
		if tx.Hash == nil && !dealNil { // 处理挖矿奖励
			morphTx = dealRewardTx(&tx)
			dealNil = true
		} else if _, ok := txSet[*tx.Hash]; !ok { // 交易尚未处理过,防止之前存了重复的交易
			txSet[*tx.Hash] = true
			morphTx = dealGeneralTx(&tx)
		}

		if morphTx != nil {
			morphTxs = append(morphTxs, morphTx)
		}
	}
	return morphTxs
}

func dealGeneralTx(tx *Transaction) *MorphTransaction {
	morphTx := NewMorphTransaction()
	morphTx.Hash = tx.Hash.Hex()
	morphTx.BlockNumber = (*big.Int)(tx.BlockNumber) // 区块奖励交易记录, 只有区块号和区块Hash
	morphTx.BlockHash = tx.BlockHash.Hex()           // 区块奖励交易记录, 只有区块号和区块Hash
	morphTx.From = tx.From.Hex()
	if tx.To != nil { // create contract
		morphTx.To = tx.To.Hex()
	}
	for _, tr := range tx.Transfer {
		morphTr := dealTransfer(tr)
		morphTx.Transfer = append(morphTx.Transfer, morphTr)
	}
	return morphTx
}

func dealRewardTx(tx *Transaction) *MorphTransaction {
	morphTx := NewMorphTransaction()
	morphTx.Hash = tx.BlockNumber.String()
	morphTx.BlockNumber = (*big.Int)(tx.BlockNumber) // 区块奖励交易记录, 只有区块号和区块Hash
	morphTx.BlockHash = tx.BlockHash.Hex()           // 区块奖励交易记录, 只有区块号和区块Hash
	for _, tr := range tx.Transfer {
		morphTr := dealTransfer(tr)
		morphTx.Transfer = append(morphTx.Transfer, morphTr)
	}
	return morphTx
}

func dealTransfer(tr *Transfer) *MorphTransfer {
	morphTx := NewMorphTransfer()
	morphTx.From = tr.From.Address
	morphTx.To = tr.To.Address
	morphTx.Type = tr.Type
	morphTx.Index = tr.Nonce
	return morphTx
}

// getTxsByBlockNumber 获取某一区块号的所有交易
func getTxsByBlockNumber(table string, blockNumber *big.Int, index int, wg *sync.WaitGroup) {
	defer wg.Done()
	sqlServer := sqlServers[index]
	rows, err := sqlServer.Query("SELECT info FROM " + table + " WHERE blockNumber=\"" + (*hexutil.Big)(blockNumber).String() + "\";")
	defer rows.Close() // 非常重要：关闭rows释放持有的数据库链接
	if err != nil {
		fmt.Println("Error: Query failed", err)
		return
	}
	// 循环读取结果集中的数据
	var txs []Transaction
	for rows.Next() {
		var (
			info string
			tx   Transaction
		)
		err = rows.Scan(&info)
		if err != nil {
			fmt.Println("Error: Scan failed", err)
			return
		}
		err = json.Unmarshal([]byte(info), &tx)
		if err != nil {
			fmt.Println("Error: Unmarshal failed", err)
			return
		}
		txs = append(txs, tx)
	}
	transactionChan <- txs
}

// getTxsByHash 获取某一交易
func getTxByHash(table string, hash string, index int, wg *sync.WaitGroup) {
	defer wg.Done()
	sqlServer := sqlServers[index]
	rows, err := sqlServer.Query("SELECT info FROM " + table + " WHERE hash=\"" + hash + "\";")
	defer rows.Close() // 非常重要：关闭rows释放持有的数据库链接
	if err != nil {
		fmt.Println("Error: Query failed", err)
		return
	}
	// 循环读取结果集中的数据
	var txs []Transaction
	for rows.Next() {
		var (
			info string
			tx   Transaction
		)
		err = rows.Scan(&info)
		if err != nil {
			fmt.Println("Error: Scan failed", err)
			return
		}
		err = json.Unmarshal([]byte(info), &tx)
		if err != nil {
			fmt.Println("Error: Unmarshal failed", err)
			return
		}
		txs = append(txs, tx)
	}
	transactionChan <- txs
}

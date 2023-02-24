package storagermw

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	_ "github.com/go-sql-driver/mysql" // 导入包但不使用, init()
	"math/big"
	"sync"
)

var txsChan chan []*Transaction

func GetTransactionByBlockNumber(blockNumber *big.Int, txsNum *int) (int, int) {
	var wg sync.WaitGroup
	numCpu := 12
	txsChan = make(chan []*Transaction, numCpu)
	wg.Add(numCpu)

	tables := tableOfTxs[0:numCpu]
	for index, table := range tables {
		go getTransactionByBlockNumber(table, (*hexutil.Big)(blockNumber).String(), index, &wg)
	}

	var (
		k = 0
		allTxs []*Transaction
	)
	for {
		if k == numCpu {
			close(txsChan)
		}
		k += 1
		txs, ok := <-txsChan
		if !ok {
			break
		}
		 allTxs = append(allTxs, txs...)
	}
	wg.Wait()
	*txsNum += len(allTxs)
	return dealAllTxs(allTxs)
}

// getTransactionByBlockNumber 获取某一区块号的区块的读写比例
func getTransactionByBlockNumber(table string, blockNumber string, index int, wg *sync.WaitGroup) {
	defer wg.Done()
	sqlServer := sqlServers[index]
	rows, err := sqlServer.Query("SELECT number, hash, infos FROM " + table + " WHERE number=\"" + blockNumber + "\";")
	defer rows.Close() // 非常重要：关闭rows释放持有的数据库链接
	if err != nil {
		fmt.Println("Error: Query failed", err)
		return
	}

	var (
		number, hash, info string
		txs []*Transaction
	)
	for rows.Next() {
		err = rows.Scan(&number, &hash, &info)
		if err != nil {
			fmt.Println("Error: Scan failed", err)
			return
		}
		txs = append(txs, newTransaction(number, hash, &info))
	}
	txsChan <- txs
}

func dealAllTxs(allTxs []*Transaction) (int, int) {
	read, rmw := 0, 0
	var rmwTmp map[common.Hash]int
	rmwTmp = make(map[common.Hash]int)
	for _, tx := range allTxs{
		trs := tx.Transfers
		//var (
		//	lastContractAddress *common.Address
		//	contractAddress []*common.Address
		//	addressLen = -1
		//)
		for _, tr := range *trs{
			if tr.GetLabel() == 1 {
				storage := tr.(*storageTransition)
				var newValue = ""
				if storage.newValue != nil{
					newValue = storage.newValue.Big().String()
				}
				if times, ok := rmwTmp[*storage.slot]; ok{
					if newValue == ""{
						rmwTmp[*storage.slot] = times + 1
					}else {
						rmw += 1
						//fmt.Println(times, storage.slot, storage.preValue.Big(), newValue)
						delete(rmwTmp, *storage.slot)
						for _, _ = range rmwTmp {
							read += 1
						}
						rmwTmp = make(map[common.Hash]int)
					}
				}else { // 未曾记录
					if newValue == ""{
						rmwTmp[*storage.slot] = 1
					}else {
						fmt.Println("直接写？", storage.slot, storage.preValue.Big(), newValue)
					}
				}
				//fmt.Println(storage.slot, storage.preValue.Big(), newValue)
			}
		}
	}
	return read, rmw
	//if rmw != 0 && read != 0{
	//	re := read * 100 / (read + rmw)
	//	rm := rmw * 100 / (read + rmw)
	//	if re + rm < 100{
	//		rm += 1
	//	}
	//	strTmp := fmt.Sprintf("%d:%d", re, rm)
	//	fmt.Println(read, rmw, strTmp)
	//}
}


package txs

import (
	"fmt"
	"math/big"
	"testing"
)

func Test(t *testing.T) {
	SetHost("202.114.7.176")
	err := OpenSqlServers()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer CloseSqlServers()
	//min, max := big.NewInt(7202102), big.NewInt(8852777)
	for i := big.NewInt(8000000); i.Cmp(big.NewInt(8450000)) == -1; i = i.Add(i, big.NewInt(1)) {
		mtxs := GetTxsByBlockNumber(i)
		for _, tx := range mtxs {

			fmt.Println(tx.Hash, tx.BlockNumber, tx.BlockHash, tx.From, tx.To, len(tx.Transfer))

		}
	}

	//mtxs := GetTxByHash("0x6de644")
	//for _, tx := range mtxs {
	//	fmt.Println(tx.Hash, tx.BlockNumber, tx.BlockHash, tx.From, tx.To, len(tx.Transfer))
	//	for _, tr := range tx.Transfer{
	//		fmt.Println("From: " + tr.From, "  To: " + tr.To)
	//	}
	//}
}

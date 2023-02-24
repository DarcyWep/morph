package storagermw

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"
)

func Test(t *testing.T) {
	//SetHost("202.114.6.243")
	SetHost("192.168.5.29")
	err := OpenSqlServers()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer CloseSqlServers()

	//OpenFile读取文件，不存在时则创建，使用追加模式
	path := "/Users/darcywep/Projects/GoProjects/morph/storagermw/storagermw.csv"
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("文件打开失败！")
	}
	defer file.Close()

	var data = []string{"number", "StateRead", "StateWrite", "StateRWRate"}
	//创建写入接口
	writerCsv := csv.NewWriter(file)
	err = writerCsv.Write(data)
	if err != nil {
		log.Println("WriterCsv写入文件失败")
	}
	writerCsv.Flush() //刷新，不刷新是无法写入的

	//min, max := big.NewInt(10039156), big.NewInt(10039157)
	min, max := big.NewInt(10039156), big.NewInt(10105118)
	txsNum := 0
	for i := min; i.Cmp(max) == -1; i = i.Add(i, big.NewInt(1)) {
		GetTransactionByBlockNumber(i, &txsNum)
		fmt.Println(txsNum)
	}
	writerCsv.Flush() //刷新，不刷新是无法写入的
}

package rwrate

import (
	"encoding/csv"
	"fmt"
	"github.com/DarcyWep/morph/storagermw"
	"log"
	"math/big"
	"os"
	"testing"
)

func Test(t *testing.T) {
	SetHost("192.168.5.29")
	err := OpenSqlServers()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer CloseSqlServers()

	storagermw.SetHost("192.168.5.29")
	err = storagermw.OpenSqlServers()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer storagermw.CloseSqlServers()

	//OpenFile读取文件，不存在时则创建，使用追加模式
	path := "/Users/darcywep/Projects/GoProjects/morph/rwrate/rw_rate.csv"
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("文件打开失败！")
	}
	defer file.Close()

	var data = []string{"Number", "Read", "Write", "RWRatio"}
	//创建写入接口
	writerCsv := csv.NewWriter(file)
	err = writerCsv.Write(data)
	if err != nil {
		log.Println("WriterCsv写入文件失败")
	}
	writerCsv.Flush() //刷新，不刷新是无法写入的

	//min, max := big.NewInt(7202102), big.NewInt(8852777)
	min, max := big.NewInt(10050000), big.NewInt(10080001)
	txsNum := 0
	for i := min; i.Cmp(max) == -1; i = i.Add(i, big.NewInt(1)) {
		rwRate := GetRWRateByBlockNumber(i)
		StorageRead, StorageWrite := storagermw.GetTransactionByBlockNumber(i, &txsNum)
		if rwRate == nil && StorageRead + StorageWrite == 0 {
			//fmt.Println(i, "StateRead", 0, "StateWrite", 0, "StorageRead", 0, "StorageWrite", 0)
			data = []string{i.String(), "0", "0", "0"}
		} else {
			read, write := rwRate.StateRead + StorageRead, rwRate.StateWrite + StorageWrite
			sum := read + write
			var rate string
			if sum > 0 {
				r, w := read*100/sum, write*100/sum
				if r+w < 100 {
					w += 1
				}
				rate = fmt.Sprintf("%d:%d", r, w)

			} else {
				rate = "0:0"
			}
			data = []string{i.String(), fmt.Sprintf("%v", read), fmt.Sprintf("%v", write), rate}
			//fmt.Println(data)
			//fmt.Println(i, "StateRead", rwRate.StateRead, "StateWrite", rwRate.StateWrite,
			//	"StorageRead", rwRate.StorageRead, "StorageWrite", rwRate.StorageWrite,
			//	"StateRWRate", StateRWRate, "StorageRWRate", StorageRWRate)
		}
		err = writerCsv.Write(data)
		if err != nil {
			log.Println("WriterCsv写入文件失败")
		}
	}

	//for i := min; i.Cmp(max) == -1; i = i.Add(i, big.NewInt(1)) {
	//	rwRate := GetRWRateByBlockNumber(i)
	//	StorageRead, StorageWrite := storagermw.GetTransactionByBlockNumber(i)
	//	if rwRate == nil &&  {
	//		//fmt.Println(i, "StateRead", 0, "StateWrite", 0, "StorageRead", 0, "StorageWrite", 0)
	//		data = []string{i.String(), "0", "0", "0"}
	//	} else {
	//		StateNum := rwRate.StateRead + rwRate.StateWrite
	//		StorageNum := rwRate.StorageRead + rwRate.StorageWrite
	//		var StateRWRate, StorageRWRate string
	//		if StateNum > 0 {
	//			r, w := rwRate.StateRead*100/StateNum, rwRate.StateWrite*100/StateNum
	//			if r+w < 100 {
	//				w += 1
	//			}
	//			StateRWRate = fmt.Sprintf("%d:%d", r, w)
	//
	//		} else {
	//			StateRWRate = "0:0"
	//		}
	//		if StorageNum != 0 {
	//			r, w := rwRate.StorageRead*100/StorageNum, rwRate.StorageWrite*100/StorageNum
	//			if r+w < 100 {
	//				w += 1
	//			}
	//			StorageRWRate = fmt.Sprintf("%d:%d", r, w)
	//		} else {
	//			StorageRWRate = "0:0"
	//		}
	//		StorageRWRate = StorageRWRate
	//		data = []string{i.String(), fmt.Sprintf("%v", rwRate.StateRead), fmt.Sprintf("%v", rwRate.StateWrite), StateRWRate}
	//		//fmt.Println(i, "StateRead", rwRate.StateRead, "StateWrite", rwRate.StateWrite,
	//		//	"StorageRead", rwRate.StorageRead, "StorageWrite", rwRate.StorageWrite,
	//		//	"StateRWRate", StateRWRate, "StorageRWRate", StorageRWRate)
	//	}
	//	err = writerCsv.Write(data)
	//	if err != nil {
	//		log.Println("WriterCsv写入文件失败")
	//	}
	//}
	fmt.Println("Number of transaction:", txsNum)
	writerCsv.Flush() //刷新，不刷新是无法写入的
}

//
//3397102
//
//4027203

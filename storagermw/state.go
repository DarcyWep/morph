package storagermw

import (
	"strconv"
	"strings"

	//"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// 以下是导出交易的相关字段需求 (注: 结构体里面的数据要给json包访问 -> 需要首字母大写)

type stateTransition struct {
	label uint8 // 0: 普通转账(state), 1: ERC20类转账(storage), 2: KECCAK256, 3: push20, 4: ContractCall

	from  *balance
	to    *balance
	value *big.Int
	// TODO type 类型(0: Storage类型)
	// type 类型(1: 合约的调用者转账给某一接收方, 可能是嵌套合约的调用; 2: 合约创建者发送到合约账户; 3: 将手续费添加给矿工, 只有To字段)
	// type 类型(4: 多扣除的手续费退还, 只有To字段; 5: 从交易发起者账户预扣除交易费, 只有From字段; 6: 合约销毁)
	// type 类型(7: 给挖出叔父区块的矿工奖励, 只有To字段; 8: 给挖出该区块的矿工奖励, 只有To字段; 9: The DAO 硬分叉相关)
	// 类型 7, 8 放在了该区块最后一个事务的交易中, 最后处理
	type_ uint8 // 一个交易中, 3 4 5 类型最多只有一个, 且需要结合起来看(5-4: 交易发起者需要扣除的手续费 != 3: 给矿工的手续费)
}

func newStateTransition(info string) *stateTransition {
	infos := strings.Split(info, ",")

	tmp := big.NewInt(0)
	type_, _ := strconv.ParseUint(infos[1], 10, 8)
	value, _ := tmp.SetString(infos[4], 10)
	var t = &stateTransition{
		label: 0,

		value: value,
		type_: uint8(type_),
	}
	if infos[2] == ""{
		t.from = nil
	}else {
		t.from = newBalance(infos[2])
	}

	if infos[3] == ""{
		t.to = nil
	}else {
		t.to = newBalance(infos[3])
	}
	return t
}

func (t *stateTransition) GetLabel() uint8 {
	return t.label
}

type balance struct {
	address  string
	beforeTx *big.Int
}

func newBalance(info string) *balance {
	infos := strings.Split(info, "~")
	tmp := big.NewInt(0)
	beforeTx, _ := tmp.SetString(infos[1], 10)

	return &balance{
		address: infos[0],
		beforeTx: beforeTx,
	}
}

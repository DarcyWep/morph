package storagermw

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
)

type storageTransition struct {
	label uint8 // 0: 普通转账(state), 1: ERC20类转账(storage), 2: KECCAK256, 3: push20, 4: ContractCall

	contractAddress *common.Address
	slot            *common.Hash // 智能合约的存储槽
	preValue        *common.Hash
	newValue        *common.Hash // newValue = nil 则是 SLOAD, 否则为 SSTORE
}

func newStorageTransition(info string) *storageTransition {
	infos := strings.Split(info, ",")
	slot := common.HexToHash(infos[1])

	tmp := big.NewInt(0)
	tmp, _ = tmp.SetString(infos[2], 10)
	preValue := common.BigToHash(tmp)
	 st := &storageTransition{
		label: 1,
		slot:     &slot,
		preValue: &preValue,
	}

	if infos[3] == ""{
		st.newValue = nil
	}else {
		tmp, _ = tmp.SetString(infos[3], 10)
		newValue := common.BigToHash(tmp)
		st.newValue = &newValue
	}
	return st
}

func (s *storageTransition) GetLabel() uint8 {
	return s.label
}

type contractCall struct {
	label uint8 // 0: 普通转账(state), 1: ERC20类转账(storage), 2: KECCAK256, 3: push20, 4: ContractCall

	address *common.Address
}

func newContractCall(info string) *contractCall {
	infos := strings.Split(info, ",")
	address := common.HexToAddress(infos[1])
	return &contractCall{
		label: 4,
		address: &address,
	}
}

func (c *contractCall) GetLabel() uint8 {
	return c.label
}
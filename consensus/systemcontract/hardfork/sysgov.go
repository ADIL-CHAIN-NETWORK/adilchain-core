package hardfork

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

var (
	addr = common.HexToAddress("0x000000000000000000000000000000000000f000")
)

const (
	code = "0x608060405234801561001057600080fd5b50600436106101215760003560e01c8063be645692116100ad578063e1607fb711610071578063e1607fb71461033c578063efd8d8e2146103cc578063f851a440146103d4578063fb48270c146103dc578063fbb847e1146103e457610121565b8063be64569214610297578063c4d66de8146102b1578063c967f90f146102d7578063db78dd28146102f6578063e08b1d381461031b57610121565b806326782247116100f4578063267822471461023c5780633656de21146102445780633a061bd3146102615780634fb9e9b7146102695780636233be5d1461028f57610121565b806305b8481014610126578063158ef93e146101dd5780631b5e358c146101f9578063232e5ffc1461021d575b600080fd5b6101496004803603602081101561013c57600080fd5b503563ffffffff166103ec565b60405180868152602001856001600160a01b03166001600160a01b03168152602001846001600160a01b03166001600160a01b0316815260200183815260200180602001828103825283818151815260200191508051906020019060200280838360005b838110156101c55781810151838201526020016101ad565b50505050905001965050505050505060405180910390f35b6101e5610556565b604080519115158252519081900360200190f35b61020161055f565b604080516001600160a01b039092168252519081900360200190f35b61023a6004803603602081101561023357600080fd5b5035610565565b005b61020161071c565b6101496004803603602081101561025a57600080fd5b503561072b565b610201610795565b61023a6004803603602081101561027f57600080fd5b50356001600160a01b031661079b565b610201610836565b61029f61083c565b60408051918252519081900360200190f35b61023a600480360360208110156102c757600080fd5b50356001600160a01b0316610849565b6102df6108c6565b6040805161ffff9092168252519081900360200190f35b6102fe6108cb565b6040805167ffffffffffffffff9092168252519081900360200190f35b6103236108d2565b6040805163ffffffff9092168252519081900360200190f35b61023a6004803603608081101561035257600080fd5b6001600160a01b0382358116926020810135909116916040820135919081019060808101606082013564010000000081111561038d57600080fd5b82018360208201111561039f57600080fd5b803590602001918460208302840111640100000000831117156103c157600080fd5b5090925090506108d9565b6102fe610c0d565b610201610c13565b61023a610c27565b61029f610ce1565b60008060008060606003805490508663ffffffff1610610448576040805162461bcd60e51b8152602060048201526012602482015271496e646578206f7574206f662072616e676560701b604482015290519081900360640190fd5b610450610ce7565b60038763ffffffff168154811061046357fe5b60009182526020918290206040805160a08101825260059093029091018054835260018101546001600160a01b039081168486015260028201541683830152600381015460608401526004810180548351818702810187019094528084529394919360808601939283018282801561052257602002820191906000526020600020906000905b825461010083900a900460f81b6001600160f81b0319168152602060019283018181049485019490930390920291018084116104e95790505b5050509190925250508151602083015160408401516060850151608090950151929c919b5099509297509550909350505050565b60005460ff1681565b61f00181565b3341146105a6576040805162461bcd60e51b815260206004820152600a6024820152694d696e6572206f6e6c7960b01b604482015290519081900360640190fd5b60005b6003548110156107185781600382815481106105c157fe5b906000526020600020906005020160000154141561071057600354600019018114610684576003805460001981019081106105f857fe5b90600052602060002090600502016003828154811061061357fe5b60009182526020909120825460059092020190815560018083015490820180546001600160a01b039283166001600160a01b031991821617909155600280850154908401805491909316911617905560038083015490820155600480830180546106809284019190610d28565b5050505b600380548061068f57fe5b6000828152602081206005600019909301928302018181556001810180546001600160a01b0319908116909155600282018054909116905560038101829055906106dc6004830182610d86565b5050905560405182907fc2946e69de813a7cede502a3b315aa221abf9fcca5c7134b0ae6b2c3857cf63d90600090a2610718565b6001016105a9565b5050565b6001546001600160a01b031681565b60008060008060606002805490508610610780576040805162461bcd60e51b8152602060048201526011602482015270125908191bd95cc81b9bdd08195e1a5cdd607a1b604482015290519081900360640190fd5b610788610ce7565b6002878154811061046357fe5b61f00081565b60005461010090046001600160a01b031633146107ec576040805162461bcd60e51b815260206004820152600a60248201526941646d696e206f6e6c7960b01b604482015290519081900360640190fd5b600180546001600160a01b0319166001600160a01b0383169081179091556040517faefcaa6215f99fe8c2f605dd268ee4d23a5b596bbca026e25ce8446187f4f1ba90600090a250565b61f00281565b6801bc16d674ec80000081565b60005460ff1615610897576040805162461bcd60e51b8152602060048201526013602482015272105b1c9958591e481a5b9a5d1a585b1a5e9959606a1b604482015290519081900360640190fd5b6000805460ff196001600160a01b0390931661010002610100600160a81b031990911617919091166001179055565b601581565b6201518081565b6003545b90565b60005461010090046001600160a01b0316331461092a576040805162461bcd60e51b815260206004820152600a60248201526941646d696e206f6e6c7960b01b604482015290519081900360640190fd5b600254610935610ce7565b6040518060a00160405280838152602001886001600160a01b03168152602001876001600160a01b031681526020018681526020018585808060200260200160405190810160405280939291908181526020018383602002808284376000920182905250939094525050600280546001810182559152825160059091027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace81019182556020808501517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5acf830180546001600160a01b039283166001600160a01b03199182161790915560408701517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad08501805491909316911617905560608501517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad1830155608085015180519596508695939450610aba937f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad290930192910190610dae565b505060038054600181018255600091909152825160059091027fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b81019182556020808501517fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85c830180546001600160a01b039283166001600160a01b03199182161790915560408701517fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85d8501805491909316911617905560608501517fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85e83015560808501518051869550610bd6937fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85f01929190910190610dae565b50506040518391507f14ca27cd9911371c77ed1cf3cee0a4320613b07478668958754713b3c880cd0f90600090a250505050505050565b61708081565b60005461010090046001600160a01b031681565b6001546001600160a01b03163314610c77576040805162461bcd60e51b815260206004820152600e60248201526d4e65772061646d696e206f6e6c7960901b604482015290519081900360640190fd5b60018054600080546001600160a01b03808416610100908102610100600160a81b0319909316929092178084556001600160a01b03199094169094556040519204909216917f7ce7ec0b50378fb6c0186ffb5f48325f6593fcb4ca4386f21861af3129188f5c91a2565b60025490565b6040518060a001604052806000815260200160006001600160a01b0316815260200160006001600160a01b0316815260200160008152602001606081525090565b82805482825590600052602060002090601f01602090048101928215610d7657600052602060002091601f016020900482015b82811115610d76578254825591600101919060010190610d5b565b50610d82929150610e4b565b5090565b50805460008255601f016020900490600052602060002090810190610dab9190610e69565b50565b82805482825590600052602060002090601f01602090048101928215610d765791602002820160005b83821115610e1557835183826101000a81548160ff021916908360f81c02179055509260200192600101602081600001049283019260010302610dd7565b8015610e425782816101000a81549060ff0219169055600101602081600001049283019260010302610e15565b5050610d829291505b6108d691905b80821115610d8257805460ff19168155600101610e51565b6108d691905b80821115610d825760008155600101610e6f56fea26469706673582212205bbe3039cd7c79c5c55d2c7251da4df29473d056a275c66c31131ccaf6870da564736f6c63430006010033"
)

type SysGov struct {
}

func (s *SysGov) GetName() string {
	return "SysGov"
}

func (s *SysGov) Execute(config *params.ChainConfig, height *big.Int, state *state.StateDB) (err error) {
	contractCode, err := hex.DecodeString(code)
	if err != nil {
		return
	}

	//write code to sys contract
	state.SetCode(addr, contractCode)
	log.Debug("Write code to system contract account", "addr", addr.String(), "code", code)

	// init contract

	return
}

package congress

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/congress/systemcontract"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"math"
	"math/big"
)

// Proposal is the system governance proposal info.
type Proposal struct {
	Id    *big.Int
	From  common.Address
	To    common.Address
	Value *big.Int
	Data  []byte
}

func (c *Congress) getPassedProposalCount(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB) (uint32, error) {

	method := "getPassedProposalCount"
	data, err := c.abi[systemcontract.SysGovContractName].Pack(method)
	if err != nil {
		log.Error("Can't pack data for getPassedProposalCount", "error", err)
		return 0, err
	}

	msg := types.NewMessage(header.Coinbase, &systemcontract.SysGovContractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)

	// use parent
	result, err := executeMsg(msg, state, header, newChainContext(chain, c), c.chainConfig)
	if err != nil {
		return 0, err
	}

	// unpack data
	ret, err := c.abi[systemcontract.SysGovContractName].Unpack(method, result)
	if err != nil {
		return 0, err
	}
	if len(ret) != 1 {
		return 0, errors.New("invalid output length")
	}
	count, ok := ret[0].(uint32)
	if !ok {
		return 0, errors.New("invalid count format")
	}

	return count, nil
}

func (c *Congress) getPassedProposalByIndex(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, idx uint32) (*Proposal, error) {

	method := "getPassedProposalByIndex"
	data, err := c.abi[systemcontract.SysGovContractName].Pack(method, idx)
	if err != nil {
		log.Error("Can't pack data for getPassedProposalByIndex", "error", err)
		return nil, err
	}

	msg := types.NewMessage(header.Coinbase, &systemcontract.SysGovContractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)

	// use parent
	result, err := executeMsg(msg, state, header, newChainContext(chain, c), c.chainConfig)
	if err != nil {
		return nil, err
	}

	// unpack data
	prop := &Proposal{}
	err = c.abi[systemcontract.SysGovContractName].UnpackIntoInterface(prop, method, result)
	if err != nil {
		return nil, err
	}

	return prop, nil
}

//finishProposalById
func (c *Congress) finishProposalById(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, id *big.Int) error {
	method := "finishProposalById"
	data, err := c.abi[systemcontract.SysGovContractName].Pack(method, id)
	if err != nil {
		log.Error("Can't pack data for getPassedProposalByIndex", "error", err)
		return err
	}

	msg := types.NewMessage(header.Coinbase, &systemcontract.SysGovContractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)

	// use parent
	_, err = executeMsg(msg, state, header, newChainContext(chain, c), c.chainConfig)
	if err != nil {
		return err
	}

	return nil
}

func (c *Congress) executeProposal(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, prop *Proposal, totalTxIndex int) (*types.Transaction, *types.Receipt, error) {
	propRLP, err := rlp.EncodeToBytes(prop)
	if err != nil {
		return nil, nil, err
	}
	//make system governance transaction
	nonce := state.GetNonce(c.validator)
	tx := types.NewTransaction(nonce, sysGovToAddr, prop.Value, header.GasLimit, new(big.Int), propRLP)
	tx, err = c.signTxFn(accounts.Account{Address: c.validator}, tx, chain.Config().ChainID)
	if err != nil {
		return nil, nil, err
	}
	//add nonce for validator
	state.SetNonce(c.validator, nonce+1)
	receipt := c.executeProposalMsg(chain, header, state, prop, totalTxIndex, tx.Hash(), common.Hash{})

	return tx, receipt, nil
}

func (c *Congress) replayProposal(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, prop *Proposal, totalTxIndex int, tx *types.Transaction) (*types.Receipt, error) {
	sender, err := types.Sender(c.signer, tx)
	if err != nil {
		return nil, err
	}
	if sender != header.Coinbase {
		return nil, errors.New("invalid sender for system governance transaction")
	}
	propRLP, err := rlp.EncodeToBytes(prop)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(propRLP, tx.Data()) {
		return nil, fmt.Errorf("data missmatch, proposalID: %s, rlp: %s, txHash:%s, txData:%s", prop.Id.String(), hexutil.Encode(propRLP), tx.Hash().String(), hexutil.Encode(tx.Data()))
	}
	//make system governance transaction
	nonce := state.GetNonce(sender)
	//add nonce for validator
	state.SetNonce(sender, nonce+1)
	receipt := c.executeProposalMsg(chain, header, state, prop, totalTxIndex, tx.Hash(), header.Hash())

	return receipt, nil
}

func (c *Congress) executeProposalMsg(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, prop *Proposal, totalTxIndex int, txHash, bHash common.Hash) *types.Receipt {
	// actually run the governance message
	msg := types.NewMessage(prop.From, &prop.To, 0, prop.Value, header.GasLimit, new(big.Int), prop.Data, false)
	state.Prepare(txHash, bHash, totalTxIndex)
	_, err := executeMsg(msg, state, header, newChainContext(chain, c), c.chainConfig)
	state.Finalise(true)

	// governance message will not actually consumes gas
	receipt := types.NewReceipt([]byte{}, err != nil, header.GasUsed)
	receipt.TxHash = txHash
	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = state.GetLogs(txHash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = state.BlockHash()
	receipt.BlockNumber = header.Number
	receipt.TransactionIndex = uint(state.TxIndex())

	return receipt
}

// Copyright (c) 2023-2024 Nibi, Inc.
package rpcapi

import (
	"context"

	gethmath "github.com/ethereum/go-ethereum/common/math"
	gethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core/types"

	"github.com/NibiruChain/nibiru/v2/eth/rpc/backend"

	"github.com/NibiruChain/nibiru/v2/eth"
	"github.com/NibiruChain/nibiru/v2/eth/rpc"
	"github.com/NibiruChain/nibiru/v2/x/evm"
)

type IEthAPI interface {
	// Getting Blocks
	//
	// Retrieves information from a particular block in the blockchain.
	BlockNumber() (hexutil.Uint64, error)
	GetBlockByNumber(ethBlockNum rpc.BlockNumber, fullTx bool) (map[string]any, error)
	GetBlockByHash(hash common.Hash, fullTx bool) (map[string]any, error)
	GetBlockTransactionCountByHash(hash common.Hash) (*hexutil.Uint, error)
	GetBlockTransactionCountByNumber(blockNum rpc.BlockNumber) (*hexutil.Uint, error)

	// Reading Transactions
	//
	// Retrieves information on the state data for addresses regardless of whether
	// it is a user or a smart contract.
	GetTransactionByHash(hash common.Hash) (*rpc.EthTxJsonRPC, error)
	GetTransactionCount(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error)
	GetTransactionReceipt(hash common.Hash) (*backend.TransactionReceipt, error)
	GetTransactionByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) (*rpc.EthTxJsonRPC, error)
	GetTransactionByBlockNumberAndIndex(blockNum rpc.BlockNumber, idx hexutil.Uint) (*rpc.EthTxJsonRPC, error)
	// eth_getBlockReceipts

	// Account Information
	//
	// Returns information regarding an address's stored on-chain data.
	Accounts() ([]common.Address, error)
	GetBalance(
		address common.Address, blockNrOrHash rpc.BlockNumberOrHash,
	) (*hexutil.Big, error)
	GetStorageAt(
		address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash,
	) (hexutil.Bytes, error)
	GetCode(
		address common.Address, blockNrOrHash rpc.BlockNumberOrHash,
	) (hexutil.Bytes, error)
	GetProof(
		address common.Address, storageKeys []string, blockNrOrHash rpc.BlockNumberOrHash,
	) (*rpc.AccountResult, error)

	// Chain Information
	//
	// Returns information on the Ethereum network and internal settings.
	ProtocolVersion() hexutil.Uint
	GasPrice() (*hexutil.Big, error)
	EstimateGas(
		args evm.JsonTxArgs, blockNrOptional *rpc.BlockNumber,
	) (hexutil.Uint64, error)
	FeeHistory(
		blockCount gethmath.HexOrDecimal64,
		lastBlock gethrpc.BlockNumber,
		rewardPercentiles []float64,
	) (*rpc.FeeHistoryResult, error)
	MaxPriorityFeePerGas() (*hexutil.Big, error)
	ChainId() (*hexutil.Big, error)

	// Getting Uncles
	//
	// Returns information on uncle blocks are which are network rejected blocks
	// and replaced by a canonical block instead.
	GetUncleByBlockHashAndIndex(
		hash common.Hash, idx hexutil.Uint,
	) map[string]any
	GetUncleByBlockNumberAndIndex(
		number, idx hexutil.Uint,
	) map[string]any
	GetUncleCountByBlockHash(hash common.Hash) hexutil.Uint
	GetUncleCountByBlockNumber(blockNum rpc.BlockNumber) hexutil.Uint

	// Other
	Syncing() (any, error)
	GetTransactionLogs(txHash common.Hash) ([]*gethcore.Log, error)
	FillTransaction(
		args evm.JsonTxArgs,
	) (*rpc.SignTransactionResult, error)
	GetPendingTransactions() ([]*rpc.EthTxJsonRPC, error)
}

var _ IEthAPI = (*EthAPI)(nil)

// EthAPI: Allows connection to a full node of the Nibiru blockchain
// network via Nibiru EVM. Developers can interact with on-chain EVM data and
// send different types of transactions to the network by utilizing the endpoints
// provided by the API.
//
// [EthAPI] contains much of the "eth_" prefixed methods in the Web3 JSON-RPC spec.
//
// The API follows a JSON-RPC standard. If not otherwise
// specified, the interface is derived from the Alchemy Ethereum API:
// https://docs.alchemy.com/alchemy/apis/ethereum
type EthAPI struct {
	ctx     context.Context
	logger  log.Logger
	backend *backend.Backend
}

// NewImplEthAPI creates an instance of the public ETH Web3 API.
func NewImplEthAPI(logger log.Logger, backend *backend.Backend) *EthAPI {
	api := &EthAPI{
		ctx:     context.Background(),
		logger:  logger.With("client", "json-rpc"),
		backend: backend,
	}

	return api
}

// --------------------------------------------------------------------------
//                           Blocks
// --------------------------------------------------------------------------

// BlockNumber returns the current block number.
func (e *EthAPI) BlockNumber() (hexutil.Uint64, error) {
	e.logger.Debug("eth_blockNumber")
	return e.backend.BlockNumber()
}

// GetBlockByNumber returns the block identified by number.
func (e *EthAPI) GetBlockByNumber(ethBlockNum rpc.BlockNumber, fullTx bool) (map[string]any, error) {
	e.logger.Debug("eth_getBlockByNumber", "blockNumber", ethBlockNum, "fullTx", fullTx)
	return e.backend.GetBlockByNumber(ethBlockNum, fullTx)
}

// GetBlockByHash returns the block identified by hash.
func (e *EthAPI) GetBlockByHash(hash common.Hash, fullTx bool) (map[string]any, error) {
	methodName := "eth_getBlockByHash"
	e.logger.Debug(methodName, "hash", hash.Hex(), "fullTx", fullTx)
	block, err := e.backend.GetBlockByHash(hash, fullTx)
	logError(e.logger, err, methodName)
	return block, err
}

// logError logs a backend error if one is present
func logError(logger log.Logger, err error, methodName string) {
	if err != nil {
		logger.Debug(methodName+" failed", "error", err.Error())
	}
}

// --------------------------------------------------------------------------
//                           Read Txs
// --------------------------------------------------------------------------

// GetTransactionByHash returns the transaction identified by hash.
func (e *EthAPI) GetTransactionByHash(hash common.Hash) (*rpc.EthTxJsonRPC, error) {
	methodName := "eth_getTransactionByHash"
	e.logger.Debug(methodName, "hash", hash.Hex())
	tx, err := e.backend.GetTransactionByHash(hash)
	logError(e.logger, err, methodName)
	return tx, err
}

// GetTransactionCount returns the number of transactions at the given address up to the given block number.
func (e *EthAPI) GetTransactionCount(
	address common.Address, blockNrOrHash rpc.BlockNumberOrHash,
) (*hexutil.Uint64, error) {
	e.logger.Debug("eth_getTransactionCount", "address", address.Hex(), "block number or hash", blockNrOrHash)
	blockNum, err := e.backend.BlockNumberFromTendermint(blockNrOrHash)
	if err != nil {
		return nil, err
	}
	return e.backend.GetTransactionCount(address, blockNum)
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (e *EthAPI) GetTransactionReceipt(
	hash common.Hash,
) (*backend.TransactionReceipt, error) {
	hexTx := hash.Hex()
	e.logger.Debug("eth_getTransactionReceipt", "hash", hexTx)
	return e.backend.GetTransactionReceipt(hash)
}

// GetBlockTransactionCountByHash returns the number of transactions in the block identified by hash.
func (e *EthAPI) GetBlockTransactionCountByHash(hash common.Hash) (*hexutil.Uint, error) {
	methodName := "eth_getBlockTransactionCountByHash"
	e.logger.Debug(methodName, "hash", hash.Hex())
	txCount, err := e.backend.GetBlockTransactionCountByHash(hash)
	logError(e.logger, err, methodName)
	return txCount, err
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block identified by number.
func (e *EthAPI) GetBlockTransactionCountByNumber(
	blockNum rpc.BlockNumber,
) (*hexutil.Uint, error) {
	methodName := "eth_getBlockTransactionCountByNumber"
	e.logger.Debug(methodName, "height", blockNum.Int64())
	txCount, err := e.backend.GetBlockTransactionCountByNumber(blockNum)
	logError(e.logger, err, methodName)
	return txCount, err
}

// GetTransactionByBlockHashAndIndex returns the transaction identified by hash and index.
func (e *EthAPI) GetTransactionByBlockHashAndIndex(
	hash common.Hash, idx hexutil.Uint,
) (*rpc.EthTxJsonRPC, error) {
	e.logger.Debug("eth_getTransactionByBlockHashAndIndex", "hash", hash.Hex(), "index", idx)
	return e.backend.GetTransactionByBlockHashAndIndex(hash, idx)
}

// GetTransactionByBlockNumberAndIndex returns the transaction identified by number and index.
func (e *EthAPI) GetTransactionByBlockNumberAndIndex(
	blockNum rpc.BlockNumber, idx hexutil.Uint,
) (*rpc.EthTxJsonRPC, error) {
	e.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "number", blockNum, "index", idx)
	return e.backend.GetTransactionByBlockNumberAndIndex(blockNum, idx)
}

// --------------------------------------------------------------------------
//                           Write Txs
// --------------------------------------------------------------------------

// SendRawTransaction send a raw Ethereum transaction.
// Allows developers to both send ETH from one address to another, write data
// on-chain, and interact with smart contracts.
func (e *EthAPI) SendRawTransaction(data hexutil.Bytes) (common.Hash, error) {
	e.logger.Debug("eth_sendRawTransaction", "length", len(data))
	return e.backend.SendRawTransaction(data)
}

// --------------------------------------------------------------------------
//                           Account Information
// --------------------------------------------------------------------------

// Accounts returns the list of accounts available to this node.
func (e *EthAPI) Accounts() ([]common.Address, error) {
	e.logger.Debug("eth_accounts")
	return e.backend.Accounts()
}

// GetBalance returns the provided account's balance up to the provided block number.
func (e *EthAPI) GetBalance(
	address common.Address, blockNrOrHash rpc.BlockNumberOrHash,
) (*hexutil.Big, error) {
	e.logger.Debug("eth_getBalance", "address", address.String(), "block number or hash", blockNrOrHash)
	return e.backend.GetBalance(address, blockNrOrHash)
}

// GetStorageAt returns the contract storage at the given address, block number, and key.
func (e *EthAPI) GetStorageAt(
	address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash,
) (hexutil.Bytes, error) {
	e.logger.Debug("eth_getStorageAt", "address", address.Hex(), "key", key, "block number or hash", blockNrOrHash)
	return e.backend.GetStorageAt(address, key, blockNrOrHash)
}

// GetCode returns the contract code at the given address and block number.
func (e *EthAPI) GetCode(
	address common.Address, blockNrOrHash rpc.BlockNumberOrHash,
) (hexutil.Bytes, error) {
	e.logger.Debug("eth_getCode", "address", address.Hex(), "block number or hash", blockNrOrHash)
	return e.backend.GetCode(address, blockNrOrHash)
}

// GetProof returns an account object with proof and any storage proofs
func (e *EthAPI) GetProof(address common.Address,
	storageKeys []string,
	blockNrOrHash rpc.BlockNumberOrHash,
) (*rpc.AccountResult, error) {
	e.logger.Debug("eth_getProof", "address", address.Hex(), "keys", storageKeys, "block number or hash", blockNrOrHash)
	return e.backend.GetProof(address, storageKeys, blockNrOrHash)
}

// --------------------------------------------------------------------------
//                           EVM/Smart Contract Execution
// --------------------------------------------------------------------------

// Call performs a raw contract call.
//
// Allows developers to read data from the blockchain which includes executing
// smart contracts. However, no data is published to the blockchain network.
func (e *EthAPI) Call(args evm.JsonTxArgs,
	blockNrOrHash rpc.BlockNumberOrHash,
	_ *rpc.StateOverride,
) (hexutil.Bytes, error) {
	e.logger.Debug("eth_call", "args", args.String(), "block number or hash", blockNrOrHash)

	blockNum, err := e.backend.BlockNumberFromTendermint(blockNrOrHash)
	if err != nil {
		return nil, err
	}
	data, err := e.backend.DoCall(args, blockNum)
	if err != nil {
		return []byte{}, err
	}

	return (hexutil.Bytes)(data.Ret), nil
}

// --------------------------------------------------------------------------
//                           Event Logs
// --------------------------------------------------------------------------
// FILTER API at ./filters/api.go

// --------------------------------------------------------------------------
//                           Chain Information
// --------------------------------------------------------------------------

// ProtocolVersion returns the supported Ethereum protocol version.
func (e *EthAPI) ProtocolVersion() hexutil.Uint {
	e.logger.Debug("eth_protocolVersion")
	return hexutil.Uint(eth.ProtocolVersion)
}

// GasPrice returns the current gas price based on Ethermint's gas price oracle.
func (e *EthAPI) GasPrice() (*hexutil.Big, error) {
	e.logger.Debug("eth_gasPrice")
	return e.backend.GasPrice()
}

// EstimateGas returns an estimate of gas usage for the given smart contract call.
func (e *EthAPI) EstimateGas(
	args evm.JsonTxArgs, blockNrOptional *rpc.BlockNumber,
) (hexutil.Uint64, error) {
	e.logger.Debug("eth_estimateGas")
	return e.backend.EstimateGas(args, blockNrOptional)
}

func (e *EthAPI) FeeHistory(blockCount gethmath.HexOrDecimal64,
	lastBlock gethrpc.BlockNumber,
	rewardPercentiles []float64,
) (*rpc.FeeHistoryResult, error) {
	e.logger.Debug("eth_feeHistory")
	return e.backend.FeeHistory(blockCount, lastBlock, rewardPercentiles)
}

// MaxPriorityFeePerGas returns a suggestion for a gas tip cap for dynamic fee
// transactions.
func (e *EthAPI) MaxPriorityFeePerGas() (*hexutil.Big, error) {
	e.logger.Debug("eth_maxPriorityFeePerGas")
	head, err := e.backend.CurrentHeader()
	if err != nil {
		return nil, err
	}
	tipcap, err := e.backend.SuggestGasTipCap(head.BaseFee)
	if err != nil {
		return nil, err
	}
	return (*hexutil.Big)(tipcap), nil
}

// ChainId is the EIP-155 replay-protection chain id for the current ethereum
// chain config.
func (e *EthAPI) ChainId() (*hexutil.Big, error) { //nolint
	e.logger.Debug("eth_chainId")
	return e.backend.ChainID(), nil
}

// --------------------------------------------------------------------------
//                           Uncles
// --------------------------------------------------------------------------

// GetUncleByBlockHashAndIndex returns the uncle identified by hash and index.
// Always returns nil.
func (e *EthAPI) GetUncleByBlockHashAndIndex(
	_ common.Hash, _ hexutil.Uint,
) map[string]any {
	return nil
}

// GetUncleByBlockNumberAndIndex returns the uncle identified by number and
// index. Always returns nil.
func (e *EthAPI) GetUncleByBlockNumberAndIndex(
	_, _ hexutil.Uint,
) map[string]any {
	return nil
}

// GetUncleCountByBlockHash returns the number of uncles in the block identified
// by hash. Always zero.
func (e *EthAPI) GetUncleCountByBlockHash(_ common.Hash) hexutil.Uint {
	return 0
}

// GetUncleCountByBlockNumber returns the number of uncles in the block
// identified by number. Always zero.
func (e *EthAPI) GetUncleCountByBlockNumber(_ rpc.BlockNumber) hexutil.Uint {
	return 0
}

// --------------------------------------------------------------------------
//                           Other
// --------------------------------------------------------------------------

// Syncing returns false in case the node is currently not syncing with the
// network. It can be up to date or has not yet received the latest block headers
// from its pears. In case it is synchronizing:
//
// - startingBlock: block number this node started to synchronize from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (e *EthAPI) Syncing() (any, error) {
	e.logger.Debug("eth_syncing")
	return e.backend.Syncing()
}

// GetTransactionLogs returns the logs given a transaction hash.
func (e *EthAPI) GetTransactionLogs(txHash common.Hash) ([]*gethcore.Log, error) {
	e.logger.Debug("eth_getTransactionLogs", "hash", txHash)

	hexTx := txHash.Hex()
	res, err := e.backend.GetTxByEthHash(txHash)
	if err != nil {
		e.logger.Debug("tx not found", "hash", hexTx, "error", err.Error())
		return nil, nil
	}

	if res.Failed {
		// failed, return empty logs
		return nil, nil
	}

	resBlockResult, err := e.backend.TendermintBlockResultByNumber(&res.Height)
	if err != nil {
		e.logger.Debug("block result not found", "number", res.Height, "error", err.Error())
		return nil, nil
	}

	// parse tx logs from events
	index := int(res.MsgIndex) // #nosec G701
	return backend.TxLogsFromEvents(resBlockResult.TxsResults[res.TxIndex].Events, index)
}

// FillTransaction fills the defaults (nonce, gas, gasPrice or 1559 fields)
// on a given unsigned transaction, and returns it to the caller for further
// processing (signing + broadcast).
func (e *EthAPI) FillTransaction(
	args evm.JsonTxArgs,
) (*rpc.SignTransactionResult, error) {
	e.logger.Debug("eth_fillTransaction")
	// Set some sanity defaults and terminate on failure
	args, err := e.backend.SetTxDefaults(args)
	if err != nil {
		return nil, err
	}

	// Assemble the transaction and obtain rlp
	tx := args.ToMsgEthTx().AsTransaction()

	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return &rpc.SignTransactionResult{
		Raw: data,
		Tx:  tx,
	}, nil
}

// GetPendingTransactions returns the transactions that are in the transaction
// pool and have a from address that is one of the accounts this node manages.
func (e *EthAPI) GetPendingTransactions() ([]*rpc.EthTxJsonRPC, error) {
	e.logger.Debug("eth_getPendingTransactions")

	txs, err := e.backend.PendingTransactions()
	if err != nil {
		return nil, err
	}

	result := make([]*rpc.EthTxJsonRPC, 0, len(txs))
	for _, tx := range txs {
		for _, msg := range (*tx).GetMsgs() {
			ethMsg, ok := msg.(*evm.MsgEthereumTx)
			if !ok {
				// not valid ethereum tx
				break
			}

			rpctx, err := rpc.NewRPCTxFromMsgEthTx(
				ethMsg,
				common.Hash{},
				uint64(0),
				uint64(0),
				nil,
				e.backend.ChainConfig().ChainID,
			)
			if err != nil {
				return nil, err
			}

			result = append(result, rpctx)
		}
	}

	return result, nil
}

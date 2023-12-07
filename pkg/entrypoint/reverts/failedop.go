package reverts

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/rpc"
)

type FailedOpRevert struct {
	OpIndex   int
	Paymaster common.Address
	Reason    string
}

func failedOp() abi.Error {
	opIndex, _ := abi.NewType("uint256", "uint256", nil)
	paymaster, _ := abi.NewType("address", "address", nil)
	reason, _ := abi.NewType("string", "string", nil)
	return abi.NewError("FailedOp", abi.Arguments{
		{Name: "opIndex", Type: opIndex},
		{Name: "paymaster", Type: paymaster},
		{Name: "reason", Type: reason},
	})
}

func NewFailedOp(err error) (*FailedOpRevert, error) {
	rpcErr, ok := err.(rpc.DataError)
	if !ok {
		return nil, fmt.Errorf(
			"failedOp: cannot assert type: error is not of type rpc.DataError, err: %s",
			err,
		)
	}

	data, ok := rpcErr.ErrorData().(string)
	if !ok {
		return nil, fmt.Errorf(
			"failedOp: cannot assert type: data is not of type string, err: %s, data: %s",
			rpcErr.Error(),
			rpcErr.ErrorData(),
		)
	}

	failedOp := failedOp()
	revert, err := failedOp.Unpack(common.Hex2Bytes(data[2:]))
	//hash the param
	//revert, err := failedOp.Unpack(crypto.Keccak256(common.Hex2Bytes(data[2:])))
	if err != nil {
		return nil, fmt.Errorf("failedOp: %s", err)
	}

	args, ok := revert.([]any)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: args is not of type []any")
	}
	if len(args) != 3 {
		return nil, fmt.Errorf("failedOp: invalid args length: expected 3, got %d", len(args))
	}

	opIndex, ok := args[0].(*big.Int)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: opIndex is not of type *big.Int")
	}

	paymaster, ok := args[1].(common.Address)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: paymaster is not of type string")
	}

	reason, ok := args[2].(string)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: reason is not of type string")
	}

	return &FailedOpRevert{
		OpIndex:   int(opIndex.Int64()),
		Paymaster: paymaster,
		Reason:    reason,
	}, nil
}

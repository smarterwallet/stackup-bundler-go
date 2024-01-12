package filter

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"math/big"
)

func filterUserOperationEvent(
	eth *ethclient.Client,
	userOpHash string,
	entryPoint common.Address,
) (*entrypoint.EntrypointUserOperationEventIterator, error) {
	println("userOpHash", userOpHash)
	println("entryPoint", entryPoint.Hex())

	ep, err := entrypoint.NewEntrypoint(entryPoint, eth)
	if err != nil {
		return nil, err
	}
	bn, err := eth.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	toBlk := big.NewInt(0).SetUint64(bn)
	startBlk := big.NewInt(0)
	sub10kBlk := big.NewInt(0).Sub(toBlk, big.NewInt(viper.GetInt64("erc4337_bundler_get_logs_gap")))
	if sub10kBlk.Cmp(startBlk) > 0 {
		startBlk = sub10kBlk
	}

	println("startBlk", startBlk.String())
	println("toBlk", toBlk.String())

	return ep.FilterUserOperationEvent(
		&bind.FilterOpts{Start: startBlk.Uint64()},
		[][32]byte{common.HexToHash(userOpHash)},
		[]common.Address{},
		[]common.Address{},
	)
}

package balance

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/rs/zerolog/log"
	"landing_rotator/contracts"
	"math/big"
)

type BalanceManager struct {
	ClientETH      *ethclient.Client
	ClientBSC      *ethclient.Client
	ClientPolygon  *ethclient.Client
	ClientContract *ethclient.Client
}

const (
	addressMint    = "0xB58A7eC97a7778760b88fd96af5B01E13021eCc4"
	addressPresale = "0xedc81ae329638B0Dd18D01249292a27546e59F5f"
)

func NewBalanceManager(clientETH *ethclient.Client, clientBSC *ethclient.Client, clientPolygon *ethclient.Client, clientContract *ethclient.Client) BalanceManager {
	return BalanceManager{ClientETH: clientETH, ClientBSC: clientBSC, ClientPolygon: clientPolygon, ClientContract: clientContract}
}

func weiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}

func (b *BalanceManager) GetBalance(wallet string) (map[string]float64, error) {
	account := common.HexToAddress(wallet)
	balances := make(map[string]float64)

	balanceETH, err := b.ClientETH.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Err(err).Msg("failed to parse ETH Balance")
		return balances, err
	}

	balanceBSC, err := b.ClientBSC.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Err(err).Msg("failed to parse BSC Balance")
		return balances, err

	}

	balancePolygon, err := b.ClientPolygon.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Err(err).Msg("failed to parse Polygon Balance")
		return balances, err
	}
	balanceFloatETH := weiToEther(balanceETH)
	balanceFloat64ETH, _ := balanceFloatETH.Float64()

	balancefloatBSC := weiToEther(balanceBSC)
	balancefloat64BSC, _ := balancefloatBSC.Float64()

	newBalancePolygon := new(big.Float).SetInt(balancePolygon)
	balancefloat64Polygon, _ := newBalancePolygon.Float64()
	balances["ethereum"] = balanceFloat64ETH
	balances["bsc"] = balancefloat64BSC
	balances["polygon"] = balancefloat64Polygon

	return balances, nil
}
func (b *BalanceManager) GetPrice() (map[string]float64, error) {

	prices := make(map[string]float64)

	addressPresaleHex := common.HexToAddress(addressPresale)
	addressMintHex := common.HexToAddress(addressMint)

	contract, err := contracts.NewNftpresale(addressPresaleHex, b.ClientContract)

	priceFirst, err := contract.GetListings(nil, addressMintHex, big.NewInt(1))
	if err != nil {
		log.Err(err).Msg("failed to parse first balance")
		return prices, err

	}
	priceSecond, err := contract.GetListings(nil, addressMintHex, big.NewInt(2))
	if err != nil {
		log.Err(err).Msg("failed to parse second balance")
		return prices, err

	}
	priceThird, err := contract.GetListings(nil, addressMintHex, big.NewInt(3))
	if err != nil {
		log.Err(err).Msg("failed to parse third balance")
		return prices, err
	}
	priceFirstInt := new(big.Float).SetInt(priceFirst)
	priceFirstFloat, _ := priceFirstInt.Float64()

	priceSecondInt := new(big.Float).SetInt(priceSecond)
	priceSecondFloat, _ := priceSecondInt.Float64()

	priceThirdInt := new(big.Float).SetInt(priceThird)
	priceThirdFloat, _ := priceThirdInt.Float64()

	prices["pichu"] = priceFirstFloat
	prices["pikachu"] = priceSecondFloat
	prices["raichu"] = priceThirdFloat

	return prices, nil
}

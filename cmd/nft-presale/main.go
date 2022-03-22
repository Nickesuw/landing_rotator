package main

import (
	"github.com/Ferluci/ip2loc"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/zerolog/log"
	"gitlab.tubecorporate.com/platform-go/core/pkg/balancer"
	"gitlab.tubecorporate.com/platform-go/core/pkg/chlog"
	"landing_rotator/internal/balance"
	"landing_rotator/internal/config"
	"landing_rotator/internal/handlers"
	_ "net/http"
)

func initClients() (*ethclient.Client, *ethclient.Client, *ethclient.Client, *ethclient.Client) {
	clientETH, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Print(err)
	}
	clientBSC, err := ethclient.Dial("https://bsc-dataseed.binance.org/")

	if err != nil {
		log.Print(err)
	}
	clientPolygon, err := ethclient.Dial("https://polygon-rpc.com")
	if err != nil {
		log.Print(err)
	}
	clientContract, err := ethclient.Dial("https://data-seed-prebsc-1-s1.binance.org:8545/")
	if err != nil {
		log.Print(err)
	}

	return clientETH, clientBSC, clientPolygon, clientContract
}

func initGeoDb(cfg *config.Config) (*ip2loc.DB, error) {
	geoDB, err := ip2loc.OpenDB(cfg.GeoDB)
	if err != nil {
		return nil, err
	}
	return geoDB, nil
}

func main() {
	log.Logger = log.With().Caller().Logger()

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Can't read  config")
	}

	geoDb, err := initGeoDb(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init geo db")
		return
	}
	server := echo.New()
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	clickHouseClusters, err := balancer.BalanceSrv(cfg.CollectorsSRV())
	if err != nil {
		log.Fatal().Err(err).Msg("Can't resolve models srv records")
	}

	eventCollector, err := chlog.New(clickHouseClusters)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't create event collector")
	}

	clientETH, clientBSC, clientPolygon, clientContract := initClients()
	userBalance := balance.NewBalanceManager(clientETH, clientBSC, clientPolygon, clientContract)
	api := handlers.NewUserHandlers(&eventCollector, geoDb, userBalance)
	api.SetRoutes(server)
	log.Info().Msg("start server on 8001 port")
	server.Logger.Fatal(server.Start(":8001"))
}

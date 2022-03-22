package handlers

import (
	"encoding/json"
	"github.com/Ferluci/ip2loc"
	"github.com/avct/uasurfer"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
	"gitlab.tubecorporate.com/platform-go/core/pkg/chlog"
	"io/ioutil"
	"landing_rotator/internal/balance"
	"landing_rotator/internal/models"
	"landing_rotator/internal/usecase"
	"net/http"
	"time"
)

type UserHandlers struct {
	userUseCase usecase.UserUseCase
	eventLogger chlog.EventLogger
	userBalance balance.BalanceManager
}

func NewUserHandlers(eventLogger *chlog.EventLogger, geoDB *ip2loc.DB, userBalance balance.BalanceManager) UserHandlers {
	userUseCase := usecase.NewUserUseCase(geoDB)
	return UserHandlers{userUseCase: userUseCase, eventLogger: *eventLogger, userBalance: userBalance}
}

func (h UserHandlers) parseUserRequest(ctx echo.Context) (*models.UserEvent, error) {
	record, err := h.userUseCase.ParseIp(ctx.RealIP())
	if err != nil {
		return nil, err
	}

	userAgent := uasurfer.Parse(ctx.Request().UserAgent())
	userEvent := models.UserEvent{
		StatsDay:       time.Now(),
		EventTime:      time.Now(),
		UserAgent:      ctx.Request().UserAgent(),
		IP:             ctx.RealIP(),
		Country:        record.CountryShort,
		ISP:            record.Isp,
		UsageType:      record.UsageType,
		AcceptLanguage: ctx.Request().Header.Get("Accept-Language"),
		Referrer:       ctx.Request().Referer(),
		DeviceType:     userAgent.DeviceType.StringTrimPrefix(),
		BrowserName:    userAgent.Browser.Name.StringTrimPrefix(),
		BrowserVersion: userAgent.Browser.Version.Major,
		OSName:         userAgent.OS.Name.StringTrimPrefix(),
		OSVersion:      userAgent.OS.Version.Major,
		Meta: map[string]string{
			"currency": "bnb",
		},
		Balances: make(map[string]float64),
	}
	return &userEvent, nil
}

func (h UserHandlers) register(ctx echo.Context) error {
	userEvent, err := h.parseUserRequest(ctx)

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Err(err).Msg("failed to read body")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to read body",
		})
	}

	err = json.Unmarshal(body, &userEvent)
	if err != nil {
		log.Err(err).Msg("failed to unmarshal json")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to unmarshal body",
		})
	}

	log.Info().Msgf("requestdata: %#v", userEvent)
	h.eventLogger.Log(userEvent)
	return ctx.NoContent(http.StatusOK)

}

func (h UserHandlers) login(ctx echo.Context) error {
	userEvent, err := h.parseUserRequest(ctx)

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Err(err).Msg("failed to read body")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to read body",
		})

	}
	err = json.Unmarshal(body, &userEvent)
	if err != nil {
		log.Err(err).Msg("failed to unmarshal json")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to unmarshal body",
		})
	}

	userEvent.Balances, err = h.userBalance.GetBalance(userEvent.Wallet)
	if err != nil {
		log.Err(err).Msg("failed to get balances")

	}

	h.eventLogger.Log(userEvent)

	log.Info().Msgf("requestdata: %#v", userEvent)
	return ctx.NoContent(http.StatusOK)

}

func (h UserHandlers) sendPrices(ctx echo.Context) error {
	prices, err := h.userBalance.GetPrice()
	if err != nil {
		log.Err(err).Msg("failed to get prices")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to get prices",
		})
	}

	return ctx.JSONPretty(http.StatusOK, prices, "")
}

func (h UserHandlers) buy(ctx echo.Context) error {
	userEvent, err := h.parseUserRequest(ctx)
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Err(err).Msg("failed to read body")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to read body",
		})
	}
	err = json.Unmarshal(body, &userEvent)
	if err != nil {
		log.Err(err).Msg("failed to unmarshal body")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to unmarshal body",
		})
	}
	h.eventLogger.Log(userEvent)
	log.Info().Msgf("requestdata: %#v", userEvent)
	return ctx.NoContent(http.StatusOK)
}

func (h UserHandlers) SetRoutes(server *echo.Echo) {
	server.POST("/email", h.register)
	server.POST("/login", h.login)
	server.POST("/buy", h.buy)
	server.GET("/prices", h.sendPrices)
}

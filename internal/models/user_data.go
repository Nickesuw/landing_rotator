package models

import (
	"gitlab.tubecorporate.com/platform-go/core/pkg/chlog"
	"gitlab.tubecorporate.com/platform-go/core/pkg/fn"
	"time"
)

type UserEvent struct {
	EventTime      time.Time          `json:"event_time"`
	StatsDay       time.Time          `json:"stats_day"`
	EventType      string             `json:"event_type"`
	RequestID      uint64             `json:"request_id"`
	IP             string             `json:"ip"`
	Country        string             `json:"country"`
	ISP            string             `json:"isp"`
	UsageType      string             `json:"usage_type"`
	AcceptLanguage string             `json:"accept_language"`
	UserAgent      string             `json:"user_agent"`
	DeviceType     string             `json:"device_type"`
	BrowserName    string             `json:"browser_name"`
	BrowserVersion int                `json:"browser_version"`
	OSName         string             `json:"os_name"`
	OSVersion      int                `json:"os_version"`
	Referrer       string             `json:"referrer"`
	Price          float64            `json:"price"`
	TokenID        uint64             `json:"token_id"`
	Amount         uint64             `json:"amount"`
	Balances       map[string]float64 `json:"balances"`
	Blockchain     string             `json:"blockchain"`
	Wallet         string             `json:"wallet"`
	Meta           map[string]string  `json:"meta"`
	Email          string             `json:"email"`
}

func (u UserEvent) BuildQuery() (string, error) {
	return chlog.PrepareQuery(
		"INSERT INTO rnd.token_presale_local (stats_day,event_time,event_type,request_id,host, ip, country,isp,usage_type,accept_language,"+
			"user_agent,device_type,browser_name, browser_version,os_name,os_version,email,url,referrer,blockchain,wallet,balances,token_id,price,amount,meta)"+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?)",
		chlog.Date(u.StatsDay), u.EventTime, u.EventType, u.RequestID, "", fn.Ip2Int(u.IP), u.Country, u.ISP, u.UsageType, u.AcceptLanguage, u.UserAgent, u.DeviceType, u.BrowserName,
		u.BrowserVersion, u.OSName, u.OSVersion, u.Email, "", u.Referrer, u.Blockchain, u.Wallet, chlog.Map(u.Balances), u.TokenID, u.Price, u.Amount, chlog.Map(u.Meta))
}

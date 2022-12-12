package http

import (
	"net/http"
	"ot-recorder/app/model"
	"ot-recorder/infrastructure/config"
	"time"
)

type PingRequest struct {
	Type  string  `json:"_type" validate:"required"`
	Tst   int64   `json:"tst" validate:"required"`
	Acc   int8    `json:"acc"`
	Alt   int8    `json:"alt"`
	Batt  int8    `json:"batt"`
	Bs    int8    `json:"bs"`
	Lat   float64 `json:"lat" validate:"required"`
	Lon   float64 `json:"lon" validate:"required"`
	M     int8    `json:"m"`
	T     string  `json:"t"`
	Tid   string  `json:"tid"`
	Vac   int8    `json:"vac"`
	Vel   int8    `json:"vel"`
	Bssid string  `json:"BSSID"`
	Ssid  string  `json:"SSID"`
}

func mapLocationRequestToModel(req *PingRequest, headers *http.Header) *model.Location {
	createdAt := parseEpochTimeToLocal(req.Tst, config.Get().App.TimeZone)

	return &model.Location{
		Username:  headers.Get("x-limit-u"),
		Device:    headers.Get("x-limit-d"),
		CreatedAt: createdAt.Format("2006-01-02 15:04:05"),
		Acc:       req.Acc,
		Alt:       req.Alt,
		Batt:      req.Batt,
		Bs:        req.Bs,
		Lat:       req.Lat,
		Lon:       req.Lon,
		M:         req.M,
		T:         req.T,
		Tid:       req.Tid,
		Vac:       req.Vac,
		Vel:       req.Vel,
		Bssid:     req.Bssid,
		Ssid:      req.Ssid,
		IP:        headers.Get("X-Real-IP"),
	}
}

func parseEpochTimeToLocal(epoch int64, tz string) time.Time {
	et := time.Unix(epoch, 0)
	loc, _ := time.LoadLocation(tz)
	return et.In(loc)
}

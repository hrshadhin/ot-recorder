package http

import (
	"net/http"
	"ot-recorder/app/model"
)

type PingRequest struct {
	Type  string  `json:"_type" validate:"required"`
	Tst   int64   `json:"tst" validate:"required"`
	Acc   int16   `json:"acc"`
	Alt   int16   `json:"alt"`
	Batt  int8    `json:"batt"`
	Bs    int8    `json:"bs"`
	Lat   float64 `json:"lat" validate:"required"`
	Lon   float64 `json:"lon" validate:"required"`
	M     int8    `json:"m"`
	T     string  `json:"t"`
	Tid   string  `json:"tid"`
	Vac   int16   `json:"vac"`
	Vel   int16   `json:"vel"`
	Bssid string  `json:"BSSID"`
	Ssid  string  `json:"SSID"`
}

func mapLocationRequestToModel(req *PingRequest, headers *http.Header) *model.Location {
	return &model.Location{
		Username:  headers.Get("x-limit-u"),
		Device:    headers.Get("x-limit-d"),
		CreatedAt: req.Tst,
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

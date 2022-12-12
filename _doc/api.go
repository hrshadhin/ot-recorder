// @title OwnTracks Recorder API
// @version 1.0
// @description Store and access data published by OwnTracks apps.
// @x-logo {"url": "https://avatars.githubusercontent.com/u/6574523?s=200&v=4", "backgroundColor": "#FFFFFF", "altText": "OwnTracks Logo"}

// @contact.name API Support
// @contact.url https://github.com/hrshadhin/ot-recorder
// @contact.email dev@hrshadhin.me

// @license.name MIT
// @license.url https://github.com/hrshadhin/ot-recorder/blob/master/LICENSE.md

// @host localhost:8000
// @BasePath /
// @schemes http
// @accept json

// Package _doc provides basic API structs for the REST services
package _doc

// RootResponse is a data structure for the / endpoint
type RootResponse struct {
	// service details message
	Message string `json:"message"`
}

// HealthResponse is a data structure for the /h34l7h endpoint
type HealthResponse struct {
	// health details
	Message string `json:"message"`
}

type successResponse struct {
	Message string `json:"message"`
}

type failedResponse struct {
	Message string `json:"message"`
}

type badReqResponse struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

type successResponseData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type pingReq struct {
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

// Root will let you know, whoami
// @Summary Service Details
// @Description will let you know, whoami
// @Tags system
// @Produce	json
// @Success	200	{object} RootResponse "return service details"
// @Failure	500	{object} failedResponse
// @Router / [get]
func Root() {}

// Health will let you know the heart beats
// @Summary Service Health
// @Description will let you know the heart beats
// @Tags system
// @Produce	json
// @Success	200	{object} HealthResponse "return service health"
// @Failure	500	{object} failedResponse
// @Router /health [get]
func Health() {}

// Ping
// @Summary Ping Location
// @Description store location ping from mobile app
// @Tags location
// @Param x-limit-u header string true "{username}"
// @Param x-limit-d header string true "{device}"
// @Param X-Real-ID header string true "{client_ip}"
// @Accept json
// @Param payload body pingReq false "Ping Payload [Details Here](https://owntracks.org/booklet/tech/http/)"
// @Produce	json
// @Success	200	{object} []string{}
// @Failure	400 {object} badReqResponse
// @Failure	500	{object} failedResponse
// @Router /api/v1/ping [post]
func Ping() {}

// LastLocation
// @Summary last location
// @Description User last ping location
// @Tags location
// @Produce	json
// @Success	200	{object} model.LocationDetails
// @Failure	404,500	{object} failedResponse
// @Router /api/v1/last-location [get]
func LastLocation() {}

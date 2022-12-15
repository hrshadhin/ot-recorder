package model

import (
	"context"
)

type Location struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	Device    string  `json:"device"`
	CreatedAt int64   `json:"created_at"`
	Acc       int16   `json:"acc"`
	Alt       int16   `json:"alt"`
	Batt      int8    `json:"batt"`
	Bs        int8    `json:"bs"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	M         int8    `json:"m"`
	T         string  `json:"t"`
	Tid       string  `json:"tid"`
	Vac       int16   `json:"vac"`
	Vel       int16   `json:"vel"`
	Bssid     string  `json:"bssid"`
	Ssid      string  `json:"ssid"`
	IP        string  `json:"ip"`
}

type LocationDetails struct {
	Username         string  `json:"username"`
	Device           string  `json:"device"`
	DateTime         string  `json:"date_time"`
	Accuracy         int16   `json:"accuracy,omitempty"`
	Altitude         int16   `json:"altitude,omitempty"`
	BatteryLevel     string  `json:"battery_level,omitempty"`
	BatteryStatus    string  `json:"battery_status,omitempty"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Mode             string  `json:"mode,omitempty"`
	VerticalAccuracy int16   `json:"vertical_accuracy,omitempty"`
	Velocity         int16   `json:"velocity,omitempty"`
	WifiName         string  `json:"wifi_name,omitempty"`
	WifiMAC          string  `json:"wifi_mac,omitempty"`
	IPAddress        string  `json:"ip_address"`
}

type BatteryStatusEnum int
type ModeEnum int

const (
	Unknown   BatteryStatusEnum = 0
	Unplugged BatteryStatusEnum = 1
	Charging  BatteryStatusEnum = 2
	Full      BatteryStatusEnum = 3

	Significant ModeEnum = 1
	Move        ModeEnum = 2

	unknownstr = "Unknown"
)

func (e BatteryStatusEnum) String() string {
	switch e {
	case Unknown:
		return unknownstr
	case Unplugged:
		return "Unplugged"
	case Charging:
		return "Charging"
	case Full:
		return "Full"
	default:
		return unknownstr
	}
}

func (e ModeEnum) String() string {
	switch e {
	case Significant:
		return "Walking or Stand By"
	case Move:
		return "Moving"
	default:
		return unknownstr
	}
}

// LocationRepository represent the locations repository contract
type LocationRepository interface {
	CreateLocation(tx context.Context, location *Location) error
	GetUserLastLocation(tx context.Context, username string) (Location, error)
}

// LocationUsecase represent the locations usecase contract
type LocationUsecase interface {
	Ping(c context.Context, l *Location) (err error)
	LastLocation(c context.Context, username string) (location *LocationDetails, err error)
}

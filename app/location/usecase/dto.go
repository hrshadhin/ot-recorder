package usecase

import (
	"fmt"
	"ot-recorder/app/model"
	"ot-recorder/infrastructure/config"
	"time"
)

const dateTimeFormat = "2006-01-02 15:04:05"

func parseEpochTimeToLocal(epoch int64, tz string) time.Time {
	et := time.Unix(epoch, 0)
	loc, _ := time.LoadLocation(tz)
	return et.In(loc)
}

func toLastLocationDetails(l *model.Location) *model.LocationDetails {
	return &model.LocationDetails{
		Username:         l.Username,
		Device:           l.Device,
		DateTime:         parseEpochTimeToLocal(l.CreatedAt, config.Get().App.TimeZone).Format(dateTimeFormat),
		Accuracy:         l.Acc,
		Altitude:         l.Alt,
		BatteryLevel:     fmt.Sprintf("%d%s", l.Batt, "%"),
		BatteryStatus:    model.BatteryStatusEnum(l.Bs).String(),
		Latitude:         l.Lat,
		Longitude:        l.Lon,
		Mode:             model.ModeEnum(l.M).String(),
		VerticalAccuracy: l.Vac,
		Velocity:         l.Vel,
		WifiName:         l.Ssid,
		WifiMAC:          l.Bssid,
		IPAddress:        l.IP,
		MapLink:          fmt.Sprintf(mapLink, l.Lat, l.Lon, l.Lat, l.Lon),
	}
}

func toTelegramMessage(l *model.Location) string {
	message := fmt.Sprintf(`Username: *%s*
Device: *%s*
DateTime: *%s*
Latitude: *%f*
Longitude: *%f*
Accuracy: *%d*
Altitude: *%d*
BatteryLevel: *%s*
Mode: *%s*
[View in map](%s)`,
		l.Username,
		l.Device,
		parseEpochTimeToLocal(l.CreatedAt, config.Get().App.TimeZone).Format(dateTimeFormat),
		l.Lat,
		l.Lon,
		l.Acc,
		l.Alt,
		fmt.Sprintf("%d%s", l.Batt, "%"),
		model.ModeEnum(l.M).String(),
		fmt.Sprintf(mapLink, l.Lat, l.Lon, l.Lat, l.Lon),
	)

	return message
}

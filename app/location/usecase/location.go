package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"ot-recorder/app/model"
	"ot-recorder/app/response"
	"time"

	"github.com/sirupsen/logrus"
)

type locationUsecase struct {
	repo           model.LocationRepository
	contextTimeout time.Duration
}

func NewLocationUsecase(repo model.LocationRepository, timeout time.Duration) model.LocationUsecase {
	return &locationUsecase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *locationUsecase) Ping(c context.Context, l *model.Location) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// store location
	err = u.repo.CreateLocation(ctx, l)
	if err != nil {
		logrus.Errorln(err)

		return response.WrapError(errors.New("internal server error, please report to admin"), http.StatusInternalServerError)
	}

	return nil
}

func (u *locationUsecase) LastLocation(
	c context.Context,
	username string,
) (location *model.LocationDetails, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	l, err := u.repo.GetUserLastLocation(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.ErrNotFound
		}

		logrus.Errorln(err)

		return nil, response.WrapError(
			errors.New("internal server error, please report to admin"),
			http.StatusInternalServerError,
		)
	}

	location = &model.LocationDetails{
		Username:         l.Username,
		Device:           l.Device,
		DateTime:         l.CreatedAt.String(),
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
	}

	return location, nil
}

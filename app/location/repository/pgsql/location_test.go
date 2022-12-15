package pgsql_test

import (
	"context"
	locationRepo "ot-recorder/app/location/repository/pgsql"
	"ot-recorder/app/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateLocation(t *testing.T) {
	l := &model.Location{
		Username:  "dev",
		Device:    "phoneAndroid",
		CreatedAt: time.Now().Unix(),
		Acc:       13,
		Alt:       -42,
		Batt:      40,
		Bs:        1,
		Lat:       23.0000000,
		Lon:       90.0000000,
		M:         1,
		Tid:       "p1",
		Vac:       1,
		Vel:       0,
		Bssid:     "c0:00:00:00:00:00",
		Ssid:      "dev-test",
		IP:        "127.0.0.1",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "INSERT INTO locations"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(l.Username,
		l.Device,
		l.CreatedAt,
		l.Acc,
		l.Alt,
		l.Batt,
		l.Bs,
		l.Lat,
		l.Lon,
		l.M,
		l.T,
		l.Tid,
		l.Vac,
		l.Vel,
		l.Bssid,
		l.Ssid,
		l.IP,
	).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ur := locationRepo.NewPgsqlLocationRepository(db)
	err = ur.CreateLocation(context.TODO(), l)
	assert.NoError(t, err)
}

func TestGetUserLastLocation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "username", "device", "created_at", "acc", "alt", "batt", "bs", "lat", "lon", "m", "t", "tid", "vac",
		"vel", "bssid", "ssid", "ip"}).
		AddRow(1, "dev", "phoneAndroid", time.Now().Unix(), 13, -42, 40, 1, 23.0000000, 90.0000000, 1, "p",
			"p1", 1, 0, "", "", "")

	query := "SELECT \\* FROM locations WHERE username = \\$1 ORDER BY created_at DESC LIMIT 1"
	mock.ExpectQuery(query).WillReturnRows(rows)

	ur := locationRepo.NewPgsqlLocationRepository(db)

	location, err := ur.GetUserLastLocation(context.TODO(), "dev")
	assert.NoError(t, err)
	assert.NotNil(t, location)
	assert.Equal(t, "phoneAndroid", location.Device)
}

package mysql

import (
	"context"
	"database/sql"
	"ot-recorder/app/model"
)

type locationRepository struct {
	db *sql.DB
}

func NewMysqlLocationRepository(db *sql.DB) model.LocationRepository {
	return &locationRepository{
		db: db,
	}
}

const createLocation = `INSERT INTO locations (
  username, device, created_at, acc, alt, batt, bs, lat, lon, m, t, tid, vac, vel, bssid, ssid, ip
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func (r *locationRepository) CreateLocation(ctx context.Context, location *model.Location) error {
	stmt, err := r.db.PrepareContext(ctx, createLocation)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.ExecContext(ctx,
		location.Username,
		location.Device,
		location.CreatedAt,
		location.Acc,
		location.Alt,
		location.Batt,
		location.Bs,
		location.Lat,
		location.Lon,
		location.M,
		location.T,
		location.Tid,
		location.Vac,
		location.Vel,
		location.Bssid,
		location.Ssid,
		location.IP,
	)

	return err
}

const getPing = `SELECT * FROM locations WHERE username = ? ORDER BY created_at DESC LIMIT 1`

func (r *locationRepository) GetUserLastLocation(ctx context.Context, username string) (model.Location, error) {
	row := r.db.QueryRowContext(ctx, getPing, username)

	var l model.Location
	err := row.Scan(
		&l.ID,
		&l.Username,
		&l.Device,
		&l.CreatedAt,
		&l.Acc,
		&l.Alt,
		&l.Batt,
		&l.Bs,
		&l.Lat,
		&l.Lon,
		&l.M,
		&l.T,
		&l.Tid,
		&l.Vac,
		&l.Vel,
		&l.Bssid,
		&l.Ssid,
		&l.IP,
	)

	return l, err
}

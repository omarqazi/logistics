package datastore

import (
	"database/sql"
	"time"
)

var getUserStatement *sql.Stmt
var insertUserStatement *sql.Stmt

type User struct {
	Id        string
	PublicKey string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeviceId  string `json:"-"` // for push notifications
	Latitude  float64
	Longitude float64
}

func GetUser(userId string) (u *User, err error) {
	query := `
	select id, public_key, createdat, updatedat, device_id,
	ST_Y(location) as longitude, ST_X(location) as latitude
	from users where id = $1
	`
	if getUserStatement == nil {
		getUserStatement, err = Postgres.Prepare(query)
		if err != nil {
			return nil, err
		}
	}

	u = &User{}
	var lat, lon sql.NullFloat64
	err = getUserStatement.QueryRow(userId).Scan(
		&u.Id,
		&u.PublicKey,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeviceId,
		&lat,
		&lon,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if lat.Valid {
		u.Latitude = lat.Float64
	} else {
		u.Latitude = 0.0
	}

	if lon.Valid {
		u.Longitude = lon.Float64
	} else {
		u.Longitude = 0.0
	}
	return
}

func (u User) Location() Point {
	return Point{
		Latitude:  u.Latitude,
		Longitude: u.Longitude,
	}
}

func (u *User) Create() (err error) {
	query := `
	insert into users (
		id, public_key, createdat, updatedat, device_id
	) VALUES (
		$1, $2, now(), now(), $3
	)
	`
	if insertUserStatement == nil {
		insertUserStatement, err = Postgres.Prepare(query)
		if err != nil {
			return err
		}
	}

	if u.Id == "" {
		u.Id = NewUUID()
	}

	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	_, err = insertUserStatement.Exec(u.Id, u.PublicKey, u.DeviceId)
	return
}

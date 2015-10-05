package datastore

import (
	"database/sql/driver"
	"fmt"
)

const SRID = 4326 // WGS-84
type Point struct {
	Latitude  float64
	Longitude float64
}

func (p Point) String() string {
	wkt := fmt.Sprintf("POINT(%v %v)", p.Longitude, p.Latitude)
	return wkt
}

func (p Point) Value() (driver.Value, error) {
	wkt := fmt.Sprintf("POINT(%v %v)", p.Longitude, p.Latitude)
	return wkt, nil
}

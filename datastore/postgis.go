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
	fmt.Println("called string")
	wkt := fmt.Sprintf("POINT(%v %v)", p.Longitude, p.Latitude)
	sql := fmt.Sprintf("ST_GeometryFromText('%s', %d)", wkt, SRID)
	return sql
}

func (p Point) Value() (driver.Value, error) {
	fmt.Println("Called value")
	wkt := fmt.Sprintf("POINT(%v %v)", p.Longitude, p.Latitude)
	sql := fmt.Sprintf("ST_GeometryFromText('%s', %d)", wkt, SRID)
	return sql, nil
}

package data

import (
	geo "github.com/kellydunn/golang-geo"
)

type CircularGeofence struct {
	Name            string
	CentralPoint    geo.Point
	RadiousInMeters float64
}

func (geofence *CircularGeofence) IncludesPosition(latitude float64, longitude float64) bool {
	point := geo.NewPoint(latitude, longitude)
	return geofence.CentralPoint.GreatCircleDistance(point)*1000 < geofence.RadiousInMeters
}

var (
	Airport = &CircularGeofence{
		Name:            "Airport",
		CentralPoint:    *geo.NewPoint(60.31146, 24.96907),
		RadiousInMeters: 2000,
	}

	Downtown = &CircularGeofence{
		Name:            "Downtown",
		CentralPoint:    *geo.NewPoint(60.16422983026082, 24.941068845053014),
		RadiousInMeters: 1700,
	}

	RailwaySquare = &CircularGeofence{
		Name:            "Railway Square",
		CentralPoint:    *geo.NewPoint(60.171285, 24.943936),
		RadiousInMeters: 150,
	}

	LauttasaariIsland = &CircularGeofence{
		Name:            "Lauttasaari island",
		CentralPoint:    *geo.NewPoint(60.158536, 24.873788),
		RadiousInMeters: 1400,
	}

	LaajasaloIsland = &CircularGeofence{
		Name:            "Laajasalo island",
		CentralPoint:    *geo.NewPoint(60.16956184470527, 25.052851825093114),
		RadiousInMeters: 2200,
	}

	KallioDistrict = &CircularGeofence{
		Name:            "Kallio district",
		CentralPoint:    *geo.NewPoint(60.18260263288996, 24.953588638997264),
		RadiousInMeters: 600,
	}
)

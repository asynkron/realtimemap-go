package contract

type Position struct {
	VehicleId string  `json:"vehicleId"`
	OrgId     string  `json:"orgId"`
	Timestamp int64   `json:"timestamp"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Heading   int32   `json:"heading"`
	DoorsOpen bool    `json:"doorsOpen"`
	Speed     float64 `json:"speed"`
}

type PositionBatch struct {
	Positions []*Position `json:"positions"`
}

type Organization struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type OrganizationDetails struct {
	Id        string      `json:"id"`
	Name      string      `json:"name"`
	Geofences []*Geofence `json:"geofences"`
}

type Geofence struct {
	Name           string   `json:"name"`
	Longitude      float64  `json:"longitude"`
	Latitude       float64  `json:"latitude"`
	RadiusInMeters float64  `json:"radiusInMeters"`
	VehiclesInZone []string `json:"vehiclesInZone"`
}

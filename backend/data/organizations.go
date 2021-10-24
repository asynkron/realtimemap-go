package data

type Organization struct {
	Id        string
	Name      string
	Geofences []*CircularGeofence
}

var AllOrganizations = map[string]*Organization{
	"0006": {
		Id:   "006",
		Name: "Oy Pohjolan Liikenne Ab",
	},
	"0012": {
		Id:        "0012",
		Name:      "Helsingin Bussiliikenne Oy",
		Geofences: []*CircularGeofence{Airport, KallioDistrict, RailwaySquare},
	},
	"0017": {
		Id:        "0017",
		Name:      "Tammelundin Liikenne Oy",
		Geofences: []*CircularGeofence{LaajasaloIsland},
	},
	"0018": {
		Id:        "0018",
		Name:      "Pohjolan Kaupunkiliikenne Oy",
		Geofences: []*CircularGeofence{KallioDistrict, LauttasaariIsland, RailwaySquare},
	},
	"0020": {
		Id:   "0020",
		Name: "Bus Travel Åbergin Linja Oy",
	},
	"0021": {
		Id:   "0021",
		Name: "Bus Travel Oy Reissu Ruoti",
	},
	"0022": {
		Id:        "0022",
		Name:      "Nobina Finland Oy",
		Geofences: []*CircularGeofence{Airport, KallioDistrict, LaajasaloIsland},
	},
	"0030": {
		Id:        "0030",
		Name:      "Savonlinja Oy",
		Geofences: []*CircularGeofence{Airport, Downtown},
	},
	"0036": {
		Id:   "0036",
		Name: "Nurmijärven Linja Oy",
	},
	"0040": {
		Id:   "0040",
		Name: "HKL-Raitioliikenne",
	},
	"0045": {
		Id:   "0045",
		Name: "Transdev Vantaa Oy",
	},
	"0047": {
		Id:   "0047",
		Name: "Taksikuljetus Oy",
	},
	"0050": {
		Id:   "0050",
		Name: "HKL-Metroliikenne",
	},
	"0051": {
		Id:   "0051",
		Name: "Korsisaari Oy",
	},
	"0054": {
		Id:   "0054",
		Name: "V-S Bussipalvelut Oy",
	},
	"0055": {
		Id:   "0055",
		Name: "Transdev Helsinki Oy",
	},
	"0058": {
		Id:   "0058",
		Name: "Koillisen Liikennepalvelut Oy",
	},
	"0060": {
		Id:   "0060",
		Name: "Suomenlinnan Liikenne Oy",
	},
	"0059": {
		Id:   "0059",
		Name: "Tilausliikenne Nikkanen Oy",
	},
	"0089": {
		Id:   "0089",
		Name: "Metropolia",
	},
	"0090": {
		Id:   "0090",
		Name: "VR Oy",
	},
	"0195": {
		Id:   "0195",
		Name: "Siuntio1",
	},
}

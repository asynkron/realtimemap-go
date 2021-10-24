package server

import (
	"sort"

	"realtimemap-go/backend/contract"
	"realtimemap-go/backend/data"
	"realtimemap-go/backend/grains"
)

func mapPositionBatch(batch *grains.PositionBatch) *contract.PositionBatch {
	result := make([]*contract.Position, 0, len(batch.Positions))

	for _, pos := range batch.Positions {
		result = append(result, &contract.Position{
			OrgId:     pos.OrgId,
			VehicleId: pos.VehicleId,
			Timestamp: pos.Timestamp,
			Longitude: pos.Longitude,
			Latitude:  pos.Latitude,
			Heading:   pos.Heading,
			DoorsOpen: pos.DoorsOpen,
			Speed:     pos.Speed,
		})
	}

	return &contract.PositionBatch{Positions: result}
}

func mapOrganization(org *data.Organization, grainResponse *grains.GetGeofencesResponse) *contract.OrganizationDetails {
	geofences := make([]*contract.Geofence, 0, len(grainResponse.Geofences))

	for _, grainGeofence := range grainResponse.Geofences {
		geofences = append(geofences, mapGeofence(grainGeofence))
	}

	sort.Slice(geofences, func(i, j int) bool {
		return geofences[i].Name < geofences[j].Name
	})

	return &contract.OrganizationDetails{
		Id:        org.Id,
		Name:      org.Name,
		Geofences: geofences,
	}
}

func mapGeofence(grainGeofence *grains.GeofenceDetails) *contract.Geofence {
	vehicles := make([]string, len(grainGeofence.VehiclesInZone))
	copy(vehicles, grainGeofence.VehiclesInZone)
	sort.Strings(vehicles)

	return &contract.Geofence{
		Name:           grainGeofence.Name,
		Longitude:      grainGeofence.Longitude,
		Latitude:       grainGeofence.Latitude,
		RadiusInMeters: grainGeofence.RadiusInMeters,
		VehiclesInZone: vehicles,
	}
}

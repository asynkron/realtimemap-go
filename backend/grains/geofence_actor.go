package grains

import (
	"realtimemap-go/backend/data"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

type geofenceActor struct {
	cluster          *cluster.Cluster
	organizationName string
	geofence         data.CircularGeofence
	vehiclesInZone   map[string]struct{}
}

func NewGeofenceActor(organizationName string, geofence *data.CircularGeofence, cluster *cluster.Cluster) actor.Actor {
	return &geofenceActor{
		cluster:          cluster,
		organizationName: organizationName,
		geofence:         *geofence,
		vehiclesInZone:   make(map[string]struct{})}
}

func (g *geofenceActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {

	case *Position:
		_, vehicleIsInZone := g.vehiclesInZone[msg.VehicleId]

		if g.geofence.IncludesPosition(msg.Latitude, msg.Longitude) {
			if !vehicleIsInZone {
				g.vehiclesInZone[msg.VehicleId] = struct{}{}

				g.cluster.MemberList.BroadcastEvent(&Notification{
					VehicleId: msg.VehicleId,
					OrgId:     msg.OrgId,
					OrgName:   g.organizationName,
					ZoneName:  g.geofence.Name,
					Event:     enter})
			}
		} else {
			if vehicleIsInZone {
				delete(g.vehiclesInZone, msg.VehicleId)

				g.cluster.MemberList.BroadcastEvent(&Notification{
					VehicleId: msg.VehicleId,
					OrgId:     msg.OrgId,
					OrgName:   g.organizationName,
					ZoneName:  g.geofence.Name,
					Event:     exit})
			}
		}

	case *GetGeofencesRequest:
		ctx.Respond(&GeofenceDetails{
			Name:           g.geofence.Name,
			RadiusInMeters: g.geofence.RadiousInMeters,
			Latitude:       g.geofence.CentralPoint.Lat(),
			Longitude:      g.geofence.CentralPoint.Lng(),
			OrgId:          msg.OrgId,
			VehiclesInZone: getMapKeys(g.vehiclesInZone),
		})
	}

}

func getMapKeys(m map[string]struct{}) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

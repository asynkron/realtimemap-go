package grains

import (
	"realtimemap-go/backend/data"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
)

type geofenceActor struct {
	organizationName string
	geofence         data.CircularGeofence
	vehiclesInZone   map[string]struct{}
}

func NewGeofenceActor(organizationName string, geofence *data.CircularGeofence) actor.Actor {
	return &geofenceActor{
		organizationName: organizationName,
		geofence:         *geofence,
		vehiclesInZone:   make(map[string]struct{})}
}

func (g *geofenceActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {

	case *Position:
		cl := cluster.GetCluster(ctx.ActorSystem())
		_, vehicleIsInZone := g.vehiclesInZone[msg.VehicleId]

		if g.geofence.IncludesPosition(msg.Latitude, msg.Longitude) {
			if !vehicleIsInZone {
				g.vehiclesInZone[msg.VehicleId] = struct{}{}

				cl.MemberList.BroadcastEvent(&Notification{
					VehicleId: msg.VehicleId,
					OrgId:     msg.OrgId,
					OrgName:   g.organizationName,
					ZoneName:  g.geofence.Name,
					Event:     GeofenceEvent_Enter}, true)
			}
		} else {
			if vehicleIsInZone {
				delete(g.vehiclesInZone, msg.VehicleId)

				cl.MemberList.BroadcastEvent(&Notification{
					VehicleId: msg.VehicleId,
					OrgId:     msg.OrgId,
					OrgName:   g.organizationName,
					ZoneName:  g.geofence.Name,
					Event:     GeofenceEvent_Exit}, true)
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

package grains

import (
	"fmt"

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

				// TODO: fix once its possible to broadcast in whole cluster
				g.cluster.ActorSystem.EventStream.Publish(&Notification{
					Message: fmt.Sprintf("%s from %s entered the zone %s", msg.VehicleId, g.organizationName, g.geofence.Name),
				})

				fmt.Printf("%s from %s entered the zone %s\n", msg.VehicleId, g.organizationName, g.geofence.Name)
			}
		} else {
			if vehicleIsInZone {
				delete(g.vehiclesInZone, msg.VehicleId)

				// TODO: fix once its possible to broadcast in whole cluster
				g.cluster.ActorSystem.EventStream.Publish(&Notification{
					Message: fmt.Sprintf("%s from %s left the zone %s", msg.VehicleId, g.organizationName, g.geofence.Name),
				})

				fmt.Printf("%s from %s left the zone %s\n", msg.VehicleId, g.organizationName, g.geofence.Name)
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

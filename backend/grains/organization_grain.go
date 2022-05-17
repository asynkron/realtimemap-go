package grains

import (
	"fmt"
	"time"

	"realtimemap-go/backend/data"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
)

type OrganizationGrain struct {
	name string
}

func (o *OrganizationGrain) Init(ctx cluster.GrainContext) {
	if organization, ok := data.AllOrganizations[ctx.Identity()]; ok {
		o.name = organization.Name
		for _, geofence := range organization.Geofences {
			o.createGeofenceActor(geofence, ctx)
		}
	}
	fmt.Printf("OrganizationGrain %v initialized\n", ctx.Identity())
}

func (o *OrganizationGrain) createGeofenceActor(geofence *data.CircularGeofence, ctx cluster.GrainContext) {
	props := actor.PropsFromProducer(func() actor.Actor {
		return NewGeofenceActor(o.name, geofence)
	})
	ctx.Spawn(props)
}

func (o *OrganizationGrain) OnPosition(position *Position, ctx cluster.GrainContext) (*Empty, error) {
	for _, geofenceActor := range ctx.Children() {
		ctx.Send(geofenceActor, position)
	}

	return &Empty{}, nil
}

func (o *OrganizationGrain) GetGeofences(_ *GetGeofencesRequest, ctx cluster.GrainContext) (*GetGeofencesResponse, error) {

	futures := make([]*actor.Future, 0, len(ctx.Children()))

	for _, child := range ctx.Children() {
		future := ctx.RequestFuture(child, &GetGeofencesRequest{OrgId: ctx.Identity()}, 5*time.Second)
		futures = append(futures, future)
	}

	geofences := make([]*GeofenceDetails, 0, len(ctx.Children()))

	for _, future := range futures {
		if res, err := future.Result(); err == nil {
			if geofence, ok := res.(*GeofenceDetails); ok {
				geofences = append(geofences, geofence)
			} else {
				fmt.Printf("%v is not a GeofenceDetails", res)
			}
		} else {
			return nil, err
		}
	}

	return &GetGeofencesResponse{Geofences: geofences}, nil
}

func (o *OrganizationGrain) Terminate(ctx cluster.GrainContext)      {}
func (o *OrganizationGrain) ReceiveDefault(ctx cluster.GrainContext) {}

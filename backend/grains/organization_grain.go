package grains

import (
	"fmt"
	"time"

	"realtimemap-go/backend/data"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

type organizationGrain struct {
	id          string
	name        string
	initialized bool
	cluster     *cluster.Cluster
}

func CreateOrganizationFactory(cluster *cluster.Cluster) func() Organization {
	return func() Organization {
		return &organizationGrain{cluster: cluster}
	}
}

func (o *organizationGrain) Init(id string) {
	o.id = id
}

func (o *organizationGrain) initializeOrgIfNeeded(ctx cluster.GrainContext) {
	if o.initialized {
		return
	}

	if organization, ok := data.AllOrganizations[o.id]; ok {
		o.name = organization.Name
		for _, geofence := range organization.Geofences {
			o.createGeofenceActor(geofence, ctx)
		}
	}

	o.initialized = true
}

func (o *organizationGrain) createGeofenceActor(geofence *data.CircularGeofence, ctx cluster.GrainContext) {
	props := actor.PropsFromProducer(func() actor.Actor {
		return NewGeofenceActor(o.name, geofence, o.cluster)
	})
	ctx.Spawn(props)
}

func (o *organizationGrain) OnPosition(position *Position, ctx cluster.GrainContext) (*Empty, error) {
	// TODO: normally this would be peformed in the Init func, but the generated code does not pass context to it
	o.initializeOrgIfNeeded(ctx)

	for _, geofenceActor := range ctx.Children() {
		ctx.Send(geofenceActor, position)
	}

	return &Empty{}, nil
}

func (o *organizationGrain) GetGeofences(_ *GetGeofencesRequest, ctx cluster.GrainContext) (*GetGeofencesResponse, error) {

	futures := make([]*actor.Future, 0, len(ctx.Children()))

	for _, child := range ctx.Children() {
		future := ctx.RequestFuture(child, &GetGeofencesRequest{OrgId: o.id}, 5*time.Second)
		futures = append(futures, future)
	}

	geofences := make([]*GeofenceDetails, 0, len(ctx.Children()))

	for _, future := range futures {
		if res, err := future.Result(); err == nil {
			if geofence, ok := res.(*GeofenceDetails); ok {
				geofences = append(geofences, geofence)
			} else {
				fmt.Println("WTF did you send me")
			}
		} else {
			return nil, err
		}
	}

	return &GetGeofencesResponse{Geofences: geofences}, nil
}

func (o *organizationGrain) Terminate()                       {}
func (o *organizationGrain) ReceiveDefault(ctx actor.Context) {}

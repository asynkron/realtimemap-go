package grains

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

const MaxPositionHistory = 200

type vehicleGrain struct {
	id              string
	positionHistory []*Position
	cluster         *cluster.Cluster
}

func CreateVehicleFactory(cluster *cluster.Cluster) func() Vehicle {
	return func() Vehicle {
		return &vehicleGrain{cluster: cluster}
	}
}

func (v *vehicleGrain) Init(id string) {
	v.id = id
	v.positionHistory = make([]*Position, 0, MaxPositionHistory)
}

func (v *vehicleGrain) OnPosition(position *Position, ctx cluster.GrainContext) (*Empty, error) {

	if len(v.positionHistory) > MaxPositionHistory {
		v.positionHistory = v.positionHistory[1:] // TODO: is this memory leak?
	}
	v.positionHistory = append(v.positionHistory, position)

	orgClient := GetOrganizationGrainClient(v.cluster, position.OrgId)
	orgClient.OnPosition(position)

	v.cluster.MemberList.BroadcastEvent(position)

	return &Empty{}, nil
}

func (v *vehicleGrain) GetPositionsHistory(*GetPositionsHistoryRequest, cluster.GrainContext) (*PositionBatch, error) {
	return &PositionBatch{Positions: v.positionHistory}, nil
}

func (v *vehicleGrain) Terminate()                       {}
func (v *vehicleGrain) ReceiveDefault(ctx actor.Context) {}

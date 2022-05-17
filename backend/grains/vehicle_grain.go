package grains

import (
	"github.com/asynkron/protoactor-go/cluster"
)

const MaxPositionHistory = 200

type VehicleGrain struct {
	positionHistory []*Position
}

func (v *VehicleGrain) Init(ctx cluster.GrainContext) {
	v.positionHistory = make([]*Position, 0, MaxPositionHistory)
}

func (v *VehicleGrain) OnPosition(position *Position, ctx cluster.GrainContext) (*Empty, error) {

	cl := cluster.GetCluster(ctx.ActorSystem())

	if len(v.positionHistory) > MaxPositionHistory {
		v.positionHistory = v.positionHistory[1:] // TODO: is this memory leak?
	}
	v.positionHistory = append(v.positionHistory, position)

	orgClient := GetOrganizationGrainClient(cl, position.OrgId)
	orgClient.OnPosition(position)

	cl.MemberList.BroadcastEvent(position, true)

	return &Empty{}, nil
}

func (v *VehicleGrain) GetPositionsHistory(*GetPositionsHistoryRequest, cluster.GrainContext) (*PositionBatch, error) {
	return &PositionBatch{Positions: v.positionHistory}, nil
}

func (v *VehicleGrain) Terminate(ctx cluster.GrainContext)      {}
func (v *VehicleGrain) ReceiveDefault(ctx cluster.GrainContext) {}

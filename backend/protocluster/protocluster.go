package protocluster

import (
	"realtimemap-go/backend/grains"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/automanaged"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
)

func StartNode() *cluster.Cluster {
	system := actor.NewActorSystem()

	vehicleKind := grains.NewVehicleKind(func() grains.Vehicle {
		return &grains.VehicleGrain{}
	}, 0)
	organizationKind := grains.NewOrganizationKind(func() grains.Organization {
		return &grains.OrganizationGrain{}
	}, 0)

	provider := automanaged.New()
	config := remote.Configure("localhost", 0)
	lookup := disthash.New()

	clusterConfig := cluster.Configure("my-cluster", provider, lookup, config, cluster.WithKinds(vehicleKind, organizationKind))
	cluster := cluster.New(system, clusterConfig)

	cluster.StartMember()

	return cluster
}

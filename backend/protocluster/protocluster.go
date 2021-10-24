package protocluster

import (
	"realtimemap-go/backend/grains"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/automanaged"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func StartNode() *cluster.Cluster {
	system := actor.NewActorSystem()

	vehicleKind := cluster.NewKind("Vehicle", actor.PropsFromProducer((func() actor.Actor {
		return &grains.VehicleActor{}
	})))
	organizationKind := cluster.NewKind("Organization", actor.PropsFromProducer((func() actor.Actor {
		return &grains.OrganizationActor{}
	})))

	provider := automanaged.New()
	config := remote.Configure("localhost", 0)

	clusterConfig := cluster.Configure("my-cluster", provider, config, vehicleKind, organizationKind)
	cluster := cluster.New(system, clusterConfig)

	grains.VehicleFactory(grains.CreateVehicleFactory(cluster))
	grains.OrganizationFactory(grains.CreateOrganizationFactory(cluster))

	cluster.Start()

	return cluster
}

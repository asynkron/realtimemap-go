package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"realtimemap-go/backend/grains"
	"realtimemap-go/backend/ingress"
	"realtimemap-go/backend/protocluster"
	"realtimemap-go/backend/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	stopOnSignals(cancel)

	cluster := protocluster.StartNode()
	time.Sleep(2 * time.Second)

	server := server.NewHttpServer(cluster, ctx)
	serverDone := server.ListenAndServe()

	ingressDone := ingress.ConsumeVehicleEvents(func(event *ingress.Event) {
		position := MapToPosition(event)
		if position != nil {
			vehicleGrainClient := grains.GetVehicleGrainClient(cluster, position.VehicleId)
			vehicleGrainClient.OnPosition(position)
		}
	}, ctx)

	<-ingressDone
	<-serverDone

	cluster.Shutdown(true)
}

func stopOnSignals(cancel func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		fmt.Println("*** STOPPING ***")
		cancel()
	}()
}

func MapToPosition(e *ingress.Event) *grains.Position {
	var payload *ingress.Payload

	if e.VehiclePosition != nil {
		payload = e.VehiclePosition
	} else if e.DoorOpen != nil {
		payload = e.DoorOpen
	} else if e.DoorClosed != nil {
		payload = e.DoorClosed
	} else {
		return nil
	}

	if !payload.HasValidPosition() {
		return nil
	}

	return &grains.Position{
		VehicleId: e.VehicleId,
		OrgId:     e.OperatorId,
		Latitude:  *payload.Latitude,
		Longitude: *payload.Longitude,
		Heading:   *payload.Heading,
		Timestamp: (*payload.Timestamp).UnixMilli(),
		Speed:     *payload.Speed,
	}
}

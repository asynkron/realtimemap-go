package server

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"realtimemap-go/backend/grains"

	"github.com/asynkron/protoactor-go/actor"
	kitlog "github.com/go-kit/log"
	"github.com/philippseith/signalr"
)

type AppHub struct {
	signalr.Hub
	initialized bool
	actorSystem *actor.ActorSystem
}

func (h *AppHub) Initialize(ctx signalr.HubContext) {
	h.Hub.Initialize(ctx)
	// Initialize will be called on first connection to the hub
	// if position is sent before that, the HubContext is nil and application crashes
	// TODO: possible bug in signalr implementation?
	h.initialized = true
}

func (h *AppHub) OnConnected(connectionID string) {
	fmt.Printf("%s connected to hub\n", connectionID)
}

func (h *AppHub) OnDisconnected(connectionID string) {
	fmt.Printf("%s disconnected from hub\n", connectionID)
	if item, loaded := h.Items().LoadAndDelete("user"); loaded {
		viewportPID, _ := item.(*actor.PID)
		h.actorSystem.Root.Stop(viewportPID)
	}
}

func (h *AppHub) SendPositionBatch(connectionID string, batch *grains.PositionBatch) {
	if h.initialized {
		h.Clients().Client(connectionID).Send("positions", mapPositionBatch(batch))
	}
}

func (h *AppHub) SendNotification(connectionID string, notification *grains.Notification) {
	if h.initialized {
		h.Clients().Client(connectionID).Send("notification", mapNotification(notification))
	}
}

func (h *AppHub) SetViewport(swLng float64, swLat float64, neLng float64, neLat float64) {
	var userPID *actor.PID

	if item, loaded := h.Items().Load(h.ConnectionID()); loaded {
		userPID, _ = item.(*actor.PID)
	} else {
		props := actor.PropsFromProducer(func() actor.Actor {
			return grains.NewUserActor(h.ConnectionID(), h.SendPositionBatch, h.SendNotification)
		})
		// spawn named so that we don't get multiple viewports for same connection id in the case of concurrency issues
		userPID, _ = h.actorSystem.Root.SpawnNamed(props, h.ConnectionID())
		h.Items().Store("user", userPID)
	}

	h.actorSystem.Root.Send(userPID, &grains.UpdateViewport{
		Viewport: &grains.Viewport{
			SouthWest: &grains.GeoPoint{Longitude: swLng, Latitude: swLat},
			NorthEast: &grains.GeoPoint{Longitude: neLng, Latitude: neLat},
		},
	})
}

func serveHub(router *http.ServeMux, actorSystem *actor.ActorSystem, ctx context.Context) *AppHub {
	hub := &AppHub{actorSystem: actorSystem}

	singnalrServer, _ := signalr.NewServer(ctx,
		signalr.UseHub(hub),
		signalr.AllowOriginPatterns([]string{"localhost:8080"}),
		signalr.Logger(kitlog.NewLogfmtLogger(os.Stdout), false))

	singnalrServer.MapHTTP(router, "/events")

	return hub
}

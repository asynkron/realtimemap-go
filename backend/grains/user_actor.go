package grains

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
)

const BatchSize = 10

type SendPositions func(connectionID string, positions *PositionBatch)
type SendNotification func(connectionID string, notification *Notification)

type userActor struct {
	connectionID     string
	viewport         Viewport
	batch            []*Position
	sendPositions    SendPositions
	sendNotification SendNotification
	subscription     *eventstream.Subscription
}

func NewUserActor(connectionID string, sendPositions SendPositions, sendNotification SendNotification) *userActor {
	return &userActor{
		connectionID:     connectionID,
		sendPositions:    sendPositions,
		sendNotification: sendNotification,
	}
}

func (u *userActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {

	case *actor.Started:
		u.batch = make([]*Position, 0, BatchSize)

		u.subscription = ctx.ActorSystem().EventStream.Subscribe(func(event interface{}) {
			// do not modify state in the callback to avoid concurrency issues, let the message pass through mailbox
			switch event.(type) {
			case *Position:
				ctx.Send(ctx.Self(), event)
			case *Notification:
				ctx.Send(ctx.Self(), event)
			}
		})

	case *Position:
		if msg.Latitude < u.viewport.SouthWest.Latitude ||
			msg.Latitude > u.viewport.NorthEast.Latitude ||
			msg.Longitude < u.viewport.SouthWest.Longitude ||
			msg.Longitude > u.viewport.NorthEast.Longitude {
			return
		}

		u.batch = append(u.batch, msg)
		if len(u.batch) >= BatchSize {
			u.sendPositions(u.connectionID, &PositionBatch{Positions: u.batch})
			u.batch = u.batch[:0]
		}

	case *Notification:
		u.sendNotification(u.connectionID, msg)

	case *UpdateViewport:
		u.viewport = *msg.Viewport
		fmt.Printf("Viewport for connection %s is now %+v\n", u.connectionID, u.viewport)

	case *actor.Stopping:
		fmt.Printf("Stopping user actor for connection %s\n", u.connectionID)
		ctx.ActorSystem().EventStream.Unsubscribe(u.subscription)
		u.subscription = nil
	}
}

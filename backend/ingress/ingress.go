package ingress

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Payload struct {
	Longitude *float64   `json:"long"`
	Latitude  *float64   `json:"lat"`
	Heading   *int32     `json:"hdg"`
	DoorState *int32     `json:"drst"`
	Timestamp *time.Time `json:"tst"`
	Speed     *float64   `json:"spd"`
}

func (p *Payload) HasValidPosition() bool {
	return p != nil && p.Latitude != nil && p.Longitude != nil && p.Heading != nil && p.Timestamp != nil && p.Speed != nil && p.DoorState != nil
}

type Event struct {
	VehiclePosition *Payload `json:"VP"`
	DoorOpen        *Payload `json:"DOO"`
	DoorClosed      *Payload `json:"DOC"`
	VehicleId       string
	OperatorId      string
}

func ConsumeVehicleEvents(onEvent func(*Event), ctx context.Context) <-chan bool {
	done := make(chan bool)
	go func() {

		var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
			var event Event

			//0/1       /2        /3             /4              /5           /6               /7            /8               /9         /10            /11        /12          /13         /14             /15       /16
			// /<prefix>/<version>/<journey_type>/<temporal_type>/<event_type>/<transport_mode>/<operator_id>/<vehicle_number>/<route_id>/<direction_id>/<headsign>/<start_time>/<next_stop>/<geohash_level>/<geohash>/<sid>/#
			topicParts := strings.Split(msg.Topic(), "/")

			if err := json.Unmarshal(msg.Payload(), &event); err != nil {
				fmt.Printf("Error unmarshalling json %v", err)
			} else {
				event.OperatorId = topicParts[7]
				event.VehicleId = topicParts[7] + "." + topicParts[8]
				onEvent(&event)
			}

		}

		mqtt.WARN = log.New(os.Stdout, "", 0)
		mqtt.ERROR = log.New(os.Stdout, "", 0)
		mqtt.CRITICAL = log.New(os.Stdout, "", 0)

		opts := mqtt.NewClientOptions()
		opts.AddBroker("ssl://mqtt.hsl.fi:8883")
		opts.SetClientID("realtimemap-go")
		opts.SetDefaultPublishHandler(f)
		opts.SetKeepAlive(2 * time.Second)
		opts.SetPingTimeout(1 * time.Second)
		opts.SetCleanSession(true)
		//opts.SetConnectRetry()

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println("CONNECTED")

		if token := client.Subscribe("/hfp/v2/journey/ongoing/vp/bus/#", 0, nil); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println("SUBSCRIBED")

		<-ctx.Done()

		if token := client.Unsubscribe("/hfp/v2/journey/ongoing/vp/bus/#"); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println("UNSUBSCRIBED")

		client.Disconnect(250)
		fmt.Println("DISCONNECTED")

		done <- true
	}()

	return done
}

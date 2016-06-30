package honeycomb

import (
	"encoding/json"
	//	"errors"
	"log"
	"net"

	"github.com/gliderlabs/logspout/router"
	"github.com/houndsh/libhound-go-private"
)

func init() {
	log.Println("init")
	router.AdapterFactories.Register(NewHoneycombAdapter, "honeycomb")
}

// HoneycombAdapter is an adapter that streams JSON to Logstash.
type HoneycombAdapter struct {
	conn  net.Conn
	route *router.Route
}

// NewHoneycombAdapter creates a HoneycombAdapter
func NewHoneycombAdapter(route *router.Route) (router.LogAdapter, error) {
	/*
		transport, found := router.AdapterTransports.Lookup(route.AdapterTransport("tls"))
		if !found {
			return nil, errors.New("unable to find adapter: " + route.Adapter)
		}

		conn, err := transport.Dial(route.Address, route.Options)
		if err != nil {
			return nil, err
		}
	*/

	libhound.Init(libhound.Config{
		WriteKey: "09f5607ab2ae0aba7fe5f38ce091feb2",
		Dataset:  "ohai",
	})
	log.Println("init libhound")

	return &HoneycombAdapter{}, nil
}

// Stream implements the router.LogAdapter interface.
func (a *HoneycombAdapter) Stream(logstream chan *router.Message) {
	for m := range logstream {
		var js []byte
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(m.Data), &data); err != nil {
			log.Println("message was JSON")
			// The message is not in JSON, make a new JSON message.
			msg := HoneycombMessage{
				Message:              m.Data,
				Stream:               m.Source,
				DockerContainerName:  m.Container.Name,
				DockerContainerID:    m.Container.ID,
				DockerHostname:       m.Container.Config.Hostname,
				DockerContainerImage: m.Container.Config.Image,
			}
			if js, err = json.Marshal(msg); err != nil {
				log.Println("logstash:", err)
				continue
			}
		} else {
			log.Println("message was not JSON")
			// The message is already in JSON, add the docker specific fields.
			data["dockerContainerName"] = m.Container.Name
			data["dockerContainerId"] = m.Container.ID
			data["dockerHostname"] = m.Container.Config.Hostname
			data["dockerImage"] = m.Container.Config.Image

			if js, err = json.Marshal(data); err != nil {
				log.Println("logstash:", err)
				continue
			}
		}

		log.Println("sending to honeycomb.")
		libhound.SendNow(js)
	}
}

// HoneycombMessage is a flat JSON object sent to the Honeycomb service
type HoneycombMessage struct {
	Message              string `json:"message"`
	Stream               string `json:"stream"`
	DockerHostname       string `json:"dockerHostname"`
	DockerContainerName  string `json:"dockerContainerName"`
	DockerContainerID    string `json:"dockerContainerId"`
	DockerContainerImage string `json:"dockerImage"`
}

package honeycomb

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/gliderlabs/logspout/router"
	"github.com/houndsh/libhound-go-private"
)

func init() {
	router.AdapterFactories.Register(NewHoneycombAdapter, "honeycomb")
}

// HoneycombAdapter is an adapter that streams JSON to Logstash.
type HoneycombAdapter struct {
	conn  net.Conn
	route *router.Route
}

// NewHoneycombAdapter creates a HoneycombAdapter
func NewHoneycombAdapter(route *router.Route) (router.LogAdapter, error) {
	writeKey := route.Options["writeKey"]
	if writeKey == "" {
		writeKey = os.Getenv("HONEYCOMB_WRITE_KEY")
	}
	if writeKey == "" {
		log.Fatal("Must provide Honeycomb WriteKey.")
	}

	dataset := route.Options["dataset"]
	if dataset == "" {
		dataset = os.Getenv("HONEYCOMB_DATASET")
	}
	if dataset == "" {
		log.Fatal("Must provide Honeycomb Dataset.")
	}

	libhound.Init(libhound.Config{
		WriteKey: writeKey,
		Dataset:  dataset,
	})

	return &HoneycombAdapter{}, nil
}

// Stream implements the router.LogAdapter interface.
func (a *HoneycombAdapter) Stream(logstream chan *router.Message) {
	for m := range logstream {

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(m.Data), &data); err != nil {
			// The message is not in JSON, make a new JSON message.
			msg := HoneycombMessage{
				Message:              m.Data,
				Stream:               m.Source,
				DockerContainerName:  m.Container.Name,
				DockerContainerID:    m.Container.ID,
				DockerHostname:       m.Container.Config.Hostname,
				DockerContainerImage: m.Container.Config.Image,
			}

			if err := libhound.SendNow(msg); err != nil {
				log.Println("error: ", err)
			}
		} else {
			// The message is already in JSON, add the docker specific fields.
			data["stream"] = m.Source
			data["logspout_container"] = m.Container.Name
			data["logspout_container_id"] = m.Container.ID
			data["logspout_hostname"] = m.Container.Config.Hostname
			data["logspout_docker_image"] = m.Container.Config.Image

			if err := libhound.SendNow(data); err != nil {
				log.Println("error: ", err)
			}
		}
	}
}

// HoneycombMessage is a flat JSON object sent to the Honeycomb service
type HoneycombMessage struct {
	Message              string `json:"message"`
	Stream               string `json:"stream"`
	DockerHostname       string `json:"logspout_hostname"`
	DockerContainerName  string `json:"logspout_container"`
	DockerContainerID    string `json:"logspout_container_id"`
	DockerContainerImage string `json:"logspout_docker_image"`
}

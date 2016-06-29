package honeycomb

import (
	"encoding/json"
	"errors"
	"log"
	"net"

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

	return &HoneycombAdapter{}, nil
}

// Stream implements the router.LogAdapter interface.
func (a *HoneycombAdapter) Stream(logstream chan *router.Message) {
	for m := range logstream {
		dockerInfo := DockerInfo{
			Name:     m.Container.Name,
			ID:       m.Container.ID,
			Image:    m.Container.Config.Image,
			Hostname: m.Container.Config.Hostname,
		}

		libhound.SendNow(m)

		/*
			var js []byte
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(m.Data), &data); err != nil {
				// The message is not in JSON, make a new JSON message.
				msg := LogstashMessage{
					Message: m.Data,
					Docker:  dockerInfo,
					Stream:  m.Source,
				}
				if js, err = json.Marshal(msg); err != nil {
					log.Println("logstash:", err)
					continue
				}
			} else {
				// The message is already in JSON, add the docker specific fields.
				data["docker"] = dockerInfo
				if js, err = json.Marshal(data); err != nil {
					log.Println("logstash:", err)
					continue
				}
			}

			if _, err := a.conn.Write(js); err != nil {
				log.Fatal("logstash:", err)
			}
		*/
	}
}

// DockerInfo  Hey, here's a comment to satisfy the linter
type DockerInfo struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Image    string `json:"image"`
	Hostname string `json:"hostname"`
}

// LogstashMessage is a simple JSON input to Logstash.
type LogstashMessage struct {
	Message string     `json:"message"`
	Stream  string     `json:"stream"`
	Docker  DockerInfo `json:"docker"`
}

package honeycomb

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/gliderlabs/logspout/router"
	"github.com/honeycombio/libhoney-go"
)

const (
	DefaultHoneycombAPIURL = "https://api.honeycomb.io"
	DefaultSampleRate      = 1
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
		return nil, errors.New("Honeycomb 'WriteKey' was not provided.")
	}

	dataset := route.Options["dataset"]
	if dataset == "" {
		dataset = os.Getenv("HONEYCOMB_DATASET")
	}
	if dataset == "" {
		log.Fatal("Must provide Honeycomb Dataset.")
		return nil, errors.New("Honeycomb 'Dataset' was not provided.")
	}

	honeycombAPIURL := route.Options["apiUrl"]
	if honeycombAPIURL == "" {
		honeycombAPIURL = os.Getenv("HONEYCOMB_API_URL")
	}
	if honeycombAPIURL == "" {
		honeycombAPIURL = DefaultHoneycombAPIURL
	}

	var sampleRate uint = DefaultSampleRate
	sampleRateString := route.Options["sampleRate"]
	if sampleRateString == "" {
		sampleRateString = os.Getenv("HONEYCOMB_SAMPLE_RATE")
	}
	if sampleRateString != "" {
		parsedSampleRate, err := strconv.ParseUint(sampleRateString, 10, 32)
		if err != nil {
			log.Fatal("Must provide Honeycomb SampleRate.")
			return nil, errors.New("Honeycomb 'SampleRate' must be an integer.")
		}
		sampleRate = uint(parsedSampleRate)
	}

	libhoney.Init(libhoney.Config{
		WriteKey:   writeKey,
		Dataset:    dataset,
		APIHost:    honeycombAPIURL,
		SampleRate: sampleRate,
	})

	return &HoneycombAdapter{}, nil
}

// Stream implements the router.LogAdapter interface.
func (a *HoneycombAdapter) Stream(logstream chan *router.Message) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("error getting hostname", err)
	}
	for m := range logstream {

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(m.Data), &data); err != nil {
			// The message is not in JSON.
			// Capture the log line and stash it in "message".
			data = make(map[string]interface{})
			data["message"] = m.Data
		}
		// The message is already in JSON, add the docker specific fields.
		data["stream"] = m.Source
		data["logspout_container"] = m.Container.Name
		data["logspout_container_id"] = m.Container.ID
		data["logspout_hostname"] = m.Container.Config.Hostname
		data["logspout_docker_image"] = m.Container.Config.Image
		data["router_hostname"] = hostname

		if err := libhoney.SendNow(data); err != nil {
			log.Println("error: ", err)
		}
	}
}

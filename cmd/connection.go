package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/cheynewallace/tabby"
	"io"
	"net/http"
	"strconv"
)

func PrintConnection(conns []Connection) {

	t := tabby.New()
	t.AddHeader("Departure", "Duration", "Arrival", "Platform")

	for _, conn := range conns {
		departureTime := UnixToHHMM(conn.Departure.Time)
		arrivalTime := UnixToHHMM(conn.Arrival.Time)
		durationInt, _ := strconv.ParseInt(conn.Duration, 10, 32)

		duration := strconv.FormatInt(durationInt/60, 10) + "m"
		track := conn.Departure.Platform
		t.AddLine(departureTime, duration, arrivalTime, track)
	}

	t.Print()
}

// GetConnections fetches the connection data from the API and returns the response body as a byte slice.
func GetConnections(stationFrom string, stationTo string, time string, arrdep string) ([]byte, error) {
	url := fmt.Sprintf("https://api.irail.be/connections/?from=%s&to=%s&timesel=departure&format=json&lang=en&typeOfTransport=automatic&alerts=false&results=6", stationFrom, stationTo)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// ParseConnections takes the response body as a byte slice and returns an array of Connection structs.
func ParseConnections(body []byte) ([]Connection, error) {
	var result ConnectionResult
	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Connection, nil
}

// Struct definitions remain the same.

type ConnectionResult struct {
	Connection []Connection `json:"connection"`
}

type Connection struct {
	ID        string              `json:"id"`
	Departure ConnectionDeparture `json:"departure"`
	Arrival   ConnectionArrival   `json:"arrival"`
	Duration  string              `json:"duration"`
	Number    string              `json:"number"`
	Vias      *Vias               `json:"vias,omitempty"`
}

type ConnectionDeparture struct {
	Station  string `json:"station"`
	Time     string `json:"time"`  // Unix
	Delay    string `json:"delay"` // seconds
	Canceled string `json:"canceled"`
	Left     string `json:"left"`
	IsExtra  string `json:"isExtra"`
	Vehicle  string `json:"vehicle"`
	Platform string `json:"platform"`
	//Stops    []Stop `json:"stops"`
	//VehicleInfo  VehicleInfo  `json:"vehicleinfo"`
	//	StationInfo  StationInfo  `json:"stationinfo"`
	//PlatformInfo PlatformInfo `json:"platforminfo"`
}

type ConnectionArrival struct {
	Station      string       `json:"station"`
	StationInfo  StationInfo  `json:"stationinfo"`
	Time         string       `json:"time"`  // Unix
	Delay        string       `json:"delay"` // seconds
	Canceled     string       `json:"canceled"`
	Left         string       `json:"left"`
	Platform     string       `json:"platform"`
	PlatformInfo PlatformInfo `json:"platforminfo"`
}

type Stop struct {
	Station string `json:"station"`
	Time    string `json:"time"`  // Unix
	Delay   string `json:"delay"` // seconds
}

type Vias struct {
	Number string    `json:"number"`
	Via    []ViaInfo `json:"via"`
}

type ViaInfo struct {
	ID          string              `json:"id"`
	Station     string              `json:"station"`
	TimeBetween string              `json:"timeBetween"`
	Vehicle     string              `json:"vehicle"`
	Departure   ConnectionDeparture `json:"departure"`
	Arrival     ConnectionArrival   `json:"arrival"`
}

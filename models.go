package tpggo

import (
	"fmt"
	"strings"
	"time"
)

const (
	timeParseLayout = "2006-01-02T15:04:05-0700"
)

// APITime wraps time.Time for handle custom
// formatting
type APITime struct {
	time.Time
}

func (t *APITime) UnmarshalJSON(b []byte) error {

	s := strings.Trim(string(b), `"`)
	if nt, err := time.Parse(timeParseLayout, s); err != nil {
		return err
	} else {
		t.Time = nt
	}

	return nil
}

func (t APITime) Equal(o time.Time) bool {
	return t.Time.Equal(o)
}

func (t APITime) String() string {
	return t.Time.String()
}

// LatLng represents a WGS84 projection of a point
type LatLng struct {
	Lat float64
	Lng float64
}

// APIError returns the information
// about the failing request
type APIError struct {
	Timestamp    APITime `json:"timestamp"`
	ErrorCode    int     `json:"errorCode"`
	ErrorMessage string  `json:"errorMessage"`
}

func (e APIError) error() string {
	return fmt.Sprintf("tpg API error %d at %s: %s", e.ErrorCode, e.Timestamp, e.ErrorMessage)
}

// LineColor contains the visual information
// about a line.
type LineColor struct {
	LineCode   string `json:"lineCode"`
	Hexa       string `json:"hexa"`
	Background string `json:"background"`
	Text       string `json:"text"`
}

// GetLinesResponse is the response to
// a get lines response.
type GetLinesResponse struct {
	Timestamp APITime     `json:"timestamp"`
	Colors    []LineColor `json:"colors"`
}

// GetStopResponse contains the response of the API
// for any GetStops methods.
type GetStopResponse struct {
	Stops     []Stop  `json:"stops"`
	Timestamp APITime `json:"timestamp"`
}

type Stop struct {
	Connections []Connection `json:"connections"`
	Code        string       `json:"stopCode"`
	Name        string       `json:"stopName"`
	Distance    int          `json:"distance"`
}

type Connection struct {
	DestinationCode string `json:"destinationCode"`
	DestinationName string `json:"destinationName"`
	LineCode        string `json:"lineCode"`
}

type GetThermometerResponse struct {
	DestinationCode string      `json:"destinationCode"`
	DestinationName string      `json:"destinationName"`
	LineCode        string      `json:"lineCode"`
	Steps           []Step      `json:"steps"`
	Stop            Stop        `json:"stop"`
	Timestamp       APITime     `json:"timestamp"`
	Disruption      Disruption  `json:"disruptions"`
	Deviations      []Deviation `json:"deviations"`
}

type Step struct {
	DepartureCode int64  `json:"departureCode,omitempty"`
	Deviation     bool   `json:"deviation"`
	DeviationCode string `json:"deviationCode"`
	Reliability   string `json:"reliability,omitempty"`
	Stop          Stop   `json:"stop"`
	Timestamp     string `json:"timestamp"`
	Visible       bool   `json:"visible"`
}

type Deviation struct {
	DeviationCode string `json:"deviationCode"`
}

type GetPhysicalStopsResponse struct {
	Stops     []PhysicalStopInfos `json:"stops"`
	Timestamp APITime             `json:"timestamp"`
}

type PhysicalStop struct {
	Connections      []Connection `json:"connections"`
	Coordinates      Coordinates  `json:"coordinates"`
	PhysicalStopCode string       `json:"physicalStopCode"`
	StopName         string       `json:"stopName"`
}

type Coordinates struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Referential string  `json:"referential"`
}

type PhysicalStopInfos struct {
	PhysicalStops []PhysicalStop `json:"physicalStops"`
	StopCode      string         `json:"stopCode"`
	StopName      string         `json:"stopName"`
}

type GetDisruptionsResponse struct {
	Disruptions []Disruption `json:"disruptions"`
	Timestamp   APITime      `json:"timestamp"`
}

type Disruption struct {
	Code        string  `json:"disruptionCode"`
	Timestamp   APITime `json:"timestamp"`
	Place       string  `json:"place"`
	Consequence string  `json:"consequence"`
	Nature      string  `json:"nature"`
	LineCode    string  `json:"lineCode"`
	StopName    string  `json:"stopName"`
}

type GetNextDeparturesResponse struct {
	Departures []NextDeparture `json:"departures"`
	Stop       Stop            `json:"stop"`
	Timestamp  APITime         `json:"timestamp"`
}

type NextDeparture struct {
	Departure
	Characteristics       string       `json:"characteristics"`
	Disruptions           []Disruption `json:"disruptions"`
	VehiculeNo            int64        `json:"vehiculeNo"`
	VehiculeType          string       `json:"vehiculeType"`
	Deviation             Deviation    `json:"deviation,omitempty"`
	ConnectionWaitingTime int          `json:"connectionWaitingTime,omitempty"`
}

type GetAllNextDeparturesResponse struct {
	Departures []Departure `json:"departures"`
	Stop       Stop        `json:"stop"`
	Timestamp  APITime     `json:"timestamp"`
}
type Departure struct {
	DepartureCode     int64      `json:"departureCode"`
	Line              Connection `json:"line"`
	Reliability       string     `json:"reliability"`
	Timestamp         string     `json:"timestamp"`
	WaitingTime       string     `json:"waitingTime"`
	WaitingTimeMillis int64      `json:"waitingTimeMillis"`
}

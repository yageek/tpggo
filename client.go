package tpggo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	apiScheme     = "https"
	apiHost       = "prod.ivtr-od.tpg.ch"
	apiPathPrefix = "/v1"

	keyPathParameterKey = "key"

	getStopsAPIPath                = "/GetStops"
	getDisruptionsAPIPath          = "/GetDisruptions"
	getGetNextDeparturesAPIPath    = "/GetNextDepartures"
	getGetAllNextDeparturesAPIPath = "/AllGetNextDepartures"
	getThermometerAPIPath          = "/GetThermometer"
	getThermometerPhysicalStops    = "/GetThermometerPhysicalStops"
	getLinesColorsAPIPath          = "/GetLinesColors"
	version                        = "1.0.0"
)

var knownHTTPErrorCodes = []int{http.StatusBadRequest, http.StatusForbidden, http.StatusGone, http.StatusNotFound, http.StatusServiceUnavailable}

// UnknownResponseError represents an unknown server error
type UnknownResponseError struct {
	answer     string
	statusCode int
}

func (e UnknownResponseError) Error() string {
	return fmt.Sprintf("Unknown response %d: %s", e.statusCode, e.answer)
}

// APIClient is an HTTP client for the
// TPG API.
type APIClient struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient create a new APIClient instance with the
// provided key as API key.
func NewClient(key string) *APIClient {
	return &APIClient{
		apiKey: key,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// NewClientWithClient offer the possibility to provide a custom
// http client to use.
func NewClientWithClient(key string, client *http.Client) *APIClient {
	return &APIClient{
		apiKey:     key,
		httpClient: client,
	}
}

func (c *APIClient) apiURL(path string, parameters map[string]string) url.URL {

	parameters[keyPathParameterKey] = c.apiKey

	var query string
	if len(parameters) > 0 {

		values := url.Values{}
		for k, v := range parameters {
			if v != "" {
				values.Add(k, v)
			}
		}
		query = values.Encode()
	}

	apiPath := apiPathPrefix + path + ".json"
	return url.URL{
		Scheme:   apiScheme,
		Host:     apiHost,
		Path:     apiPath,
		RawQuery: query,
	}
}

func (c *APIClient) getAPIData(path string, parameters map[string]string, v interface{}) error {

	URL := c.apiURL(path, parameters)

	req, err := http.NewRequest("GET", URL.String(), nil)
	if err != nil {
		return nil
	}

	userAgentValue := fmt.Sprintf("tpggo %s", version)
	req.Header.Set("User-Agent", userAgentValue)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode == http.StatusOK {
		return decoder.Decode(v)
	}

	for _, code := range knownHTTPErrorCodes {
		if resp.StatusCode == code {

			apiError := APIError{}
			if err := decoder.Decode(&apiError); err != nil {

				return UnknownResponseError{
					statusCode: resp.StatusCode,
					answer:     "Unknown content",
				}
			}
		}
	}

	element, contentError := ioutil.ReadAll(resp.Body)
	var answer string
	if contentError != nil {
		answer = string(element)
	} else {
		answer = contentError.Error()
	}

	return UnknownResponseError{
		statusCode: resp.StatusCode,
		answer:     answer,
	}
}

func (c *APIClient) getStops(stopCode, stopName, line string, pos LatLng) (GetStopResponse, error) {

	params := map[string]string{
		"stopCode": stopCode,
		"stopName": stopName,
		"line":     line,
	}

	if (pos != LatLng{}) {
		params["longitude"] = strconv.FormatFloat(pos.Lng, 'b', -1, 32)
		params["latitude"] = strconv.FormatFloat(pos.Lat, 'b', -1, 32)
	}

	response := GetStopResponse{}
	err := c.getAPIData(getStopsAPIPath, params, &response)
	return response, err
}

// GetStops returns the list of all the stops.
// Results are send by ascending order regarding the stop code.
func (c *APIClient) GetStops() (GetStopResponse, error) {
	return c.getStops("", "", "", LatLng{})
}

func sortedJoinList(list []string) string {

	sort.Strings(list)
	codeParams := strings.Join(list, ",")
	return codeParams
}

// GetStopsFromCodes returns the list of stops whose stopCode is contained in the `codes` list.
// Results are send by ascending order regarding the stop code.
func (c *APIClient) GetStopsFromCodes(codes []string) (GetStopResponse, error) {
	return c.getStops(sortedJoinList(codes), "", "", LatLng{})
}

// GetStopsByName returns the list of stops whose name include the substring `name`.
// Results are send by ascending order regarding the stop code.
func (c *APIClient) GetStopsByName(name string) (GetStopResponse, error) {
	return c.getStops("", name, "", LatLng{})
}

// GetStopsByLine returns the list of stops on the provided line.
// Results are send by ascending order regarding the stop code.
func (c *APIClient) GetStopsByLine(name string) (GetStopResponse, error) {
	return c.getStops("", name, "", LatLng{})
}

// GetStopsClosedToLatLng returns all the stops located inside a 500 meter range
// from the provided position.
// Results are send by descending distance from the provided position.
func (c *APIClient) GetStopsClosedToLatLng(pos LatLng) (GetStopResponse, error) {
	return c.getStops("", "", "", pos)
}

func (c *APIClient) getPhysicalStops(stopCode, stopName string) (GetPhysicalStopsResponse, error) {

	params := map[string]string{
		"stopCode": stopCode,
		"stopName": stopName,
	}

	response := GetPhysicalStopsResponse{}
	err := c.getAPIData(getStopsAPIPath, params, &response)
	return response, err
}

// GetPhysicalStopsFromCodes returns the list of  physicalstops whose stopCode is contained in the `codes` list.
// Results are send by ascending order regarding the stop code.
func (c *APIClient) GetPhysicalStopsFromCodes(codes []string) (GetPhysicalStopsResponse, error) {
	return c.getPhysicalStops(sortedJoinList(codes), "")
}

// GetPhysicalStopsByName returns the list of physical stops whose stopCode is contained in the `codes` list.
// Results are send by ascending order regarding the stop code.
func (c *APIClient) GetPhysicalStopsByName(name string) (GetPhysicalStopsResponse, error) {
	return c.getPhysicalStops("", name)
}

func (c *APIClient) getNextDepartures(stopCode, departureCode, linesCode, destinationsCode string) (GetNextDeparturesResponse, error) {

	params := map[string]string{
		"stopCode":         stopCode,
		"departureCode":    departureCode,
		"linesCode":        linesCode,
		"destinationsCode": destinationsCode,
	}

	response := GetNextDeparturesResponse{}
	err := c.getAPIData(getGetNextDeparturesAPIPath, params, &response)
	return response, err
}

// GetNextDeparturesForStop returns the list of departures for the providedstop.
// If departure code is an empty value, all the next departure corresponding to the
// stop code will be returned.
func (c *APIClient) GetNextDeparturesForStop(stopCode, departureCode string) (GetNextDeparturesResponse, error) {
	return c.getNextDepartures(stopCode, departureCode, "", "")
}

// GetNextDeparturesForLines returns the list of departures corresponding to the
// provided lines and destinations codes.
func (c *APIClient) GetNextDeparturesForLines(linesCodes, destinationsCodes []string) (GetNextDeparturesResponse, error) {
	return c.getNextDepartures("", "", sortedJoinList(linesCodes), sortedJoinList(destinationsCodes))
}

// GetAllNextDepartures returns the next departures.
func (c *APIClient) GetAllNextDepartures(stopCode, lineCode, destinationCode string) (GetAllNextDeparturesResponse, error) {
	params := map[string]string{
		"stopCode":        stopCode,
		"lineCode":        lineCode,
		"destinationCode": destinationCode,
	}

	response := GetAllNextDeparturesResponse{}
	err := c.getAPIData(getGetAllNextDeparturesAPIPath, params, &response)
	return response, err
}

// GetThermometer returns the list of stops with
// the disruptions and deviations.
func (c *APIClient) GetThermometer(departureCode string) (GetThermometerResponse, error) {
	params := map[string]string{
		"departureCode": departureCode,
	}

	response := GetThermometerResponse{}
	err := c.getAPIData(getThermometerAPIPath, params, &response)
	return response, err
}

// GetThermometerPhysicalStops returns the list of stops with
// the disruptions and deviations.
func (c *APIClient) GetThermometerPhysicalStops(departureCode string) (GetThermometerResponse, error) {
	params := map[string]string{
		"departureCode": departureCode,
	}

	response := GetThermometerResponse{}
	err := c.getAPIData(getThermometerPhysicalStops, params, &response)
	return response, err
}

// GetLinesColor returns visual informations about the lines
func (c *APIClient) GetLinesColor() (GetLinesResponse, error) {

	response := GetLinesResponse{}
	err := c.getAPIData(getLinesColorsAPIPath, map[string]string{}, &response)
	return response, err
}

// GetDisruptions returns the list of the disruptions on the network.
func (c *APIClient) GetDisruptions() (GetDisruptionsResponse, error) {
	response := GetDisruptionsResponse{}
	err := c.getAPIData(getLinesColorsAPIPath, map[string]string{}, &response)
	return response, err
}

package tpggo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestURLSerialization(t *testing.T) {

	fakeKey := "foo"
	client := NewClient(fakeKey)
	expected := fmt.Sprintf("https://prod.ivtr-od.tpg.ch/v1/GetDisruptions.json?key=%s", fakeKey)
	URL := client.apiURL(getDisruptionsAPIPath, map[string]string{})
	assert.Equal(t, expected, URL.String(), "The formatted url does not have the exepected values")

	expected = fmt.Sprintf("https://prod.ivtr-od.tpg.ch/v1/GetDisruptions.json?key=%s&key1=value1", fakeKey)
	URL = client.apiURL(getDisruptionsAPIPath, map[string]string{"key1": "value1"})
	assert.Equal(t, expected, URL.String(), "The formatted url does not have the exepected values")
}

func TestParseTimeLayout(t *testing.T) {
	zone := time.FixedZone("UTC-8", 60*60)
	expected := time.Date(2018, 12, 14, 8, 34, 36, 0, zone)
	computed, err := time.Parse(timeParseLayout, "2018-12-14T08:34:36+0100")
	if err != nil {
		t.Errorf("Invalid layout for time parsing: %s", err)
		t.FailNow()
	}

	assert.True(t, expected.Equal(computed), "Times do not match")

}

func TestUnmarshallTime(t *testing.T) {

	zone := time.FixedZone("UTC-8", 60*60)
	expected := time.Date(2018, 12, 14, 8, 34, 36, 0, zone)
	inStruct := struct {
		Time APITime `json:"time"`
	}{}
	js := `{
		"time": "2018-12-14T08:34:36+0100"
	}`

	err := json.NewDecoder(bytes.NewBufferString(js)).Decode(&inStruct)
	if err != nil {
		t.Errorf("JSON invalid: %s", err)
		t.FailNow()
	}
	assert.True(t, expected.Equal(inStruct.Time.Time), "Time should be equal")
}

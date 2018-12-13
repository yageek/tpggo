package tpggo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLSerialization(t *testing.T) {

	fakeKey := "foo"
	client := NewClient(fakeKey)
	expected := "https://prod.ivtr-od.tpg.ch/v1/GetDisruptions.json?key=foo"
	URL := client.apiURL(getDisruptionsAPIPath, map[string]string{})
	assert.Equal(t, expected, URL.String(), "The formatted url does not have the exepected values")

	expected = "https://prod.ivtr-od.tpg.ch/v1/GetDisruptions.json?key=foo&key1=value1"
	URL = client.apiURL(getDisruptionsAPIPath, map[string]string{"key1": "value1"})
	assert.Equal(t, expected, URL.String(), "The formatted url does not have the exepected values")
}

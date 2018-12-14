package tpggo

import (
	"fmt"
	"testing"

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

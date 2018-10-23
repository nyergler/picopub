package mf_test

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nyergler/picopub/mf"
)

func TestCreateEntryFromURL(t *testing.T) {
	body := url.Values{
		"h":        []string{"entry"},
		"title":    []string{"Testing"},
		"content":  []string{"Hello world"},
		"category": []string{"recipes"},
	}

	request := httptest.NewRequest("POST", fmt.Sprintf("/micropub?%s", body.Encode()), strings.NewReader(""))
	result := new(mf.MicroformatObject)

	err := mf.ParseRequest(request, result)
	require.Nil(t, err)

	assert.Equal(t, []string{"Testing"}, result.Properties["title"].Content())
	assert.Equal(t, []string{"Hello world"}, result.Properties["content"].Content())
	assert.Equal(t, []string{"recipes"}, result.Properties["category"].Content())
	assert.Equal(t, mf.Entry, result.Type)
}

func TestCreateErrorsOnUnknownH(t *testing.T) {
	body := url.Values{
		"h": []string{"blarf"},
	}

	request := httptest.NewRequest("POST", fmt.Sprintf("/micropub?%s", body.Encode()), strings.NewReader(""))
	result := new(mf.MicroformatObject)

	err := mf.ParseRequest(request, result)
	assert.NotNil(t, err)
}

func TestArrayHandling(t *testing.T) {
	body := url.Values{
		"h":          []string{"entry"},
		"title":      []string{"Testing"},
		"content":    []string{"Hello world"},
		"category[]": []string{"recipes"},
	}

	request := httptest.NewRequest("POST", fmt.Sprintf("/micropub?%s", body.Encode()), strings.NewReader(""))
	result := new(mf.MicroformatObject)

	err := mf.ParseRequest(request, result)
	require.Nil(t, err)

	assert.Equal(t, []string{"recipes"}, result.Properties["category"].Content())
	assert.Equal(t, mf.Entry, result.Type)
}

func TestCreateEntryMultipart(t *testing.T) {
	var buf bytes.Buffer
	mpW := multipart.NewWriter(&buf)
	mpW.WriteField("h", "entry")
	mpW.WriteField("content", "Hello Multipart")
	mpW.Close()

	request := httptest.NewRequest("POST", "/micropub", &buf)
	request.Header.Set("Content-Type", mpW.FormDataContentType())
	result := new(mf.MicroformatObject)

	err := mf.ParseRequest(request, result)
	require.Nil(t, err)

	assert.Equal(t, []string{"Hello Multipart"}, result.Properties["content"].Content())
	assert.Equal(t, mf.Entry, result.Type)
}

func TestCreateEntry(t *testing.T) {
	payload := `
	{
		"type": ["h-entry"],
		"properties": {
			"content": ["hello world"],
			"photo": ["https://photos.example.com/592829482876343254.jpg"]
		}
	}`

	request := httptest.NewRequest("POST", "/micropub", strings.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	result := new(mf.MicroformatObject)

	err := mf.ParseRequest(request, result)
	require.Nil(t, err)

	assert.Equal(t, mf.Entry, result.Type)
	assert.Equal(t, []string{"hello world"}, result.Properties["content"].Content())
}

func TestNestedObjects(t *testing.T) {
	payload := `
	  {
		"type": ["h-entry"],
		"properties": {
		  "summary": [
			"Weighed 70.64 kg"
		  ],
		  "weight": [
			{
			  "type": ["h-measure"],
			  "properties": {
				"num": ["70.64"],
				"unit": ["kg"]
			  }
			}
		  ],
		  "bodyfat": [
			{
			  "type": ["h-measure"],
			  "properties": {
				"num": ["19.83"],
				"unit": ["%"]
			  }
			}
		  ]
		}
	  }
	`

	request := httptest.NewRequest("POST", "/micropub", strings.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	result := new(mf.MicroformatObject)

	err := mf.ParseRequest(request, result)
	require.Nil(t, err)

	assert.Equal(t, mf.Entry, result.Type)
	assert.Equal(t, []string{"Weighed 70.64 kg"}, result.Properties["summary"].Content())
	assert.Equal(t, []string{"kg"},
		result.Properties["weight"].Object().Properties["unit"].Content(),
	)
}

package picopub

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestEntryStore struct {
	entries []*MicroformatObject
}

func (s *TestEntryStore) Length() int {
	return len(s.entries)
}

func (s *TestEntryStore) Get(i int) *MicroformatObject {
	return s.entries[i]
}

func (s *TestEntryStore) Append(p *MicroformatObject) {
	s.entries = append(s.entries, p)
}

func TestCreateEntryFromURL(t *testing.T) {
	store := &TestEntryStore{}
	handler := NewHandler(store)
	body := url.Values{
		"h":        []string{"entry"},
		"title":    []string{"Testing"},
		"content":  []string{"Hello world"},
		"category": []string{"recipes"},
	}

	request := httptest.NewRequest("POST", fmt.Sprintf("/micropub?%s", body.Encode()), strings.NewReader(""))
	response := &httptest.ResponseRecorder{}
	handler.ServeHTTP(response, request)

	assert.Equal(t, 201, response.Code)
	assert.Equal(t, store.Length(), 1)
	assert.Equal(t, []string{"Testing"}, store.Get(0).Properties["title"].Content())
	assert.Equal(t, []string{"Hello world"}, store.Get(0).Properties["content"].Content())
	assert.Equal(t, []string{"recipes"}, store.Get(0).Properties["category"].Content())
	assert.Equal(t, Entry, store.Get(0).Type)
}

func TestCreateErrorsOnUnknownH(t *testing.T) {
	store := &TestEntryStore{}
	handler := NewHandler(store)
	body := url.Values{
		"h": []string{"blarf"},
	}

	request := httptest.NewRequest("POST", fmt.Sprintf("/micropub?%s", body.Encode()), strings.NewReader(""))
	response := &httptest.ResponseRecorder{}
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestArrayHandling(t *testing.T) {
	store := &TestEntryStore{}
	handler := NewHandler(store)
	body := url.Values{
		"h":          []string{"entry"},
		"title":      []string{"Testing"},
		"content":    []string{"Hello world"},
		"category[]": []string{"recipes"},
	}

	request := httptest.NewRequest("POST", fmt.Sprintf("/micropub?%s", body.Encode()), strings.NewReader(""))
	response := &httptest.ResponseRecorder{}
	handler.ServeHTTP(response, request)

	assert.Equal(t, 201, response.Code)
	assert.Equal(t, store.Length(), 1)
	assert.Equal(t, []string{"recipes"}, store.Get(0).Properties["category"].Content())
	assert.Equal(t, Entry, store.Get(0).Type)

}

func TestCreateEntryMultipart(t *testing.T) {
	store := &TestEntryStore{}
	handler := NewHandler(store)

	var buf bytes.Buffer
	mpW := multipart.NewWriter(&buf)
	mpW.WriteField("h", "entry")
	mpW.WriteField("content", "Hello Multipart")
	mpW.Close()

	request := httptest.NewRequest("POST", "/micropub", &buf)
	request.Header.Set("Content-Type", mpW.FormDataContentType())
	response := &httptest.ResponseRecorder{}
	handler.ServeHTTP(response, request)

	assert.Equal(t, 201, response.Code)
	assert.Equal(t, store.Length(), 1)
	assert.Equal(t, []string{"Hello Multipart"}, store.Get(0).Properties["content"].Content())
	assert.Equal(t, Entry, store.Get(0).Type)
}

func TestCreateEntry(t *testing.T) {
	store := &TestEntryStore{}
	handler := NewHandler(store)
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
	response := &httptest.ResponseRecorder{}

	handler.ServeHTTP(response, request)

	assert.Equal(t, 201, response.Code)
	assert.Equal(t, store.Length(), 1)
	assert.Equal(t, Entry, store.Get(0).Type)
	assert.Equal(t, []string{"hello world"}, store.Get(0).Properties["content"].Content())
}

func TestNestedObjects(t *testing.T) {
	store := &TestEntryStore{}
	handler := NewHandler(store)
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
	response := &httptest.ResponseRecorder{}

	handler.ServeHTTP(response, request)

	assert.Equal(t, 201, response.Code)
	assert.Equal(t, store.Length(), 1)
	assert.Equal(t, Entry, store.Get(0).Type)
	assert.Equal(t, []string{"Weighed 70.64 kg"}, store.Get(0).Properties["summary"].Content())
	assert.Equal(t, []string{"kg"},
		store.Get(0).Properties["weight"].Object().Properties["unit"].Content(),
	)
}

// func TestParseUpdate(t *testing.T) {

// }

// func TestParseUpdateFromJOSN(t *testing.T) {

// }

// func TestParseDelete(t *testing.T) {

// }

// func TestParseDeleteFromJSON(t *testing.T) {

// }

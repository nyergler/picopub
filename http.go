package picopub

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type EntryStore interface {
	Length() int
	Get(i int) *MicroformatObject
	Append(*MicroformatObject)
}

type Handler struct {
	store EntryStore
}

// NewHandler returns an HTTP Handler which serves micropub
func NewHandler(store EntryStore) http.Handler {

	return &Handler{store: store}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// q := r.URL.Query()

	switch r.Method {
	case http.MethodGet:
		break
	case http.MethodPost:
		// create
		entry := new(MicroformatObject)
		contentType := r.Header.Get("Content-Type")
		if contentType == "application/json" {
			dec := json.NewDecoder(r.Body)
			if err := dec.Decode(entry); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
		} else {
			if strings.HasPrefix(contentType, "multipart/form-data") {
				err := r.ParseMultipartForm(http.DefaultMaxHeaderBytes)
				if err != nil {
					log.Fatal(err)
				}
				formFields := (url.Values)(r.MultipartForm.Value)
				ParseForm(&formFields, entry)
			} else {
				r.ParseForm()
				ParseForm(&r.Form, entry)
			}
		}
		if entry.Type == Unknown {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.store.Append(entry)
		w.WriteHeader(http.StatusCreated)
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

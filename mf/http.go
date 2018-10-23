package mf

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func ParseRequest(r *http.Request, entry *MicroformatObject) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(entry); err != nil {
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

	return nil
}

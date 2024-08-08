package core

import (
	"bytes"
	"encoding/base64"
	ics "github.com/arran4/golang-ical"
	"io"
	"net/http"
)

type httpOption func(*http.Request)

func withAuth(username string, password string) httpOption {

	b := []byte{}
	base64.RawStdEncoding.Encode([]byte(username+password), b)

	return func(r *http.Request) {
		r.Header.Add("Authorization", string(b))
	}
}

func parseICS(r io.Reader) (*ics.Calendar, error) {
	return ics.ParseCalendar(r)
}

func parseICSURI(uri string, opts ...httpOption) (*ics.Calendar, error) {
	req, err := http.NewRequest(http.MethodGet, uri, &bytes.Buffer{})
	if err != nil {
		return nil, err
	}
    
    _, err = parseICS(&bytes.Buffer{})

	req.Header.Add("Content-Type", "text/calendar")

	for _, opt := range opts {
		opt(req)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	cal, err := parseICS(res.Body)
	if err != nil {
		return nil, err
	}

	return cal, nil
}

func serialize(cal ics.Calendar) []byte {
	buf := bytes.Buffer{}
	cal.SerializeTo(&buf)
	return buf.Bytes()
}

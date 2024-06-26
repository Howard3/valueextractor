package valueextractor

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// ErrNotFound is an error that is returned when a key is not found
var ErrNotFound = errors.New("key not found")

// extractors
type ValueExtractor interface {
	Get(key string) (val string, err error)
}

// MapExtractor is a value extractor that extracts values from a map
type MapExtractor map[string]string

// Get returns the value of a key from the map
func (m MapExtractor) Get(key string) (string, error) {
	value, ok := m[key]

	if !ok {
		return "", ErrNotFound
	}

	return value, nil
}

// QueryExtractor is a value extractor that extracts values from a http request's query parameters
type QueryExtractor struct {
	Query url.Values
}

// Get returns the value of a query parameter from the request
func (qe QueryExtractor) Get(key string) (string, error) {
	value := qe.Query.Get(key)
	if value == "" {
		return "", ErrNotFound
	}
	return value, nil
}

// ErrRequestNil is an error that is returned when the request is nil
var ErrRequestNil = errors.New("request is nil")
var ErrRequestParseForm = errors.New("error parsing form")

// FormExtractor is a value extractor that extracts values from a http request's form
type FormExtractor struct {
	Request *http.Request
	parsed  bool
	getter  func(string) string
}

func (fe *FormExtractor) isMultipart() bool {
	ctype := fe.Request.Header.Get("Content-Type")
	return strings.HasPrefix(ctype, "multipart/form-data")
}

// ensureParsed ensures that the form has been parsed
// for multipart forms, it parses the form using ParseMultipartForm and sets the getter to FormValue
// for urlencoded forms, it parses the form using ParseForm and sets the getter to Get
func (fe *FormExtractor) ensureParsed() error {
	if fe.parsed {
		return nil
	}

	fe.getter = fe.Request.FormValue

	if fe.isMultipart() {
		if err := fe.Request.ParseMultipartForm(0); err != nil {
			return errors.Join(ErrRequestParseForm, err)
		}
		fe.getter = fe.Request.FormValue
	} else if err := fe.Request.ParseForm(); err != nil {
		return errors.Join(ErrRequestParseForm, err)
	}

	fe.parsed = true

	return nil
}

// Get returns the value of a form parameter from the Request
func (fe *FormExtractor) Get(key string) (string, error) {
	if fe.Request == nil {
		return "", ErrRequestNil
	}

	if err := fe.ensureParsed(); err != nil {
		return "", err
	}

	value := fe.getter(key)
	if value == "" {
		return "", ErrNotFound
	}

	return value, nil
}

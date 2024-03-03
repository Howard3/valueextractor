package valueextractor

import (
	"errors"
	"fmt"
)

// Extractor is a value extractor that can be used to extract values from a request
// and type-convert them to the desired type, collecting errors along the way
type Extractor struct {
	extractor ValueExtractor
	errors    []error
}

// With taks a key and a converter and extracts the value from the request
func (ec *Extractor) With(key string, converter Converter) {
	str, err := ec.extractor.Get(key)
	if err != nil {
		ec.AddError(key, err)
		return
	}

	if err := converter(ec, str); err != nil {
		ec.AddError(key, err)
	}
}

// WithOptional is a function that ignores if the error is ErrNotFound
func (ec *Extractor) WithOptional(key string, converter Converter) {
	str, err := ec.extractor.Get(key)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return
		}

		ec.AddError(key, err)
		return
	}

	if err := converter(ec, str); err != nil {
		ec.AddError(key, err)
	}
}

// AddError adds an error to the chain
func (ec *Extractor) AddError(key string, err error) *Extractor {
	ec.errors = append(ec.errors, fmt.Errorf("error extracting key %s: %w", key, err))
	return ec
}

// Using creates a new Extractor with the given value extractor
// A value extractor is a function that takes a key and returns a value and an error, if any
func Using(extractor ValueExtractor) *Extractor {
	return &Extractor{extractor: extractor}
}

// Errors returns an error if there are any errors in the parser
func (ec *Extractor) Errors() error {
	if len(ec.errors) == 0 {
		return nil
	}

	var errMsgs error
	for _, err := range ec.errors {
		errMsgs = errors.Join(errMsgs, err)
	}

	return errMsgs
}

package valueextractor

import (
	"errors"
	"fmt"
)

// Extractor is a value extractor that can be used to extract values from a request
// and type-convert them to the desired type, collecting errors along the way
type Extractor struct {
	extractor    ValueExtractor
	errors       []*Error
	optionalKeys []string
}

func WithOptionalKeys(keys ...string) func(*Extractor) {
	return func(ex *Extractor) {
		ex.optionalKeys = keys
	}
}

func isOptional(ex *Extractor, key string) bool {
	for _, k := range ex.optionalKeys {
		if k == key {
			return true
		}
	}

	return false
}

// With taks a key and a converter and extracts the value from the request
func (ec *Extractor) With(key string, converter Converter) {
	str, err := ec.extractor.Get(key)
	if err != nil {
		if errors.Is(err, ErrNotFound) && isOptional(ec, key) {
			return
		}

		ec.AddExtractError(key, err)
		return
	}

	if err := converter(ec, str); err != nil {
		ec.AddConvertError(key, err)
	}
}

// AddExtractError adds an error to the chain
func (ec *Extractor) AddExtractError(key string, err error) {
	ec.errors = append(ec.errors, NewExtractError(key, err))
}

// AddConvertError adds an error to the chain
func (ec *Extractor) AddConvertError(key string, err error) {
	ec.errors = append(ec.errors, NewConvertError(key, err))
}

// Using creates a new Extractor with the given value extractor
// A value extractor is a function that takes a key and returns a value and an error, if any
func Using(extractor ValueExtractor, options ...func(*Extractor)) *Extractor {
	ex := &Extractor{extractor: extractor}
	for _, option := range options {
		option(ex)
	}

	return ex
}

// Errors returns an error if there are any errors in the parser
func (ec *Extractor) Errors() []*Error {
	if len(ec.errors) == 0 {
		return nil
	}

	return ec.errors
}

// JoinedErrors returns a single error with all the errors JoinedErrors
func (ec *Extractor) JoinedErrors() error {
	if len(ec.errors) == 0 {
		return nil
	}

	var err error
	for _, e := range ec.errors {
		err = errors.Join(err, fmt.Errorf("%s: %w", e.key, e.err))
	}

	return err
}

// ResultConverter defines a wrapped converter with input argument as a reference that returns
// a converter function. It's intended to be used with the Result function
type ResultConverter[T any] func(*T) Converter

// Result is a function that extracts a value from the request and converts it to the desired type
// It offers a simpler API than the With function
// NOTE: this is a convenience function but is slower than using the With function directly
// the performance is about on par with the traditional struct+reflection approach
func Result[T any](ex *Extractor, key string, converter ResultConverter[T]) T {
	var result T
	ex.With(key, converter(&result))
	return result
}

// ResultPtr is a function that extracts a value from the request and converts it to the desired type
// It offers a simpler API than the With function, but returns a pointer to the result instead of the result itself
// NOTE: this is a convenience function but is slower than using the With function directly
// the performance is about on par with the traditional struct+reflection approach
func ResultPtr[T any](ex *Extractor, key string, converter ResultConverter[T]) *T {
	var result T
	ex.With(key, converter(&result))
	return &result
}

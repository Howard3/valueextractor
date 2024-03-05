package valueextractor

import (
	"fmt"
	"strconv"
)

// Converter is a function that takes an Extractor and a key and returns a value and an error
type Converter func(ec *Extractor, value string) error

// DirectReturnType is a function that takes an Extractor and a key and returns a value
// this is a more performant alternative to the Result generic.
type DirectReturnType func(ec *Extractor, key string) interface{}

func ReturnString(ec *Extractor, key string) *string {
	var s string
	ec.With(key, AsString(&s))
	return &s
}

// AsString maintains the value as a string, just allowing extraction
func AsString(ref *string) Converter {
	return func(ec *Extractor, value string) error {
		*ref = value
		return nil
	}
}

// AsUint is a function that converts a string to a uint64
func AsUint64(ref *uint64) Converter {
	return func(ec *Extractor, key string) error {
		parsed, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid uint value: %v", err)
		}

		*ref = parsed
		return nil
	}
}

// ReturnUint64 is a function that returns a uint64
func ReturnUint64(ec *Extractor, key string) *uint64 {
	var i uint64
	ec.With(key, AsUint64(&i))
	return &i
}

// AsInt64 is a function that converts a string to an int64
func AsInt64(ref *int64) Converter {
	return func(ec *Extractor, value string) error {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int value: %v", err)
		}

		*ref = parsed
		return nil
	}
}

// ReturnInt64 is a function that returns an int64
func ReturnInt64(ec *Extractor, key string) *int64 {
	var i int64
	ec.With("age", AsInt64(&i))
	return &i
}

// AsFloat64 is a function that converts a string to a float64
func AsFloat64(ref *float64) Converter {
	return func(ec *Extractor, value string) error {
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value: %v", err)
		}

		*ref = parsed
		return nil
	}
}

// ReturnFloat64 is a function that returns a float64
func ReturnFloat64(ec *Extractor, key string) *float64 {
	var i float64
	ec.With("age", AsFloat64(&i))
	return &i
}

// AsBool is a function that converts a string to a bool
func AsBool(ref *bool) Converter {
	return func(ec *Extractor, value string) error {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid bool value: %v", err)
		}

		*ref = parsed
		return nil
	}
}

// ReturnBool is a function that returns a bool
func ReturnBool(ec *Extractor, key string) *bool {
	var i bool
	ec.With("age", AsBool(&i))
	return &i
}

package valueextractor

import (
	"fmt"
	"strconv"
)

// Converter is a function that takes an Extractor and a key and returns a value and an error
type Converter func(ec *Extractor, value string) error

// AsString maintains the value as a string, just allowing extraction
func AsString(ref *string) Converter {
	return func(ec *Extractor, value string) error {
		*ref = value
		return nil
	}
}

// AsUint is a function that converts a string to a uint64
func AsUint64(ref *uint64) Converter {
	return func(ec *Extractor, value string) error {
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid uint value: %v", err)
		}

		*ref = parsed
		return nil
	}
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

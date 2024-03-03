package valueextractor

type ErrorType int

const (
	ExtractError ErrorType = iota
	ConvertError
)

// Error is a custom error type that is used to represent errors that occur during the extraction and conversion of values.
type Error struct {
	eType ErrorType
	key   string
	err   error
}

// Error returns the error message.
func (e Error) Error() string {
	return e.err.Error()
}

// IsExtractError returns true if the error is of type ExtractError.
func (e Error) IsExtractError() bool {
	return e.eType == ExtractError
}

// IsConvertError returns true if the error is of type ConvertError.
func (e Error) IsConvertError() bool {
	return e.eType == ConvertError
}

// Key returns the key that was used to extract the value.
func (e Error) Key() string {
	return e.key
}

// NewExtractError creates a new ExtractError.
func NewExtractError(key string, err error) *Error {
	return &Error{eType: ExtractError, key: key, err: err}
}

// NewConvertError creates a new ConvertError.
func NewConvertError(key string, err error) *Error {
	return &Error{eType: ConvertError, key: key, err: err}
}

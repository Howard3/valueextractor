# Value Extractor

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/howard3/valueextractor)
[![Go Report Card](https://goreportcard.com/badge/github.com/Howard3/valueextractor)](https://goreportcard.com/report/github.com/Howard3/valueextractor)
![GitHub license](https://img.shields.io/github/license/Howard3/valueextractor)
![GitHub stars](https://img.shields.io/github/stars/Howard3/valueextractor?style=social)
![GitHub forks](https://img.shields.io/github/forks/Howard3/valueextractor?style=social)
[![Go Reference](https://pkg.go.dev/badge/github.com/Howard3/valueextractor.svg)](https://pkg.go.dev/github.com/Howard3/valueextractor)

This package provides a flexible system for extracting values from various sources (like HTTP requests, maps, etc.) and converting them into the desired Go types with comprehensive error handling. It's designed to make it easier to work with dynamic data sources in a type-safe manner.

## Features

- **Extensible Value Extraction**: Supports extracting values from maps, HTTP request query parameters, and form data.
- **Type Conversion**: Convert extracted strings into specific Go types (`string`, `uint64`, `int64`), including custom type conversions.
- **Error Handling**: Collects and aggregates errors throughout the extraction and conversion process for robust error reporting.
- **No external dependencies** - only the Go standard library
- **Fast** - run `go test -bench .`. It's about **20% faster** than the idiomatic struct tag + reflection
- **Lighweight** - Less than 300 LOC 
- **Easy to read** - Only 1 `if err != nil` for all of your conversions.
 
## Getting Started

### Installation

```sh
go get github.com/Howard3/valueextractor
```

### Basic Usage

Here's a quick example to get you started:

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/Howard3/valueextractor"
)

func main() {
    // Example: Extracting values from a query parameter
    req, _ := http.NewRequest("GET", "/?id=123&name=John", nil)
    queryExtractor := valueextractor.QueryExtractor{Query: req.URL.Query()}

    var id uint64
    var name string
    extractor := valueextractor.Using(queryExtractor)
    extractor.With("id", valueextractor.AsUint64(&id))
    extractor.With("name", valueextractor.AsString(&name))

    if err := extractor.Errors(); err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("Extracted values - ID: %d, Name: %s\n", id, name)
    }
}
```

Alternatively this can be written with fewer lines using the `Result` function.
```go
package main

import (
    "fmt"
    "net/http"
    "github.com/Howard3/valueextractor"
)

func main() {
    // Example: Extracting values from a query parameter
    req, _ := http.NewRequest("GET", "/?id=123&name=John", nil)
    queryExtractor := valueextractor.QueryExtractor{Query: req.URL.Query()}

    ex := valueextractor.Using(queryExtractor)
    id := valueextractor.Result(ex, "id", valueextractor.AsUint64)
    name := valueextractor.Result(ex, "name", valueextractor.AsString)

    if err := ex.Errors(); err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("Extracted values - ID: %d, Name: %s\n", id, name)
    }
}
```
## Using `WithOptional` for Optional Values

The `WithOptional` function is designed to streamline the handling of optional values. When extracting a value that may or may not be present (and its absence is not considered an error), use `WithOptional` instead of `With`. This function behaves similarly to `With`, attempting to extract and convert a value, but it will ignore `ErrNotFound` errors. This is particularly useful for working with HTTP requests where certain parameters may be optional.

### Example Usage of `WithOptional`

```go
var age uint64
extractor.WithOptional("age", valueextractor.AsUint64(&age))
```

In this example, if the "age" parameter is missing from the data source, `WithOptional` will not add an `ErrNotFound` to the error chain, allowing the application to proceed without treating the absence of "age" as an error.

By leveraging `WithOptional`, you can write cleaner, more concise code for handling optional parameters without cluttering your error handling logic with checks for missing but non-essential data.


## Extensibility

The system is designed with extensibility in mind. You can extend it by implementing custom `ValueExtractor` interfaces or `Converter` functions.

### Implementing a Custom ValueExtractor

To extract values from a new data source, implement the `ValueExtractor` interface:

```go
type CustomExtractor struct {}

func (ce CustomExtractor) Get(key string) (string, error) {
    // Logic to extract and return the value based on the key
}
```

### Adding a New Converter Function

To support converting to a new type:

```go
func AsCustomType(ref *CustomType) valueextractor.Converter {
    return func(ec *valueextractor.Extractor, value string) error {
        // Convert value to CustomType and assign to ref
        
        *ref = newVal // assign the new value

        return nil
    }
}
```

## Error Handling

Errors are collected throughout the extraction and conversion process. Use the `Errors` method to retrieve any accumulated errors:

```go
if err := extractor.Errors(); err != nil {
    // Handle error
}
```

In addition to the general error handling mechanism provided by the `Extractor` system, there are specific errors that users should be aware of when working with the value extraction and conversion system. Understanding these errors can help in diagnosing and handling common issues that may arise during the extraction and conversion process.

#### Defined Errors

- **ErrNotFound**: This error is returned when the specified key is not found within the source (e.g., map, query parameters, form). It indicates that the requested value for conversion does not exist.

#### Common Errors During Conversion

- **Invalid uint value**: This error occurs when attempting to convert a string to a `uint64` and the string does not represent a valid unsigned integer. It indicates a format or value error in the source string.
- **Invalid int value**: Similar to the "Invalid uint value" error, this error is returned when a string cannot be successfully converted to an `int64` due to formatting issues or value constraints.

#### Handling Specific Errors

When using the `Extractor` system, it's important to check for these specific errors where appropriate. You can use Go's `errors.Is` function to check for a specific error type. Here’s how you might handle `ErrNotFound` to differentiate between missing values and other errors:

```go
var id uint64
extractor := valueextractor.Using(queryExtractor)
extractor.With("id", valueextractor.AsUint64(&id))

if err := extractor.Errors(); err != nil {
    if errors.Is(err, valueextractor.ErrNotFound) {
        fmt.Println("The specified key was not found.")
    } else {
        fmt.Println("Error:", err)
    }
}
```

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

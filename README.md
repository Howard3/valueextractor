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
- **Type Conversion**: Convert extracted strings into specific Go types (`string`, `uint64`, `int64`, `float64`, `bool`), including custom type conversions.
- **Direct Return Types**: Provides direct return functions (`ReturnString`, `ReturnUint64`, etc.) for performance optimization.
- **Error Handling**: Collects and aggregates errors throughout the extraction and conversion process for robust error reporting.
- **No external dependencies** - only the Go standard library.
- **Fast** - run `go test -bench .`. It's between **20-40% faster** than the idiomatic struct tag + reflection.
- **Lighweight** - Less than 300 LOC.
- **Easy to read** - Only 1 `if err != nil` for all of your conversions.
 
## Getting Started

### Installation

```sh
go get github.com/Howard3/valueextractor
```

### Basic Usage

Here are quick examples to get you started:

There are three main ways to use this library:
- References
- Direct value return objects
- Return Generics

#### Using References
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

##### Using `WithOptional` for Optional Values
The `WithOptional` function is designed to streamline the handling of optional values. When extracting a value that may or may not be present (and its absence is not considered an error), use `WithOptional` instead of `With`. This function behaves similarly to `With`, attempting to extract and convert a value, but it will ignore `ErrNotFound` errors. This is particularly useful for working with HTTP requests where certain parameters may be optional.
Similarly the Return generics offer a method, ReturnOptional.

```go
var age uint64
extractor.WithOptional("age", valueextractor.AsUint64(&age))
```

In this example, if the "age" parameter is missing from the data source, `WithOptional` will not add an `ErrNotFound` to the error chain, allowing the application to proceed without treating the absence of "age" as an error.

By leveraging `WithOptional`, you can write cleaner, more concise code for handling optional parameters without cluttering your error handling logic with checks for missing but non-essential data.



#### Using Direct Return Functions

Direct return functions like `ReturnString` and `ReturnUint64` offer a more straightforward way to extract and convert values, especially when used within conditional statements or benchmarks.

```go
package main

import (
    "fmt"
    "github.com/Howard3/valueextractor"
)

func main() {
    ex := valueextractor.NewExtractor(...) // Assume ex is properly initialized

    if *valueextractor.ReturnString(ex, "name") != "John" {
        fmt.Println("Name not parsed correctly")
    }

    if *valueextractor.ReturnUint64(ex, "age") != 30 {
        fmt.Println("Age not parsed correctly")
    }
}
```

#### Using Return Generics
Return generics are easy to work with but have approximately the same performance as the struct+reflection approach.
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


## Extensibility

The system is designed with extensibility in mind. You can extend it by implementing custom `ValueExtractor` interfaces or `Converter` functions. See the documentation for details on creating custom extractors and converters.
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

Errors are collected throughout the extraction and conversion process. Use the `Errors` method to retrieve any accumulated errors. The package defines `ErrNotFound` for missing keys and provides detailed error types for extract and convert errors, allowing precise error handling and diagnostics.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


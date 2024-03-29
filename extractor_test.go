package valueextractor

import (
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

type Bench1 struct {
	Name string `query:"name"`
	Age  uint64 `query:"age"`
}

func TestOptionalKeys(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8080?name=John", nil)

	var name string
	var age uint64

	ex := Using(QueryExtractor{Query: req.URL.Query()}, WithOptionalKeys("age"))
	ex.With("name", AsString(&name))
	ex.With("age", AsUint64(&age))
	err := ex.Errors()

	switch {
	case err != nil:
		t.Fatal(err)
	case ex.optionalKeys == nil:
		t.Fatal("Optional keys not set")
	}
}

func TestResultGeneric(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8080?name=John&age=30", nil)

	ex := Using(QueryExtractor{Query: req.URL.Query()})
	name := Result(ex, "name", AsString)
	age := Result(ex, "age", AsUint64)
	err := ex.Errors()

	switch {
	case err != nil:
		t.Fatal(err)
	case name != "John":
		t.Fatal("Name not parsed correctly")
	case age != 30:
		t.Fatal("Age not parsed correctly")
	}
}

func BenchmarkParamsParser(b *testing.B) {
	// construct a request with sample query data
	req, _ := http.NewRequest("GET", "http://localhost:8080?name=John&age=30", nil)

	for i := 0; i < b.N; i++ {
		res := Bench1{}

		ex := Using(QueryExtractor{Query: req.URL.Query()})
		ex.With("name", AsString(&res.Name))
		ex.With("age", AsUint64(&res.Age))
		err := ex.Errors()

		switch {
		case err != nil:
			b.Fatal(err)
		case res.Name != "John":
			b.Fatal("Name not parsed correctly")
		case res.Age != 30:
			b.Fatal("Age not parsed correctly")
		}
	}
}

func BenchmarkTypeReturn(b *testing.B) {
	req, _ := http.NewRequest("GET", "http://localhost:8080?name=John&age=30", nil)

	for i := 0; i < b.N; i++ {
		ex := Using(QueryExtractor{Query: req.URL.Query()})

		if *ReturnString(ex, "name") != "John" {
			b.Fatal("Name not parsed correctly")
		}

		if *ReturnUint64(ex, "age") != 30 {
			b.Fatal("Age not parsed correctly")
		}

		if err := ex.Errors(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParamsParserResultsValue(b *testing.B) {
	// construct a request with sample query data
	req, _ := http.NewRequest("GET", "http://localhost:8080?name=John&age=30", nil)

	for i := 0; i < b.N; i++ {
		ex := Using(QueryExtractor{Query: req.URL.Query()})

		if Result(ex, "name", AsString) != "John" {
			b.Fatal("Name not parsed correctly")
		}

		if Result(ex, "age", AsUint64) != 30 {
			b.Fatal("Age not parsed correctly")
		}

		if err := ex.Errors(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParamsParserResultsPtr(b *testing.B) {
	// construct a request with sample query data
	req, _ := http.NewRequest("GET", "http://localhost:8080?name=John&age=30", nil)

	for i := 0; i < b.N; i++ {
		ex := Using(QueryExtractor{Query: req.URL.Query()})

		if *ResultPtr(ex, "name", AsString) != "John" {
			b.Fatal("Name not parsed correctly")
		}

		if *ResultPtr(ex, "age", AsUint64) != 30 {
			b.Fatal("Age not parsed correctly")
		}

		if err := ex.Errors(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParamsParserNoStruct(b *testing.B) {
	req, _ := http.NewRequest("GET", "http://localhost:8080?name=John&age=30", nil)

	for i := 0; i < b.N; i++ {
		name := ""
		age := uint64(0)

		ex := Using(QueryExtractor{Query: req.URL.Query()})
		ex.With("name", AsString(&name))
		ex.With("age", AsUint64(&age))
		err := ex.Errors()

		switch {
		case err != nil:
			b.Fatal(err)
		case name != "John":
			b.Fatal("Name not parsed correctly")
		case age != 30:
			b.Fatal("Age not parsed correctly")
		}
	}
}

func BenchmarkWithReflection(b *testing.B) {
	// construct a request with sample query data
	req, _ := http.NewRequest("GET", "http://localhost:8080?name=John&age=30", nil)

	for i := 0; i < b.N; i++ {
		bench1 := Bench1{}
		v := reflect.ValueOf(&bench1).Elem()
		t := v.Type()
		query := req.URL.Query()

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			tag := t.Field(i).Tag.Get("query")
			value := query.Get(tag)

			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Int:
				intVal, _ := strconv.Atoi(value)
				field.SetUint(uint64(intVal))
			}
		}

	}
}

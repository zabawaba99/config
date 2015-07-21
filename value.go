package config

import (
	"log"
	"reflect"
)

type value struct {
	// UserEnv is a flag that to determine whether or not to use
	// Env over Flag when resolving the value
	UseEnv bool
	// Flag is the pointer value retrieved and parsed from `flag.<Type>()`
	Flag interface{}
	// Env is the value retrieved and parsed from `os.Getenv`
	Env interface{}
	// Default is the value that is returned if both Flag and Env are empty
	Default interface{}
	// Type is the string representation of the type
	// of value that `Flag` and `Env` are
	// 		i.e: `uint`, `int`, `float64`, `string`, etc
	Type string
}

// resolve returns the appropriate `interface{}` value
func (v value) resolve() interface{} {
	if v.UseEnv {
		return v.Env
	}

	// check to see if the flag was even set
	val := reflect.ValueOf(v.Flag)
	if !val.IsValid() {
		return v.Default
	}

	// check to see if the flag was set but has an
	// empty default
	val = val.Elem()
	if isZeroValue(val.Interface()) {
		return v.Default
	}
	return val.Interface()
}

func isZeroValue(v interface{}) bool {
	switch v.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return v.(uint64) == 0
	case float32, float64:
		return v.(float64) == 0
	case string:
		return v.(string) == ""
	case nil:
		return true
	default:
		log.Printf("que? %T\n", v)
	}

	return false
}

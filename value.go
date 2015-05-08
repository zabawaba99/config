package config

import "reflect"

type value struct {
	// UserEnv is a flag that to determine whether or not to use
	// Env over Flag when resolving the value
	UseEnv bool
	// Flag is the pointer value retrieved and parsed from `flag.<Type>()`
	Flag interface{}
	// Env is the value retrieved and parsed from `os.Getenv`
	Env interface{}
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

	// get the value of the flag pointer
	return reflect.ValueOf(v.Flag).Elem().Interface()
}

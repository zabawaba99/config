package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

var errNoEnvVar = errors.New("no env var")
var values = map[string]value{}

func init() {
	var filename string
	var parser func(string) (map[string]argument, error)
	if _, err := os.Stat("config.json"); err == nil {
		filename, parser = "config.json", parseJSON
	}

	// TODO: parse yaml

	if filename == "" {
		return
	}

	config, err := parser(filename)
	if err != nil {
		println("err", err.Error())
		// do something
		return
	}

	load(config)
}

func load(config map[string]argument) {
	for name, c := range config {
		var (
			f      interface{}
			e      interface{}
			useEnv bool
		)

		var err error
		if c.FlagName != "" {
			f, err = loadFlag(c)
			if err != nil {
				println(err.Error())
				continue
			}
		}

		if c.EnvName != "" {
			e, err = loadEnv(c)
			useEnv = err != errNoEnvVar
			if useEnv && err != nil {
				println(err.Error())
				continue
			}
		}

		if c.Type == "" {
			c.Type = "string"
		}
		values[name] = value{Flag: f, Env: e, Type: c.Type, Fallback: c.Default, UseEnv: useEnv}
	}
	flag.Parse()
}

func loadFlag(a argument) (interface{}, error) {
	var v interface{}
	switch a.Type {
	case "uint", "uint8", "uint16", "uint32", "uint64":
		ii, ok := a.Default.(float64)
		if !ok {
			return nil, errors.New("invalid")
		}

		v = flag.Uint64(a.FlagName, uint64(ii), a.Description)
	default:
		v = flag.String(a.FlagName, fmt.Sprintf("%s", a.Default), a.Description)
	}
	return v, nil
}

func loadEnv(a argument) (interface{}, error) {
	envVal := os.Getenv(a.EnvName)
	if envVal == "" {
		return nil, errNoEnvVar
	}

	var v interface{}
	switch a.Type {
	case "uint", "uint8", "uint16", "uint32", "uint64":
		i, err := strconv.ParseUint(envVal, 10, 64)
		if err != nil {
			return nil, err
		}
		v = uint64(i)
	default:
		v = envVal
	}
	return v, nil
}

func Load(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if k := val.Kind(); k != reflect.Struct {
		return fmt.Errorf("can only load config into a struct not: %v", k)
	}

	typ := val.Type()
	// loop through the struct's fields and set the map
	for i := 0; i < val.NumField(); i++ {
		fv := val.Field(i)
		if !fv.CanSet() {
			continue
		}

		f := typ.Field(i)
		configName := f.Tag.Get("config")
		if configName == "" {
			configName = f.Name
		}

		c, ok := values[configName]
		if !ok {
			// not in config file
			continue
		}

		// TODO: validate that the field type
		// is the same as the type we are setting
		newVal := c.resolve()
		switch c.Type {
		case "uint", "uint8", "uint16", "uint32", "uint64":
			fv.SetUint(newVal.(uint64))
		default:
			fv.SetString(newVal.(string))
		}
	}
	return nil
}

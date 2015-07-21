package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

var errNoEnvVar = errors.New("no env var")
var errNotPtr = errors.New("cannot only load onto pointer")
var values = map[string]value{}

func init() {
	var filename string
	var parser func(string) (map[string]argument, error)
	if _, err := os.Stat("config.json"); err == nil {
		filename, parser = "config.json", parseJSON
	}

	// TODO: parse yaml

	if filename == "" {
		log.Fatal("missing config file")
	}

	config, err := parser(filename)
	if err != nil {
		log.Fatal("Could not parse config file", err)
	}

	errs := load(config)
	if len(errs) != 0 {
		printUsageAndExit(config, errs)
	}
}

func printUsageAndExit(config map[string]argument, errs []argError) {
	logger := log.New(os.Stderr, "", 0)
	logger.Print("Cannot start application:")
	for _, err := range errs {
		logger.Printf("  %s", err)
	}
	logger.Print("\nAttempted to use the following values:")
	for k, v := range config {
		var flagVal interface{} = ""
		envVal := os.Getenv(v.EnvName)

		if val, ok := values[k]; ok {
			flagVal = reflect.ValueOf(val.Flag).Elem().Interface()
		}
		logger.Printf("-\n  Flag:\t%s: \"%v\"", v.FlagName, flagVal)
		logger.Printf("  Env:\t%s : %q", v.EnvName, envVal)
	}
	logger.Println()
	logger.Println("Usage")
	for k, v := range config {
		logger.Printf("  %s:\t\t%s\n", k, v.Description)
	}
	logger.Println()
	logger.Fatal("Exiting...")
}

func load(config map[string]argument) (errs []argError) {
	// try and load the flag/env for every config in the json file
	for name, c := range config {
		v := value{Type: c.Type, Default: c.Default}

		var err error
		if c.FlagName != "" {
			v.Flag, err = loadFlag(c)
			if err != nil {
				println(err.Error())
				continue
			}
		}

		if c.EnvName != "" {
			v.Env, err = loadEnv(c)
			v.UseEnv = err != errNoEnvVar
			if v.UseEnv && err != nil {
				println(err.Error())
				continue
			}
		}
		values[name] = v
	}
	flag.Parse()

	// check if any of the required configurations are missing
	for name, arg := range config {
		if !arg.Required {
			continue
		}

		val, ok := values[name]
		if !ok {
			errs = append(errs, argError{couldNotParse, arg})
			continue
		}

		if isZeroValue(val.resolve()) {
			errs = append(errs, argError{argMissing, arg})
			continue
		}
	}
	return errs
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
		def := fmt.Sprintf("%s", a.Default)
		if isZeroValue(a.Default) {
			def = ""
		}
		v = flag.String(a.FlagName, def, a.Description)
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

// Load marshals the loaded configuration into the given interface.
// If the interface is not a struct pointer, and error will be returned.
func Load(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return errNotPtr
	}
	val = val.Elem()

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

		// TODO: make sure we can actually set the newVal
		// onto the struct field
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

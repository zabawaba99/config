package config

import "fmt"

type argErrType int

const (
	couldNotParse argErrType = iota
	argMissing
)

type argError struct {
	Error    argErrType
	Argument argument
}

func (ce argError) String() string {
	var errType string
	switch ce.Error {
	case couldNotParse:
		errType = "could not be parsed"
	case argMissing:
		errType = "was not specified"
	}

	return fmt.Sprintf("Argument %s/%s %s", ce.Argument.FlagName, ce.Argument.EnvName, errType)
}

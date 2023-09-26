// params.go
package config

import (
	"flag"
	"os"
)

// Params holds the parsed command-line parameters.
type Params struct {
	ProtoFilePath string
	Debug         bool
	Uglify        bool
}

// ParseParams parses command-line parameters and environment variables.
func ParseParams() Params {
	var params Params

	// Parse command-line parameters
	flag.StringVar(&params.ProtoFilePath, "proto", "", "Path to the protobuf schema definition file")
	flag.BoolVar(&params.Debug, "debug", false, "Enable debug mode to show detailed field information")
	flag.BoolVar(&params.Uglify, "uglify", false, "Uglify JSON output everywhere")
	flag.Parse()

	// Check environment variables for debug and uglify flags
	if debugEnv := os.Getenv("DEBUG"); debugEnv == "true" {
		params.Debug = true
	}

	if uglifyEnv := os.Getenv("UGLIFY"); uglifyEnv == "true" {
		params.Uglify = true
	}

	return params
}

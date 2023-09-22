package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"gopkg.in/yaml.v2"
)

type FieldOverrideConfig struct {
	Service string            `yaml:"service"`
	Method  string            `yaml:"method"`
	Fields  map[string]string `yaml:"fields"`
}

var configs []FieldOverrideConfig

func readConfigs() {
	data, err := ioutil.ReadFile("overrides.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %s\n", err)
	}

	var configWrapper struct {
		Overrides []FieldOverrideConfig `yaml:"overrides"`
	}
	if err := yaml.Unmarshal(data, &configWrapper); err != nil {
		log.Fatalf("Error parsing YAML file: %s\n", err)
	}

	configs = configWrapper.Overrides
}

func generateExampleValue(field *desc.FieldDescriptor) interface{} {
	if field.IsRepeated() {
		// Handle repeated fields by generating a list of example values
		return []interface{}{generateSingleFieldValue(field)}
	} else {
		return generateSingleFieldValue(field)
	}
}

func generateSingleFieldValue(field *desc.FieldDescriptor) interface{} {
	if field.GetMessageType() != nil {
		// For google.protobuf.Timestamp special handling
		if field.GetMessageType().GetFullyQualifiedName() == "google.protobuf.Timestamp" {
			return time.Now().Format(time.RFC3339Nano)
		}

		// Handle nested messages
		nestedExample := make(map[string]interface{})
		for _, nestedField := range field.GetMessageType().GetFields() {
			nestedValue := generateExampleValue(nestedField)
			nestedExample[nestedField.GetName()] = nestedValue
		}
		return nestedExample
	}

	// Handle basic types, for example:
	switch field.GetType().String() {
	case "TYPE_DOUBLE":
		return 1.7976931348623157e+308
	case "TYPE_FLOAT":
		return float32(3.402823466e+38)
	case "TYPE_INT64":
		return int64(9223372036854775807)
	case "TYPE_UINT64":
		return "18446744073709551615"
	case "TYPE_INT32":
		return "2147483647"
	case "TYPE_FIXED64":
		return "18446744073709551615"
	case "TYPE_FIXED32":
		return "4294967295"
	case "TYPE_BOOL":
		return "true"
	case "TYPE_STRING":
		return "abcdefghijklmnopqrstuvwxyzABCD" // 32 bytes
	case "TYPE_GROUP", "TYPE_MESSAGE", "TYPE_BYTES":
		return "NOT_SUPPORTED"
	case "TYPE_UINT32":
		return "4294967295"
	case "TYPE_ENUM":
		return "ENUM_VALUE_MAX"
	case "TYPE_SFIXED32", "TYPE_SFIXED64", "TYPE_SINT32", "TYPE_SINT64":
		return "NOT_SUPPORTED"
	default:
		return "UNKNOWN_TYPE"
	}
}

func customValueGenerator(service, method, fieldName string) (interface{}, bool) {
	// Check for universal rules first
	for _, config := range configs {
		if config.Service == "*" && config.Method == "*" {
			for key, customValue := range config.Fields {
				matched, _ := regexp.MatchString(key, fieldName)
				if matched {
					if customValue == "uuid" {
						return uuid.New().String(), true
					} else {
						// Try to detect if the custom value is an integer
						if intValue, err := strconv.Atoi(customValue); err == nil {
							return intValue, true
						}
						// Add more type detections here if needed
						return customValue, true
					}
				}
			}
		}
	}

	// Then check for service/method specific rules
	for _, config := range configs {
		if config.Service == service && config.Method == method {
			for key, customValue := range config.Fields {
				matched, _ := regexp.MatchString(key, fieldName)
				if matched {
					if customValue == "uuid" {
						return uuid.New().String(), true
					} else {
						// Try to detect if the custom value is an integer
						if intValue, err := strconv.Atoi(customValue); err == nil {
							return intValue, true
						}
						// Add more type detections here if needed
						return customValue, true
					}
				}
			}
		}
	}

	return nil, false
}
func generateFields(service, method string, fields []*desc.FieldDescriptor, debug bool) map[string]interface{} {
	example := make(map[string]interface{})
	for _, field := range fields {
		exampleValue := generateExampleValue(field)
		if debug {
			fmt.Printf("Field: %s (%s): %v\n", field.GetName(), field.GetType().String(), exampleValue)
		}

		// Use customValueGenerator for specific overrides
		if customValue, ok := customValueGenerator(service, method, field.GetName()); ok {
			exampleValue = customValue
			if debug {
				fmt.Printf("Field: %s (%s): %v (custom)\n", field.GetName(), field.GetType().String(), customValue)
			}
		}

		example[field.GetName()] = exampleValue
	}
	return example
}

func main() {
	readConfigs()
	var protoFilePath string
	var debug bool

	flag.StringVar(&protoFilePath, "proto", "", "Path to the protobuf schema definition file")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode to show detailed field information")
	flag.Parse()

	if debugEnv := os.Getenv("DEBUG"); debugEnv == "true" {
		debug = true
	}

	parser := protoparse.Parser{}
	fds, err := parser.ParseFiles(protoFilePath)
	if err != nil {
		fmt.Printf("Failed to parse .proto file: %v\n", err)
		os.Exit(1)
	}

	lookupTable := make(map[string]string)
	serviceIdx := 1
	for _, fd := range fds {
		for _, sd := range fd.GetServices() {
			fmt.Printf("%d. Service: %s\n", serviceIdx, sd.GetName())
			for mIdx, md := range sd.GetMethods() {
				fmt.Printf("      %d/%d. %s/%s\n", serviceIdx, mIdx+1, sd.GetName(), md.GetName())
				lookupTable[fmt.Sprintf("%d/%d", serviceIdx, mIdx+1)] = sd.GetName() + "/" + md.GetName()
			}
			serviceIdx++
		}
	}

	var input string
	fmt.Println("Select a method to generate examples for (e.g., 1/1 or Service/Method):")
	fmt.Scanln(&input)

	var service, method string
	parts := strings.Split(input, "/")
	if len(parts) == 2 {
		if _, err := strconv.Atoi(parts[0]); err == nil {
			if translated, ok := lookupTable[input]; ok {
				parts = strings.Split(translated, "/")
			}
		}
		service = parts[0]
		method = parts[1]
	} else {
		fmt.Println("FAIL: Invalid input format.")
		os.Exit(1)
	}

	found := false
	for _, fd := range fds {
		for _, sd := range fd.GetServices() {
			if sd.GetName() == service {
				for _, md := range sd.GetMethods() {
					if md.GetName() == method {
						found = true

						if debug {
							fmt.Println("Generating request fields:")
						}
						requestExample := generateFields(service, method, md.GetInputType().GetFields(), debug)

						if debug {
							fmt.Println("Generating response fields:")
						}
						responseExample := generateFields(service, method, md.GetOutputType().GetFields(), debug)

						requestJSON, err := json.Marshal(requestExample)
						if err != nil {
							fmt.Println("Failed to generate request JSON:", err)
							os.Exit(1)
						}

						responseJSON, err := json.Marshal(responseExample)
						if err != nil {
							fmt.Println("Failed to generate response JSON:", err)
							os.Exit(1)
						}

						fmt.Println("Request example:")
						fmt.Println(string(requestJSON))
						fmt.Println("gRPCurl call example:")
						fmt.Printf("grpcurl -d '%s' -plaintext HOST:PORT %s/%s\n", string(requestJSON), service, method)
						fmt.Println("Response example:")
						fmt.Println(string(responseJSON))
					}
				}
			}
		}
	}

	if !found {
		fmt.Println("FAIL: Service or Method does not exist.")
		os.Exit(1)
	}
}

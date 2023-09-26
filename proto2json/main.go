package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"

	"github.com/mb/proto2json/config"
)

var configs []config.FieldOverrideConfig

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
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Date generator
	dateGenerator := func() string {
		now := time.Now()
		return now.Format("2023-01-02")
	}

	// Time generator
	timeGenerator := func() string {
		now := time.Now()
		return now.Format("15:04:05.999")
	}

	// DateTime generator
	dateTimeGenerator := func() string {
		now := time.Now()
		return now.Format("2023-01-02T15:04:05.999Z07:00")
	}

	// Phone number generator
	phoneNumberGenerator := func() string {
		return "+44 (0) 123 456 7890"
	}

	// Country name generator
	countryNameGenerator := func() string {
		countries := []string{"France", "United Kingdom", "Germany", "Spain"}
		randomIndex := random.Intn(len(countries))
		return countries[randomIndex]
	}

	// Country code generator
	countryCodeGenerator := func() string {
		countryCodes := []string{"FR", "UK", "DE", "ES"}
		randomIndex := random.Intn(len(countryCodes))
		return countryCodes[randomIndex]
	}

	// City generator
	cityGenerator := func() string {
		cities := []string{"Paris", "London", "Berlin", "Madrid"}
		randomIndex := random.Intn(len(cities))
		return cities[randomIndex]
	}

	// First name generator
	firstNameGenerator := func() string {
		firstNames := []string{"John", "Jane", "Alice", "Bob", "Eve"}
		randomIndex := random.Intn(len(firstNames))
		return firstNames[randomIndex]
	}

	// Last name generator
	lastNameGenerator := func() string {
		lastNames := []string{"Smith", "Doe", "Johnson", "Brown", "Wilson"}
		randomIndex := random.Intn(len(lastNames))
		return lastNames[randomIndex]
	}

	// Email generator
	emailGenerator := func() string {
		firstName := firstNameGenerator()
		lastName := lastNameGenerator()
		return fmt.Sprintf("%s.%s@example.com", strings.ToLower(firstName), strings.ToLower(lastName))
	}

	// imo generator
	imoGenerator := func() int {
		min := 9000000
		max := 9999999
		return min + random.Intn(max-min+1)
	}

	// uuid generator
	uuidGenerator := func() string {
		return uuid.New().String()
	}

	// Check for universal rules and service/method specific rules
	for _, config := range configs {
		if (config.Service == "*" && config.Method == "*") || (config.Service == service && config.Method == method) {
			for key, customValue := range config.Fields {
				matched, _ := regexp.MatchString(key, fieldName)
				if matched {
					switch customValue {
					case "uuid":
						return uuidGenerator(), true
					case "imo":
						return imoGenerator(), true
					case "first_name":
						return firstNameGenerator(), true
					case "last_name":
						return lastNameGenerator(), true
					case "iso_date":
						return dateGenerator(), true
					case "iso_time":
						return timeGenerator(), true
					case "iso_datetime":
						return dateTimeGenerator(), true
					case "email":
						return emailGenerator(), true
					case "phone":
						return phoneNumberGenerator(), true
					case "country_name":
						return countryNameGenerator(), true
					case "country_code":
						return countryCodeGenerator(), true
					case "city":
						return cityGenerator(), true
					default:
						// Add more custom values here as needed
						return customValue, true
					}
				}
			}
		}
	}

	return nil, false
}

func generateFields(service, method string, fields []*desc.FieldDescriptor, debug bool) interface{} {
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

		// Recursively apply custom value replacements to nested fields
		if field.GetMessageType() != nil {
			nestedFields := field.GetMessageType().GetFields()
			nestedExample := generateFields(service, method, nestedFields, debug)
			if _, isMap := exampleValue.(map[string]interface{}); isMap {
				for key, value := range nestedExample.(map[string]interface{}) {
					exampleValue.(map[string]interface{})[key] = value
				}
			} else if _, isSlice := exampleValue.([]interface{}); isSlice {
				// Handle repeated fields with nested messages
				// This assumes that if it's a slice, it's a repeated message field
				// and expects a list of maps.
				for _, item := range exampleValue.([]interface{}) {
					if itemMap, isItemMap := item.(map[string]interface{}); isItemMap {
						for key, value := range nestedExample.(map[string]interface{}) {
							itemMap[key] = value
						}
					}
				}
			}
		}
	}
	return example
}

func beautifyJSON(jsonData []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, jsonData, "", "    "); err != nil {
		return string(jsonData) // Return the original JSON if beautifying fails
	}
	return prettyJSON.String()
}

func printBlock(title, content string) {
	fmt.Printf("%s:\n%s\n\n", title, content)
}

func main() {
	var err error
	// todo: allow overrides.yaml to be the default, but configurable.
	configs, err = config.ReadConfigs("overrides.yaml")
	if err != nil {
		log.Fatalf("Error reading or parsing YAML file: %s\n", err)
	}

	// params
	var protoFilePath string
	flag.StringVar(&protoFilePath, "proto", "", "Path to the protobuf schema definition file")
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Enable debug mode to show detailed field information")
	var uglify bool
	flag.BoolVar(&uglify, "uglify", false, "Uglify JSON output everywhere")

	// parse params
	flag.Parse()

	if debugEnv := os.Getenv("DEBUG"); debugEnv == "true" {
		debug = true
	}

	if uglifyEnv := os.Getenv("UGLIFY"); uglifyEnv == "true" {
		uglify = true
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

						// Optionally beautify or uglify JSON based on the flag
						var requestOutput, responseOutput string
						if uglify {
							requestOutput = string(requestJSON)
							responseOutput = string(responseJSON)
						} else {
							requestOutput = beautifyJSON(requestJSON)
							responseOutput = beautifyJSON(responseJSON)
						}

						// Escape double quotes for the grpcurl command
						escapedRequestJSON := strings.ReplaceAll(requestOutput, `"`, `\"`)
						// Escape the opening curly brace to prevent a new line
						escapedRequestJSON = strings.ReplaceAll(escapedRequestJSON, "{", "{\\")

						// Format the grpcurl command with proper line breaks and indentation
						grpcurlCommand := fmt.Sprintf("grpcurl -d \"'%s'\" -H \"Authorization: Bearer ${TOKEN}\" -plaintext ${HOST}:${PORT} ${API_PROTO_SERVICE_VERSION}.%s/%s", escapedRequestJSON, service, method)

						printBlock("Request example", requestOutput)
						printBlock("gRPCurl call example", grpcurlCommand)
						printBlock("Response example", responseOutput)
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

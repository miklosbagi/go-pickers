package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jhump/protoreflect/desc/protoparse"
)

func generateExampleValue(fieldType string) string {
	switch fieldType {
	case "TYPE_DOUBLE":
		return "1.7976931348623157E+308"
	case "TYPE_FLOAT":
		return "3.402823466E+38"
	case "TYPE_INT64":
		return "9223372036854775807"
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

func customValueGenerator(fieldName string) (string, bool) {
	if matched, _ := regexp.MatchString("[a-z]+_id$", fieldName); matched {
		return uuid.New().String(), true
	}
	return "", false
}

func main() {
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

	requestExample := make(map[string]interface{})
	responseExample := make(map[string]interface{})
	found := false

	for _, fd := range fds {
		for _, sd := range fd.GetServices() {
			if sd.GetName() == service {
				for _, md := range sd.GetMethods() {
					if md.GetName() == method {
						found = true

						for fIdx, field := range md.GetInputType().GetFields() {
							var exampleValue interface{}
							if stringValue, ok := customValueGenerator(field.GetName()); ok {
								exampleValue = stringValue
							} else {
								exampleValue = generateExampleValue(field.GetType().String())
							}

							if debug {
								fmt.Printf("Request Field %d: %s (%s): %v\n", fIdx+1, field.GetName(), field.GetType().String(), exampleValue)
							}

							requestExample[field.GetName()] = exampleValue
						}

						for fIdx, field := range md.GetOutputType().GetFields() {
							var exampleValue interface{}
							if stringValue, ok := customValueGenerator(field.GetName()); ok {
								exampleValue = stringValue
							} else {
								exampleValue = generateExampleValue(field.GetType().String())
							}

							if debug {
								fmt.Printf("Response Field %d: %s (%s): %v\n", fIdx+1, field.GetName(), field.GetType().String(), exampleValue)
							}

							responseExample[field.GetName()] = exampleValue
						}
					}
				}
			}
		}
	}

	if !found {
		fmt.Println("FAIL: Service or Method does not exist.")
		os.Exit(1)
	}

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

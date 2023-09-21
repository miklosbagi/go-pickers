# proto2json
This is go code to read up a defined proto file, and attempt to put out example JSON requests and responses. The point of this small tool is to save time on manual analysis, and have a copy-paste solution at hand for any proto you throw at this, for testing purposes.

### Usage
```
go mod tidy; go run main.go --proto <path-to-proto-file>
```
There is a `--debug` optional flag to spit out more info in case something go sideways.

### Example run
```
go mod tidy; go build; ./proto2json --proto ~/Documents/test.proto
1. Service: MyService
      1/1. MyService/List
      1/2. MyService/Add
      1/3. MyService/Remove
2. Service: YourService
      2/1. YourService/List
      2/2. YourService/Add
      2/3. YourService/Remove

Select a method to generate examples for (e.g., 1/1 or Service/Method):
1/2

Request example:
{"id":"2107f934-ac08-4fa8-8b09-e8bce4701df6","days":1,"add_id":"0fcf7143-292c-4992-9f31-6e2d2384166d"}

gRPCurl call example:
grpcurl -d '{"id":"2107f934-ac08-4fa8-8b09-e8bce4701df6","days":1,"add_id":"0fcf7143-292c-4992-9f31-6e2d2384166d"}' -plaintext HOST:PORT TokensService/Generate

Response example:
{"blah":"9b16e97d-5e28-4c10-bc2c-19fe04060d42","expires":"2023-09-21T21:08:46.082502+02:00","id":"998c53cb-7f62-48b3-9f4a-c69004fd75aa"}
```

### Configuration
Proto files can define a few types of fields, therefore proto2json uses detaults so avoid having to set some value to everything. These defaults are:
- DOUBLE: 1.7976931348623157e+308
- FLOAT: float32(3.402823466e+38)
- INT64: int64(9223372036854775807)
- UINT64: "18446744073709551615"
- INT32: "2147483647"
- FIXED64: "18446744073709551615"
- FIXED32: "4294967295"
- BOOL: "true"
- STRING": "abcdefghijklmnopqrstuvwxyzABCD"
- UINT32": "4294967295"
- ENUM": "ENUM_VALUE_MAX"

Unsupported (returning "NOT_SUPPORTED"):
- GROUP: "NOT_SUPPORTED"
- MESSAGE: "NOT_SUPPORTED"
- BYTES: "NOT_SUPPORTED"
- SFIXED32: "NOT_SUPPORTED"
- SFIXED64: "NOT_SUPPORTED"
- SINT32: "NOT_SUPPORTED"
- SINT64: "NOT_SUPPORTED"
		return "NOT_SUPPORTED"

In all other cases:
- default: "UNKNOWN_TYPE"

### Configuration overrides via YAML
The `overrides.yaml` file allows for the above default values to be overridden with static output, or reference to supported code.

Example:
```
overrides:
  # overriding this field for all services and methods
  - service: "*"
    method: "*"
    fields:
      # uuid is replaced by uuid code here automatically: cc325791-84ef-4269-b492-8515e5a88520
      ".*_id": "uuid"  # This will match any field ending with "_id"
      # uuid is replaced by uuid code here automatically: cc325791-84ef-4269-b492-8515e5a88520
      id: uuid
      # integer value
      expire_days: 1

  # overriding a specific service
  - service: TokensService
    method: Generate
    fields:
      code: uuid
```

Code ovverides:
- all "uuid" values are replaced by generated uuid, e.g: cc325791-84ef-4269-b492-8515e5a88520
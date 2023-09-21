# proto2json
This is go code to read up a defined proto file, and attempt to put out example JSON requests and responses. The point of this small tool is to save time on manual analysis, and have a copy-paste solution at hand for any proto you throw at this, for testing purposes.

### Usage
```
go mod tidy; go run main.go --proto <path-to-proto-file>
```
There is a `--debug` optional flag to spit out more info in case something go sideways.

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
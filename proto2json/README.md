# proto2json
This is go code to read up a defined proto file, and attempt to put out example JSON requests and responses. The point of this small tool is to save time on manual analysis, and have a copy-paste solution at hand for any proto you throw at this, for testing purposes.
This provides a similar output like postman, when example data fills in a request, with the following differences:

- gRPC only (no REST)
- Examples aim to be better than the default
- Creates exmaples for both request and response, including grpcurl call.
- Not upselling a subscription package


### Usage
```
go mod tidy; go run main.go --proto <path-to-proto-file>
```

Params:
- `--proto`: full / relative path to proto file

Flags:
- `--debug`: print out more info
- `--uglify`: single line json output for all examples

### Example run
We are using [nested proto test data](./test/test-data/nested proto) for this.

```
go mod tidy; go build; ./proto2json --proto ./test/test-data/nested proto
1. Service: OrderService
      1/1. OrderService/CreateOrder
      1/2. OrderService/GetOrderInfo
Select a method to generate examples for (e.g., 1/1 or Service/Method):
```

```
Request example:
{
    "customer_info": {
        "address": {
            "city": "abcdefghijklmnopqrstuvwxyzABCD",
            "contacts": [
                {
                    "email": "abcdefghijklmnopqrstuvwxyzABCD",
                    "phone": {
                        "country_code": "abcdefghijklmnopqrstuvwxyzABCD",
                        "number": "abcdefghijklmnopqrstuvwxyzABCD"
                    }
                }
            ],
            "postal_code": "abcdefghijklmnopqrstuvwxyzABCD",
            "state": "abcdefghijklmnopqrstuvwxyzABCD",
            "street": "abcdefghijklmnopqrstuvwxyzABCD"
        },
        "first_name": "abcdefghijklmnopqrstuvwxyzABCD",
        "last_name": "abcdefghijklmnopqrstuvwxyzABCD"
    },
    "items": [
        {
            "item_id": "abcdefghijklmnopqrstuvwxyzABCD",
            "name": "abcdefghijklmnopqrstuvwxyzABCD",
            "price": 1.7976931348623157e+308,
            "quantity": "2147483647"
        }
    ],
    "order_id": "799e6325-9697-4fd9-826b-0d0341101a81"
}
```

```
gRPCurl call example:
grpcurl -d "'{\
    \"customer_info\": {\
        \"address\": {\
            \"city\": \"abcdefghijklmnopqrstuvwxyzABCD\",
            \"contacts\": [
                {\
                    \"email\": \"abcdefghijklmnopqrstuvwxyzABCD\",
                    \"phone\": {\
                        \"country_code\": \"abcdefghijklmnopqrstuvwxyzABCD\",
                        \"number\": \"abcdefghijklmnopqrstuvwxyzABCD\"
                    }
                }
            ],
            \"postal_code\": \"abcdefghijklmnopqrstuvwxyzABCD\",
            \"state\": \"abcdefghijklmnopqrstuvwxyzABCD\",
            \"street\": \"abcdefghijklmnopqrstuvwxyzABCD\"
        },
        \"first_name\": \"abcdefghijklmnopqrstuvwxyzABCD\",
        \"last_name\": \"abcdefghijklmnopqrstuvwxyzABCD\"
    },
    \"items\": [
        {\
            \"item_id\": \"abcdefghijklmnopqrstuvwxyzABCD\",
            \"name\": \"abcdefghijklmnopqrstuvwxyzABCD\",
            \"price\": 1.7976931348623157e+308,
            \"quantity\": \"2147483647\"
        }
    ],
    \"order_id\": \"799e6325-9697-4fd9-826b-0d0341101a81\"
}'" -H "Authorization: Bearer ${TOKEN}" -plaintext ${HOST}:${PORT} ${API_PROTO_SERVICE_VERSION}.OrderService/CreateOrder
```

```
Response example:
{
    "customer_info": {
        "address": {
            "city": "abcdefghijklmnopqrstuvwxyzABCD",
            "contacts": [
                {
                    "email": "abcdefghijklmnopqrstuvwxyzABCD",
                    "phone": {
                        "country_code": "abcdefghijklmnopqrstuvwxyzABCD",
                        "number": "abcdefghijklmnopqrstuvwxyzABCD"
                    }
                }
            ],
            "postal_code": "abcdefghijklmnopqrstuvwxyzABCD",
            "state": "abcdefghijklmnopqrstuvwxyzABCD",
            "street": "abcdefghijklmnopqrstuvwxyzABCD"
        },
        "first_name": "abcdefghijklmnopqrstuvwxyzABCD",
        "last_name": "abcdefghijklmnopqrstuvwxyzABCD"
    },
    "items": [
        {
            "item_id": "abcdefghijklmnopqrstuvwxyzABCD",
            "name": "abcdefghijklmnopqrstuvwxyzABCD",
            "price": 1.7976931348623157e+308,
            "quantity": "2147483647"
        }
    ],
    "order_id": "509fa97c-3c56-4ff0-a5a5-d946a13acddd"
}
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

# Known issues
- Uglify has a printout bug with `\{\\`.
- Proto files with `import` statements are only supported only relative to run path.
- item_id not replaced with uuid.
- `repeated type` gets only a single item (e.g: repeated string: "names" should generate {"name":"Steve","name":"George"} as example.

# TODO
- Generators
  - Firstname
  - Lastname
  - Email
  - Phone Number
  - Country name
  - Country code
  - City
  - FlagIso2Code
  - Dates
  - Times
- Tests
- Makefile, packaging, etc
- Feature
  - ENTER on selections secreen == All services in loop
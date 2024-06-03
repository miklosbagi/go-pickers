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

JSON Request example:
```json
{
    "customer_info": {
        "address": {
            "city": "Berlin",
            "contacts": [
                {
                    "email": "bob.brown@example.com",
                    "phone": "+44 (0) 123 456 7890"
                }
            ],
            "postal_code": "abcdefghijklmnopqrstuvwxyzABCD",
            "state": "abcdefghijklmnopqrstuvwxyzABCD",
            "street": "abcdefghijklmnopqrstuvwxyzABCD"
        },
        "first_name": "Alice",
        "last_name": "Brown"
    },
    "expire_days": "1",
    "first_name": "Eve",
    "imo": 9076580,
    "items": [
        {
            "item_id": "3c4b6ca8-ca06-4b91-b179-2cf1307f61b2",
            "name": "abcdefghijklmnopqrstuvwxyzABCD",
            "price": 1.7976931348623157e+308,
            "quantity": "2147483647"
        }
    ],
    "order_id": "6abb2568-55fe-49d0-948f-247f10d04b67"
}
```

YAML Request Example:
```yaml
customer_info:
  address:
    city: Berlin
    contacts:
      - email: bob.brown@example.com
        phone: +44 (0) 123 456 7890
    postal_code: abcdefghijklmnopqrstuvwxyzABCD
    state: abcdefghijklmnopqrstuvwxyzABCD
    street: abcdefghijklmnopqrstuvwxyzABCD
  first_name: Alice
  last_name: Brown
expire_days: "1"
first_name: Eve
imo: 9076580
items:
  - item_id: 3c4b6ca8-ca06-4b91-b179-2cf1307f61b2
    name: abcdefghijklmnopqrstuvwxyzABCD
    price: 1.7976931348623157e+308
    quantity: "2147483647"
order_id: 6abb2568-55fe-49d0-948f-247f10d04b67
```

```
gRPCurl call example:
grpcurl -d "'{\
    \"customer_info\": {\
        \"address\": {\
            \"city\": \"Berlin\",
            \"contacts\": [
                {\
                    \"email\": \"bob.brown@example.com\",
                    \"phone\": \"+44 (0) 123 456 7890\"
                }
            ],
            \"postal_code\": \"abcdefghijklmnopqrstuvwxyzABCD\",
            \"state\": \"abcdefghijklmnopqrstuvwxyzABCD\",
            \"street\": \"abcdefghijklmnopqrstuvwxyzABCD\"
        },
        \"first_name\": \"Alice\",
        \"last_name\": \"Brown\"
    },
    \"expire_days\": \"1\",
    \"first_name\": \"Eve\",
    \"imo\": 9076580,
    \"items\": [
        {\
            \"item_id\": \"3c4b6ca8-ca06-4b91-b179-2cf1307f61b2\",
            \"name\": \"abcdefghijklmnopqrstuvwxyzABCD\",
            \"price\": 1.7976931348623157e+308,
            \"quantity\": \"2147483647\"
        }
    ],
    \"order_id\": \"6abb2568-55fe-49d0-948f-247f10d04b67\"
}'" -H "Authorization: Bearer ${TOKEN}" -plaintext ${HOST}:${PORT} ${API_PROTO_SERVICE_VERSION}.OrderService/CreateOrder
```

JSON Response example:
```json
{
    "customer_info": {
        "address": {
            "city": "Berlin",
            "contacts": [
                {
                    "email": "jane.brown@example.com",
                    "phone": "+44 (0) 123 456 7890"
                }
            ],
            "postal_code": "abcdefghijklmnopqrstuvwxyzABCD",
            "state": "abcdefghijklmnopqrstuvwxyzABCD",
            "street": "abcdefghijklmnopqrstuvwxyzABCD"
        },
        "first_name": "Jane",
        "last_name": "Brown"
    },
    "expire_days": "1",
    "first_name": "Bob",
    "imo": 9285517,
    "items": [
        {
            "item_id": "554b9674-71d7-40c5-a40f-c615f0b210bf",
            "name": "abcdefghijklmnopqrstuvwxyzABCD",
            "price": 1.7976931348623157e+308,
            "quantity": "2147483647"
        }
    ],
    "order_id": "8899559d-01ee-4e1f-ad5a-f845ece48a3c"
}
```

### Configuration
Proto files can define a few types of fields, therefore proto2json uses detaults so avoid having to set some value to everything. These defaults are:
- DOUBLE: return 1.7976931348623157e+308
- FLOAT: float32(3.402823466e+38)
- INT64: int64(9223372036854775807)
- UINT64: "18446744073709551615"
- INT32: "2147483647"
- FIXED64": 18446744073709551615"
- FIXED32: "4294967295"
- BOOL: "true"
- STRING: "abcdefghijklmnopqrstuvwxyzABCD" // 32 bytes
- UINT32: "4294967295"
- ENUM: "ENUM_VALUE_MAX"

Unsupported types:
- GROUP, MESSAGE, BYTES, SFIXED32, SFIXED64 SINT32 SINT64: "UNSUPPORTED_TYPE"

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
      "^[a-z]+_id$": "uuid"  # This will match any field ending with "_id"
      id: uuid
      imo: imo
      IMONumber: imo
      first_name: first_name
      last_name: last_name
      expire_days: 1
      city: city
      country: country
      email: email
      phone: phone
  - service: TokensService
    method: Generate
    fields:
      code: uuid
```

Code ovverides:
- uuid: replaced by generated uuid, e.g: cc325791-84ef-4269-b492-8515e5a88520
- imo and imoNumber: replaced by generated imo number, e.g: 9682510
- firstName: replaced by generated first name, e.g: John
- lastName: replaced by generated last name, e.g: Doe
- date: replaced by generated date, e.g: 2021-01-01
- time: replaced by generated time, e.g: 12:00:00
- dateTime: replaced by generated date time, e.g: 2021-01-01T12:00:00Z
- email: replaced by generated email, e.g: firstname.lastname@example.com
- phoneNumber: replaced by generated phone number, e.g: +44 (0) 123 456 7890
- countryName: replaced by generated country name, e.g: Spain
- countryCode: replaced by generated country code, e.g: GR
- city: replaced by generated city, e.g: London


# Known issues
- Uglify has a printout bug with `\{\\`.
- Proto files with `import` statements are only supported only relative to run path.
- `repeated type` gets only a single item (e.g: repeated string: "names" should generate {"name":"Steve","name":"George"} as example.
- mapped values are not handled (e.g: map<string, string> = {"key": "value"}).
- Enum Type is wrong, need to parse enum BlahTypes for example to have a value 0-x, generate as integer.
- Some int32 fields are generated as "fieldName": { "value": 12345 }, despite they are a stright int32 field, with no map. This breaks the grpc response.
- SIUnit is not handled properly
-

# TODO
- String fields would require to include field name, for better identification for debugging. All abcdefgh.... won't help much.
- Tests
- Makefile, packaging, etc
- Feature
  - ENTER on selections secreen == All services in loop

# proto2json
This is go code to read up a defined proto file, and attempt to put out example JSON requests and responses. The point of this small tool is to save time on manual analysis, and have a copy-paste solution at hand for any proto you throw at this, for testing purposes.

### Usage
```
go mod tidy; go run main.go --proto <path-to-proto-file>
```
There is a `--debug` optional flag to spit out more info in case something go sideways.

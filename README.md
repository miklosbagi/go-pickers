# go-pickers
Some scripts to play around with go


## proto2json
This is go code to read up a defined proto file, and attempt to put out example JSON requests and responses. The point of this small tool is to save time on manual analysis, and have a copy-paste solution at hand for any proto you throw at this, for testing purposes.

### Usage
```
go mod tidy; go run main.go --proto <path-to-proto-file>
```
There is a `--debug` optional flag to spit out more info in case something go sideways.

## unit-test-extractor
`get-tests.sh` is a small tool that aims to analyze go code and put out unit tests, so we have a single view at the coverage (not counting %, this is just info).
To use this, include the directory in your $PATH, and navigate to the root of go code you need information about, and run `get-tests.sh` without any params.
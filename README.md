# go-pickers
Some scripts to play around with go


## proto2json
This tool helps generating JSON request and response examples from protobuf files. See more in [proto2json](proto2json/README.md)

## unit-test-extractor
`get-tests.sh` is a small tool that aims to analyze go code and put out unit tests, so we have a single view at the coverage (not counting %, this is just info).
To use this, include the directory in your $PATH, and navigate to the root of go code you need information about, and run `get-tests.sh` without any params.
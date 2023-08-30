#!/bin/bash

# Function to extract and print the test names and their descriptions
extract_and_print() {
    local file="$1"
    # Get test function names
    local funcs=$(awk '/func Test[^(]*/ {print $2}' "$file" | sed 's/(.*//')

    # For each function, print its name and then any descriptions associated with it
    echo "$funcs" | while read -r fn; do
        if [[ ! -z "$fn" ]]; then
            echo "- $fn"

            # Use awk to extract descriptions between the function start and the next function start or file end
            awk "/func ${fn}/{flag=1; next} /func Test/{flag=0} flag" "$file" | \
                grep 'description:' | \
                awk -F\" '{print "  -", $(NF-1)}'
        fi
    done
}

# UNIT TESTS
echo "UNIT_TESTS"
echo "=========="
# Recursively find all *_test.go files excluding *_integration_test.go files
find . -name '*_test.go' ! -name '*_integration_test.go' | while read -r file; do
    echo "$file"
    extract_and_print "$file"
    echo ""  # Add a blank line for readability
done

# INTEGRATION TESTS
echo "INTEGRATION_TESTS"
echo "================"
# Recursively find all *_integration_test.go files
find . -name '*_integration_test.go' | while read -r file; do
    echo "$file"
    extract_and_print "$file"
    echo ""  # Add a blank line for readability
done

#!/bin/bash

# Read the YAML content from stdin and save it to a temporary file
input_file=$(mktemp)
cat > "$input_file"

# Split the large YAML into separate files for each document
csplit -f split -b "%d.yaml" -z "$input_file" "/^---$/" "{*}"

# Go through each split file
for split_file in split*.yaml; do
    # Extract the 'kind' field using yq
    resource_kind=$(yq e '.kind' "$split_file" | tr '[:upper:]' '[:lower:]')
    resource_name=$(yq e '.metadata.name' "$split_file" | tr '[:upper:]' '[:lower:]')

    # If the 'resource_kind' field is non-empty, rename the file
    if [ -n "$resource_kind" ]; then
        mv "$split_file" "./chart/templates/${resource_name}_${resource_kind}.yaml"
    fi
done

# Clean up the temporary file
rm "$input_file"

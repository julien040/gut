#!/bin/bash

# Clone the repository into /tmp
rm -rf /tmp/gitattributes
git clone https://github.com/gitattributes/gitattributes.git /tmp/gitattributes

if [ $? -ne 0 ]; then
    echo "Failed to clone repository"
    exit 1
fi

# Change to the repository directory
cd /tmp/gitattributes || {
    echo "Failed to change directory"
    exit 1
}

# Start the Go map
echo ""
echo "var GitAttributes = map[string]string{"

# Find all gitattributes files and process them
find . -name "*.gitattributes" -type f | while read -r file; do
    # Get filename without path and extension
    filename=$(basename "$file" .gitattributes)

    # Format the content for Go string
    content=$(cat "$file" | sed 's/`/`+"`"+`/g' | sed ':a;N;$!ba;s/\n/\\n/g')

    # Output the map entry
    echo "    \"$filename\": \`$content\`,"
done

# Close the map
echo "}"

# Clean up
rm -rf /tmp/gitattributes

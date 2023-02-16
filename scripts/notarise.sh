# Check if the script is running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    echo "This script is only for macOS"
    exit 1
fi

# Check if ENV variables are set
if [[ -z "$NOTARISE_PASSWORD" || -z "$NOTARISE_USERNAME" || -z "$NOTARISE_APPLICATION_IDENTITY" ]]; then
    echo "Please set the following environment variables: NOTARISE_PASSWORD, NOTARISE_USERNAME, NOTARISE_APPLICATION_IDENTITY"
    exit 1
fi

# Check if the file exists
if [[ ! -f $1 ]]; then
    echo "File $1 does not exist"
    exit 1
fi

# Run the notarisation for amd64
xcrun altool --notarize-app --primary-bundle-id "dev.gut-cli.app" --username "$NOTARISE_USERNAME" --password "$NOTARISE_PASSWORD" --file $1

echo Building for $TARGETPLATFORM and output to $OUTPUT
# split format linux/amd64 into linux amd64 and pass as args to go to compile
OS=$(echo $TARGETPLATFORM | cut -d "/" -f 1)
ARCH=$(echo $TARGETPLATFORM | cut -d "/" -f 2)
GOOS=$OS GOARCH=$ARCH go build -o $OUTPUT .
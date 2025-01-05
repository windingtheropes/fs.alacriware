FROM golang

ARG TARGETPLATFORM

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

# Local modules
COPY auth ./auth
COPY based ./based
COPY logger ./logger

# split format linux/amd64 into linux amd64 and pass as args to go to compile
RUN GOOS=$(echo $TARGETPLATFORM | cut -d "/" -f 1 | read ouput; echo $ouput); GOARCH=$(echo $TARGETPLATFORM | cut -d "/" -f 2 | read ouput; echo $ouput); go build -o /fs-alacriware
RUN chmod a+x /fs-alacriware

EXPOSE 3030
ENV GIN_MODE=release
ENV PUBDIR="/public"

CMD ["/fs-alacriware"]
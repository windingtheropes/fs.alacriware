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

# build script
COPY gobuild.sh ./

# split format linux/amd64 into linux amd64 and pass as args to go to compile
RUN TARGETPLATFORM=$TARGETPLATFORM OUTPUT=/fs-alacriware sh gobuild.sh
RUN chmod a+x /fs-alacriware

EXPOSE 3030
ENV GIN_MODE=release
ENV PUBDIR="/public"

CMD ["/fs-alacriware"]
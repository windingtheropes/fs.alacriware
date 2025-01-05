FROM golang

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

# Local modules
COPY auth ./auth
COPY based ./based
COPY logger ./logger

RUN GOOS=linux GOARCH=amd64 go build -o /fs-alacriware
RUN chmod a+x /fs-alacriware

EXPOSE 3030
ENV GIN_MODE=release
ENV PUBDIR="/public"

CMD ["/fs-alacriware"]
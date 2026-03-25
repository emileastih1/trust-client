FROM golang:1.22

# Install git.
# Git is required for fetching the dependencies.
RUN apt update
RUN apt install -y git

WORKDIR /go/src/bulletin-board-api
COPY . .

RUN go get github.com/google/wire/cmd/wire
RUN go install github.com/google/wire/cmd/wire

RUN go get -d -v ./...
RUN go install -v ./...
RUN go vet ./...
RUN go mod download
RUN wire ./...

# build
RUN go build -ldflags "-s -w" -o /build/server ./main.go

ENTRYPOINT ["/build/server"]

EXPOSE 7000
EXPOSE 7001
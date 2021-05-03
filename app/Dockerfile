FROM golang:1.15.7

ENV WORKSPACE /var/www/
WORKDIR $WORKSPACE

COPY go.mod go.sum $WORKSPACE
RUN go mod download 

COPY cmd $WORKSPACE/cmd/
COPY internal $WORKSPACE/internal/

EXPOSE 8080

RUN go build cmd/server/main.go
ENTRYPOINT ["./main"]

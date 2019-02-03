FROM golang:1.11.5-alpine

RUN apk add --no-cache git protobuf

RUN mkdir -p otf-workout-grpc
WORKDIR /otf-workout-grpc
ADD . /otf-workout-grpc

RUN go get -u google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/aws/aws-sdk-go

RUN protoc --version
RUN protoc -I definition/ definition/workout.proto --go_out=plugins=grpc:definition
RUN go build ./server/main.go

EXPOSE 8080

CMD ["./main"]
FROM golang:latest

#Changes workdir to the local of golang

WORKDIR /go/src/NewPhoto

#Copies all the working files

COPY . .

RUN go build main.go

ENTRYPOINT . ./credentials.sh && ./main



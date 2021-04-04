FROM golang:latest

ARG pass_mysql
ARG pass_redis

#Changes workdir to the local of golang
ENV runAddr=":8082"
ENV aiAddr="newphotoai_server:9089"

#Redis credentils
ENV redisAddr="newphoto_redis"
ENV redisPassword=${pass_redis}
ENV redisUsername="default"

#Mysql credentials
ENV mysqlPassword=${pass_mysql}
ENV mysqlUsername="yariksvitlitskiy"
ENV mysqlTable="newphotos"
ENV mysqlAddr="newphoto_mysql"


WORKDIR /go/src/NewPhoto

#Copies all the working files

COPY . .

RUN go build main.go

ENTRYPOINT ./main



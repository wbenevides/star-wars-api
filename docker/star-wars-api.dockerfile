FROM golang
LABEL author="Wallace Benevides"
ADD . /go/src/github.com/wallacebenevides/star-wars-api

RUN go get github.com/gorilla/mux github.com/sirupsen/logrus gopkg.in/mgo.v2
RUN go install github.com/wallacebenevides/star-wars-api
ENTRYPOINT /go/bin/star-wars-api
EXPOSE 8080

FROM golang
LABEL author="Wallace Benevides"
ADD . /go/src/github.com/wallacebenevides/star-wars-api
RUN go get -d -v github.com/gorilla/mux github.com/sirupsen/logrus go.mongodb.org/mongo-driver/mongo

RUN go install github.com/wallacebenevides/star-wars-api
ENTRYPOINT /go/bin/star-wars-api
EXPOSE 8080

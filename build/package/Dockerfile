FROM golang
LABEL author="Wallace Benevides"

ADD . /go/src/github.com/wallacebenevides/star-wars-api
RUN go get -d -v github.com/gorilla/mux github.com/sirupsen/logrus go.mongodb.org/mongo-driver/mongo github.com/spf13/viper

RUN go install github.com/wallacebenevides/star-wars-api/cmd/starwars
ENTRYPOINT /go/bin/starwars
EXPOSE 8080
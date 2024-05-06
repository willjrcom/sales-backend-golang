FROM golang:1.21

ENV API=/go/src/sales-backend-golang \
	GO111MODULE=on \
	environment=dev \
	GOPRIVATE="github.com/willjrcom"

WORKDIR $API
COPY . $API/
RUN go mod download
RUN go install

RUN apt-get install apt-transport-https
RUN apt-get update

CMD ["/httpserver"]
EXPOSE 8080
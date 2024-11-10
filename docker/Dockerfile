FROM golang:1.23.3-alpine

WORKDIR /code

RUN apk --no-cache --update add git
RUN go install github.com/air-verse/air@latest

ENV PATH /go/bin:$PATH

EXPOSE 80

CMD ["air"]

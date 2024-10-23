FROM golang:1.23

RUN mkdir /gin

ADD . /gin

WORKDIR /gin

RUN go mod tidy

RUN go build -o gin .

RUN chmod +x /gin/gin

ENTRYPOINT ["/gin/gin"]
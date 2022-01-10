
FROM golang:1-alpine AS builder
WORKDIR /go/src/github.com/soyuka/incolore/
RUN apk --no-cache add git
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o incolore .

FROM scratch
COPY --from=builder /go/src/github.com/soyuka/incolore/incolore .
CMD ["./incolore"]

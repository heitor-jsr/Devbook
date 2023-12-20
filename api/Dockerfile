FROM golang:latest as builder

WORKDIR /go/src/app
COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags '-w -s' -o devbook .

FROM scratch
COPY --from=builder /go/src/app/devbook .
CMD ["./devbook"]
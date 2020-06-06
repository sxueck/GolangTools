FROM golang:1.13-buster as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...
RUN go build -o /go/bin/sxrehttps

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/sxrehttps /
EXPOSE 80
CMD ["/sxrehttps"]
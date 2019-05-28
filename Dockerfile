FROM golang:1.11.5 as build
WORKDIR /go/src/github.com/mitchell/selfpass
COPY . .
ENV GO111MODULE=on
RUN make gen-certs-go
RUN make build

FROM debian:stable-20190506-slim
WORKDIR /usr/bin
COPY --from=build /go/src/github.com/mitchell/selfpass/bin/server .
ENTRYPOINT ["server"]
EXPOSE 8080

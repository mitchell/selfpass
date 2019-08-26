FROM golang:1.12 as build
WORKDIR /go/src/github.com/mitchell/selfpass/services
COPY . .
ENV GO111MODULE on
RUN make build

FROM debian:stable-20190812-slim
COPY --from=build /go/src/github.com/mitchell/selfpass/services/bin/server /usr/bin/server
RUN groupadd -r selfpass && useradd --no-log-init -r -g selfpass selfpass
USER selfpass
WORKDIR /home/selfpass
ENTRYPOINT ["server"]
EXPOSE 8080

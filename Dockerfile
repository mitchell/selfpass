FROM golang:1.11.5 as build
WORKDIR /go/src/github.com/mitchell/selfpass
COPY . .
RUN go get -u golang.org/x/tools/cmd/goimports
ENV GO111MODULE=on
RUN make build

FROM debian:stable-20190326-slim
RUN printf "deb http://httpredir.debian.org/debian stretch-backports main non-free\ndeb-src http://httpredir.debian.org/debian stretch-backports main non-free" > /etc/apt/sources.list.d/backports.list
RUN apt-get update && apt-get install -t stretch-backports -y --no-install-recommends redis-server=5:5.0.3-3~bpo9+2 \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*
WORKDIR /usr/bin/selfpass/
COPY --from=build /go/src/github.com/mitchell/selfpass/bin/server .
COPY --from=build /go/src/github.com/mitchell/selfpass/redis.conf .
COPY --from=build /go/src/github.com/mitchell/selfpass/db/dump.rdb ./db/dump.rdb
COPY --from=build /go/src/github.com/mitchell/selfpass/certs/ca.pem ./certs/ca.pem
COPY --from=build /go/src/github.com/mitchell/selfpass/certs/server.pem ./certs/server.pem
COPY --from=build /go/src/github.com/mitchell/selfpass/certs/server-key.pem ./certs/server-key.pem
COPY --from=build /go/src/github.com/mitchell/selfpass/dual-entry ./dual-entry
ENTRYPOINT ./dual-entry
EXPOSE 8080

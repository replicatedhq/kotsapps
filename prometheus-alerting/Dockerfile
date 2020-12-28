FROM golang:1.15 as build
ADD . /go/src/github.com/replicatedhq/kotsapps/prometheus-alerts
RUN go install /go/src/github.com/replicatedhq/kotsapps/prometheus-alerts/cmd/example-app

FROM debian:stretch-slim
COPY --from=build /go/bin/example-app /go/bin/example-app
EXPOSE 3000
CMD ["/go/bin/example-app"]

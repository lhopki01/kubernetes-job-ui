FROM alpine:3.10

ADD kubernetes-job-ui /usr/local/bin

CMD ["kubernetes-job-ui", "serve"]

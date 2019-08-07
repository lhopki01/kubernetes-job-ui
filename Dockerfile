FROM alpine:3.10

RUN mkdir /app

ADD kubernetes-job-ui /app
ADD templates /app/templates
ADD static /app/static

CMD ["/app/kubernetes-job-ui", "serve"]

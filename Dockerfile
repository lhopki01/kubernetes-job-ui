FROM alpine:3.10

RUN mkdir /app

ADD kubernetes-job-ui /app

CMD ["/app/kubernetes-job-ui"]

FROM ubuntu:18.04
#FROM golang:latest
WORKDIR app
ADD build/benchmark-app /app
ADD files/test.txt /app/files/
CMD ["/app/benchmark-app"]
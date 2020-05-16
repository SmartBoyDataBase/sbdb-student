FROM golang:1.12-alpine as builder
RUN apk add git
COPY . /go/src/sbdb-student
ENV GO111MODULE on
WORKDIR /go/src/sbdb-student
RUN go get && go build

FROM alpine
MAINTAINER longfangsong@icloud.com
COPY --from=builder /go/src/sbdb-student/sbdb-student /
WORKDIR /
CMD ./sbdb-student
ENV PORT 8000
EXPOSE 8000
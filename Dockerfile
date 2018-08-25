ARG ARCH
FROM golang as build 

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY ./ /go/src/github.com/innovate-technologies/iammobile
WORKDIR /go/src/github.com/innovate-technologies/iammobile

RUN dep ensure -v

ARG GO_ARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH="${GO_ARCH}" go build -a -installsuffix cgo ./

ARG ARCH
FROM multiarch/alpine:${ARCH}-edge

RUN apk add --no-cache ca-certificates

COPY --from=build /go/src/github.com/innovate-technologies/iammobile/iammobile /opt/iammobile/iammobile
COPY --from=build /go/src/github.com/innovate-technologies/iammobile/html /opt/iammobile/html

ENV GH_USERNAME=""
ENV GH_TOKEN=""

WORKDIR /opt/iammobile/
CMD [ "/opt/iammobile/iammobile" ]

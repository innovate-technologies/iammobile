ARG ARCH
FROM golang as build 

COPY ./ /go/src/github.com/innovate-technologies/iammobile
WORKDIR /go/src/github.com/innovate-technologies/iammobile

ARG GO_ARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH="${GO_ARCH}" go build -a -installsuffix cgo ./

ARG ARCH
FROM multiarch/alpine:${ARCH}-edge

RUN apk add --no-cache ca-certificates

COPY --from=build /go/src/github.com/innovate-technologies/iammobile/iammobile /opt/iammobile/iammobile
COPY --from=build /go/src/github.com/innovate-technologies/iammobile/html /opt/iammobile/html

ENV GH_USERNAME=""
ENV GH_TOKEN=""

CMD [ "/opt/iammobile/iammobile" ]

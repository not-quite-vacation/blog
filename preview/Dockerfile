FROM golang:alpine
ARG pkg=not-quite-vacation/preview
RUN apk add --no-cache ca-certificates
COPY . $GOPATH/src/$pkg
WORKDIR $GOPATH/src/$pkg
RUN go install 
CMD echo "run preview if you want to see the blog"; exit 1

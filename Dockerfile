FROM alpine:3.5
RUN apk add --no-cache ca-certificates
COPY gopath/bin/blog /bin/blog
CMD echo "run blog if you want to see the blog"; exit 1

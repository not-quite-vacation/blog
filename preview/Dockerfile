FROM ubuntu:16.04
RUN apt-get update -y && apt-get install -y curl 
COPY preview /bin/preview
ENTRYPOINT ["/bin/preview"]

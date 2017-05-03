tag := foolusion/nqv:$(shell date +%y%m%d%H%M%S)

all: docker/push

docker/build: main.go bucket-cache.go
	docker build -t ${tag} .

docker/push: docker/build
	docker push ${tag}

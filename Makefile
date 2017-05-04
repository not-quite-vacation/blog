tag := foolusion/nqv:$(shell date +%y%m%d%H%M%S)

all: kube.yaml

preview: main.go bucket-cache.go ../blog/static.go
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .

docker/build: main.go bucket-cache.go
	docker build -t ${tag} .

docker/push: docker/build
	docker push ${tag}

kube.yaml: docker/push
	sed "s#IMAGE#${tag}#;" k8s/deployment.yaml > kube.yaml

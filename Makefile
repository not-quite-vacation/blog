image := gcr.io/notquitevacation/blog
tag := cloud-build-test

all: kube.yaml

kube.yaml: 
	sed "s#IMAGE#${image}:${tag}#;" k8s/deployment.yaml > kube.yaml

kube/deploy: kube.yaml
	kubectl apply -f kube.yaml

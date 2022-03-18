SERVICE := am-stats-legacy
NAMESPACE := aftermath-services
# REGISTRY := ghcr.io/byko-dev
REGISTRY := docker.io/vkouzin
# 
VERSION = $(shell git rev-parse --short HEAD)
TAG := ${REGISTRY}/${SERVICE}

echo:
	@echo "Tag:" ${TAG}

pull:
	git pull

build:
	docker build -t ${TAG}:${VERSION} -t ${TAG}:latest .
	docker image prune -f

push:
	docker push ${TAG}:latest

apply:
	kubectl apply -f .kube/

restart:
	kubectl rollout restart deployment/${SERVICE} -n ${NAMESPACE}

ctx:
	kubectl config set-context --current --namespace=${NAMESPACE}
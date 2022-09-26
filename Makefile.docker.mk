

docker-login:
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}

docker-build:
	docker build -t fsrg/$(PROJECT_NAME):$(APP_VERSION) .

docker-push:
	docker push fsrg/$(PROJECT_NAME):$(APP_VERSION)

docker-deploy: docker-login docker-build docker-push

docker-run:
	docker run -p 7070:7070 fsrg/$(PROJECT_NAME):$(APP_VERSION)

.PHONY: docker-login docker-build docker-push docker-deploy docker-run
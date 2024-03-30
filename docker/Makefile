build-image:
	DOCKER_BUILDKIT=1 docker build --progress=plain -t nftm/market-server-pre -f $(ROOT_PATH)/docker/Dockerfile-pre.buildkit $(ROOT_PATH)
	DOCKER_BUILDKIT=1 docker build --progress=plain -t ${SERVER_IMAGE} -f $(ROOT_PATH)/docker/Dockerfile.buildkit $(ROOT_PATH)

push-image:
	docker push ${SERVER_IMAGE}

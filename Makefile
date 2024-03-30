SLASH:=/
REPLACE_SLASH:=\/
export ROOT_PATH ?= ${PWD}
export MODE ?= prod
ROOT_PATH_REPLACE=$(subst $(SLASH),$(REPLACE_SLASH),$(ROOT_PATH))

export ARCH:=$(shell arch)
DOCKER=`which docker`

export CODE_PATH=${PWD}/src
export TEST_PATH=${PWD}/test
export REPORT_PATH?=/tmp/report
export VERSION ?= latest
-include .makerc/help
-include .makerc/docker

.EXPORT_ALL_VARIABLES:

check:
	cd ${ROOT_PATH}/src/listener && staticcheck .
	cd ${ROOT_PATH}/src/server && staticcheck .

build-image: ##@Docker build docker images
	@$(MAKE) -C docker build-image

push-image: ##Docker push docker images
	@$(MAKE) -C docker push-image

clean-image:
	docker rmi ${SERVER_IMAGE}
	#docker rmi -f `docker images -f "dangling=true" -q`

alldefconfig:
	python3 kconfig-lib/alldefconfig.py

menuconfig:
	MENUCONFIG_STYLE=aquatic python3 kconfig-lib/menuconfig.py

start:
	if [ "$(CONFIG_DOCKER_COMPOSE)" = "y" ]; then \
  		make -C bootup/docker-compose start; \
  	fi

stop:
	if [ "$(CONFIG_DOCKER_COMPOSE)" = "y" ]; then \
  		make -C bootup/docker-compose stop; \
  	fi

HELP_FUN = \
	%help; \
	while(<>) { push @{$$help{$$2 // 'options'}}, [$$1, $$3] if /^([a-zA-Z\-]+)\s*:.*\#\#(?:@([a-zA-Z\-]+))?\s(.*)$$/ }; \
	print "usage: make [target]\n\n"; \
	for (sort keys %help) { \
	print "${WHITE}$$_:${RESET}\n"; \
	for (@{$$help{$$_}}) { \
	$$sep = " " x (32 - length $$_->[0]); \
	print "  ${YELLOW}$$_->[0]${RESET}$$sep${GREEN}$$_->[1]${RESET}\n"; \
	}; \
	print "\n"; }

help: ##@other Show this help.
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)



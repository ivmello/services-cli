.PHONY: all services_cli create_volume
all: services_cli
build:
	docker build -t services_cli .
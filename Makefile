.PHONY: setup up down clean help

## setup: Install k3d cluster with local registry
setup:
	k3d cluster create hikmah-cluster \
		-p "8080:80@loadbalancer" \
		-p "5432:5432@loadbalancer" \
		--agents 1

## up: Launch Tilt (The Dev Loop)
up:
	tilt up

## down: Stop Tilt and remove resources
down:
	tilt down

## clean: Delete the entire k3d cluster
clean:
	k3d cluster delete hikmah-cluster

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

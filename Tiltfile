# -*- mode: Python -*-

default_registry("localhost:5005")

# 1. Infrastructure - Load everything via Kustomize
# This includes postgres, redpanda, kratos, maktba, and gateway
k8s_yaml(local("kustomize build deploy/overlays/dev"))

# 2. Resources Configuration
k8s_resource("postgres", port_forwards=5432)

k8s_resource("kratos", port_forwards=["4433:4433", "4434:4434"])

# 3. Maktba Build
docker_build(
    "maktba-image",
    context=".",
    dockerfile="./src/maktba/Dockerfile",
    live_update=[
        sync("./src/maktba", "/src/src/maktba"),
    ],
)

# We define the resource here to add the port forward and dependency
k8s_resource("maktba", port_forwards=5000, resource_deps=["postgres", "redpanda"])

# 4. Gateway Build
docker_build(
    "gateway-image",
    context="./src/madkhal",
    dockerfile="./src/madkhal/Dockerfile",
)
# resource_deps ensures it doesn't stay 'Red' while maktba is still building
k8s_resource("gateway", port_forwards=8080, resource_deps=["maktba", "kratos"])

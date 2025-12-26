# -*- mode: Python -*-

# 1. Infrastructure (Backing Services)
# specific to dev environment
k8s_yaml(local("kustomize build deploy/overlays/dev"))

# Expose Redpanda and Postgres to localhost for convenience (optional)
k8s_resource("postgres", port_forwards=5432)
k8s_resource("redpanda", port_forwards=["9094:9094", "8081:8081"])
k8s_resource('grafana', port_forwards='3000:3000')
k8s_resource('kratos', port_forwards=['4433:4433', '4434:4434'])


# 2. Maktba (Catalog - .NET 10)
docker_build(
    "maktba-image",
    '.',
    dockerfile="./src/maktba/Dockerfile",
    live_update=[
        sync("./src/maktba", "/app"),
    ],
    # entrypoint is usually defined in Dockerfile, but for dev we might override to 'dotnet watch'
)

# For now, let's comment out the service resource until we have the manifest.
k8s_resource("maktba", port_forwards=5000)

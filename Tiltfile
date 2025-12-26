# -*- mode: Python -*-

# 1. Backing Services (Infrastructure)
# We deploy these first as static YAMLs or Helm charts
k8s_yaml('deploy/manifests/db.yaml')
k8s_yaml('deploy/manifests/redpanda.yaml')

# 2. Khazin (.NET 9)
docker_build('al-khazin-image', './al-khazin', 
    dockerfile='./al-khazin/Dockerfile',
    live_update=[
        sync('./al-khazin', '/app') # Syncs code for fast feedback
    ])
k8s_yaml('deploy/manifests/al-khazin.yaml')
k8s_resource('al-khazin', port_forwards=5000)

# 3. Warraq (Go)
docker_build('al-warraq-image', './al-warraq',
    dockerfile='./al-warraq/Dockerfile',
    live_update=[
        sync('./al-warraq', '/app'),
        run('cd /app && go build -o main ./cmd/main.go', trigger='./al-warraq/cmd')
    ])
k8s_yaml('deploy/manifests/al-warraq.yaml')
k8s_resource('al-warraq', port_forwards=8081)

# 4. Gateway (KrakenD)
docker_build('gateway-image', './gateway', dockerfile='./gateway/Dockerfile')
k8s_yaml('deploy/manifests/gateway.yaml')
k8s_resource('gateway', port_forwards=8080)

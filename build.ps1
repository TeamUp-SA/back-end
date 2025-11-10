# === SECTION 1: CLEAN UP OLD RESOURCES ===
Write-Host "--- (1/5) Deleting old Kubernetes resources..."
# We pipe errors to $null to suppress "not found" messages, which are expected.
kubectl delete -f k8s-full-stack.yaml 2>$null
kubectl delete secret postgres-secret user-service-secret auth-service-secret notification-service-secret app-service-secret 2>$null
kubectl delete configmap krakend-config-cm postgres-init-sql-cm 2>$null
kubectl delete pvc postgres-pvc redis-pvc 2>$null
Write-Host "--- Cleanup complete."

# === SECTION 2: BUILD DOCKER IMAGES ===
Write-Host "--- (2/5) Building local Docker images..."
docker build -t user-service:latest ./user-service
docker build -t auth-service:latest -f auth-service/Dockerfile .
docker build -t notification-service:latest ./notification-service
docker build -t app-service:latest ./app-service
Write-Host "--- Docker images built."

# === SECTION 3: CREATE CONFIGS & SECRETS ===
Write-Host "--- (3/5) Creating ConfigMaps and Secrets..."
kubectl create configmap krakend-config-cm --from-file=krakend.json=./krakend.json
kubectl create configmap postgres-init-sql-cm --from-file=init.sql=./init.sql

kubectl create secret generic postgres-secret --from-env-file=postgres.env
kubectl create secret generic user-service-secret --from-env-file=user-service.env
kubectl create secret generic auth-service-secret --from-env-file=auth-service.env
kubectl create secret generic notification-service-secret --from-env-file=notification-service.env
kubectl create secret generic app-service-secret --from-env-file=app-service.env
Write-Host "--- Configs and Secrets created."

# === SECTION 4: DEPLOY TO KUBERNETES ===
Write-Host "--- (4/5) Deploying k8s-full-stack.yaml..."
kubectl apply -f k8s-full-stack.yaml
Write-Host "--- Deployment applied."

# === SECTION 5: WATCH STATUS ===
Write-Host "--- (5/5) Watching pod status... (Press Ctrl+C to exit)"
kubectl get pods -w
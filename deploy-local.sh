#!/bin/bash

# Define the namespace
NAMESPACE="my-microservices"

# Function to check for errors and exit
check_error() {
    if [ $? -ne 0 ]; then
        echo "❌ ERROR: $1"
        exit 1
    fi
}

echo "--- 🚀 Starting Microservices Deployment Pipeline (Namespace: $NAMESPACE) ---"

# === 1. Namespace Check/Creation ===
if ! kubectl get namespace "$NAMESPACE" &> /dev/null; then
    echo "--- Creating namespace '$NAMESPACE'..."
    kubectl create namespace "$NAMESPACE"
    check_error "Failed to create namespace $NAMESPACE."
fi
echo "--- Target namespace set."

# === 2. CLEAN UP OLD RESOURCES ===
echo "--- (1/5) Deleting old Kubernetes resources..."
# We pipe stderr to /dev/null to suppress "not found" messages.
kubectl delete -f k8s-full-stack.yaml -n "$NAMESPACE" 2>/dev/null
kubectl delete secret postgres-secret user-service-secret auth-service-secret notification-service-secret app-service-secret search-service-secret -n "$NAMESPACE" 2>/dev/null
kubectl delete configmap krakend-config-cm postgres-init-sql-cm -n "$NAMESPACE" 2>/dev/null
kubectl delete pvc postgres-pvc redis-pvc -n "$NAMESPACE" 2>/dev/null

# Wait a moment for resources to clean up (optional but helpful)
sleep 2 
echo "--- Cleanup complete."

# === 3. BUILD DOCKER IMAGES ===
echo "--- (2/5) Building local Docker images..."
docker build -t user-service:latest ./user-service
check_error "User Service image build failed."
docker build -t auth-service:latest -f auth-service/Dockerfile .
check_error "Auth Service image build failed."
docker build -t notification-service:latest ./notification-service
check_error "Notification Service image build failed."
docker build -t app-service:latest ./app-service
check_error "App Service image build failed."
docker build -t search-service:latest ./search-service
check_error "Search Service image build failed."
echo "--- Docker images built."

# === 4. CREATE CONFIGS & SECRETS ===
echo "--- (3/5) Creating ConfigMaps and Secrets in $NAMESPACE..."
kubectl create configmap krakend-config-cm --from-file=krakend.json=./krakend.json -n "$NAMESPACE"
check_error "ConfigMap krakend-config-cm failed."
kubectl create configmap postgres-init-sql-cm --from-file=init.sql=./init.sql -n "$NAMESPACE"
check_error "ConfigMap postgres-init-sql-cm failed."

kubectl create secret generic postgres-secret --from-env-file=postgres.env -n "$NAMESPACE"
check_error "Secret postgres-secret failed."
kubectl create secret generic user-service-secret --from-env-file=user-service.env -n "$NAMESPACE"
check_error "Secret user-service-secret failed."
kubectl create secret generic auth-service-secret --from-env-file=auth-service.env -n "$NAMESPACE"
check_error "Secret auth-service-secret failed."
kubectl create secret generic notification-service-secret --from-env-file=notification-service.env -n "$NAMESPACE"
check_error "Secret notification-service-secret failed."
kubectl create secret generic app-service-secret --from-env-file=app-service.env -n "$NAMESPACE"
check_error "Secret app-service-secret failed."
kubectl create secret generic search-service-secret --from-env-file=search-service.env -n "$NAMESPACE"
check_error "Secret search-service-secret failed."
echo "--- Configs and Secrets created."

# === 5. DEPLOY TO KUBERNETES ===
echo "--- (4/5) Deploying k8s-full-stack.yaml to $NAMESPACE..."
kubectl apply -f k8s-full-stack.yaml -n "$NAMESPACE"
check_error "Deployment via k8s-full-stack.yaml failed."
echo "--- Deployment applied."

# === 6. WATCH STATUS ===
echo "--- (5/5) Watching pod status in $NAMESPACE... (Press Ctrl+C to exit)"
kubectl get pods -n "$NAMESPACE" -w
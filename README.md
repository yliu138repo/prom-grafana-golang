# Setting Up KinD Kubernetes Cluster with Prometheus, Grafana, and K6 for Monitoring and Stress Testing a Go Application

In this guide, we'll walk through setting up a Kubernetes cluster using **KIND** (Kubernetes IN Docker) and configuring **Prometheus** and **Grafana** for monitoring a sample Go application in local environment. This application is built with the **CHI Router**, **PostgreSQL (pgx)**, and the **Viper** library, demonstrating how to monitor any web app deployed on Kubernetes. Additionally, we will use K6 to perform stress testing on the application, allowing us to observe its performance under load and effectively utilize the monitoring setup.

## Prerequisites

Ensure you have the following tools installed on your machine:

- **Docker** - [Installation Guide](https://docs.docker.com/engine/install/)
- **KinD** - [Installation Guide](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- **Helm** - [Installation Guide](https://helm.sh/docs/intro/install/)
- **Go** - [Installation Guide](https://go.dev/doc/install)

## Step 1: Create a KIND Cluster

To set up a Kubernetes cluster with two nodes, switch to kindk8 folder and run the following command:
(you can adjust the host and container part as per your local system)
```bash
kind create cluster --config kind-config.yaml
```

## Step 2: Verify Cluster Information
Get information about your cluster to confirm it was created successfully:
```bash
kubectl cluster-info --context kind-kind
```

## Step 3: Install Prometheus and Grafana with Node Metrics Exporter
Add the necessary Helm repositories and install Prometheus and Grafana with the node metrics exporter on your cluster:
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add stable https://charts.helm.sh/stable
helm repo update
kubectl create namespace monitoring
helm install kind-prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --set prometheus.service.nodePort=30000 \
  --set prometheus.service.type=NodePort \
  --set grafana.service.nodePort=31000 \
  --set grafana.service.type=NodePort \
  --set alertmanager.service.nodePort=32000 \
  --set alertmanager.service.type=NodePort \
  --set prometheus-node-exporter.service.nodePort=32001 \
  --set prometheus-node-exporter.service.type=NodePort
```
Verify the services in the monitoring namespace:
```bash
kubectl get svc -n monitoring
kubectl get namespaces
```
## Step 4: Label Nodes
Label the nodes in your cluster for easy identification:
```bash
kubectl label nodes kind-worker node=worker1
kubectl label nodes kind-worker2 node=worker2
kubectl get nodes --show-labels
```
## Step 5: Build and Load Your Application Docker Image
Build the Docker image for your Go application and load it into the KIND cluster:
```bash
docker build -t dsi/gok8-app:latest .
kind load docker-image dsi/gok8-app:latest --name kind
```
## Step 6: Deploy Your Application
Apply the `deployment.yaml` file to create a pod for your application:
(Make sure you use your PostgreSQL database configuration in the `deployment.yam` `env` section and your PostgreSQL is running and accessible from KinD cluster)
```bash
kubectl apply -f deployment.yaml
```
Verify that the pod has been created:
```bash
kubectl get pods
```
### Additional Debugging Commands
Describe a specific pod:

```bash
kubectl describe pod <pod-name>
```
View logs of a specific pod:

```bash
kubectl logs <pod-name>
```
Open a shell in a pod:

```bash
kubectl exec -it <pod-name> -- /bin/sh
```
Debug a node:
```bash
kubectl debug node/<node-name> -it --image=busybox --namespace=kube-system
```
## Step 7: Expose Your Application with a Service
Create a service for your application by applying the `service.yaml` file:

```bash
kubectl apply -f service.yaml
```
Verify that the service and pods are running:

```bash
kubectl get pods
kubectl get svc
```
## Step 8: Forward Service Port to Localhost
Forward the service port to access your application locally:

```bash
kubectl port-forward svc/go-app 8080:8080 &
```
Now, your application should be accessible at http://localhost:8080.

## Step 9: Forward Grafana and Prometheus Ports to Localhost
Forward the Grafana and Prometheus services to access them on your local machine:

```bash
kubectl port-forward svc/kind-prometheus-kube-prome-prometheus -n monitoring 9090:9090 --address=0.0.0.0 &
kubectl port-forward svc/kind-prometheus-grafana -n monitoring 31000:80 --address=0.0.0.0 &
```
- Prometheus URL: http://localhost:9090
- Grafana URL: http://localhost:31000
Note: The default Grafana login credentials are:
Username: admin
Password: prom-operator

## Step 10: Add Application Metrics to Prometheus
Add your Go app's metrics to Prometheus by applying the go-app-podmonitor.yaml file:

```bash
kubectl apply -f go-app-podmonitor.yaml
```
Verify that your app is listed under Prometheus targets by accessing Prometheus at http://localhost:9090/targets.

## Step 11: Import Grafana Dashboard
Switch to grafana folder and import the k8-go-app-monitoring-dashboard.json file into your Grafana dashboard:

- Log in to Grafana at http://localhost:31000.
- Click on the "+" icon on the left panel and select "Import".
- Upload the k8-go-app-monitoring-dashboard.json file.
- Replace the existing pod name with your pod name in the dashboard settings if necessary.

## Step 12: Run Load Testing with k6
Before running stress test, make sure PostgreSQL database is running with configuration that you have setup in deployment.yaml file.
Switch to K6 folder and run a load test on your API using k6:

```bash
k6 run k6/api_test.js
```

This will simulate traffic to your application, allowing you to see real-time metrics in Prometheus and Grafana.

**If you want to run Prometheus and Grafana in a Docker environment, you can try `docker-compose-full.yml`**


Happy Coding !
# K3d Deployment

These manifests mirror `docker-compose.yml` for a local k3d cluster.

Build and import the API image:

```bash
docker compose build api
k3d image import golang_vs_api_demo-api:latest -c demo
```

Deploy:

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/
kubectl -n golang-vs-api-demo rollout status deploy/postgres
kubectl -n golang-vs-api-demo rollout status deploy/api
```

Access the API locally:

```bash
kubectl -n golang-vs-api-demo port-forward svc/api 8080:8080
curl http://localhost:8080/healthz
```

Delete:

```bash
kubectl delete -f k8s/
```

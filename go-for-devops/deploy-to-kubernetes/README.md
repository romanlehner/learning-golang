# Deploy to Kubernetes with Golang
In this exercise we create a k8s cluster and create a nginx deployement, service and ingress resource. The code was written to clean up k8s resources, but the cluster needs to be created beforehand and deleted afterwards manually. 

Notes for myself:
- shut down app gracefully to clean up resources
- handling context

## Create a k8s cluster with kind
Create the cluster and verify the kubeconfig has been created:
```bash
kind create cluster --name workloads --config kind-config.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

cat ~/.kube/config
```
Run the script:
```bash
go run main.go
```
It should clean up the resources by itself. 

## Test the app
The logs of each pod should now be streamed to the console output. Check the resources being created:
```bash
kubectl -n nginx get svc,pod,ing
```

You can either open the app in the browser or run a curl against the ingress. The access logs should be reflected in the log stream.
```bash
curl localhost:8080/hello
10.244.0.7 - - [xx/Feb/20xx:20:15:28 +0000] "GET /hello HTTP/1.1" 200 12142 "-" "curl/7.68.0" "172.18.0.1"
```
Shut the app down with `ctrl+C`. Resources should be gracefully removed:
```bash
Deleting ingress "nginx".
Deleting service "nginx".
Deleting deployment "nginx".
Deleting namespace "nginx".
```

## Delete the k8s cluster
```bash
kind delete cluster --name workloads
```


go mod tidy
make generate
make manifests
eval $(minikube -p minikube docker-env)
make docker-build
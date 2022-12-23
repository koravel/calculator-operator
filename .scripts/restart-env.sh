minikube stop
minikube delete
minikube start
eval $(minikube -p minikube docker-env)
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.6.1/aio/deploy/recommended.yaml
kubectl apply -f config/dashboard/role-binding.yaml
kubectl apply -f config/dashboard/account.yaml
#operator-sdk olm install
kubectl get service -A
kubectl get deployment -A
kubectl get pods -A
kubectl apply -f config/samples/calc_v1alpha1_calculator.yaml
kubectl describe crd calculators.calc.example.com
kubectl describe calculator/calculator-sample
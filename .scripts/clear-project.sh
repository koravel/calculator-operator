mkdir .tmp
cp api/v1alpha1/calculator_types.go .tmp/
cp -r config/dashboard/ .tmp/
cp -r config/samples/calc_v1alpha1_calculator.yaml .tmp/
cp controllers/calculator_controller.go .tmp/
cp controllers/calculator_controller_test.go .tmp/

rm -rf api/
sudo rm -rf  bin/
rm -rf  config/
rm -rf  controllers/
rm -rf  hack/
rm .dockerignore
rm .gitignore
rm Dockerfile
rm go.mod
rm go.sum
rm cover.out
rm main.go
rm Makefile
rm PROJECT
rm README.md
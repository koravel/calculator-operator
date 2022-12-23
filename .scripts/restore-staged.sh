cp .tmp/calculator_types.go api/v1alpha1/
cp -r .tmp/dashboard/ config/
cp -r .tmp/calc_v1alpha1_calculator.yaml config/samples/
cp .tmp/calculator_controller.go controllers/
cp .tmp/calculator_controller_test.go controllers/

rm -rf .tmp/
operator-sdk init --domain example.com --repo github.com/calculator-operator
operator-sdk create api --group calc --version v1alpha1 --kind Calculator --plugins="deploy-image/v1-alpha" --image=memcached:1.4.36-alpine --image-container-command="memcached,-m=64,modern,-v" --run-as-user="1001"
sed -i 's/IMG ?= controller:latest/IMG ?= $(IMAGE_TAG_BASE):$(VERSION)/g' Makefile
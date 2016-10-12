test:
	go test -cover
	go vet

install-glide:
	curl https://glide.sh/get | sh

deps:
	install-glide
	GO15VENDOREXPERIMENT=1 glide install --cache

deps-update:
	install-glide
	rm -rf ./vendor
	GO15VENDOREXPERIMENT=1 glide update --cache

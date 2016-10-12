test:
	go test -cover
	go vet

deps:
	GO15VENDOREXPERIMENT=1 glide install --cache

deps-update:
	rm -rf ./vendor
	GO15VENDOREXPERIMENT=1 glide update --cache

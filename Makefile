install: bindata
	go install .

build: bindata
	go get github.com/mitchellh/gox
	gox -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" .

bindata:
	go-bindata -o assets/views.go -pkg assets assets/...

.PNONY: install build bindata

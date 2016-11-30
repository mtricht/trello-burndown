install: bindata
	go install .

build: bindata
	go get github.com/mitchellh/gox
	gox -osarch="linux/amd64 linux/386" -cgo -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" .

bindata:
	go-bindata -o assets/views.go -pkg assets assets/...

docker:
	docker build --no-cache -t trello-burndown .

.PNONY: install build bindata docker

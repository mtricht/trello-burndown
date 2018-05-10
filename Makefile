install: bindata
	go install .

build: bindata
	go get github.com/mitchellh/gox
	gox -osarch="linux/amd64" -cgo -output="build/trello-burndown_{{.OS}}_{{.Arch}}" ./cmd/

bindata:
	go-bindata -o assets/views.go -pkg assets assets/...

docker:
	docker build --no-cache -t trello-burndown .

run: install
	trello-burndown

.PNONY: install build bindata docker run

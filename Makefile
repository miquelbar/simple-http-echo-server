project = simple-http-echo-server

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./bin/$(project)-linux .

build-osx:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./bin/$(project)-osx .

build-with-docker:
	@docker build -t $(project) -f Dockerfile .

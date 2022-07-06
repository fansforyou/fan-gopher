build:
	go build -o fan-gopher main.go

release: release-linuxx64 release-macx64 release-winx64

release-linuxx64:
	env GOOS=linux GOARCH=amd64 go build -o dist/fan-gopher-linuxx64-$(version) main.go

release-macx64:
	env GOOS=darwin GOARCH=amd64 go build -o dist/fan-gopher-macx64-$(version) main.go

release-winx64:
	env GOOS=windows GOARCH=amd64 go build -o dist/fan-gopher-winx64-$(version).exe main.go
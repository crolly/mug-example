bin/% : functions/%/main.go
		env GOOS=linux go build -ldflags="-s -w" -o $@

build: vendor 

vendor: Gopkg.toml
	    dep ensure
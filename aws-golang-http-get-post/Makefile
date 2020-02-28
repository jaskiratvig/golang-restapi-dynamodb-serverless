.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/getArtist getArtist/getArtist.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getAllArtists getAllArtists/getAllArtists.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/createArtist createArtist/createArtist.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/editArtist editArtist/editArtist.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/deleteArtist deleteArtist/deleteArtist.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

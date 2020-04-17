.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/getArtist routers/crud/getArtist/getArtist.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getAllArtists routers/crud/getAllArtists/getAllArtists.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/createArtist routers/crud/createArtist/createArtist.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/editArtist routers/crud/editArtist/editArtist.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/deleteArtist routers/crud/deleteArtist/deleteArtist.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/home routers/auth0/home/home.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/login routers/auth0/login/login.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/callback routers/auth0/callback/callback.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

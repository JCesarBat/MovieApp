movie-gen: 
	$env:GOOS="linux"; $env:CGO_ENABLED="0"; go build -o main movie/cmd/main.
metadata-gen:
	$env:GOOS="linux"; $env:CGO_ENABLED="0"; go build -o main metadata/cmd/main.go
rating-gen:
		

movie-build:
	docker build -f movie/Dockerfile -t movie .

metadata-build:
	docker build -f movie/Dockerfile -t metadata .

rating-build:
	docker build -f rating/Dockerfile -t rating .

mock-metadata-controller:
	 mockgen -package=repository -source=metadata/internal/controller/metadata/controller.go
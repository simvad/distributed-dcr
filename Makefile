# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

build:
	$(GOBUILD) -o bin/server cmd/server/main.go
	$(GOBUILD) -o bin/cli cmd/cli/main.go

clean:
	$(GOCLEAN)
	rm -rf bin/*
	

run:
	$(GOBUILD) -o bin/server cmd/server/main.go
	$(GOBUILD) -o bin/cli cmd/cli/main.go
	./bin/server
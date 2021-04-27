binary=bin/couture

all: fmt clean build

install: all
	@echo installing
	@cp $(binary) $(HOME)/bin/

build:
	@echo building
	@go build -o $(binary) cmd/couture.go

clean:
	@echo cleaning
	@rm -f $(binary)

fmt:
	@echo formatting
	@find cmd internal -name \*.go -exec go fmt {} \;

tidy:
	@echo tidying
	@go mod tidy

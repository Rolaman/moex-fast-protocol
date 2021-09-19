.PHONY:	deps run

deps:
	export GO111MODULE=on && go mod vendor

run: deps
	go run -mod=vendor cmd/main.go

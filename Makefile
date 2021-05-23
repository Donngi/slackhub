.PHONY: build, deploy, tidy

build:
	cd "$(PWD)/initial" && make build
	cd "$(PWD)/interactive" && make build
	cd "$(PWD)/register" && make build
	cd "$(PWD)/editor" && make build
	cd "$(PWD)/catalog" && make build
	cd "$(PWD)/eraser" && make build
	cd "$(PWD)/examples/go" && make build

pr-prep:
	cd "$(PWD)/initial" && go test
	cd "$(PWD)/interactive" && go test
	cd "$(PWD)/register" && go test
	cd "$(PWD)/editor" && go test
	cd "$(PWD)/catalog" && go test
	cd "$(PWD)/eraser" && go test
	cd "$(PWD)/examples/go" && go test
	
deploy: build
	cd "$(PWD)/terraform/env/example" && terraform init && terraform apply

tidy:
	cd "$(PWD)/initial" && make tidy
	cd "$(PWD)/interactive" && make tidy
	cd "$(PWD)/register" && make tidy
	cd "$(PWD)/editor" && make tidy
	cd "$(PWD)/catalog" && make tidy
	cd "$(PWD)/eraser" && make tidy
	cd "$(PWD)/examples/go" && make tidy

update-dependencies:
	cd "$(PWD)/initial" && go get -u
	cd "$(PWD)/interactive" && go get -u
	cd "$(PWD)/register" && go get -u
	cd "$(PWD)/editor" && go get -u
	cd "$(PWD)/catalog" && go get -u
	cd "$(PWD)/eraser" && go get -u
	cd "$(PWD)/examples/go" && go get -u

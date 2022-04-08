#	Copyright (c) 2022 Cyber Home Lab authors
#
#	Licensed under the Apache License, Version 2.0 (the "License");
#	you may not use this file except in compliance with the License.
#	You may obtain a copy of the License at
#
#		http://www.apache.org/licenses/LICENSE-2.0
#
#	Unless required by applicable law or agreed to in writing, software
#	distributed under the License is distributed on an "AS IS" BASIS,
#	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#	See the License for the specific language governing permissions and
#	limitations under the License.

COPYRIGHT=Copyright (c) 2022 Cyber Home Lab authors
DEBUG_MSG=$$(tput setaf 6)DEBUG  : \t$$(tput sgr0)
INFO_MSG=$$(tput setaf 2)INFO   : \t$$(tput sgr0)
WARNING_MSG=$$(tput setaf 3)WARNING: \t$$(tput sgr0)
ERROR_MSG=$$(tput setaf 1)ERROR  : \t$$(tput sgr0)
SHELL=/bin/bash
GO111MODULE=on
BINARY_NAME=core
CYBERHOMELAB_CONFIG=$(shell echo $$(pwd)/config.toml)

GIT_TAG?=$(shell git describe --tags --match "[0-9]*")
LDFLAGS="-X 'main.Version=${GIT_TAG}'"
GO_BUILD=CGO_ENABLED=0 go build -ldflags=$(LDFLAGS)

export GO111MODULE
export CYBERHOMELAB_CONFIG

.DEFAULT_GOAL := help

.ONESHELL:
init:
	@git config --local core.hooksPath .githooks/
	go version || ( echo "$(ERROR_MSG)go not found, please install it https://go.dev/doc/install" ; exit 2 )
	golangci-lint version || ( echo "$(ERROR_MSG)golangci-lint not found, please install it https://golangci-lint.run/usage/install/" ; exit 2 )
	if ! grep -q CYBERHOMELAB_CONFIG ~/.bashrc
	then
		echo "export CYBERHOMELAB_CONFIG=$(CYBERHOMELAB_CONFIG)" >> ~/.bashrc
		export CYBERHOMELAB_CONFIG=$(CYBERHOMELAB_CONFIG)
	fi

clean:
	@go clean
	rm -f "$(BINARY_NAME)" __debug_bin nohup.out
	rm -rf /tmp/{{file{1,2,3},source,destination}.txt,source,destination,pre-push.{err,log}}
	echo -e "$(INFO_MSG)Cleanup was executed"

build-linux:
	@mkdir -p bin
	GOOS=linux		GOARCH=amd64 $(GO_BUILD) -o bin/$(BINARY_NAME) main.go

build-release:
	mkdir -p bin
	GOOS=linux		GOARCH=amd64 $(GO_BUILD) -o bin/$(BINARY_NAME)-linux-x86_64 main.go
	GOOS=linux		GOARCH=arm64 $(GO_BUILD) -o bin/$(BINARY_NAME)-linux-arm64 main.go
	GOOS=darwin		GOARCH=amd64 $(GO_BUILD) -o bin/$(BINARY_NAME)-darwin-x86_64 main.go
	GOOS=darwin		GOARCH=arm64 $(GO_BUILD) -o bin/$(BINARY_NAME)-darwin-arm64 main.go
	GOOS=windows	GOARCH=amd64 $(GO_BUILD) -o bin/$(BINARY_NAME)-windows-x86_64.exe main.go

build: build-linux

run-go:
	@go run main.go

run-binary: build
	@./bin/$(BINARY_NAME)

run: run-go

tests:
	@go test -v ./... ; echo $?

check:
	@.githooks/pre-push
	echo -e "$(INFO_MSG)All check(s) passed"

check-pre-push:
	@cat /tmp/pre-push.err /tmp/pre-push.log

check-coverage:
	@go test -cover ./... -coverprofile=coverage.out >/dev/null
	go tool cover -func=coverage.out | grep -vE '^total:|init|GetLastMessage|(5|6|7|8|9|10)[0-9].[0-9]%$$'
	go tool cover -html=coverage.out -o coverage.html
	echo -e "\n * For more info, check coverage.html or visit https://go.dev/blog/cover\n"

check-todos:
	@find . -type f -name '*.go' -exec grep -n TODO {} +

.ONESHELL:
check-copyright:
	@for file in $$(find . -type f -name "*.go")
	do
		if head -3 "$${file}" | tail -1 | grep -v $(COPYRIGHT)
		then
			echo -e "$(ERROR_MSG)File $${file} doesn't have the copyright set"
		else
			echo -e "$(INFO_MSG)File $${file} has the copyright set"
		fi
	done

.ONESHELL:
help:
	@echo -e "
	================================================
	\t\tCyber Home Lab
	\tThe core of Cyber Home Lab platform
	================================================
	
	Commands available:
	\tmake init
	\tmake clean
	\tmake build
	\tmake check
	\tmake check-pre-push
	\tmake check-coverage
	\tmake check-todos
	\tmake check-copyright
	\tmake tests
	\tmake run
	\tmake run-go
	\tmake run-binary
	"

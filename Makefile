SHELL := /bin/bash

test:
	go test ./

VERSION := $(shell cat VERSION | tr -d '\n')
push-tag:
	git tag v$(VERSION)
	git push origin v$(VERSION)

publish-version:
	GOPROXY=proxy.golang.org go list -m github.com/gelmium/graceful-shutdown@v$(VERSION)
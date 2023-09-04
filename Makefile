SHELL := /bin/bash

test:
	go test ./

VERSION:=0.0.1
push-tag:
	git tag v$(VERSION)
	git push origin v$(VERSION)

publish-version:
	GOPROXY=proxy.golang.org go list -m github.com/gelmium/graceful-shutdown@v$(VERSION)
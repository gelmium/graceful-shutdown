SHELL := /bin/bash

test:
	go test ./

VERSION := $(shell cat VERSION | tr -d '\n')
push-tag:
	git tag v$(VERSION)
	git push origin v$(VERSION)

publish-version:
	GOPROXY=proxy.golang.org go list -m github.com/gelmium/graceful-shutdown@v$(VERSION)

# NEW_VERSION is the next patch version
NEW_VERSION := $(shell echo $(VERSION) | awk -F. '{$$NF = $$NF + 1;} 1' | sed 's/ /./g')
PRERELEASE_VERSION := $(NEW_VERSION)-$(shell git rev-parse --short HEAD)
echo-version:
	@echo "Current version: $(VERSION)"
	@echo "New version: $(NEW_VERSION)"
	@echo "Prerelease version: $(PRERELEASE_VERSION)"
bump-version:
	echo -n $(NEW_VERSION) > VERSION
	git add VERSION

BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
# if env not availlable then try to use git to find out
SOURCE_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
SOURCE_COMMIT ?= $(shell git rev-parse --verify HEAD)
SHORT_COMMIT := $(shell v='$(SOURCE_COMMIT)'; echo "$${v::7}")
.ci-helper-git-config-user:
	git config user.name "$(shell git log -n 1 --pretty=format:%an)"
	git config user.email "$(shell git log -n 1 --pretty=format:%ae)"
.ci-helper-gh-bump-version-commit-with-pr:
	git checkout -b bump/v$(NEW_VERSION)
	make bump-version
	git commit -m "AUTOBUMP-$(NEW_VERSION) [skip ci]"
	git push origin bump/v$(NEW_VERSION)
	echo "Create new PR to bump version"
	gh pr create --head bump/v$(NEW_VERSION) --base $(SOURCE_BRANCH) --fill
	echo "bump/v$(NEW_VERSION)" > .ci-helper-gh-bump-version-commit-with-pr
GH_BRANCH ?= $(SOURCE_BRANCH)
.ci-helper-gh-auto-merge-pr-of-branch:
	echo "Approve and Merge the submitted PR automatically"
	# TODO: gh pr review $(GH_BRANCH) --approve
	gh pr merge $(GH_BRANCH) --auto -d -s
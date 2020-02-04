release_version:= v0.4

export GO111MODULE=on

.PHONY: bin
bin:
	go build -o bin/l2i github.com/mhausenblas/l2i

.PHONY: release
release:
	curl -sL https://git.io/goreleaser | bash -s -- --rm-dist --config .goreleaser.yml

.PHONY: publish
publish:
	git tag ${release_version}
	git push --tags
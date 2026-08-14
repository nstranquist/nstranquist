.PHONY: render check test verify

render:
	go run ./tools/render --root .

check:
	go run ./tools/render --root . --check

test:
	go test ./tools/render

verify: test render check

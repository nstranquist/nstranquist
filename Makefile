.PHONY: render check test verify refresh-products

NICO_TOOLS ?= ../nicos-tools

refresh-products:
	go -C $(NICO_TOOLS)/nicos-dev run ./cmd/nico-portfolio-sync sync-products --repo-root .. --portfolio-root ../apps/nico-portfolio --profile-root ../../nstranquist --refresh-github --verify-access

render:
	go run ./tools/render --root .

check:
	go run ./tools/render --root . --check

test:
	go test ./tools/render

verify: test render check

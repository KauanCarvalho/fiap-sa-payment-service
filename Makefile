GO ?= go

.DEFAULT_GOAL := help
.PHONY: help deps setup-git-hooks generate-changelog

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  help            Show this help message"
	@echo "  deps            Install dependencies"
	@echo "  setup-git-hooks Install Git hooks using Lefthook"

deps:
	@echo "Installing dependencies..."
	$(GO) mod download

setup-git-hooks: deps
	@echo "Installing Git hooks with Lefthook..."
	$(GO) tool lefthook install

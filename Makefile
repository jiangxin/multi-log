define message
	@echo "### $(1)"
endef

test: $(TARGET) golint
	$(call message,Testing goconfig using golint for coding style)
	@golint
	$(call message,Testing goconfig for unit tests)
	@go test ./...

golint:
	@if ! type golint >/dev/null 2>&1; then \
		go get golang.org/x/lint/golint; \
	fi

.PHONY: test golint

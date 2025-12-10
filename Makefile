
.PHONY: cover clean

RUNS=100

cover:
	@rm -f coverage.out coverage.html
	@echo "Running tests $(RUNS) times..."
	@go test ./... -count=$(RUNS) -coverprofile=coverage.out >/dev/null 2>&1
	@go tool cover -html=coverage.out -o coverage.html
	@xdg-open coverage.html 2>/dev/null || open coverage.html 2>/dev/null || start coverage.html

clean:
	rm -f coverage.out coverage.html

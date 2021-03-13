.PHONY: examples test

examples:
	go run . -pkg main ./examples/echo/echo.graphql

test:
	go test ./generator/. -count=1

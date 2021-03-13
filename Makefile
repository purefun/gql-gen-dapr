.PHONY: examples test

examples:
	go run . ./examples/echo/echo.graphql

test:
	go test ./generator/. -count=1

.PHONY: examples test

examples:
	go run . -s echo -pkg echo -f ./examples/echo/echo.graphql -o ./examples/echo

test:
	go test ./generator/. -count=1

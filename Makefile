.PHONY: examples

examples:
	go run . -s echo -pkg echo -f ./examples/echo/echo.graphql -o ./examples/echo

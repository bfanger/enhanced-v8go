
dev:
	find . -name "*.go" -o -name "*.js"| entr go run examples/console/console-example.go

test:
	go test ./module/

dev_test:
	gow test ./module/

dev:
	find . -name "*.go" -o -name "*.js"| entr go run examples/svelte/svelte-example.go
	# find . -name "*.go" -o -name "*.js"| entr go run examples/console/console-example.go

test:
	go test ./js/

dev_test:
	gow test -v ./js/
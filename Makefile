cover: 
	go test -v -covermode=count -coverprofile=.cover ./...

html:
	go tool cover -html=.cover

run:
	go run main.go
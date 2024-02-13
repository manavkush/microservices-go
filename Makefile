consul:
	@docker compose up -d

movie:
	go run movie/cmd/main.go

rating:
	go run rating/cmd/main.go

metadata:
	go run metadata/cmd/main.go

hello:
	echo "hello"

.PHONY: movie rating metadata
	# go run movie/cmd/main.go
	# go run rating/cmd/main.go
	# go run metadata/cmd/main.go


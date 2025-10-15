hot-reload:
	air --build.cmd "go build -o bin/backend-api cmd/main.go" --build.bin "./bin/backend-api" --build.exclude_dir "templates,build"
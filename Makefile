package:
	docker buildx build --network=host --platform linux/amd64 -f docker/Dockerfile -t go-static-file-server:latest .

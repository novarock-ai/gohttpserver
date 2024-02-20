package:
	docker buildx build --network=host --platform linux/amd64 -f docker/Dockerfile -t registry.dp.tech/mlops/launching-static-file-server:latest . --push
deploy:
	kubectl rollout restart deploy -n project-launching launching-static-file-server

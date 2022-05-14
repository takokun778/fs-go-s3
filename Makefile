minio:
	@docker run \
		-p 9000:9000 \
		-p 9001:9001 \
		-e "MINIO_ROOT_USER=minio" \
		-e "MINIO_ROOT_PASSWORD=miniominio" \
		--name fs-go-s3 \
		--rm -it -d \
		minio/minio server /data --console-address ":9001"
localstack:
	@docker-compose up -d

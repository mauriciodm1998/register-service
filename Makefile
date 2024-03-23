boiler-plate: test-build-bake docker-push

run-db:
	docker-compose -f build/docker-compose.yaml up -d

run-local-stack:
	docker run --rm -it -p 4566:4566 localstack/localstack

create-queue:
	awslocal kinesis list-streams
	awslocal sqs create-queue --queue-name clock_in_register
	awslocal sqs create-queue --queue-name month_report_requests

test-build-bake:
	docker build -t docker.io/mauricio1998/register-service . -f build/Dockerfile

docker-push:
	docker push docker.io/mauricio1998/register-service
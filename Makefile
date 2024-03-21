run-db:
	docker-compose -f build/docker-compose.yaml up -d

run-local-stack:
	docker run --rm -it -p 4566:4566 localstack/localstack

create-queue:
	awslocal kinesis list-streams
	awslocal sqs create-queue --queue-name clock_in_register
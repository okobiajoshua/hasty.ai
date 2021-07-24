# we will put our integration testing in this path
E2E_TEST_PATH?=./e2e

# this command will start a docker components that we set in docker-compose.yml
docker.start:
	docker-compose up -d --remove-orphans;

# shutting down docker components
docker.stop:
	docker-compose down;

# this command will trigger e2e test
test.e2e:
	cd ./backend; \
	go test $(E2E_TEST_PATH) -count=1 -v; \
	cd ../

.PHONY: build start generate test

test:
	go mod vendor
	docker-compose up --abort-on-container-exit
	docker-compose down
	rm -rf vendor
	
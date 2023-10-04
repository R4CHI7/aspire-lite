.PHONY: test

integration-test:
	@docker compose -f e2e/docker-compose.yml down -v
	@docker compose -f e2e/docker-compose.yml up --build --abort-on-container-exit --remove-orphans --force-recreate
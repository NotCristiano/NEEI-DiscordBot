up:
	docker compose -f docker/docker-compose.yml -f docker/docker-compose.override.yml up --build

down:
	docker compose -f docker/docker-compose.yml -f docker/docker-compose.override.yml down

prod-up:
	docker compose -f docker/docker-compose.yml -f docker/docker-compose.prod.yml up -d

prod-down:
	docker compose -f docker/docker-compose.yml -f docker/docker-compose.prod.yml down

prod-logs:
	docker compose -f docker/docker-compose.yml -f docker/docker-compose.prod.yml logs -f
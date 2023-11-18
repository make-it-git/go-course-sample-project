up:
	docker compose -f docker-compose.yml --project-name common up -d
	docker compose -f driver-service/build/docker-compose.yml --project-name driver up -d
	docker compose -f rider-service/build/docker-compose.yml --project-name rider up -d
	docker compose -f ride-service/build/docker-compose.yml --project-name ride up -d

stop:
	docker compose -f docker-compose.yml --project-name common stop
	docker compose -f driver-service/build/docker-compose.yml --project-name driver stop
	docker compose -f rider-service/build/docker-compose.yml --project-name rider stop
	docker compose -f ride-service/build/docker-compose.yml --project-name ride stop

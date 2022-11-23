up:
	docker-compose build; docker-compose up

down:
	docker-compose down

clean:
	docker rm $(docker ps -a);
	docker volume rm -f $(docker volume ls -q);
	rm -rf ./temp/*

BASEDIR = $(shell pwd)
PORT_FE=8080
PORT_DB=3306
PORT_API=9000

dev: db redis api fe

fe: cleanfe
	cd frontend && docker build -t todo-html .
	docker run --name todo-html --expose $(PORT_FE) -p $(PORT_FE):80 \
	-v $(BASEDIR)/frontend/www:/usr/share/nginx/html todo-html

db: cleandb
	cd database && docker build -t todo-mysql .
	docker run --name todo-mysql -p $(PORT_DB):$(PORT_DB) \
	-e MYSQL_ROOT_PASSWORD=password -e MYSQL_ROOT_HOST=% -d todo-mysql

api: cleanapi
	cd middleware && docker build -t todo-goapi .
	docker run --name todo-goapi --expose $(PORT_API) \
	-p $(PORT_API):$(PORT_API)  -e PORT=$(PORT_API) -e todo_user=root \
	-e todo_pass=password -e todo_host=host.docker.internal -e todo_name=todo  \
	-e REDISPORT=6379 -e REDISHOST=host.docker.internal -d todo-goapi	

cleanfe:
	-docker stop todo-html
	-docker rm todo-html

cleandb:
	-docker stop todo-mysql
	-docker rm todo-mysql

cleanapi:
	-docker stop todo-goapi
	-docker rm todo-goapi		


redis: cleanredis
	docker run --name todo-redis -p 6379:6379 -d redis	

cleanredis:
	-docker stop todo-redis
	-docker rm todo-redis	
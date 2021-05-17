build:
	docker build -t cron-job-vue-go .

run:
	docker run --rm -i -t -p 3000:3000 --dns 8.8.8.8 cron-job-vue-go

all: build run
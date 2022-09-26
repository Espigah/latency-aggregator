run:
	docker run -d -p 7070:7070 fsrg/latency-aggregator:0.0.1
	sleep 5
	go run main.go
	sleep 5
	curl -X GET http://localhost:7070/metrics | grep order_
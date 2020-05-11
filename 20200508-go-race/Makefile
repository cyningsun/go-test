run:
	go run main.go

racerun:
	go run main.go 

benchmark:
	wrk -t12 -c20 -d30s http://127.0.0.1:8080/read &
	wrk -t12 -c20 -d30s http://127.0.0.1:8080/write &
 

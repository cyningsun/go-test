run:
	go run main.go &> normal.log

racerun:
	go run -race main.go &> race.log 

benchmark:
	wrk -t1 -c2 -d30s http://127.0.0.1:8080/read &
	wrk -t1 -c2 -d30s http://127.0.0.1:8080/write &
 

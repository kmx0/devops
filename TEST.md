
RUN APPs for Testing
//run agent 
go run agent/main.go -k "hashkey"

//run server
go run server/main.go -d "postgres://postgres:postgres@localhost:5432/metrics" -i 10s -k "hashkey"


curl -sK -v http://localhost:8080/debug/pprof/heap > heap.out

go tool pprof -http=":9090" -seconds=30 heap.out
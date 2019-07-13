
.PHONY: bench plot pprof

bench:
	go test -bench . | ./bench.py

plot:
	./bench.py --plot

pprof: cpu.prof
	go tool pprof -web cpu.prof

cpu.prof: *.go
	go test -cpuprofile cpu.prof -memprofile mem.prof -bench .

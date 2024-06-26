.DEFAULT_GOAL := run

callee.o: callee.c
	gcc -fPIC -c callee.c

libcallee.so: callee.o
	gcc -shared -o $@ $^

caller: main.go libcallee.so
	go build main.go

clean:
	rm -f caller libcallee.so callee.o

run: caller
	./main
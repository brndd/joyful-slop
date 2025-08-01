DESTDIR=/
PREFIX=${HOME}/
BINDIR=bin/

.PHONY: default build test install uninstall clean

default: build

build:
	go build -o build/ ./...

test:
	go test ./...

install:
	mkdir -p ${DESTDIR}/${PREFIX}/${BINDIR}
	cp build/evinfo ${DESTDIR}/${PREFIX}/${BINDIR}/
	cp build/joyful ${DESTDIR}/${PREFIX}/${BINDIR}/

uninstall:
	rm ${DESTDIR}/${PREFIX}/${BINDIR}/evinfo
	rm ${DESTDIR}/${PREFIX}/${BINDIR}/joyful

clean:
	rm -rf build/
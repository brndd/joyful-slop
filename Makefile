DESTDIR=/
PREFIX=${HOME}/
BINDIR=bin/

default:
	go build -o build/ ./...

install:
	mkdir -p ${DESTDIR}/${PREFIX}/${BINDIR}
	cp build/evinfo ${DESTDIR}/${PREFIX}/${BINDIR}/
	cp build/joyful ${DESTDIR}/${PREFIX}/${BINDIR}/

uninstall:
	rm ${DESTDIR}/${PREFIX}/${BINDIR}/evinfo
	rm ${DESTDIR}/${PREFIX}/${BINDIR}/joyful

clean:
	rm -rf build/
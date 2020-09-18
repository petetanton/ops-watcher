lint:
	golangci-lint run

build:
	for targetos in 'darwin' 'linux'; do \
  		for arch in 'amd64'; do \
  			env GOOS=$$targetos GOARCH=$$arch go build -o ops-watcher-$$targetos-$$arch main.go ; \
  			zip ops-watcher-$$targetos-$$arch.zip ops-watcher-$$targetos-$$arch face.png ; \
		done \
	done


run:
	./ops-watcher-darwin-amd64

install:
	cp ops-watcher-darwin-amd64 /usr/local/bin/ops-watcher

remove:
	rm /usr/local/bin/ops-watcher

clean:
	rm -rf ops-watcher-*

all: clean lint build run

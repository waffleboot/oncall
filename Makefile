.PHONY: build
build:
	go build -o bin/oncall ./cmd/oncall

.PHONY: run
run: build
	rm debug.log
	./bin/oncall

.PHONY: zip
zip:
	rm ../oncall.zip && zip -r ../oncall.zip . -x '.*' 'bin/*' '*.log' '*.json' 'journal.txt'

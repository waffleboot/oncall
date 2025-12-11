.PHONY: build
build:
	go build -o bin/oncall ./cmd/oncall

.PHONY: run
run: build
	./bin/oncall

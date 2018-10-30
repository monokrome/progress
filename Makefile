PREFIX   ?= /usr/local
BIN_PATH ?= bin/

target=$(BIN_PATH)/prg

all: $(BIN_PATH)prg
.PHONY: all


install: all
	mkdir -p $(PREFIX)/bin/
	cp $(BIN_PATH)* $(PREFIX)/bin/
.PHONY: install


clean:
	rm -rf $(target)
.PHONY: clean


$(BIN_PATH)prg: $(wildcard cmd/prg/*.go)
	mkdir -p $(@D)
	go build -o $@ -i $^
	cd $(dir $<) && go install ./...

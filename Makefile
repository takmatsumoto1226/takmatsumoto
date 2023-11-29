build: clean
	CGO_ENABLED=0 go build -ldflags "-w -s" -o ./build/linux/mongodb .

clean:
	@rm -rf ./build
	@rm -rf ./release

.PHONY: clean

## ex : go get -insecure gitlab01.mitake.com.tw/RD1/GO/mitake-tcp.git
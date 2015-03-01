NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
UNAME := $(shell uname -s)
ifeq ($(UNAME),Darwin)
ECHO=echo
else
ECHO=/bin/echo -e
endif

all: 
	@$(ECHO) "$(OK_COLOR)==> Building$(NO_COLOR)"
	go get -v ./...
	go test -v

bin: 
	@$(ECHO) "$(OK_COLOR)==> Building$(NO_COLOR)"
	go build

clean:
	@rm -rf dist/ post-processor-consul

format:
	go fmt ./...

dist:
	@$(ECHO) "$(OK_COLOR)==> Building Packages...$(NO_COLOR)"
	@gox -osarch="darwin/386 darwin/amd64 linux/386 linux/amd64 freebsd/386 freebsd/amd64 openbsd/386 openbsd/amd64 windows/386 windows/amd64 netbsd/386 netbsd/amd64"
	@mv post-processor-consul_darwin_386 post-processor-consul; tar cvfz post-processor-consul.darwin-i386.tar.gz post-processor-consul; rm post-processor-consul
	@mv post-processor-consul_darwin_amd64 post-processor-consul; tar cvfz post-processor-consul.darwin-amd64.tar.gz post-processor-consul; rm post-processor-consul
	@mv post-processor-consul_freebsd_386 post-processor-consul; tar cvfz post-processor-consul.freebsd-i386.tar.gz post-processor-consul; rm post-processor-consul
	@mv post-processor-consul_freebsd_amd64 post-processor-consul; tar cvfz post-processor-consul.freebsd-amd64.tar.gz post-processor-consul; rm post-processor-consul
	@mv post-processor-consul_linux_386 post-processor-consul; tar cvfz post-processor-consul.linux-i386.tar.gz post-processor-consul; rm post-processor-consul
	@mv post-processor-consul_linux_amd64 post-processor-consul; tar cvfz post-processor-consul.linux-amd64.tar.gz post-processor-consul; rm post-processor-consul
	@mv post-processor-consul_netbsd_386 post-processor-consul; tar cvfz post-processor-consul.netbsd-i386.tar.gz post-processor-consul; rm post-processor-consul
	@mv post-processor-consul_netbsd_amd64 post-processor-consul; tar cvfz post-processor-consul.netbsd-amd64.tar.gz post-processor-consul; rm post-processor-consul
	@mv post-processor-consul_openbsd_386 post-processor-consul; tar cvfz post-processor-consul.openbsd-i386.tar.gz post-processor-consul; rm post-processor-consul
	@mv post-processor-consul_openbsd_amd64 post-processor-consul; tar cvfz post-processor-consul.openbsd-amd64.tar.gz post-processor-consul; rm post-processor-consul
	@mv post-processor-consul_windows_386.exe post-processor-consul.exe; zip post-processor-consul.windows-i386.zip post-processor-consul.exe; rm post-processor-consul.exe
	@mv post-processor-consul_windows_amd64.exe post-processor-consul.exe; zip post-processor-consul.windows-amd64.zip post-processor-consul.exe; rm post-processor-consul.exe
	@mkdir -p dist/
	@mv post-processor-consul* dist/.

.PHONY: all clean deps format test updatedeps

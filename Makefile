cmds=test-cbor reflect-cbor test-segmentheader test-headerhash test-handshake test-blockfetch playground

export build = $(abspath ./build)

generator=cbor-type

pkg=$(shell find . -name \*.go)

dsts=$(addprefix $(build)/,$(cmds))

.SECONDEXPANSION:

all: $(generator) $(dsts)

$(generator): $$(shell find ./cmd/cbor-type -name \*.go) | $(build)
	cd $(dir $<); \
	go build -o $(abspath $@)


$(dsts): $$(shell find ./cmd/$$(notdir $$@) -name \*.go) $(pkg) | $(build)
	export CGO_ENABLE=1; \
	go generate; \
	cd $(dir $<); \
	go build -o $(abspath $@) $(build_flags)

$(build):
	mkdir -p $@


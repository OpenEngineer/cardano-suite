generator=cbor-type

.SECONDEXPANSION:

$(generator): $$(shell find ./cmd/cbor-type -name \*.go) | $(build)
	cd $(dir $<); \
	go build -o $(abspath $@)

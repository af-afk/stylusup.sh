
.PHONY: clean

popcon.stylusup.sh: $(shell find . -name '*.go')
	@go build

bootstrap.zip: popcon.stylusup.sh
	@cp popcon.stylusup.sh bootstrap
	@zip bootstrap.zip bootstrap

clean:
	@rm -f bootstrap bootstrap.zip popcon.stylusup.sh

-include .env
export

.PHONY: run
run:
	go run ./cmd/nft-presale/main.go

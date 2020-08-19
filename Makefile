test:
	@TEST_TARGET=dev go test -count=1 -v ./device ./user ./group
	#@go test -count=1 ./device ./auth ./group ./alert ./user ./parser

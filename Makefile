test_feature:
	@TEST_TARGET=feature go test -count=1 ./device ./user ./group ./wave ./alert ./thingplug

test_dev:
	@TEST_TARGET=dev go test -count=1 ./device ./user ./group ./wave ./alert ./thingplug ./parser

test:
	@go test -count=1 ./device ./user ./group ./wave ./alert ./thingplug ./parser

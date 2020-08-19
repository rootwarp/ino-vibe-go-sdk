test_feature:
	@TEST_TARGET=feature go test -count=1 ./device ./user ./group ./wave ./alert

test_def:
	@TEST_TARGET=dev go test -count=1 ./device ./user ./group ./wave ./alert

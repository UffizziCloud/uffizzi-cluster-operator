
.PHONY: list-test-suites
list-test-suites:
	$(eval TEST_SUITES := $(shell find tests/e2e -name Makefile | grep  -E -o '(\w|\-)*' | grep -v 'test' | grep -v 'Makefile' | grep -v 'e2e'))

.PHONE: e2e-test-suites
e2e-test-suites: list-test-suites
	@echo "Test suites to run"
	@for test_suite in $(TEST_SUITES) ; do \
		echo -e "\t $$test_suite" ; \
	done


test-verbose: create-temp ## Run verbose all tests
	go test -v -count=1 -race -coverprofile=tmp/coverage-report-unit.cov -covermode=atomic ./... 

test-unit: create-temp ## run tests	
	go test -count=1 -race -coverprofile=tmp/coverage-report-unit.cov -covermode=atomic --tags=unit ./... 

test-unit-verbose: create-temp ## Run unit tests
	go test -v -count=1 -race -coverprofile=tmp/coverage-report-unit.cov -covermode=atomic --tags=unit ./... 

test: create-temp ## create temp folder 
	go test -count=1 -race -coverprofile=tmp/coverage-report.cov -covermode=atomic  ./... 

test-clear-cache: ## clear test cache 
	go clean -testcache

test-race:
	@for i in `seq 1 100`; \
		do echo ==========================$$i================== >> log_test_running.txt && \
			make test >> log_test_running.txt; \
	done;
	
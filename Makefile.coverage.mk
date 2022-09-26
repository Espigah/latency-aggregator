
create-temp:
	@mkdir -p tmp

coverage-merge:
	@rm -rf tmp/coverage-report.cov
	@find tmp/ -name *.cov | xargs 3rd/gocovmerge/bin/gocovmerge > tmp/coverage-report.cov

coverage-percent: ## Print coverage percent. run test-unit before
	@make coverage-cover --silent | grep total | awk '{print substr($$3, 1, length($$3)-1)}' # $$ only works in makefile

coverage-cover: 
	@go tool cover -func tmp/coverage-report.cov

coverage-show: ## Open coverage report. run test-unit before
	@go tool cover -html=tmp/coverage-report.cov

coverage-to-html: ## Open coverage report. run test-unit before
	@go tool cover -html=tmp/coverage-report.cov -o tmp/coverage-report.html


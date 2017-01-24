design_package := github.com/apiarian/ipfs-pinbase/design

gen: ## the normal workflow goa regenration process
	goagen app -d $(design_package)
	@if [ -a _scaffolds ]; then rm -f _scaffolds/*; fi;
	goagen main -d $(design_package) -o _scaffolds
	goagen client -d $(design_package)
	# goagen js -d $(design_package)
	goagen swagger -d $(design_package)
	bootprint openapi swagger/swagger.json api-doc

clean: ## remove generated stuff except for the main package
	rm -rf app/
	rm -rf client/
	rm -rf schema/
	rm -rf swagger/
	rm -rf js/
	rm -rf _scaffolds/
	rm -rf tool/
	rm -rf api-doc/

gen-bootstrap: ## initial bootstrap (run this first!)
	goagen bootstrap -d $(design_package)

.PHONY: help

prep: ## install required packages
	npm install -g bootprint
	npm install -g bootprint-openapi

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help


# self-documenting makefile:
# http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
# interesting goa-centric makefile:
# https://github.com/kkeuning/cug/blob/master/examples/adder2/Makefile

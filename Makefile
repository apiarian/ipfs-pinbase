api_dir := cmd/ipfs-pinbase
design_package := github.com/apiarian/ipfs-pinbase/$(api_dir)/design

gen: ## the normal workflow goa regenration process
	cd $(api_dir); goagen app -d $(design_package)
	@if [ -a $(api_dir)/_scaffolds ]; then rm -f $(api_dir)/_scaffolds/*; fi;
	cd $(api_dir); goagen main -d $(design_package) -o _scaffolds
	cd $(api_dir); goagen client -d $(design_package)
	# cd $(api_dir); goagen js -d $(design_package)
	cd $(api_dir); goagen swagger -d $(design_package)
	cd $(api_dir); bootprint openapi swagger/swagger.json api-doc

clean: ## remove generated stuff except for the main package
	rm -rf $(api_dir)/app/
	rm -rf $(api_dir)/client/
	rm -rf $(api_dir)/schema/
	rm -rf $(api_dir)/swagger/
	rm -rf $(api_dir)/js/
	rm -rf $(api_dir)/_scaffolds/
	rm -rf $(api_dir)/tool/
	rm -rf $(api_dir)/api-doc/

gen-bootstrap: ## initial bootstrap (run this first!)
	cd $(api_dir); goagen bootstrap -d $(design_package)

prep: ## install required packages
	npm install -g bootprint
	npm install -g bootprint-openapi

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help


# self-documenting makefile:
# http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
# interesting goa-centric makefile:
# https://github.com/kkeuning/cug/blob/master/examples/adder2/Makefile

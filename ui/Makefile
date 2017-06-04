ELMFILES := Main.elm PinCss.elm Edit.elm Party.elm Party/Pin.elm


build-elm-debug: ## build a debug version of the elm app
	elm-make $(ELMFILES) --debug --output build/elm.js

build-elm: ## build a production version of the elm app
	elm-make $(ELMFILES) --output build/elm.js


install-purecss: ## installs pure.css
	npm install purecss --save

copy-all-purecss: ## copies all of the pure.css files into the build directory
	cp -R node_modules/purecss/build build/purecss


install-elm-css: ## installs elm-css helper command
	npm install -g elm-css

build-css:
	elm-css Stylesheets.elm --output build


clean: ## removes junk
	rm -f build/*.js
	rm -rf build/purecss
	rm -f build/pinbase.css


rebuild-debug: clean build-elm-debug build-css copy-all-purecss ## just does the rebuild for debug, without starting the server again

debug: rebuild-debug ## run the debug verion of the app, serving the build directory
	cd build && python -m SimpleHTTPServer 8000


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help


# self-documenting makefile:
# http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

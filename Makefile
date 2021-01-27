help:
	@grep -E '^[a-zA-Z_0-9]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## sudo cleans the tmp folder
	sudo rm -rf /tmp/ruby-gems-crawl/tmp
	mkdir /tmp/ruby-gems-crawl/tmp
	mkdir /tmp/ruby-gems-crawl/tmp/gems
	mkdir /tmp/ruby-gems-crawl/tmp/gems/uncompressed

all: ## runs all steps
	step_1 step_2 step_3 step_4

step_1: ## Fetches the latest version of all ruby gems from rubygems.org and converts them to a csv file for golang to process
	ruby src/entry/latest_spec_processor.rb

step_2: ## Downloads all of the latest ruby gems and unpacks them
	go run ./src/entry/gem_download.go

step_3: ## Corrects read and write permissions of all of the files
	find /tmp/ruby-gems-crawl/tmp/gems -type d -exec chmod u=rx,go=rx {} \;
	find /tmp/ruby-gems-crawl/tmp/gems -type f -exec chmod u=r,go=r {} \;

step_4: ## Scan the files for vulnerabilities
	go run ./src/entry/gem_scan.go scan /tmp/ruby-gems-crawl/tmp/gems
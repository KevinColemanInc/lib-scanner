![lib-scanner logo](./assets/logo.png)
# lib-scanner

Downloads the latest ruby gems and scans for malicious code!

## Useful for
- Checking if a ruby gem has malicious code.
- Fast auditing of ruby gems for malicious code that may run on install or malicious

## Requirements

- ruby 2.7+
- golang 1.15+

## Quick Start
Run `make help` gives you a quick overview of how this project works.
```
$ make help
all                            runs all steps
clean                          sudo cleans the tmp folder
step_1                         Fetches the latest version of all ruby gems from rubygems.org and converts them to a csv file for golang to process
step_2                         Downloads all of the latest ruby gems and unpacks them
step_3                         Corrects read and write permissions of all of the files
step_4                         Scan the files for vulnerabilities
```

### Step 1 - Download latest gem metadata

https://rubygems.org/latest_specs.4.8.gz is a file containing Marshalled ruby objects, I wrote a [ruby script](/master/script/entry/latest_spec_processor.rb) to fetch the file and convert the gems into a csv file with their latest versions.

### Step 2 - Downloads all of the latest ruby gems and unpacks them

Golang script to download and unpack latest gems. This takes a long time (+2 hours) depending on your CPU and network connection

### Step 3 - Corrects read and write permissions of all of the files

Files in ruby gems have weird permissions (like no read access). This forcefully makes all files read-only. Because there are ~20m files, it takes a long time to run this.

### Step 4 - Scan the files for vulnerabilities

This runs regexes against all `.rb` files to identify vulnerabilities and writes the results to STOUT.
I recommend appending a `> results.csv` to the cmd to save them to a file. `.csv` is used because you can use [q](http://harelba.github.io/q/) or dump into a database for fast and grouping of results. I use the `≫` to delimit the columns because the last column prints out the entire line of ruby code the regex matched against. `≫` is a character that safely handles splitting the csv file.

## Scan your ruby-gems folder

You can use this tool to scan the gems that you are actively using in your projects instead of scanning the latest versions of everything.

```
~/project-dir $ bundle show --paths
```
Make a note of where `bundler` installs your gems

Then scan the directory
```
$ go run ./src/entry/gem_scan.go scan /Users/kevin/.rvm/gems/ruby-2.5.7/gems > results.csv
➜  lib-crawl git:(fcf5996) ✗ go run ./src/entry/gem_scan.go scan /Users/kevin/.rvm/gems/ruby-2.5.7/gems
Started!
Found paths: 42868
has http≫/Users/kevin/.rvm/gems/ruby-2.5.7/gems/autoprefixer-rails-9.1.0/spec/compass/config.rb≫http_path       = "/"
has http≫/Users/kevin/.rvm/gems/ruby-2.5.7/gems/aws-sdk-core-3.100.0/lib/seahorse/client/net_http/patches.rb≫require 'net/http'
has http≫/Users/kevin/.rvm/gems/ruby-2.5.7/gems/aws-sdk-core-3.100.0/lib/seahorse/client/net_http/handler.rb≫require 'net/https'
```

## Launch Blog

Long-form article about the how and why
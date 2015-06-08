require File.expand_path('../../base/base.rb', __FILE__)

class HelloWorld < Base
    def run
        puts "Hello World"
    end
end

if __FILE__ == $0
    HelloWorld.new.run
end
require 'socket'
require 'json'

class Base
    attr_accessor :dry_run
    attr_accessor :hostname
    attr_accessor :data

    def initialize
        @dry_run = true
        @hostname = Socket.gethostname
        @data = {}
        _read_data
    end

    def current_dir
        File.expand_path(File.dirname(__FILE__))
    end

    def init; end
    def run; end

    protected
    def _read_data
        data_dir = File.join(current_dir, 'data')
        unless File.exist?(data_dir)
            return
        end

        Dir.entries(data_dir).each do |filename|
            puts filename
            if filename.downcase.start_with?('readme') || filename == '.' || filename == '..'
                next
            end

            full_filename = File.join(data_dir, filename)
            if full_filename.end_with?('.json')
                @data[full_filename] = JSON.parse(File.read(full_filename))
            else
                @data[full_filename] = File.read(full_filename)
            end
        end
    end
end
class Base
    def initialize
        @dry_run = true
        @data = {}
        _read_data
    end

    def init; end
    def run; end

    protected
    def _read_data
    end
end
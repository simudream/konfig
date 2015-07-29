This directory serves as a complete example of resourced-stacks project.

## Testing commands
```
# Given that hostname is didip.mac-mini.local,

# Example on how to match exact hostname.
go run main.go -root=./tests/project -cmd=run -stack=helloworld -conditions='name == "didip-mac-mini.local"'

# Example on how to match hostname using regex.
go run main.go -root=./tests/project -cmd=run -stack=helloworld -conditions='name.match(/^didip/i)'
```

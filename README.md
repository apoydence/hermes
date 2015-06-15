# hermes
A P2P over HTTP API

# Usage

### Upload File
```
curl -v localhost:8080 -X POST -H "key:some-key" --data-binary "@/tmp/example.file"
````

### Download File
```
curl -v localhost:8080 --cookie "key=some-key" -o output.file
```


# for mac
GOOS=darwin GOARCH=amd64 go build -o bin/qurandw-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o bin/qurandw-darwin-arm64

# for linux
GOOS=linux GOARCH=amd64 go build -o bin/qurandw-linux-amd64
GOOS=linux GOARCH=386 go build -o bin/qurandw-linux-x86

# for windows
GOOS=windows GOARCH=amd64 go build -o bin/qurandw-win32-amd64.exe 
GOOS=windows GOARCH=386 go build -o bin/qurandw-win32-x86.exe

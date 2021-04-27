cleanup: tidy fmt
fmt:; @find pkg internal -name \*.go -exec go fmt {} \;   
tidy:; @go mod tidy

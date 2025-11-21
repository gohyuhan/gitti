package constant

const (
	APPNAME = "gitti"
)

// this will be injected during build
// exmaple) go build -ldflags "-X gitti/constant.APPVERSION=v1.x.x" -o main
var APPVERSION = "v0.1.0"

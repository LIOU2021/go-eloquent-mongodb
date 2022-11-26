package orm

import (
	"fmt"
	"os"
	"runtime"
)

func getUri() (uri string) {
	user := os.Getenv("mongodb_user")
	password := os.Getenv("mongodb_password")
	host := os.Getenv("mongodb_host")
	port := os.Getenv("mongodb_port")

	domain := ""
	if user == "" {
		domain = fmt.Sprintf("%s:%s", host, port)
	} else {
		domain = fmt.Sprintf("%s:%s@%s:%s", user, password, host, port)
	}

	uri = fmt.Sprintf("mongodb://%s/?retryWrites=true&w=majority", domain)
	return
}

func getLogTitle(collection string) string {
	return fmt.Sprintf("[collection - %s] : ", collection)
}

func getCurrentFuncInfo() string {
	pc, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("\nPC:%s\nFILE:%s\nLINE:%d", runtime.FuncForPC(pc).Name(), file, line)
}

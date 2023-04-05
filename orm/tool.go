package orm

import (
	"errors"
	"fmt"
	"runtime"
)

func getUri() (uri string) {
	user := conf.User
	password := conf.Password
	host := conf.Host
	port := conf.Port

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

func getCurrentFuncInfo(skip int) string {
	pc, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("\nPC:%s\nFILE:%s\nLINE:%d\n", runtime.FuncForPC(pc).Name(), file, line)
}

func (e *Eloquent[T]) errMsg(msg ...any) (err error) {
	concatMsg := fmt.Sprintln(msg...)
	message := fmt.Sprintln(e.logTitle, concatMsg, getCurrentFuncInfo(2))
	err = errors.New(message)
	return
}

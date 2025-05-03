package dd

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"strings"
)

func Dump(vars ...interface{}) {
	file, line := caller()
	fmt.Printf("\n── dump @ %s:%d ──\n", file, line)
	for i, v := range vars {
		fmt.Printf("#%d ⇒ ", i+1)
		spew.Config.Dump(v)
	}
	fmt.Println("──────────────────")
}

func D(vars ...interface{}) {
	Dump(vars...)
	panic("execution stopped by dd.D()")
}

func Web(ctx *gin.Context, vars ...interface{}) {
	if ctx == nil {
		D(vars...)
		return
	}
	Dump(vars...)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"debug": "execution stopped by dd.Web()",
	})
}

func caller() (string, int) {
	for skip := 2; ; skip++ {
		if pc, file, line, ok := runtime.Caller(skip); ok {
			fn := runtime.FuncForPC(pc)
			if !strings.Contains(fn.Name(), "pkg/dd.") {
				if i := strings.LastIndex(file, "/"); i >= 0 {
					file = file[i+1:]
				}
				return file, line
			}
		}
		return "?", 0
	}
}

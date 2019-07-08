package common

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	"github.com/apache/dubbo-go/common/logger"
)

// checkArgs check concurrency and total request count.
func CheckArgs(c, n int) (int, int, error) {
	if c < 1 {
		fmt.Printf("c < 1 and reset c = 1")
		c = 1
	}
	if n < 1 {
		fmt.Printf("n < 1 and reset n = 1")
		n = 1
	}
	if c > n {
		return c, n, errors.New("c must be set <= n")
	}
	return c, n, nil
}

// pprof
func InitProfiling(port string) {
	go func() {
		logger.Info(http.ListenAndServe(":"+port, nil))
	}()
}

func GetString(size int) string {
	argBuf := new(bytes.Buffer)
	for i := 0; i < size; i++ {
		// size: 300
		argBuf.WriteString("击鼓其镗，踊跃用兵。土国城漕，我独南行。从孙子仲，平陈与宋。不我以归，忧心有忡。爰居爰处？爰丧其马？于以求之？于林之下。死生契阔，与子成说。执子之手，与子偕老。于嗟阔兮，不我活兮。于嗟洵兮，不我信兮。")
	}
	return argBuf.String()
}

func InitSignal(survivalTimeout int) {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP,
		syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			go time.AfterFunc(time.Duration(survivalTimeout), func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})

			fmt.Println("app exit now...")
			return
		}
	}
}

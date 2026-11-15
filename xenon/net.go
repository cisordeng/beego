package xenon

import (
	"fmt"
	"net"
	"time"
)

func GetAvailablePort() int {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		PanicNotNilError(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()
	PanicNotNilError(err)
	return port
}

func CheckAvailablePort(port int) bool {
	bindChannel := make(chan bool, 1)
	timeout := time.NewTimer(time.Second)
	expired := false
	defer timeout.Stop()

	go func() {
		checkPortCount := 5
		checkPort: conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), time.Millisecond * 50)
		if expired {
			return
		}
		if err == nil {
			if conn != nil {
				_ = conn.Close()
			}
			if checkPortCount > 0 {
				checkPortCount = checkPortCount - 1
				time.Sleep(time.Millisecond * 50)
				goto checkPort
			}
			bindChannel <- true
			close(bindChannel)
			return
		}
		bindChannel <- false
		close(bindChannel)
		return
	}()

	select {
	case bind := <- bindChannel:
		return bind
	case <- timeout.C:
		expired = true
		return false
	}
}
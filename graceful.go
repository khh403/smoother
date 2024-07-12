package smoother

//smoother listeners and connections allow graceful
//restarts by tracking when all connections from a listener
//have been closed

import (
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func newSmootherListener(l net.Listener) *smootherListener {
	return &smootherListener{
		Listener:     l,
		closeByForce: make(chan bool),
	}
}

// gracefully closing net.Listener
type smootherListener struct {
	net.Listener
	closeError   error
	closeByForce chan bool
	wg           sync.WaitGroup
}

func (l *smootherListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.(*net.TCPListener).AcceptTCP()
	if err != nil {
		return nil, err
	}
	conn.SetKeepAlive(true)                  // see http.tcpKeepAliveListener
	conn.SetKeepAlivePeriod(3 * time.Minute) // see http.tcpKeepAliveListener
	uconn := smootherConn{
		Conn:   conn,
		wg:     &l.wg,
		closed: make(chan bool),
	}
	go func() {
		//connection watcher
		select {
		case <-l.closeByForce:
			uconn.Close()
		case <-uconn.closed:
			//closed manually
		}
	}()
	l.wg.Add(1)
	return uconn, nil
}

// non-blocking trigger close
func (l *smootherListener) release(timeout time.Duration) {
	//stop accepting connections - release fd
	l.closeError = l.Listener.Close()
	if l.closeError != nil {
		l.debugf("[2] close smootherListener.net.Listener failed, errmsg: %+v", l.closeError.Error())
	} else {
		l.debugf("[2] close smootherListener.net.Listener success")
	}
	//start timer, close by force if deadline not met
	waited := make(chan bool)
	go func() {
		l.wg.Wait()
		l.debugf("listener's all connections closed!\n")
		waited <- true
	}()
	go func() {
		select {
		case <-time.After(timeout):
			close(l.closeByForce)
		case <-waited:
			//no need to force close
			l.debugf("no need to force close!\n")
		}
	}()
}

// blocking wait for close
func (l *smootherListener) Close() error {
	l.debugf("[1] close net.Listener, and wait all connections close")
	l.wg.Wait()
	l.debugf("[5] close net.Listener, and all connections closed!!!!")
	return l.closeError
}

func (l *smootherListener) File() *os.File {
	// returns a dup(2) - FD_CLOEXEC flag *not* set
	tl := l.Listener.(*net.TCPListener)
	fl, _ := tl.File()
	return fl
}

// notifying on close net.Conn
type smootherConn struct {
	net.Conn
	wg     *sync.WaitGroup
	closed chan bool
}

func (o smootherConn) Close() error {
	err := o.Conn.Close()
	if err == nil {
		o.wg.Done()
		o.closed <- true
	}
	o.debugf("connections closed")
	return err
}

func (l *smootherListener) debugf(f string, args ...interface{}) {
	log.Printf("[smoother slave] "+f, args...)
}

func (o *smootherConn) debugf(f string, args ...interface{}) {
	log.Printf("[smoother slave] "+f, args...)
}

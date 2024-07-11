package smoother

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
<<<<<<< HEAD
=======
	"errors"
>>>>>>> develop
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var tmpBinPath = filepath.Join(os.TempDir(), "smoother-"+token()+extension())

// a smoother master process
type master struct {
	*Config
	slaveID             int
	slaveCmd            *exec.Cmd
	slaveExtraFiles     []*os.File
	binPath, tmpBinPath string
	binPerms            os.FileMode
	binHash             []byte
	restartMux          sync.Mutex
	restarting          bool
	restartedAt         time.Time
	restarted           chan bool
	awaitingUSR1        bool
	descriptorsReleased chan bool
	signalledAt         time.Time
	printCheckUpdate    bool
}

// run 是主进程执行的方法，监听 USR2 信号，并一直在 fork loop 中死循环
func (mp *master) run() error {
<<<<<<< HEAD
	mp.debugf("run")
=======
	mp.debugf("process master run")
>>>>>>> develop
	if err := mp.checkBinary(); err != nil {
		return err
	}
	if mp.Config.Fetcher != nil {
		if err := mp.Config.Fetcher.Init(); err != nil {
			mp.warnf("fetcher init failed (%s). fetcher disabled.", err)
			mp.Config.Fetcher = nil
		}
	}
<<<<<<< HEAD
	mp.setupSignalling()
	if err := mp.retreiveFileDescriptors(); err != nil {
		return err
	}
=======
	// 设置 master 进程需要相应的信号处理函数
	mp.setupSignalling()

	if err := mp.retreiveFileDescriptors(); err != nil {
		return err
	}
	mp.debugf("master process retrive file descriptors success, count: %d, detail: %+v",
		len(mp.slaveExtraFiles), mp.slaveExtraFiles)

>>>>>>> develop
	if mp.Config.Fetcher != nil {
		mp.printCheckUpdate = true
		mp.fetch()
		go mp.fetchLoop()
	}
<<<<<<< HEAD
=======

>>>>>>> develop
	return mp.forkLoop()
}

func (mp *master) checkBinary() error {
	//get path to binary and confirm its writable
	binPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find binary path (%s)", err)
	}
	mp.binPath = binPath
<<<<<<< HEAD
=======
	mp.debugf("check bin path: %s", mp.binPath)
>>>>>>> develop
	if info, err := os.Stat(binPath); err != nil {
		return fmt.Errorf("failed to stat binary (%s)", err)
	} else if info.Size() == 0 {
		return fmt.Errorf("binary file is empty")
	} else {
		//copy permissions
		mp.binPerms = info.Mode()
	}
	f, err := os.Open(binPath)
	if err != nil {
		return fmt.Errorf("cannot read binary (%s)", err)
	}
	//initial hash of file
	hash := sha1.New()
	io.Copy(hash, f)
	mp.binHash = hash.Sum(nil)
	f.Close()
	//test bin<->tmpbin moves
	if mp.Config.Fetcher != nil {
		if err := move(tmpBinPath, mp.binPath); err != nil {
			return fmt.Errorf("cannot move binary (%s)", err)
		}
		if err := move(mp.binPath, tmpBinPath); err != nil {
			return fmt.Errorf("cannot move binary back (%s)", err)
		}
	}
	return nil
}

func (mp *master) setupSignalling() {
<<<<<<< HEAD
=======
	//process master 会接受所有系统的信号，并一一处理他们
	//setup master process's signals
>>>>>>> develop
	//updater-forker comms
	mp.restarted = make(chan bool)
	mp.descriptorsReleased = make(chan bool)
	//read all master process signals
	signals := make(chan os.Signal)
	signal.Notify(signals)
	go func() {
		for s := range signals {
			mp.handleSignal(s)
		}
	}()
}

func (mp *master) handleSignal(s os.Signal) {
	if s == mp.RestartSignal {
<<<<<<< HEAD
		//user initiated manual restart  USR2 信号
=======
		//user initiated manual restart
		//这是用户手动给主进程发送的信号，一般默认为 USR2 信号，比如 kill -USR2 <master process id>
>>>>>>> develop
		go mp.triggerRestart()
	} else if s.String() == "child exited" {
		// will occur on every restart, ignore it
	} else
	//**during a restart** a SIGUSR1 signals
	//to the master process that, the file
	//descriptors have been released
	if mp.awaitingUSR1 && s == SIGUSR1 {
<<<<<<< HEAD
		mp.debugf("signaled, sockets ready")
=======
		mp.debugf("SIGUSR1 signaled, sockets ready")
>>>>>>> develop
		mp.awaitingUSR1 = false
		mp.descriptorsReleased <- true
	} else
	//while the slave process is running, proxy
	//all signals through
	if mp.slaveCmd != nil && mp.slaveCmd.Process != nil {
<<<<<<< HEAD
		mp.debugf("proxy signal (%s)", s)
=======
		//mp.debugf("proxy signal (%s)", s)
>>>>>>> develop
		mp.sendSignal(s)
	} else
	//otherwise if not running, kill on CTRL+c
	if s == os.Interrupt {
		mp.debugf("interupt with no slave")
		os.Exit(1)
	} else {
		mp.debugf("signal discarded (%s), no slave process", s)
	}
}

func (mp *master) sendSignal(s os.Signal) {
	if mp.slaveCmd != nil && mp.slaveCmd.Process != nil {
		if err := mp.slaveCmd.Process.Signal(s); err != nil {
			mp.debugf("signal failed (%s), assuming slave process died unexpectedly", err)
			os.Exit(1)
		}
	}
}

func (mp *master) retreiveFileDescriptors() error {
	mp.slaveExtraFiles = make([]*os.File, len(mp.Config.Addresses))
	for i, addr := range mp.Config.Addresses {
		a, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
<<<<<<< HEAD
			return fmt.Errorf("Invalid address %s (%s)", addr, err)
=======
			return fmt.Errorf("invalid address %s (%s)", addr, err)
>>>>>>> develop
		}
		l, err := net.ListenTCP("tcp", a)
		if err != nil {
			return err
		}
		f, err := l.File()
		if err != nil {
<<<<<<< HEAD
			return fmt.Errorf("Failed to retreive fd for: %s (%s)", addr, err)
		}
		if err := l.Close(); err != nil {
			return fmt.Errorf("Failed to close listener for: %s (%s)", addr, err)
=======
			return fmt.Errorf("failed to retreive fd for: %s (%s)", addr, err)
		}
		if err := l.Close(); err != nil {
			return fmt.Errorf("failed to close listener for: %s (%s)", addr, err)
>>>>>>> develop
		}
		mp.slaveExtraFiles[i] = f
	}
	return nil
}

// fetchLoop is run in a goroutine
func (mp *master) fetchLoop() {
	min := mp.Config.MinFetchInterval
	time.Sleep(min)
	for {
		t0 := time.Now()
		mp.fetch()
		//duration fetch of fetch
		diff := time.Now().Sub(t0)
		if diff < min {
			delay := min - diff
			//ensures at least MinFetchInterval delay.
			//should be throttled by the fetcher!
			time.Sleep(delay)
		}
	}
}

func (mp *master) fetch() {
	if mp.restarting {
		return //skip if restarting
	}
	if mp.printCheckUpdate {
<<<<<<< HEAD
		mp.debugf("checking for updates...")
=======
		//mp.debugf("checking for updates...")
>>>>>>> develop
	}
	reader, err := mp.Fetcher.Fetch()
	if err != nil {
		mp.debugf("failed to get latest version: %s", err)
		return
	}
	if reader == nil {
		if mp.printCheckUpdate {
<<<<<<< HEAD
			mp.debugf("no updates")
=======
			//mp.debugf("no updates")
>>>>>>> develop
		}
		mp.printCheckUpdate = false
		return //fetcher has explicitly said there are no updates
	}
	mp.printCheckUpdate = true
<<<<<<< HEAD
	mp.debugf("streaming update...")
=======
	//mp.debugf("streaming update...")
>>>>>>> develop
	//optional closer
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}
	tmpBin, err := os.OpenFile(tmpBinPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		mp.warnf("failed to open temp binary: %s", err)
		return
	}
	defer func() {
		tmpBin.Close()
		os.Remove(tmpBinPath)
	}()
	//tee off to sha1
	hash := sha1.New()
	reader = io.TeeReader(reader, hash)
	//write to a temp file
	_, err = io.Copy(tmpBin, reader)
	if err != nil {
		mp.warnf("failed to write temp binary: %s", err)
		return
	}
	//compare hash
	newHash := hash.Sum(nil)
	if bytes.Equal(mp.binHash, newHash) {
<<<<<<< HEAD
		mp.debugf("hash match - skip")
=======
		//mp.debugf("hash match - skip")
>>>>>>> develop
		return
	}
	//copy permissions
	if err := chmod(tmpBin, mp.binPerms); err != nil {
		mp.warnf("failed to make temp binary executable: %s", err)
		return
	}
	if err := chown(tmpBin, uid, gid); err != nil {
		mp.warnf("failed to change owner of binary: %s", err)
		return
	}
	if _, err := tmpBin.Stat(); err != nil {
		mp.warnf("failed to stat temp binary: %s", err)
		return
	}
	tmpBin.Close()
	if _, err := os.Stat(tmpBinPath); err != nil {
		mp.warnf("failed to stat temp binary by path: %s", err)
		return
	}
	if mp.Config.PreUpgrade != nil {
		if err := mp.Config.PreUpgrade(tmpBinPath); err != nil {
			mp.warnf("user cancelled upgrade: %s", err)
			return
		}
	}
<<<<<<< HEAD
	//smoother sanity check, dont replace our good binary with a non-executable file
=======
	//smoother sanity check, do not replace our good binary with a non-executable file
>>>>>>> develop
	tokenIn := token()
	cmd := exec.Command(tmpBinPath)
	cmd.Env = append(os.Environ(), []string{envBinCheck + "=" + tokenIn}...)
	cmd.Args = os.Args
	returned := false
	go func() {
		time.Sleep(5 * time.Second)
		if !returned {
			mp.warnf("sanity check against fetched executable timed-out, check smoother is running")
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		}
	}()
	tokenOut, err := cmd.CombinedOutput()
	returned = true
	if err != nil {
		mp.warnf("failed to run temp binary: %s (%s) output \"%s\"", err, tmpBinPath, tokenOut)
		return
	}
	if tokenIn != string(tokenOut) {
		mp.warnf("sanity check failed")
		return
	}
	//overwrite!
	if err := overwrite(mp.binPath, tmpBinPath); err != nil {
		mp.warnf("failed to overwrite binary: %s", err)
		return
	}
	mp.debugf("upgraded binary (%x -> %x)", mp.binHash[:12], newHash[:12])
	mp.binHash = newHash
	//binary successfully replaced
	if !mp.Config.NoRestartAfterFetch {
		mp.triggerRestart()
	}
	//and keep fetching...
	return
}

func (mp *master) triggerRestart() {
<<<<<<< HEAD
	if mp.restarting {
=======
	// 父进程响应用户 kill -USR2 重启信号
	if mp.restarting {
		// 如果已经完成重启动作，name就忽略
>>>>>>> develop
		mp.debugf("already graceful restarting")
		return //skip
	} else if mp.slaveCmd == nil || mp.restarting {
		mp.debugf("no slave process")
		return //skip
	}
<<<<<<< HEAD
	mp.debugf("graceful restart triggered")
	mp.restarting = true
=======

	mp.debugf("graceful restart triggered")
	mp.restarting = true
	// master 主进程执行 restart 动作后要等待子进程会发的 USR1 信号!!!
>>>>>>> develop
	mp.awaitingUSR1 = true
	mp.signalledAt = time.Now()
	mp.sendSignal(mp.Config.RestartSignal) //ask nicely to terminate
	select {
	case <-mp.restarted:
		//success
		mp.debugf("restart success")
	case <-time.After(mp.TerminateTimeout):
		//times up mr. process, we did ask nicely!
		mp.debugf("graceful timeout, forcing exit")
		mp.sendSignal(os.Kill)
	}
}

// not a real fork
func (mp *master) forkLoop() error {
	//loop, restart command
	for {
<<<<<<< HEAD
		if err := mp.fork(); err != nil {
			return err
		}
=======
		// 这是一个死循环，一般正常工作下有一个master进程和一个slave在执行，但是当进程重启时，会创建一个新的slave进程，这样子就会有两个slave进程在执行了, 等老进程所有请求都处理完释放连接后，父进程会关闭老进程，
		// 正常情况下父进程会卡在这里，直到重启!!!
		//
		mp.debugf("fork begin")
		if err := mp.fork(); err != nil {
			return err
		}
		mp.debugf("fork finish")
>>>>>>> develop
	}
}

func (mp *master) fork() error {
<<<<<<< HEAD
	mp.debugf("starting %s", mp.binPath)
=======
	mp.debugf("fork new slave process, starting %s", mp.binPath)
>>>>>>> develop
	cmd := exec.Command(mp.binPath)
	//mark this new process as the "active" slave process.
	//this process is assumed to be holding the socket files.
	mp.slaveCmd = cmd
	mp.slaveID++
	//provide the slave process with some state
	e := os.Environ()
	e = append(e, envBinID+"="+hex.EncodeToString(mp.binHash))
	e = append(e, envBinPath+"="+mp.binPath)
	e = append(e, envSlaveID+"="+strconv.Itoa(mp.slaveID))
	e = append(e, envIsSlave+"=1")
	e = append(e, envNumFDs+"="+strconv.Itoa(len(mp.slaveExtraFiles)))
	cmd.Env = e
	//inherit master args/stdfiles
	cmd.Args = os.Args
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//include socket files
	cmd.ExtraFiles = mp.slaveExtraFiles
<<<<<<< HEAD
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start slave process: %s", err)
=======
	//os/exec exec.Command(mp.binPath).Start() 返回一个 *exec.Cmd 对象，然后调用 cmd.Start() 启动进程，
	//然后返回一个 err 值，如果 err 值不为 nil，则说明启动进程失败，否则说明启动进程成功，
	//然后通过 cmd.Wait() 方法等待进程退出，然后返回一个 err 值，如果 err 值不为 nil，则说明进程退出失败，否则说明进程退出成功，
	//然后通过 cmd.Process.Kill() 方法杀死进程，然后返回一个 err 值，如果 err 值不为 nil，则说明杀死进程失败，否则说明杀死进程成功，
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start slave process: %s", err)
>>>>>>> develop
	}
	//was scheduled to restart, notify success
	if mp.restarting {
		mp.restartedAt = time.Now()
		mp.restarting = false
		mp.restarted <- true
	}
	//convert wait into channel
<<<<<<< HEAD
	cmdwait := make(chan error)
	go func() {
		cmdwait <- cmd.Wait()
	}()
	//wait....
	select {
	case err := <-cmdwait:
		//program exited before releasing descriptors
		//proxy exit code out to master
		code := 0
		if err != nil {
			code = 1
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
=======
	cmdWait := make(chan error)
	go func() {
		cmdWait <- cmd.Wait()
	}()
	//wait....
	select {
	case err := <-cmdWait:
		//program exited before releasing descriptors
		//proxy exit code out to master
		code := 0
		if err != nil { // err 对应上面 cmd.Start() 的返回值
			code = 1
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
>>>>>>> develop
					code = status.ExitStatus()
				}
			}
		}
		mp.debugf("prog exited with %d", code)
		//if a restarts are disabled or if it was an
		//unexpected crash, proxy this exit straight
		//through to the main process
		if mp.NoRestart || !mp.restarting {
			os.Exit(code)
		}
	case <-mp.descriptorsReleased:
		//if descriptors are released, the program
		//has yielded control of its sockets and
		//a parallel instance of the program can be
		//started safely. it should serve state.Listeners
		//to ensure downtime is kept at <1sec. The previous
		//cmd.Wait() will still be consumed though the
		//result will be discarded.
<<<<<<< HEAD
=======
		mp.debugf("all connections released, slave process quit normally")
>>>>>>> develop
	}
	return nil
}

func (mp *master) debugf(f string, args ...interface{}) {
	if mp.Config.Debug {
		log.Printf("[smoother master] "+f, args...)
	}
}

func (mp *master) warnf(f string, args ...interface{}) {
	if mp.Config.Debug || !mp.Config.NoWarn {
		log.Printf("[smoother master] "+f, args...)
	}
}

func token() string {
	buff := make([]byte, 8)
	rand.Read(buff)
	return hex.EncodeToString(buff)
}

// On Windows, include the .exe extension, noop otherwise.
func extension() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}

	return ""
}

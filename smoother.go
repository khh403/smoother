// Package smoother implements daemonizable
// self-upgrading binaries in Go (golang).
package smoother

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/khh403/smoother/fetcher"
)

const (
<<<<<<< HEAD
	envSlaveID        = "OVERSEER_SLAVE_ID"
	envIsSlave        = "OVERSEER_IS_SLAVE"
	envNumFDs         = "OVERSEER_NUM_FDS"
	envBinID          = "OVERSEER_BIN_ID"
	envBinPath        = "OVERSEER_BIN_PATH"
	envBinCheck       = "OVERSEER_BIN_CHECK"
	envBinCheckLegacy = "GO_UPGRADE_BIN_CHECK"
=======
	envSlaveID  = "SMOOTHER_SLAVE_ID"  // specify slave process id
	envIsSlave  = "SMOOTHER_IS_SLAVE"  // check current process is master or slave
	envNumFDs   = "SMOOTHER_NUM_FDS"   // number of fds
	envBinID    = "SMOOTHER_BIN_ID"    // binary file id
	envBinPath  = "SMOOTHER_BIN_PATH"  // binary file path
	envBinCheck = "SMOOTHER_BIN_CHECK" // binary file check
>>>>>>> develop
)

// Config defines smoother's run-time configuration
type Config struct {
<<<<<<< HEAD
	//Required will prevent smoother from fallback to running
	//running the program in the main process on failure.
	Required bool
	//Program's main function
	Program func(state State)
	//Program's zero-downtime socket listening address (set this or Addresses)
	Address string
	//Program's zero-downtime socket listening addresses (set this or Address)
	Addresses []string
	//RestartSignal will manually trigger a graceful restart. Defaults to SIGUSR2.
=======
	//Required will prevent smoother from fallback to
	//running the program in the main process on failure.
	//如果设置为 true，在 smoother 失败时不会回退到主进程运行程序，确保程序始终在 smoother 的控制下运行
	//一般在使用中，要通过 smoother 来管理程序，这个设置为 false
	Required bool
	//Program's main function
	//应用程序的主入口点函数，smoother 将管理这个函数
	Program func(state State)
	//Program's zero-downtime socket listening address (set this or Addresses)
	//程序的零停机时间 socket 监听地址（设置此字段或 Addresses）
	Address string
	//Program's zero-downtime socket listening addresses (set this or Address)
	//程序的零停机时间 socket 监听地址列表（设置此字段或 Address）
	Addresses []string
	//RestartSignal will manually trigger a graceful restart. Defaults to SIGUSR2.
	//手动触发优雅重启的信号。默认是 SIGUSR2，支持手动指定不同的信号
>>>>>>> develop
	RestartSignal os.Signal
	//TerminateTimeout controls how long smoother should
	//wait for the program to terminate itself. After this
	//timeout, smoother will issue a SIGKILL.
<<<<<<< HEAD
=======
	//控制 smoother 等待程序自行终止的时间。在此超时后，smoother 将发送 SIGKILL 信号
>>>>>>> develop
	TerminateTimeout time.Duration
	//MinFetchInterval defines the smallest duration between Fetch()s.
	//This helps to prevent unwieldy fetch.Interfaces from hogging
	//too many resources. Defaults to 1 second.
<<<<<<< HEAD
	MinFetchInterval time.Duration
	//PreUpgrade runs after a binary has been retrieved, user defined checks
	//can be run here and returning an error will cancel the upgrade.
	PreUpgrade func(tempBinaryPath string) error
	//Debug enables all [smoother] logs.
	Debug bool
	//NoWarn disables warning [smoother] logs.
	NoWarn bool
	//NoRestart disables all restarts, this option essentially converts
	//the RestartSignal into a "ShutdownSignal".
	NoRestart bool
	//NoRestartAfterFetch disables automatic restarts after each upgrade.
	//Though manual restarts using the RestartSignal can still be performed.
	NoRestartAfterFetch bool
	//Fetcher will be used to fetch binaries.
=======
	//定义两次 Fetch() 之间的最小时间间隔。默认是 1 秒,调整这个值以控制更新拉取的频率，防止资源占用过多
	MinFetchInterval time.Duration
	//PreUpgrade runs after a binary has been retrieved, user defined checks
	//can be run here and returning an error will cancel the upgrade.
	//在升级前进行自定义检查或验证逻辑，以确保新的二进制文件适合使用,在二进制文件被检索后运行，用户可以在这里定义检查，返回错误将取消升级
	PreUpgrade func(tempBinaryPath string) error
	//Debug enables all [smoother] logs.
	//启用所有 [smoother] 日志,设置为 true 以启用详细日志记录，便于调试
	Debug bool
	//NoWarn disables warning [smoother] logs.
	//禁用 [smoother] 警告日志,设置为 true 以抑制警告日志
	NoWarn bool
	//NoRestart disables all restarts, this option essentially converts
	//the RestartSignal into a "ShutdownSignal".
	//禁用所有重启，此选项实际上将 RestartSignal 转换为 "ShutdownSignal",使用这个选项来防止任何重启，仅允许关闭
	NoRestart bool
	//NoRestartAfterFetch disables automatic restarts after each upgrade.
	//Though manual restarts using the RestartSignal can still be performed.
	//禁用每次升级后的自动重启，但仍然可以使用 RestartSignal 手动重启,如果你希望在每次升级后手动处理重启，可以设置为 true
	NoRestartAfterFetch bool
	//Fetcher will be used to fetch binaries.
	//用于获取二进制文件的接口。 实现并分配一个 fetcher，用于检索应用程序的更新二进制文件
>>>>>>> develop
	Fetcher fetcher.Interface
}

func validate(c *Config) error {
	//validate
	if c.Program == nil {
		return errors.New("smoother.Config.Program required")
	}
	if c.Address != "" {
		if len(c.Addresses) > 0 {
<<<<<<< HEAD
			return errors.New("smoother.Config.Address and Addresses cant both be set")
		}
		c.Addresses = []string{c.Address}
	} else if len(c.Addresses) > 0 {
=======
			return errors.New("smoother.Config.Address and Addresses can not both be set")
		}
		//如果 Addresses 没有被设置，默认设置为 Address
		c.Addresses = []string{c.Address}
	} else if len(c.Addresses) > 0 {
		//如果 Address 没有被设置，默认设置为 Addresses[0]
>>>>>>> develop
		c.Address = c.Addresses[0]
	}
	if c.RestartSignal == nil {
		c.RestartSignal = SIGUSR2
	}
	if c.TerminateTimeout <= 0 {
<<<<<<< HEAD
		c.TerminateTimeout = 30 * time.Second
=======
		c.TerminateTimeout = 600 * time.Second
>>>>>>> develop
	}
	if c.MinFetchInterval <= 0 {
		c.MinFetchInterval = 1 * time.Second
	}
	return nil
}

// RunErr allows manual handling of any
// smoother errors.
func RunErr(c Config) error {
	return runErr(&c)
}

// Run executes smoother, if an error is
<<<<<<< HEAD
// encountered, smoother fallsback to running
=======
// encountered, smoother fallback to running
>>>>>>> develop
// the program directly (unless Required is set).
func Run(c Config) {
	err := runErr(&c)
	if err != nil {
<<<<<<< HEAD
		if c.Required {
=======
		// if running runErr error, check Required flag
		if c.Required {
			// if Required is set, then exit with error
>>>>>>> develop
			log.Fatalf("[smoother] %s", err)
		} else if c.Debug || !c.NoWarn {
			log.Printf("[smoother] disabled. run failed: %s", err)
		}
		c.Program(DisabledState)
		return
	}
	os.Exit(0)
}

// sanityCheck returns true if a check was performed
func sanityCheck() bool {
	//sanity check
	if token := os.Getenv(envBinCheck); token != "" {
		fmt.Fprint(os.Stdout, token)
		return true
	}
<<<<<<< HEAD
	//legacy sanity check using old env var
	if token := os.Getenv(envBinCheckLegacy); token != "" {
		fmt.Fprint(os.Stdout, token)
		return true
	}
=======

>>>>>>> develop
	return false
}

// SanityCheck manually runs the check to ensure this binary
// is compatible with smoother. This tries to ensure that a restart
// is never performed against a bad binary, as it would require
// manual intervention to rectify. This is automatically done
// on smoother.Run() though it can be manually run prior whenever
// necessary.
func SanityCheck() {
	if sanityCheck() {
		os.Exit(0)
	}
}

// abstraction over master/slave
var currentProcess interface {
	triggerRestart()
	run() error
}

func runErr(c *Config) error {
	//os not supported
	if !supported {
		return fmt.Errorf("os (%s) not supported", runtime.GOOS)
	}
	if err := validate(c); err != nil {
		return err
	}
	if sanityCheck() {
		return nil
	}
	//run either in master or slave mode
	if os.Getenv(envIsSlave) == "1" {
<<<<<<< HEAD
		currentProcess = &slave{Config: c}
	} else {
		// 最开始刚启动的时候 envIsSlave 是空的，默认是 master 启动，运行 proc_master.go 里的 run()
=======
		//当master进程fork出 slave 进程，slave 进程会运行 proc_slave.go 里的 run()
		currentProcess = &slave{Config: c}
	} else {
		//最开始刚启动的时候 envIsSlave 是空的，默认是 master 启动，运行 proc_master.go 里的 run()
>>>>>>> develop
		currentProcess = &master{Config: c}
	}
	return currentProcess.run()
}

// Restart programmatically triggers a graceful restart. If NoRestart
// is enabled, then this will essentially be a graceful shutdown.
func Restart() {
	if currentProcess != nil {
		currentProcess.triggerRestart()
	}
}

// IsSupported returns whether smoother is supported on the current OS.
func IsSupported() bool {
	return supported
}

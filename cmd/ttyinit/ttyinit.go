package ttyinit

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rancher/k3os/config"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	running      = true
	processes    = map[int]*os.Process{}
	processLock  = sync.Mutex{}
	agettyCmd    = "/sbin/agetty -a rancher -J -p %s linux"
	expectedList = []string{
		"tty1",
		"tty2",
		"tty3",
		"tty4",
		"tty5",
		"tty6",
		"ttyS0",
		"ttyS1",
		"ttyS2",
		"ttyS3",
		"ttyAMA0",
	}
)

func Main() {
	// TODO: add logs
	//logrus.Initlogrusger()
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
	app := cli.NewApp()

	app.Name = os.Args[0]
	app.Usage = fmt.Sprintf("%s K3OS(%s)", app.Name, config.OSBuildDate)
	app.Version = config.OSVersion
	app.Author = "Rancher Labs, Inc."

	app.Flags = []cli.Flag{}
	app.Action = run

	app.Run(os.Args)
}

func setupSigterm() {
	sigtermChan := make(chan os.Signal)
	signal.Notify(sigtermChan, syscall.SIGTERM)
	go func() {
		for range sigtermChan {
			termPids()
		}
	}()
}

func run(c *cli.Context) error {
	setupSigterm()

	doneChannel := make(chan string, len(expectedList))

	for _, tty := range expectedList {
		if isAvailable(tty) {
			go execute(tty, doneChannel)
		}
	}

	for i := 0; i < len(expectedList); i++ {
		tty := <-doneChannel
		//logrus.Infof("%s has been stopped", tty)
		fmt.Println(tty, "has been stopped")
	}
	return nil
}

func addProcess(process *os.Process) {
	processLock.Lock()
	defer processLock.Unlock()
	processes[process.Pid] = process
}

func removeProcess(process *os.Process) {
	processLock.Lock()
	defer processLock.Unlock()
	delete(processes, process.Pid)
}

func termPids() {
	running = false
	processLock.Lock()
	defer processLock.Unlock()

	for _, process := range processes {
		//logrus.Infof("sending SIGTERM to %d", process.Pid)
		process.Signal(syscall.SIGTERM)
	}
}

func execute(tty string, doneChannel chan string) {
	defer func() { doneChannel <- tty }()

	start := time.Now()
	count := 0

	args := strings.Split(fmt.Sprintf(agettyCmd, tty), " ")

	for {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}

		if err := cmd.Start(); err == nil {
			addProcess(cmd.Process)
			if err = cmd.Wait(); err != nil {
				if !strings.Contains(err.Error(), syscall.SIGTERM.String()) {
					logrus.Errorf("wait cmd to exit: %s, err: %v", tty, err)
				}
			}
			removeProcess(cmd.Process)
		} else {
			logrus.Errorf("start cmd: %s, err: %v", tty, err)
		}

		if !running {
			//logrus.Infof("%s : not restarting, exiting", tty)
			break
		}

		count++

		if count > 10 {
			if time.Now().Sub(start) <= (1 * time.Second) {
				logrus.Errorf("%s : restarted too fast, not executing", tty)
				break
			}

			count = 0
			start = time.Now()
		}
	}
}

func isAvailable(name string) bool {
	if f, err := os.Open(fmt.Sprintf("/dev/%s", name)); err == nil {
		return terminal.IsTerminal(int(f.Fd()))
	}
	return false
}

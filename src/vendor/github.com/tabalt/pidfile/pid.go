package pidfile

import (
	"os"
	"syscall"
)

type Pid struct {
	Id int
}

func NewPid(id int) *Pid {
	return &Pid{
		Id: id,
	}
}

func (p *Pid) ProcessExist() bool {
	process, err := os.FindProcess(p.Id)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}

	return true
}

package pidfile

import (
	"errors"
	"os"
	"strconv"
)

const (
	PidFileMinId         = 0
	PidFileTmpPathSuffix = ".tmp"
)

type PidFile struct {
	Pid *Pid
	File *File
	TmpFile *File
}

func NewPidFile(path string) *PidFile {
	return &PidFile{
		Pid: NewPid(os.Getpid()),
		File: NewFile(path),
		TmpFile: NewFile(path + PidFileTmpPathSuffix),
	}
}

func (pf *PidFile) Create() error {
	file := pf.File
	pid, err := pf.ReadPidFromFile(file)
	if err == nil && pid.ProcessExist() {
		file = pf.TmpFile
	}

	if err := pf.WritePidToFile(file, pf.Pid); err != nil {
		return err
	}

	return nil
}

func (pf *PidFile) Clear() error {
	pid, err := pf.ReadPidFromFile(pf.File)
	tmpPid, tmpErr := pf.ReadPidFromFile(pf.TmpFile)
	
	if err != nil && tmpErr != nil {
		return errors.New("clear pid error: " + err.Error() +", clear tmp pid error: "+ tmpErr.Error())
	}

	if err == nil && pf.Pid.Id == pid.Id {
		if err = pf.File.Remove(); err != nil {
			return err
		}

		if tmpErr == nil {
			if err = pf.TmpFile.Rename(pf.File.Path); err != nil {
				return err
			}
		}
	} else {
		if tmpErr == nil && pf.Pid.Id == tmpPid.Id {
			if err = pf.TmpFile.Remove(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (pf *PidFile) ReadPidFromFile(file *File) (*Pid, error) {
	fb, err := file.Read()
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(string(fb))
	if err != nil || id <= PidFileMinId {
		return nil, errors.New("pid file data error")
	}

	return NewPid(id), nil
}

func (pf *PidFile)  WritePidToFile(file *File, pid *Pid) error {
	fb := []byte(strconv.Itoa(pid.Id))
	return file.Write(fb)
}

func CreatePidFile(path string) (*PidFile, error) {
	pf := NewPidFile(path)
	err := pf.Create()
	return pf, err
}

func ClearPidFile(pf *PidFile) error {
	return pf.Clear()
}

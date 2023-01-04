package utilities

import (
	"github.com/prometheus/procfs"
)

var Proc procfs.FS

func InitProcfs() {
	proc, err := procfs.NewFS("/host/proc")
	if err != nil {
		panic(err)
	}
	Proc = proc
}

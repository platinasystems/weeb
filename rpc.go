package weeb

import (
	"github.com/platinasystems/vnet/srpc"

	"io"
)

type Rpc struct {
	srpc.Server
}

func NewRpc(rwc io.ReadWriteCloser, regs ...interface{}) (r *Rpc) {
	r = &Rpc{}
	r.Init(rwc, regs...)
	return
}

package main

import (
	"github.com/platinasystems/elib/elog"
	"github.com/platinasystems/weeb"

	"fmt"
	"log"
	"os/exec"
	"time"
)

type Listener struct {
	rpc *weeb.Rpc
}

func (l *Listener) Exec(args *[]string, result *string) (err error) {
	var b []byte
	b, err = exec.Command((*args)[0], (*args)[1:]...).Output()
	if err == nil {
		*result = string(b)
	} else {
		*result = ""
	}
	return
}

var ackReply = "ack"
var count = 0
var printEvery = 100
var sleep = 0 * time.Second

func (l *Listener) Hello(msg []byte, ack *string) error {
	if printEvery != 0 && count%printEvery == 0 {
		log.Println("Hello: ", string(msg))
	}
	count++
	*ack = ackReply
	return nil
}

func HelloRpcClient(r *weeb.Rpc) {
	i := 0
	for {
		msg := fmt.Sprintf("Hello %d", i)
		i++
		var reply string
		err := r.Call("Listener.Hello", []byte(msg), &reply)
		if err != nil {
			log.Printf("Call %v", err)
			return
		}
		if reply != ackReply {
			log.Printf("Reply: %v", reply)
		}
		if sleep != 0 {
			time.Sleep(sleep)
		}
	}
}

func (l *Listener) GetEventView(args struct{}, ack *elog.View) (err error) {
	v := elog.NewView()
	*ack = *v
	return
}

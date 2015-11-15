// +build !js

package main

import (
	"github.com/platinasystems/elib/elog"
	"github.com/platinasystems/weeb"
	"golang.org/x/net/websocket"

	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

//go:generate weebgen -url /js/foundation_deps.min.js -no-inline-data internal/js/foundation_deps/foundation_deps.min.js
import _ "github.com/platinasystems/weeb/example/internal/js/foundation_deps"

//go:generate weebgen  -url /js/foundation.min.js -no-inline-data internal/js/foundation/foundation.min.js
import _ "github.com/platinasystems/weeb/example/internal/js/foundation"

//go:generate weebgen -url /css/eg.min.css -no-inline-data internal/css/eg/eg.min.css
import _ "github.com/platinasystems/weeb/example/internal/css/eg"

//go:generate sh -c "gopherjs build -o js.min.js github.com/platinasystems/weeb/example && weebgen -url /js/js.min.js -no-inline-data -package main js.min.js"

func handle_content(res http.ResponseWriter, req *http.Request) {
	c, found := weeb.ContentByPath[req.URL.Path]

	if !found {
		http.NotFound(res, req)
		return
	}

	res.Header().Set("Content-Type", c.ContentType)
	if len(c.ContentEncoding) != 0 {
		res.Header().Set("Content-Encoding", c.ContentEncoding)
	}
	var r io.ReadSeeker
	if len(c.Data) > 0 {
		r = bytes.NewReader(c.Data)
	} else {
		var err error
		r, err = os.Open(c.FilePath)
		if err != nil {
			http.NotFound(res, req)
			return
		}
	}
	http.ServeContent(res, req, "", time.Unix(c.UnixTimeLastModified, 0), r)
}

func root(w http.ResponseWriter, r *http.Request) {
	p, d, _ := mySite.Match(r.URL.Path)
	if p != nil {
		d.Reset()
		d.Body = p.PageBody(r.URL.Path, d)
		io.WriteString(w, d.Markup())
	} else {
		http.NotFound(w, r)
	}
}

func handle_ws(ws *websocket.Conn) {
	// Gob based RPC requires binary web socket frames.
	ws.PayloadType = websocket.BinaryFrame
	l := &Listener{}
	r := weeb.NewRpc(ws, l)
	l.rpc = r
	// go HelloRpcClient(r)
	err := r.Serve()
	if err != io.EOF {
		log.Printf("%v", err)
	}
	ws.Close()
}

func main() {
	http.HandleFunc("/js/", handle_content)
	http.HandleFunc("/css/", handle_content)
	http.HandleFunc("/", root)
	http.Handle("/ws/rpc/", websocket.Handler(handle_ws))

	elog.Enable(true)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

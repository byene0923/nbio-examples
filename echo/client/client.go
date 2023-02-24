package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/lesismal/nbio"
	"github.com/lesismal/nbio/logging"
)

func main() {
	var (
		ret    []byte
		buf    = make([]byte, 1024*1024*4)
		addr   = "localhost:8888"
		ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	)

	logging.SetLevel(logging.LevelInfo)

	rand.Read(buf)

	engine := nbio.NewEngine(nbio.Config{})

	done := make(chan int)
	engine.OnData(func(c *nbio.Conn, data []byte) {
		ret = append(ret, data...)
		if len(ret) == len(buf) {
			if bytes.Equal(buf, ret) {
				close(done)
			}
		}
	})

	err := engine.Start()
	if err != nil {
		fmt.Printf("Start failed: %v\n", err)
	}
	defer engine.Stop()

	// net.Dial also can be used here
	c, err := nbio.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	engine.AddConn(c)
	c.Write(buf)

	select {
	case <-ctx.Done():
		logging.Error("timeout")
	case <-done:
		logging.Info("success")
	}
}

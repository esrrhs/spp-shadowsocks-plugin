package main

import (
	"flag"
	"fmt"
	"github.com/esrrhs/go-engine/src/common"
	"github.com/esrrhs/go-engine/src/conn"
	"github.com/esrrhs/go-engine/src/loggo"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	usage := "proto=tcp/rudp/kcp/quic"

	flag.Usage = func() {
		fmt.Println(usage)
	}
	vpn := flag.Bool("V", false, "")
	flag.Parse()

	loggo.Ini(loggo.Config{
		Level:     loggo.LEVEL_INFO,
		Prefix:    "spp",
		MaxDay:    3,
		NoLogFile: true,
		NoPrint:   false,
	})

	log_init()

	loggo.Info("start plugin")

	loggo.Info("vpn mode %v", *vpn)

	loggo.Info("parse env")
	opts, err := parseEnv()
	if err != nil {
		loggo.Error("parseEnv fail %v", err)
		os.Exit(-11)
	}

	remoteaddr, _ := opts.Get(strings.ToLower("remoteaddr"))
	localaddr, _ := opts.Get(strings.ToLower("localaddr"))
	loggo.Info("remoteaddr %v", remoteaddr)
	loggo.Info("localaddr %v", localaddr)

	var protos []string
	if value, b := opts.Get("proto"); b {
		protos = append(protos, value)
	} else {
		protos = append(protos, "tcp")
	}
	loggo.Info("protos %v", protos)

	if *vpn {
		registerControlFunc()
	}

	go parentMonitor(3)

	ty, _ := opts.Get("type")
	if len(ty) > 0 && ty == "server" {
		loggo.Info("start server")

		c, err := conn.NewConn(protos[0])
		if err != nil {
			loggo.Error("NewConn fail %v", err)
			os.Exit(-21)
		}

		lc, err := c.Listen(remoteaddr)
		if err != nil {
			loggo.Error("Listen fail %v", err)
			os.Exit(-22)
		}

		loggo.Info("start server ok")

		for {
			sc, err := lc.Accept()
			if err != nil {
				loggo.Error("Accept fail %v", err)
				time.Sleep(time.Second)
				continue
			}
			go process("tcp", localaddr, sc)
		}
	} else {
		loggo.Info("start client")

		c, err := conn.NewConn("tcp")
		if err != nil {
			loggo.Error("NewConn fail %v", err)
			os.Exit(-31)
		}

		lc, err := c.Listen(localaddr)
		if err != nil {
			loggo.Error("Listen fail %v", err)
			os.Exit(-32)
		}

		loggo.Info("start client ok")

		for {
			sc, err := lc.Accept()
			if err != nil {
				loggo.Error("Accept fail %v", err)
				time.Sleep(time.Second)
				continue
			}
			go process(protos[0], remoteaddr, sc)
		}
	}
}

func process(proto string, remoteaddr string, c conn.Conn) {

	defer common.CrashLog()

	loggo.Info("process begin %s %s %s", proto, remoteaddr, c.Info())

	c, err := conn.NewConn(proto)
	if err != nil {
		loggo.Error("process NewConn fail %v", err)
		os.Exit(-41)
	}

	pc, err := c.Dial(remoteaddr)
	if err != nil {
		loggo.Error("process Dial fail %v", err)
		return
	}

	errCh := make(chan error, 2)
	go proxy(c, pc, c.Info(), pc.Info(), errCh)
	go proxy(pc, c, pc.Info(), c.Info(), errCh)

	for i := 0; i < 2; i++ {
		<-errCh
	}

	c.Close()
	pc.Close()

	loggo.Info("process end %s %s %s", proto, remoteaddr, c.Info())
}

func proxy(destination io.Writer, source io.Reader, dst string, src string, errCh chan error) {
	loggo.Info("proxy begin from %s -> %s", src, dst)
	n, err := io.Copy(destination, source)
	errCh <- err
	loggo.Info("proxy end from %s -> %s %v %v", src, dst, n, err)
}

func parentMonitor(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	pid := os.Getppid()
	for {
		select {
		case <-ticker.C:
			curpid := os.Getppid()
			if curpid != pid {
				os.Exit(-51)
			}
		}
	}
}

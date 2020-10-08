package main

import (
	"flag"
	"fmt"
	"github.com/esrrhs/go-engine/src/common"
	"github.com/esrrhs/go-engine/src/loggo"
	"github.com/esrrhs/go-engine/src/proxy"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func main() {
	usage := ""

	config := proxy.DefaultConfig()
	ss := reflect.ValueOf(config).Elem()
	typeOfT := ss.Type()
	for i := 0; i < ss.NumField(); i++ {
		name := typeOfT.Field(i).Name
		if ss.Field(i).Kind() == reflect.Int {
			usage += fmt.Sprintf("%v = %v\n", strings.ToLower(name), ss.Field(i).Int())
		} else if ss.Field(i).Kind() == reflect.String {
			usage += fmt.Sprintf("%v = %v\n", strings.ToLower(name), ss.Field(i).String())
		} else if ss.Field(i).Kind() == reflect.Bool {
			usage += fmt.Sprintf("%v = %v\n", strings.ToLower(name), ss.Field(i).Bool())
		}
	}

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

	loggo.Info("set config")
	for i := 0; i < ss.NumField(); i++ {
		name := typeOfT.Field(i).Name
		if value, b := opts.Get(strings.ToLower(name)); b {
			if ss.Field(i).IsValid() && ss.Field(i).CanSet() {
				if ss.Field(i).Kind() == reflect.Int {
					x, _ := strconv.Atoi(value)
					if !ss.Field(i).OverflowInt(int64(x)) {
						ss.Field(i).SetInt(int64(x))
						loggo.Info("%v = %v", name, x)
					}
				} else if ss.Field(i).Kind() == reflect.String {
					ss.Field(i).SetString(value)
					loggo.Info("%v = %v", name, value)
				} else if ss.Field(i).Kind() == reflect.Bool {
					x := value == "true"
					ss.Field(i).SetBool(x)
					loggo.Info("%v = %v", name, x)
				}
			}
		}
	}

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
		_, err := proxy.NewServer(config, protos, []string{remoteaddr})
		if err != nil {
			loggo.Error("NewServer fail %v", err)
			os.Exit(-12)
		}
		loggo.Info("start server ok")
	} else {
		loggo.Info("start client")
		_, err := proxy.NewClient(config, protos[0], remoteaddr, common.UniqueId(),
			"ss_proxy",
			[]string{"tcp"}, []string{localaddr}, []string{""})
		if err != nil {
			loggo.Error("NewClient fail %v", err)
			os.Exit(-13)
		}
		loggo.Info("start client ok")
	}

	for {
		time.Sleep(time.Hour)
	}
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
				os.Exit(1)
			}
		}
	}
}

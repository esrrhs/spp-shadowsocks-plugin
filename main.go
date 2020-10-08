package main

import (
	"flag"
	"fmt"
	"github.com/esrrhs/go-engine/src/common"
	"github.com/esrrhs/go-engine/src/loggo"
	"github.com/esrrhs/go-engine/src/proxy"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func main() {
	log_init()

	log.Printf("start plugin")

	loggo.Ini(loggo.Config{
		Level:     loggo.LEVEL_ERROR,
		Prefix:    "spp",
		MaxDay:    3,
		NoLogFile: true,
		NoPrint:   false,
	})

	usage := ""

	config := proxy.DefaultConfig()
	ss := reflect.ValueOf(config).Elem()
	typeOfT := ss.Type()
	for i := 0; i < ss.NumField(); i++ {
		name := typeOfT.Field(i).Name
		if ss.Field(i).Kind() == reflect.Int {
			usage += fmt.Sprintf("%v = %v", strings.ToLower(name), ss.Field(i).Int())
		} else if ss.Field(i).Kind() == reflect.String {
			usage += fmt.Sprintf("%v = %v", strings.ToLower(name), ss.Field(i).String())
		} else if ss.Field(i).Kind() == reflect.Bool {
			usage += fmt.Sprintf("%v = %v", strings.ToLower(name), ss.Field(i).Bool())
		}
	}

	flag.Usage = func() {
		fmt.Printf(usage)
	}
	v := flag.Bool("V", false, "")
	log.Printf("parse args")
	flag.Parse()
	log.Printf("vpn mode %v", *v)

	log.Printf("parse env")
	opts, err := parseEnv()
	if err != nil {
		log.Printf("parseEnv fail %v", err)
		os.Exit(-11)
	}

	remoteaddr, _ := opts.Get(strings.ToLower("remoteaddr"))
	localaddr, _ := opts.Get(strings.ToLower("localaddr"))
	log.Printf("remoteaddr %v", remoteaddr)
	log.Printf("localaddr %v", localaddr)

	log.Printf("set config")
	for i := 0; i < ss.NumField(); i++ {
		name := typeOfT.Field(i).Name
		if value, b := opts.Get(strings.ToLower(name)); b {
			if ss.Field(i).IsValid() && ss.Field(i).CanSet() {
				if ss.Field(i).Kind() == reflect.Int {
					x, _ := strconv.Atoi(value)
					if !ss.Field(i).OverflowInt(int64(x)) {
						ss.Field(i).SetInt(int64(x))
						log.Printf("%v = %v", name, x)
					}
				} else if ss.Field(i).Kind() == reflect.String {
					ss.Field(i).SetString(value)
					log.Printf("%v = %v", name, value[0])
				} else if ss.Field(i).Kind() == reflect.Bool {
					x := value == "true"
					ss.Field(i).SetBool(x)
					log.Printf("%v = %v", name, x)
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

	go parentMonitor(3)

	ty, _ := opts.Get("type")
	if len(ty) > 0 && ty == "server" {
		log.Printf("start server")
		_, err := proxy.NewServer(config, protos, []string{remoteaddr})
		if err != nil {
			log.Printf("NewServer fail %v", err)
			os.Exit(-12)
		}
		log.Printf("start server ok")
	} else {
		log.Printf("start client")
		_, err := proxy.NewClient(config, protos[0], remoteaddr, common.UniqueId(),
			"ss_proxy",
			[]string{"tcp"}, []string{localaddr}, []string{""})
		if err != nil {
			log.Printf("NewClient fail %v", err)
			os.Exit(-13)
		}
		log.Printf("start client ok")
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

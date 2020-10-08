package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/esrrhs/go-engine/src/common"
	"github.com/esrrhs/go-engine/src/loggo"
	"github.com/esrrhs/go-engine/src/proxy"
	"net/http"
	_ "net/http/pprof"
	"os"
	"reflect"
	sort "sort"
	"strconv"
	"strings"
	"time"
)

type Args map[string][]string

func (args Args) Get(key string) (value string, ok bool) {
	if args == nil {
		return "", false
	}
	vals, ok := args[key]
	if !ok || len(vals) == 0 {
		return "", false
	}
	return vals[0], true
}

func (args Args) Add(key, value string) {
	args[key] = append(args[key], value)
}

func parseArgs() Args {
	opts := make(Args)

	ss_remote_host := os.Getenv("SS_REMOTE_HOST")
	ss_remote_port := os.Getenv("SS_REMOTE_PORT")
	ss_local_host := os.Getenv("SS_LOCAL_HOST")
	ss_local_port := os.Getenv("SS_LOCAL_PORT")
	if len(ss_remote_host) == 0 {
		loggo.Info("need ss_remote_host")
		os.Exit(-1)
	}
	if len(ss_remote_port) == 0 {
		loggo.Info("need ss_remote_port")
		os.Exit(-2)
	}
	if len(ss_local_host) == 0 {
		loggo.Info("need ss_local_host")
		os.Exit(-3)
	}
	if len(ss_local_port) == 0 {
		loggo.Info("need ss_local_port")
		os.Exit(-4)
	}

	opts.Add("remoteAddr", ss_remote_host)
	opts.Add("remotePort", ss_remote_port)
	opts.Add("localAddr", ss_local_host)
	opts.Add("localPort", ss_local_port)
	loggo.Info("remoteAddr = %v", ss_remote_host)
	loggo.Info("remotePort = %v", ss_remote_port)
	loggo.Info("localAddr = %v", ss_local_host)
	loggo.Info("localPort = %v", ss_local_port)

	ss_plugin_options := os.Getenv("SS_PLUGIN_OPTIONS")
	if len(ss_plugin_options) > 0 {
		other_opts, err := parsePluginOptions(ss_plugin_options)
		if err != nil {
			loggo.Info("parse SS_PLUGIN_OPTIONS fail %v", err)
			os.Exit(-5)
		}
		for k, v := range other_opts {
			opts[k] = v
		}
	}

	return opts
}

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
		fmt.Printf(usage)
	}
	flag.Parse()

	loggo.Ini(loggo.Config{
		Level:     loggo.LEVEL_INFO,
		Prefix:    "spp",
		MaxDay:    3,
		NoLogFile: true,
		NoPrint:   false,
	})

	opts := parseArgs()

	for i := 0; i < ss.NumField(); i++ {
		name := typeOfT.Field(i).Name
		value := opts[strings.ToLower(name)]
		if len(value) > 0 && ss.Field(i).IsValid() && ss.Field(i).CanSet() {
			if ss.Field(i).Kind() == reflect.Int {
				x, _ := strconv.Atoi(value[0])
				if !ss.Field(i).OverflowInt(int64(x)) {
					ss.Field(i).SetInt(int64(x))
					loggo.Info("%v = %v", name, x)
				}
			} else if ss.Field(i).Kind() == reflect.String {
				ss.Field(i).SetString(value[0])
				loggo.Info("%v = %v", name, value[0])
			} else if ss.Field(i).Kind() == reflect.Bool {
				x := value[0] == "true"
				ss.Field(i).SetBool(x)
				loggo.Info("%v = %v", name, x)
			}
		}
	}

	var protos []string
	proto := opts["proto"]
	if len(proto) > 0 {
		protos = append(protos, proto[0])
	} else {
		protos = append(protos, "tcp")
	}

	ty := opts["type"]
	if len(ty) > 0 && ty[0] == "server" {
		_, err := proxy.NewServer(config, protos, []string{opts["remoteAddr"][0] + ":" + opts["remotePort"][0]})
		if err != nil {
			loggo.Info("NewServer fail %v", err)
			os.Exit(-8)
		}
	} else {
		_, err := proxy.NewClient(config, protos[0], opts["remoteAddr"][0]+":"+opts["remotePort"][0], common.UniqueId(),
			"ss_proxy",
			[]string{"tcp"}, []string{opts["localAddr"][0] + ":" + opts["localPort"][0]}, []string{""})
		if err != nil {
			loggo.Info("NewClient fail %v", err)
			os.Exit(-10)
		}
		loggo.Info("Client start")
	}

	profile := opts["profile"]
	if len(profile) > 0 {
		go http.ListenAndServe("0.0.0.0:"+profile[0], nil)
	}

	for {
		time.Sleep(time.Hour)
	}
}

func parsePluginOptions(s string) (opts Args, err error) {
	opts = make(Args)
	if len(s) == 0 {
		return
	}
	i := 0
	for {
		var key, value string
		var offset, begin int

		if i >= len(s) {
			break
		}
		begin = i
		// Read the key.
		offset, key, err = indexUnescaped(s[i:], []byte{'=', ';'})
		if err != nil {
			return
		}
		if len(key) == 0 {
			err = fmt.Errorf("empty key in %q", s[begin:i])
			return
		}
		i += offset
		// End of string or no equals sign?
		if i >= len(s) || s[i] != '=' {
			opts.Add(key, "1")
			// Skip the semicolon.
			i++
			continue
		}
		// Skip the equals sign.
		i++
		// Read the value.
		offset, value, err = indexUnescaped(s[i:], []byte{';'})
		if err != nil {
			return
		}
		i += offset
		opts.Add(key, value)
		// Skip the semicolon.
		i++
	}
	return opts, nil
}

// Escape backslashes and all the bytes that are in set.
func backslashEscape(s string, set []byte) string {
	var buf bytes.Buffer
	for _, b := range []byte(s) {
		if b == '\\' || bytes.IndexByte(set, b) != -1 {
			buf.WriteByte('\\')
		}
		buf.WriteByte(b)
	}
	return buf.String()
}

// Encode a name–value mapping so that it is suitable to go in the ARGS option
// of an SMETHOD line. The output is sorted by key. The "ARGS:" prefix is not
// added.
//
// "Equal signs and commas [and backslashes] must be escaped with a backslash."
func encodeSmethodArgs(args Args) string {
	if args == nil {
		return ""
	}

	keys := make([]string, 0, len(args))
	for key := range args {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	escape := func(s string) string {
		return backslashEscape(s, []byte{'=', ','})
	}

	var pairs []string
	for _, key := range keys {
		for _, value := range args[key] {
			pairs = append(pairs, escape(key)+"="+escape(value))
		}
	}

	return strings.Join(pairs, ",")
}

func indexUnescaped(s string, term []byte) (int, string, error) {
	var i int
	unesc := make([]byte, 0)
	for i = 0; i < len(s); i++ {
		b := s[i]
		// A terminator byte?
		if bytes.IndexByte(term, b) != -1 {
			break
		}
		if b == '\\' {
			i++
			if i >= len(s) {
				return 0, "", fmt.Errorf("nothing following final escape in %q", s)
			}
			b = s[i]
		}
		unesc = append(unesc, b)
	}
	return i, string(unesc), nil
}

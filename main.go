package main

import (
	"bytes"
	"fmt"
	"github.com/esrrhs/go-engine/src/common"
	"github.com/esrrhs/go-engine/src/loggo"
	"github.com/esrrhs/go-engine/src/proxy"
	"net/http"
	_ "net/http/pprof"
	"os"
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
		os.Exit(-1)
	}
	if len(ss_remote_port) == 0 {
		os.Exit(-2)
	}
	if len(ss_local_host) == 0 {
		os.Exit(-3)
	}
	if len(ss_local_port) == 0 {
		os.Exit(-4)
	}

	opts.Add("remoteAddr", ss_remote_host)
	opts.Add("remotePort", ss_remote_port)
	opts.Add("localAddr", ss_local_host)
	opts.Add("localPort", ss_local_port)

	ss_plugin_options := os.Getenv("SS_PLUGIN_OPTIONS")
	if len(ss_plugin_options) > 0 {
		other_opts, err := parsePluginOptions(ss_plugin_options)
		if err != nil {
			os.Exit(-5)
		}
		for k, v := range other_opts {
			opts[k] = v
		}
	} else {
		os.Exit(-6)
	}

	return opts
}

func main() {
	opts := parseArgs()

	loggo.Ini(loggo.Config{
		Level:     loggo.LEVEL_ERROR,
		Prefix:    "spp",
		MaxDay:    3,
		NoLogFile: true,
		NoPrint:   true,
	})

	config := proxy.DefaultConfig()
	compress := opts["compress"]
	if len(compress) > 0 {
		config.Compress, _ = strconv.Atoi(compress[0])
	}
	key := opts["key"]
	if len(key) > 0 {
		config.Key = key[0]
	}
	encrypt := opts["encrypt"]
	if len(encrypt) > 0 {
		config.Encrypt = encrypt[0]
	}
	config.ShowPing = false
	maxclient := opts["maxclient"]
	if len(maxclient) > 0 {
		config.MaxClient, _ = strconv.Atoi(maxclient[0])
	}
	maxconn := opts["maxconn"]
	if len(maxclient) > 0 {
		config.MaxSonny, _ = strconv.Atoi(maxconn[0])
	}
	var protos []string
	proto := opts["proto"]
	if len(proto) > 0 {
		protos = append(protos, proto[0])
	}

	ty := opts["type"]
	if len(ty) > 0 && ty[0] == "server" {
		var listenaddrs []string
		listenaddr := opts["listenaddr"]
		if len(listenaddr) > 0 {
			listenaddrs = append(listenaddrs, listenaddr[0])
		} else {
			os.Exit(-7)
		}
		_, err := proxy.NewServer(config, protos, listenaddrs)
		if err != nil {
			os.Exit(-8)
		}
	} else if len(ty) > 0 && ty[0] == "client" {
		server := opts["server"]
		if len(server) <= 0 {
			os.Exit(-9)
		}
		_, err := proxy.NewClient(config, protos[0], server[0], common.UniqueId(),
			"proxy_client",
			[]string{"tcp"}, []string{opts["localAddr"][0] + ":" + opts["localPort"][0]}, []string{opts["remoteAddr"][0] + ":" + opts["remotePort"][0]})
		if err != nil {
			os.Exit(-10)
		}
		loggo.Info("Client start")
	} else {
		os.Exit(-11)
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

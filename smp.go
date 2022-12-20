package smp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type LoadOpts struct {
	path string
}

type option func(*LoadOpts)

func Path(p string) option {
	return func(l *LoadOpts) {
		l.path = p
	}
}

func Load(path ...option) string {

	home, _ := os.UserHomeDir()
	p := &LoadOpts{
		path: filepath.Join(home, ".ssh", "config"),
	}
	for _, opt := range path {
		opt(p)
	}
	f, _ := os.Open(p.path)
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}
	return string(data)
}

type Host struct {
	Name    string
	Options []Option
}

type Option struct {
	Name  string
	Value string
}

func Parse(s string) []Host {
	var hosts []Host
	reg := "\r\n|\n"
	tmp := regexp.MustCompile(reg).Split(s, -1)
	r := regexp.MustCompile(`^[a-zA-Z]`)
	for _, v := range tmp {
		if r.MatchString(v) {
			arr := strings.Fields(v)
			if arr[0] == "Host" {
				hosts = append(hosts, Host{Name: arr[1]})
			}
			if arr[0] == "Include" {
				home, _ := os.UserHomeDir()
				files, _ := filepath.Glob(filepath.Join(home, ".ssh", arr[1]))
				for _, f := range files {
					s := Load(Path(f))
					hosts = append(hosts, Parse(s)...)
				}
			}
		} else {
			i := strings.TrimSpace(v)
			if len(i) == 0 {
				continue
			}
			if strings.HasPrefix(i, "#") {
				continue
			}
			s := strings.Fields(i)
			if len(s) == 0 {
				continue
			}
			hosts[len(hosts)-1].Options = append(hosts[len(hosts)-1].Options, Option{Name: s[0], Value: s[1]})
		}
	}
	fmt.Printf("%+v\n", hosts)
	return hosts
}

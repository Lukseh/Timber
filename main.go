package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v4"
	"golang.org/x/mod/semver"
)

type BuildMode string

const (
	ModeStandard BuildMode = "standard"
	ModeWatch    BuildMode = "watch"
)

type BuildEntry struct {
	Entryfile   string
	Outname     string
	Options     string
	mode        BuildMode
	Dockerfile  string         // Dockerfile location
	Env         map[string]any // Env mapping
	CmdOverrite string         // instead of `go` use for example `wails`
}

type Gopherfile struct {
	Name        string
	Version     string
	Go          string
	Description string
	Build       map[string]BuildEntry
}

var (
	ConfigFileName string
	args           []string
)

func init() {
	flag.StringVar(&ConfigFileName, "f", "Gopherfile", "Specify Gopherfile path")
	flag.Parse()
	args = flag.Args()
}

func main() {
	if len(args) > 0 && args[0] == "init" {
		f, err := os.OpenFile("Gopherfile", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		helloworld := Gopherfile{
			Name:        "helloworld",
			Version:     "1.0.0",
			Go:          "1.25.0",
			Description: "example gopherfile",
			Build: map[string]BuildEntry{
				"release": {
					Entryfile:  "helloworld.go",
					Outname:    "helloworld",
					Options:    "",
					mode:       ModeStandard,
					Dockerfile: "",
				},
			},
		}
		c, err := yaml.Marshal(helloworld)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		if _, err := f.Write(c); err != nil {
			fmt.Println(err)
		}
		f.Sync()
		f.Close()
		os.Exit(0)
	}
	var (
		gf  []byte
		err error
	)
	if gf, err = os.ReadFile(ConfigFileName); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	gc := Gopherfile{}
	if err = yaml.Unmarshal(gf, &gc); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if len(args) < 1 {
		fmt.Println("No arguments passed")
		os.Exit(1)
	}

	out, err := exec.Command("go", "version").Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	rg := regexp.MustCompile(`go version go([\d.]+) `)
	matches := rg.FindSubmatch(out)
	if matches == nil {
		fmt.Println("No go found on device")
		os.Exit(-1)
	}
	goVersion := string(matches[1])
	if !goVersionAtLeast(goVersion, gc.Go) {
		fmt.Println("Your go version does not match minimum project version.")
		os.Exit(-1)
	}

	switch args[0] {
	case "build":
		if len(args) == 1 {
			fmt.Println("Building for release")
			build(&gc, "release")
			break
		}
		build(&gc, args[1])
	case "info":
		os.Exit(info(&gc))
	default:
		fmt.Println("Unsupported action")
		os.Exit(1)
	}
}

func goVersionAtLeast(actual, minimum string) bool {
	a := "v" + actual
	m := "v" + minimum
	return semver.Compare(a, m) >= 0
}

func info(gf *Gopherfile) int {
	if c, err := yaml.Marshal(struct{ Name, Version, Go, Description string }{Name: gf.Name, Version: gf.Version, Go: gf.Go, Description: gf.Description}); err != nil {
		fmt.Println(err)
		return -1
	} else {
		fmt.Println(string(c))
		return 0
	}
}

func build(gf *Gopherfile, name string) {
	var v BuildEntry
	var ok bool
	if v, ok = gf.Build[name]; !ok {
		fmt.Println("No script with that name found in Gopherfile")
		os.Exit(1)
	}
	if name != "release" {
		fmt.Println("Building", name)
	}
	if v.mode == "" {
		v.mode = "standard"
	}
	switch v.mode {
	case ModeStandard:
		if v.CmdOverrite != "" {
			c := strings.Split(v.CmdOverrite, " ")
			cmd := exec.Command(c[0], c[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		goArgs := []string{"build"}
		if v.Outname != "" {
			goArgs = append(goArgs, "-o", v.Outname)
		} else {
			goArgs = append(goArgs, "-o", gf.Name)
		}
		if v.Options != "" {
			goArgs = append(goArgs, strings.Fields(v.Options)...)
		}
		if v.Entryfile != "" {
			goArgs = append(goArgs, v.Entryfile)
		} else {
			goArgs = append(goArgs, "main.go")
		}
		cmd := exec.Command("go", goArgs...)
		cmd.Env = append(os.Environ(), maptoslice(v.Env)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func maptoslice(s map[string]any) []string {
	var r []string
	for k, v := range s {
		r = append(r, fmt.Sprintf("%s=%v", k, v))
	}
	return r
}

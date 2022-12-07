package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zigdon/adventofcode/common"
)

type fs struct {
	Dirs  map[string][]string
	Files map[string]int
}

func (f *fs) String() string {
	res := []string{}
	for d, files := range f.Dirs {
		res = append(res, d+":")
		for _, file := range files {
			fname := d + "/" + file
			if d == "" {
				fname = "/" + fname
			}
			res = append(res, fmt.Sprintf("%s: %d", fname, f.Files[fname]))
		}
	}

	return strings.Join(res, "\n")
}

func (f *fs) Du(dir string) int {
	res := 0
	for fname, size := range f.Files {
		if strings.HasPrefix(fname, dir) {
			res += size
		}
	}

	return res
}

func newFS() *fs {
	return &fs{
		Dirs:  make(map[string][]string),
		Files: make(map[string]int),
	}
}

func parseChdir(dir, cmd string) string {
	cmd = strings.TrimPrefix(cmd, "$ cd ")
	if cmd == "/" {
		return ""
	}
	if cmd == ".." {
		path := strings.Split(dir, "/")
		return strings.Join(path[0:len(path)-1], "/")
	}
	return dir + "/" + cmd
}

func one(data *fs) int {
	res := 0
	for d := range data.Dirs {
		size := data.Du(d)
		if size < 100000 {
			res += size
		}
	}

	return res
}

func two(data *fs) int {
	total := 70000000
	need := 30000000
	best := data.Du("")
	free := total - best
	for d := range data.Dirs {
		size := data.Du(d)
		if size < best && free+size > need {
			best = size
		}
	}
	return best
}

func readFile(path string) (*fs, error) {
	join := func(a, b string) string {
		if a == "/" {
			a = ""
		}
		return a + "/" + b
	}
	res := common.AsStrings(common.ReadTransformedFile(path, common.IgnoreBlankLines))

	f := newFS()
	cwd := ""
	for _, l := range res {
		// log.Printf("<== %s", l)
		if l == "$ ls" {
			continue
		}
		if strings.HasPrefix(l, "$ cd") {
			cwd = parseChdir(cwd, l)
			// log.Printf("chdir => %s", cwd)
			continue
		}
		words := strings.SplitN(l, " ", 2)
		if words[0] == "dir" {
			// log.Printf("newdir %s/%s", cwd, words[1])
			f.Dirs[join(cwd, words[1])] = []string{}
			continue
		}
		// log.Printf("newfile %s/%s: %s", cwd, words[1], words[0])
		f.Files[join(cwd, words[1])] = common.MustInt(words[0])
		f.Dirs[cwd] = append(f.Dirs[cwd], words[1])
	}

	return f, nil
}

func main() {
	data, err := readFile(os.Args[1])
	if err != nil {
		log.Fatalf("%v", err)
	}

	res := one(data)
	fmt.Printf("%v\n", res)

	res = two(data)
	fmt.Printf("%v\n", res)
}

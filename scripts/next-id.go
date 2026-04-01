//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func main() {
	max := 0
	re := regexp.MustCompile(`^(\d+)-`)
	for _, lane := range []string{"backlog", "active", "done"} {
		entries, _ := os.ReadDir(filepath.Join(".agents/work", lane))
		for _, e := range entries {
			if m := re.FindStringSubmatch(e.Name()); m != nil {
				if n, _ := strconv.Atoi(m[1]); n > max {
					max = n
				}
			}
		}
	}
	fmt.Printf("%03d\n", max+1)
}

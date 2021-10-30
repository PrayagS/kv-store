package main

import (
	"fmt"

	"github.com/PrayagS/kv-store/pkg/kvstore"
)

func main() {
	s := kvstore.New()
	s.Set("lmao", 420)
	s.Set("nice", 69)
	x, ok := s.Get("lmao")
	if ok == nil {
		fmt.Printf("%s", x)
	}
	s.Set("lmao", 42069)
	x, ok = s.Get("lmao")
	if ok == nil {
		fmt.Printf("%v", x)
	}
	y, ok := s.Get("gg")
	fmt.Printf("%v %s", y, ok)
}

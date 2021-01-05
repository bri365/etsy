package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

// write 500GB of random data
const (
	dirs  int = 512
	files int = 1024
	size  int = 1024 * 1024 // 1MB files
)

func main() {
	for i := 0; i < dirs; i++ {
		dir := fmt.Sprintf("dir%04d", i)
		os.Mkdir(dir, 0755)
		for j := 0; j < files; j++ {
			data := make([]byte, size)
			rand.Read(data)
			name := fmt.Sprintf("%s/temp%04d.dat", dir, j)
			err := ioutil.WriteFile(name, data, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

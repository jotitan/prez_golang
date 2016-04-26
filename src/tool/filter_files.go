package main

import (
	"os"
	"strconv"
)

func main(){
	dir := os.Args[1]
	threashold,_ := strconv.ParseInt(os.Args[2],10,32)

	f,_ := os.Open(dir)
	files,_ := f.Readdir(-1)
	for _,file := range files {
		if file.Size() > threashold {

		}
	}
}


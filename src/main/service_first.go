package main

import (
	"net/http"
	"fmt"
	"time"
	"os"
	"math/rand"
)

func main(){
	port := os.Args[1]
	prefix := os.Args[2]
	s := http.NewServeMux()
	s.HandleFunc("/get",func(w http.ResponseWriter,r *http.Request){
		value := r.FormValue("value")
		waitTime := time.Millisecond * (time.Duration(700 + (300 - rand.Int()%600)))
		fmt.Println("Receive",value,"and wait",waitTime,"ms")
		<-time.Tick(waitTime)
		w.Write([]byte(fmt.Sprintf("%s_transform_%s",prefix,value)))
	})
	fmt.Println("Launch service on",port)
	http.ListenAndServe(":" + port,s)
}
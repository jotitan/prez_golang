package main

import (
	"math"
	"fmt"
	"crypto/md5"
	"encoding/base64"
	"net/http"
	"io/ioutil"
	"time"
	"runtime"
	"os"
	"strconv"
)

func main() {
	n := 10
	if len(os.Args) > 1 {
		if val,err := strconv.ParseInt(os.Args[1],10,32) ; err == nil {
			n = int(val)
		}
	}
	fmt.Println("Lancement des tests avec n =",n)

	begin := time.Now()
	runSync(n)
	fmt.Println("SYNC :", time.Now().Sub(begin))

	begin = time.Now()
	runAsync(n, 1)
	fmt.Println("ASYNC (1) :", time.Now().Sub(begin))

	begin = time.Now()
	runAsync(n, 4)
	fmt.Println("ASYNC (4) :", time.Now().Sub(begin))

	begin = time.Now()
	runAsync(n, 10)
	fmt.Println("ASYNC (10) :", time.Now().Sub(begin))
}

func runSync(n int)[]Result{
	results := make([]Result,n)
	for i := 0 ; i < n ; i++ {
		value := fmt.Sprintf("value_%d",i)
		result := Result{position:i,original:value}
		result.valueUrl1 = callUrl(url1,value)
		result.valueUrl2 = callUrl(url2,value)
		result.valueFunc = localHeavyFunction(i)
		results[i] = result
	}
	return results
}

func runAsync(n,nbThread int)[]Result{
	runtime.GOMAXPROCS(nbThread)
	chResults := make(chan Result)
	for i := 0 ; i < n ; i++ {
		go func(intValue int, value string){
			channels := []chan string{make(chan string),make(chan string),make(chan string)}
			go func(c chan string){c<- callUrl(url1,value)}(channels[0])
			go func(c chan string){c<- callUrl(url2,value)}(channels[1])
			go func(c chan string){c<- localHeavyFunction(intValue)}(channels[2])

			result := Result{position:intValue,original:value}
			result.valueUrl1 = <- channels[0]
			result.valueUrl2 = <- channels[1]
			result.valueFunc = <- channels[2]
			chResults <- result
		}(i,fmt.Sprintf("value_%d",i))
	}
	results := make([]Result,n)
	for i := 0 ; i < n ; i++ {
		r := <- chResults
		results[r.position] = r
	}
	return results
}

type Result struct {
	position int
	original string
	valueUrl1 string
	valueUrl2 string
	valueFunc string
}

const (
	url1 = "http://localhost:10001/get"
	url2 = "http://localhost:10002/get"
)

func callUrl(url,value string)string{
	resp,_ := http.DefaultClient.Get(url + "?value=" + value)
	data,_ := ioutil.ReadAll(resp.Body)
	return string(data)
}

func localFunction(i int)string{
	m := md5.New()
	m.Write([]byte(fmt.Sprintf("Value_%d",math.Pow(float64(i*5),float64(5)))))
	return base64.URLEncoding.EncodeToString(m.Sum(nil))
}

func localHeavyFunction(i int)string{
	m := md5.New()
	n := 500000
	for i := 0 ; i < n ; i++ {
		m.Write([]byte(fmt.Sprintf("Value_%d",math.Pow(float64(i*5),float64(5)))))
	}
	return base64.URLEncoding.EncodeToString(m.Sum(nil))
}

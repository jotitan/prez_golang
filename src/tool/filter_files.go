package main

import (
	"os"
	"strconv"
	"fmt"
	"path/filepath"
	"math"
	"time"
	"sort"
)

type FileResults []FileResult
func (frs FileResults) Len() int{return len(frs)}
func (frs FileResults) Less(i, j int) bool {
	return frs[i].size < frs[j].size
}
func (frs FileResults) Swap(i, j int){frs[i],frs[j] = frs[j],frs[i]}


type FileResult struct {
	path string
	size int64
}

func (fr FileResult)String()string{
	return fmt.Sprintf("%s : %s\n",fr.path,formatSize(float64(fr.size)))
}

var sizeFormaters = []string{"o","Ko","Mo","Go","To","Po","Eo,Zo,Yo"}

func formatSize(size float64)string{
	exp := math.Floor(math.Log(size)/math.Log(float64(1024)))
	finalValue := size / math.Pow(1024,exp)
	return fmt.Sprintf("%.2f %s",finalValue,sizeFormaters[int(exp)])
}

func main(){
	begin := time.Now()
	dir := os.Args[1]
	threashold,_ := strconv.ParseInt(os.Args[2],10,32)
	recursive :=false
	if len(os.Args) > 3 {
		recursive = os.Args[3] == "-R"
	}
	results := check(dir,int64(threashold),recursive)

	sort.Sort(FileResults(results))
	fmt.Println(results)
	fmt.Println("Results in",time.Now().Sub(begin))
}

func check(dir string,threashold int64,recursive bool)[]FileResult {
	f,_ := os.Open(dir)
	defer f.Close()
	files,_ := f.Readdir(-1)
	results := make([]FileResult,0,len(files))
	for _,file := range files {
		if !file.IsDir() {
			if file.Size() > threashold {
				results = append(results, FileResult{filepath.Join(dir, file.Name()), file.Size()})
			}
		}else{
			if recursive {
				results = append(results,check(filepath.Join(dir,file.Name()),threashold,recursive)...)
			}
		}
	}
	return results
}
// checkurl project main.go
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var MAX_CHECK_NUM = 10

var succ_txt_name = "succ.txt"
var source_txt_name = "url.txt"
var lines []string
var check_word = ""

func read_source(file_name string) {
	file, err := os.Open(file_name)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return
}

func check_url(url string, limit_ch chan int) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if strings.Index(string(body), check_word) != -1 {
		file, err := os.OpenFile(succ_txt_name, os.O_APPEND|os.O_CREATE, os.ModeAppend)
		if err != nil {
			return
		}
		fmt.Println(url)
		file.WriteString(url + "\n")
		file.Close()
	}
	<-limit_ch
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("*************************************************************************")
	fmt.Println("*\t将要搜索的URL全放在跟这个程序同一个目录下的url.txt里一行一个\t*")
	fmt.Println("*\t国外的网站很有可能中间抓取会出错建议多运行几次之后去下重复\t*")
	fmt.Println("*\t\t在这个窗口的标题栏点右键选编辑里面可以粘贴\t\t*")
	fmt.Println("*************************************************************************")
	fmt.Println("输入关键词")
	data, _, _ := reader.ReadLine()
	check_word = string(data)
	read_source(source_txt_name)
	var limit_ch = make(chan int, MAX_CHECK_NUM)
	for _, url := range lines {
		limit_ch <- 1
		go check_url(url, limit_ch)
	}
}

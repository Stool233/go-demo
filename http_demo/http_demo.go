package http_demo

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)
import "github.com/valyala/fasthttp"

type DictType int

const (
	SPLIT_DICT DictType = 0
	STOP_WORD  DictType = 1
	SYNONYM    DictType = 2
)

func StartServer() {
	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/syncDictionaries":
			dictionaryHandler(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	if err := fasthttp.ListenAndServe(":10080", m); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func dictionaryHandler(ctx *fasthttp.RequestCtx) {
	GetDictionaries()
	ctx.SetContentType("text/plain; charset=utf8")
	ctx.SetBody([]byte("ok"))
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func GetDictionaries() {

	// http 获取远程词库
	url := "http://dev.api.tinya.huya.com:8080/dictionary/all"
	resp := new(Result)
	e := get(url, resp)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(resp)

	dictionaries := resp.Data

	// 分成三个词库
	splitDicts := make([]string, 0)
	stopwordDicts := make([]string, 0)
	synonymDicts := make([]string, 0)

	for _, dict := range dictionaries {
		if dict.Type == SPLIT_DICT {
			splitDicts = append(splitDicts, dict.Content)
		} else if dict.Type == STOP_WORD {
			stopwordDicts = append(stopwordDicts, dict.Content)
		} else if dict.Type == SYNONYM {
			synonymDicts = append(synonymDicts, dict.Content)
		}
	}

	fmt.Println(splitDicts)
	fmt.Println(stopwordDicts)
	fmt.Println(synonymDicts)

	// 转化格式后写到文件里
	isExist, err := PathExists("./tmp")
	if err != nil {
		fmt.Println(err)
	}
	if !isExist {
		err := os.MkdirAll("./tmp", os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	write2File("./tmp/splitDicts", strings.Join(splitDicts, "\n"))
	write2File("./tmp/stopword", strings.Join(stopwordDicts, "\n"))
	write2File("./tmp/synonym", strings.Join(synonymDicts, "\n"))

	// 复制到原文件
	CopyFile("./splitDicts", "./tmp/splitDicts")
	CopyFile("./stopword", "./tmp/stopword")
	CopyFile("./synonym", "./tmp/synonym")
}

func CopyFile(dst string, src string) (written int64, err error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func write2File(filePath string, data string) {
	var d1 = []byte(data + "\n")
	err := ioutil.WriteFile(filePath, d1, 0777) //写入文件(字节数组)
	if err != nil {
		fmt.Println(err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func get(url string, resp interface{}) error {
	var (
		code int
		body []byte
		e    error
	)
	if url == "" {
		return fmt.Errorf("invalid url")
	}
	if code, body, e = getRaw(url, 0); e != nil {
		return e
	}
	if !HTTPStatusOk(code) {
		return fmt.Errorf("%d:%s", code, body)
	}
	e = json.Unmarshal(body, resp)
	return e
}

func getRaw(url string, timeout time.Duration) (code int, body []byte, e error) {
	if timeout > 0 {
		return fasthttp.GetTimeout(nil, url, timeout)
	}
	return fasthttp.Get(nil, url)
}

func HTTPStatusOk(code int) bool {
	return fasthttp.StatusOK == code
}

type Result struct {
	Code      int          `json:"code"`
	RequestId string       `json:"requestId"`
	Message   string       `json:"message"`
	Data      []Dictionary `json:"data"`
}

type Dictionary struct {
	Id          int      `json:"id"`
	PlatformKey string   `json:"platformKey"`
	Content     string   `json:"content"`
	Note        string   `json:"note"`
	Type        DictType `json:"type"`
}

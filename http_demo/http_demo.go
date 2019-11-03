package http_demo

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
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

type Config struct {
	RemoteUrl         string
	Port              int16
	SplitDictFilePath string
	SplitDictFileName string
	StopWordFilePath  string
	StopWordFileName  string
	SynonymFilePath   string
	SynonymFileName   string
}

func StartServer(cfg *Config) {
	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/syncDictionaries":
			dictionaryHandler(ctx, cfg)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	if err := fasthttp.ListenAndServe(":10080", m); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func dictionaryHandler(ctx *fasthttp.RequestCtx, cfg *Config) {
	GetDictionaries(cfg)
	ctx.SetContentType("text/plain; charset=utf8")
	ctx.SetBody([]byte("ok"))
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func GetDictionaries(cfg *Config) {

	// http 获取远程词库
	url := cfg.RemoteUrl
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

	logrus.Info(splitDicts)
	logrus.Info(stopwordDicts)
	logrus.Info(synonymDicts)

	// 转化格式后写到文件里
	isExist, err := PathExists("./tmp")
	if err != nil {
		fmt.Println(err)
	}
	if !isExist {
		err := os.MkdirAll("./tmp", os.ModePerm)
		if err != nil {
			logrus.Error(err)
			return
		}
	}

	write2File("./tmp/"+cfg.SplitDictFileName, strings.Join(splitDicts, "\n"))
	write2File("./tmp/"+cfg.StopWordFileName, strings.Join(stopwordDicts, "\n"))
	write2File("./tmp/"+cfg.SynonymFileName, strings.Join(synonymDicts, "\n"))

	// 复制到原文件
	CopyFile(cfg.SplitDictFilePath, "./tmp/"+cfg.SplitDictFileName)
	CopyFile(cfg.StopWordFilePath, "./tmp/"+cfg.StopWordFileName)
	CopyFile(cfg.SynonymFilePath, "./tmp/"+cfg.SynonymFileName)
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

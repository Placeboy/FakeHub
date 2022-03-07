package ldc

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const COMPILEBASEURL = "http://kubernetes.tinylink.cn/linklab/compilev2/"

type CompileType string
type BoardType string

// BoardType
// esp32
// sky
// STM32F103C8

// CompileType
// esp32duino-virtual
// contiki-ng-virtual
// stm32-virtual

type CompileInfo struct {
	BoardType string
	CompileType string
	Filepath string
}

type CompilePara struct {
	Filehash string 				`json:"filehash"`
	BoardType string				`json:"boardType"`
	CompileType string 				`json:"compileType"`
}

// SHA1New_file 对文件进行SHA1加密
func SHA1New_file(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// 获取编译信息
func getCompileInfo() CompileInfo {
	var cInfo CompileInfo
	fmt.Println(flag.Args())
	for _, arg := range flag.Args() {
		arglist := strings.Split(arg, "=")
		if len(arglist) != 2 {
			arglist = strings.Split(arg, " ")
		}
		if len(arglist) != 2 {
			fmt.Printf(" unknown arg:%v\n", arg)
			continue
		}
		arg0 := strings.ToLower(arglist[0])
		if arg0 == "f" || arg0 == "F" {
			cInfo.Filepath = arglist[1]
		} else if arg0 == "b" || arg0 == "B" {
			cInfo.BoardType = arglist[1]
		} else if arg0 == "c" || arg0 == "C" {
			cInfo.CompileType = arglist[1]
		}
	}
	//fmt.Errorf("File name shoule be specified.\n")
	return cInfo
}

// 编译接口
// curl -v --request POST
// --form 'parameters={"filehash":"98bca5e26f43055315c81dc79cda22d29950f3d2", "boardType":"developerkit", "compileType":"alios"};type=application/json'
// --form "file=@alios.zip;type=application/octet-stream"
// http://kubernetes.tinylink.cn/linklab/compilev2/api/compile
func compileFile(cInfo CompileInfo) {
	// 构建parameters字符串
	para := CompilePara{
		Filehash: SHA1New_file(cInfo.Filepath),
		BoardType: cInfo.BoardType,
		CompileType: cInfo.CompileType,
	}
	para_json, err := json.Marshal(&para)
	if err != nil {
		log.Fatal("json Marshal fail", err.Error())
	}
	parameters := string(para_json)
	fmt.Println(parameters)

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	fileWriter, _ := bodyWriter.CreateFormFile("file", cInfo.Filepath)
	file, _ := os.Open(cInfo.Filepath)
	defer file.Close()
	io.Copy(fileWriter, file)

	// other form data
	extraParams := map[string]string{
		"parameters": parameters,
		"type":      "application/octet-stream",
	}
	for key, value := range extraParams {
		_ = bodyWriter.WriteField(key, value)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	//fmt.Println(contentType)
	resp, _ := http.Post(COMPILEBASEURL + "api/compile", contentType, bodyBuffer)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("parse body fail", err.Error())
	}
	fmt.Println(string(body))
}

// 获取编译后的文件(非阻塞)
// 测试通过
func downloadCompiledFile(cInfo CompileInfo, savedFile string) {
	//var file *os.File
	client := &http.Client{}
	req, err := http.NewRequest("GET", COMPILEBASEURL + "api/compile/nonblock", nil)
	if err != nil {
		log.Fatal("http request fail ", err.Error())
	}
	//req.Header.Set("Authorization", token)
	filehash := SHA1New_file(cInfo.Filepath)
	q := req.URL.Query()
	q.Add("filehash", filehash)
	q.Add("boardtype", cInfo.BoardType)
	q.Add("compiletype", cInfo.CompileType)
	req.URL.RawQuery = q.Encode()
	//fmt.Println(req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("get response fail ", err.Error())
	}
	//fmt.Println(resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("parse body fail", err.Error())
	}
	fmt.Println(string(body))

	err = ioutil.WriteFile(savedFile, body, 0666)
	if err != nil {
		log.Fatal("write file fail", err.Error())
	}
}
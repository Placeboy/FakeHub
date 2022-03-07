package ldc

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/Placeboy/FakeHub/vs"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"os/signal"
	"time"

	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const LDCBASEURL = "http://kubernetes.tinylink.cn/linklab/device-control-v2/"

// gid+index 与 deviceid的映射
var gdm map[string][]string
//const WEBSOCKETURL = "wss://kubernetes.wss.tinylink.cn/linklab/device-control-v2/user-service/api/ws"
// 测试版: ws://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/ws
//const LOCALHOSTURL = "http://localhost/"

// 登录接口
func login() (string, error) {
	formData := make(map[string]string)
	formData["id"] = "UserTest"
	formData["password"] = "6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918"
	jsonStr, err := json.Marshal(formData)
	if err != nil {
		log.Fatal("json Marshal fail", err.Error())
		return "", err
	}
	reader := strings.NewReader(string(jsonStr))
	resp, err := http.Post(LDCBASEURL + "login-authentication/user/login","raw", reader)
	if err != nil {
		log.Fatal("post form fail", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("parse body fail", err.Error())
		return "", err
	}
	//fmt.Println(string(body))
	var result LoginResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal("json Unmarshal fail ", err.Error())
		return "", err
	}
	//fmt.Println(result)
	//token :=
	fmt.Println("--- Login success!")
	return result.Data["token"], nil
}

// 显示系统支持的设备种类
func listBoard(token string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", LDCBASEURL + "user-service/api/board/list", nil)
	if err != nil {
		log.Fatal("http request fail ", err.Error())
	}
	req.Header.Set("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("get response fail ", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("parse body fail", err.Error())
	}

	//fmt.Println(string(body))
	var result ListBoardResult
	err = json.Unmarshal(body, &result)
	//fmt.Println(result.Code, result.Msg)
	if err != nil {
		log.Fatal("json Unmarshal fail ", err.Error())
	}
	boards := result.Data["boards"]
	fmt.Printf("%-30s\t%-20s\n", "Boardname", "Boardtype")
	for _, val := range boards {
		fmt.Printf("%-30s\t%-20s\n", val["boardname"], val["boardtype"])
	}
	//fmt.Println(result.Data["boards"])
	//fmt.Println(result)
	//token :=
	//return result.Data["token"]
}

// 烧写文件上传
// curl -v --request POST \
//  --header "Authorization: d6b804c0169fdbc0952dc8ef54a2a147d059438e70ec03eee05762913801fd9d" \
//  --form 'parameters={"boardname":"ESP32DevKitC"};type=application/json' \
//  --form "file=@bin/ESP32DevKitC.bin;type=application/octet-stream" \
//  http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file

func uploadFile(boardname, filepath, token string) string {
	// 构建parameters字符串
	para := UploadPara{
		BoardName: boardname,
	}
	parameters, _ := stringifyJsonStruct(para)
	//fmt.Println(parameters)

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	fileWriter, _ := bodyWriter.CreateFormFile("file", filepath)
	file, _ := os.Open(filepath)
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
	//token, _ := login()
	//fmt.Println(contentType)
	//resp, _ := http.Post(LDCBASEURL + "file-cache/api/file", contentType, bodyBuffer)
	client := &http.Client{}
	req, err := http.NewRequest("POST", LDCBASEURL + "file-cache/api/file", bodyBuffer)
	if err != nil {
		log.Fatal("http request fail ", err.Error())
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", contentType)
	resp, err := client.Do(req)
	//client.Head("")
	//client.Post(LDCBASEURL + "file-cache/api/file", contentType, bodyBuffer)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("parse body fail", err.Error())
	}
	//fmt.Println(string(body))
	var result UploadFileResult
	err = json.Unmarshal(body, &result)
	//fmt.Println(result.Code, result.Msg)
	if err != nil {
		log.Fatal("json Unmarshal fail ", err.Error())
	}
	fmt.Println("--- Upload File success!")
	return result.Data["filehash"]
	//fmt.Println(result.Data["boards"])
}

// 烧写任务提交
//# 不包含pid字段的调用方式
//curl -v --request POST \
//--header "Authorization: faff67ef90c9f2181a58f4bd983fe4dbfd38e8dc32035aad2396ea2ef98b21c3" \
//--data '{"tasks":[{"boardname":"ArduinoMega2560","deviceid":"/dev/ArduinoMega2560-8","runtime":30,"filehash":"eb3920b037e505b19c9a0ce0d8f28ae56f5ed28d9f70830ed22e11fd07d01c82","clientid":"ClientTest","taskindex":1},{"boardname":"ESP32DevKitC","deviceid":"/dev/ESP32DevKitC-0","runtime":30,"filehash":"6ec0d4238b7164784a62f4b163c712c480b45b2a931ed6c2c6b00e4c66890ca1","clientid":"ClientTest","taskindex":2}]}' \
//http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/device/burn

func submitTasks(taskNum, runtime int, boardname, filehash, token string) string {
	tasks := make([]*Task, taskNum)
	for i := 0; i < taskNum; i++ {
		task := NewTask()
		task.BoardName = boardname
		task.FileHash = filehash
		task.Runtime = runtime
		task.TaskIndex = i+1
		tasks[i] = task
	}
	//task_str, _ := stringifyJsonStruct(tasks)
	//fmt.Println(task_str)
	para := SubmitPara{
		Tasks: tasks,
	}
	parameters, err := json.Marshal(para)
	if err != nil {
		log.Fatal("json Marshal fail", err.Error())
	}
	//fmt.Println(string(parameters))
	client := &http.Client{}
	req, err := http.NewRequest("POST", LDCBASEURL + "user-service/api/device/burn", bytes.NewBuffer(parameters))
	if err != nil {
		log.Fatal("http request fail ", err.Error())
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", token)
	resp, err := client.Do(req)
	//client.Head("")
	//client.Post(LDCBASEURL + "file-cache/api/file", contentType, bodyBuffer)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("parse body fail", err.Error())
	}
	//fmt.Println(string(body))
	var result SubmitTaskResult
	err = json.Unmarshal(body, &result)
	//fmt.Println(result.Code, result.Msg)
	if err != nil {
		log.Fatal("json Unmarshal fail ", err.Error())
	}
	fmt.Println("--- Submit tasks success!")
	return result.Data["groupid"]
}

func runWebSocket(gid string, taskNum int, token string) {
	done := make(chan struct{})
	var addr = flag.String("addr", "kubernetes.wss.tinylink.cn", "http service address")
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/linklab/device-control-v2/user-service/api/ws"}
	log.Printf("connecting to %s", u.String())

	header := http.Header{}
	header.Set("Authorization", token)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	//done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			//log.Printf("recv: %s", message)
			HandleWebSocketMsg(gid, message, done, taskNum)
		}
	}()

	//ticker := time.NewTicker(time.Second)
	//defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		//case t := <-ticker.C:
		//	err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		//	if err != nil {
		//		log.Println("write:", err)
		//		return
		//	}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func stopWebSocket(done chan<- struct{}) {
	done <- struct{}{}
	fmt.Println("--- Stop WebSocket")
}

// json结构体转json字符串
func stringifyJsonStruct(jsonStruct interface{}) (string, error) {
	jsonBytes, err := json.Marshal(&jsonStruct)
	if err != nil {
		log.Fatal("json Marshal fail", err.Error())
		return "", err
	}
	return string(jsonBytes), nil
}

//var cur_gid string
// 只处理cur_gid的消息
func HandleWebSocketMsg(cur_gid string, message []byte, done chan<- struct{}, taskNum int) {

	//log.Printf("recv: %s", message)
	msgBody := NewWebSocketMsg()
	//var msgBody WebSocketMsg
	json.Unmarshal(message, msgBody)
	//fmt.Println(msgBody.Data)
	gid := msgBody.Data.GroupID
	if gid != cur_gid {
		// 直接忽略掉
		return
	}
	idx := msgBody.Data.TaskIndex
	switch msgBody.Data.Type {
	case "TaskAllocateMsg":
		// 此时,Hub应当有groupid 与 taskindex的配置信息
		// 将其对应替换为deviceid

		// 使CreateConfig不再阻塞
		//gid := msgBody.Data.GroupID
		//idx := msgBody.Data.TaskIndex
		//cur_gid = gid
		deviceid := msgBody.Data.Data["deviceid"]
		fmt.Printf("--- Allocated deviceid = %s\n", deviceid)
		if _, ok := gdm[gid]; !ok {
			gdm[gid] = make([]string, taskNum)
		}
		gdm[gid][idx-1] = deviceid
		//gdm[gid][idx] = deviceid
		//CreateDeviceConfig(gid, idx, dc)
	case "TaskLogMsg":
		msg := msgBody.Data.Msg
		timestamp := msgBody.TimeStamp
		if msg[:13] != "LinkLab-2022-" {
			HandleUserOutput(msg, timestamp)
			return
		}
		res, err := vs.HandleDeviceRequest(msg, gdm[gid][idx-1])
		if err != nil {
			HandleUserOutput(msg, timestamp)
			return
		}
		log.Printf("--- Command = %s", msg) //后面把这里替换为向LDC发送消息
		log.Printf("--- Result = %s", res) //后面把这里替换为向LDC发送消息

	case "TaskEndRunMsg":
		//time.Sleep(time.Second) // 给其他任务结束的时间
		stopWebSocket(done)
	default:
		//HandleUserOutput(message)
	}
	//if msgBody.Data["type"] == "TaskEndRunMsg" {
	//	time.Sleep(time.Second) // 给其他任务结束的时间
	//	// 清除当前gid的信息
	//
	//	stopWebSocket(done)
	//}
	//fmt.Println(msgBody.Data)
}

// 用户输出交给WebIDE
func HandleUserOutput(message string, timestamp int64) {
	tm := time.Unix(timestamp/1e9, 0)
	//fmt.Println(tm.Format("2006-01-02 03:04:05 PM"))
	log.Printf("--- %s: UserMsg --- %s", tm.Format("2006-01-02 15:04:05"), message)
}


func Init() {
	gdm = make(map[string][]string)
}

//func main() {
//	fmt.Println(SHA1New_file("led.zip"))
//}

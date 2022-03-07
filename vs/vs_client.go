package vs

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/Placeboy/FakeHub/rsa"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)
const VSBASEURL = "http://localhost:8080/"

type CreateRequest struct {
	DeviceID string `form:`
}

type DataRequest struct {
	DeviceID string `form:"deviceid"`
	DataType string `form:"datatype"`
}

type DeviceLog struct {
	DeviceID string `form:"id"`
	LogMsg	string `form:"msg"`
}

type DataResult struct {
	Data int `json:"data"`
}
// 发送读取数据的请求，已测试通过
func GetData(id, datatype string) int {
	var dr DataRequest
	dr.DataType = datatype
	dr.DeviceID = id

	client := &http.Client{}
	req, err := http.NewRequest("GET", VSBASEURL + "getSensorData", nil)
	if err != nil {
		log.Fatal("http request fail ", err.Error())
	}

	q := req.URL.Query()
	q.Add("deviceid", id)
	q.Add("datatype", datatype)
	req.URL.RawQuery = q.Encode()

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
	var result DataResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal("json Unmarshal fail ", err.Error())
	}
	return result.Data
	//fmt.Println(string(body))
	//return string(body)
	// 发送数据给LDC，这里用打印数据代替
}

// 功能：接收LDC发来的设备日志，进行解密操作
func HandleDeviceRequest(msg, deviceid string) (string, error) {
	//return msg, nil
	//var dl DeviceLog
	var data int
	//c.Bind(&dl)
	//id := dl.DeviceID
	//msg := dl.LogMsg
	// 对msg进行解密
	//prefix := msg[:13]
	//fmt.Println(prefix)
	//if prefix != "LinkLab-2022-" {
	//	c.JSON(http.StatusForbidden, gin.H{"error": "Not Valid Request!"})
	//	return nil
	//}
	cryptograph := msg[13:]
	//text := Decrypt(cryptograph)
	decodeBytes, err := base64.StdEncoding.DecodeString(cryptograph)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	origData, _ := rsa.RsaDecrypt(decodeBytes)
	//fmt.Println(string(origData))
	text := string(origData)
	return text, nil
	switch text {
	case "ReadLight":
		// 调用读取光照的函数
		data = GetData(deviceid, "light")
	case "ReadTemperature":
		data = GetData(deviceid, "temperature")
	case "ReadHumidity":
		data = GetData(deviceid, "humidity")
	default:
		//c.JSON(http.StatusForbidden, gin.H{"error": "Not Valid Request!"})
		return "", errors.New("not valid request")
	}
	//c.JSON(http.StatusOK, gin.H{"data": data})
	// 发送一个request
	res := strconv.Itoa(data)
	return res, nil
}

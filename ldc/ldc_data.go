package ldc

type WebSocketMsg struct {
	Code int						`json:"code"`
	Type string 					`json:"type"`
	TimeStamp int64 				`json:"timestamp"`
	Data *WebSocketMsgData			`json:"data"`
}

type WebSocketMsgData struct {
	GroupID string					`json:"groupid"`
	TaskIndex int					`json:"taskindex"`
	Type string						`json:"type"`
	Msg string						`json:"msg"`
	Data map[string]string			`json:"data"`
}

func NewWebSocketMsgData() *WebSocketMsgData {
	return &WebSocketMsgData{
		Data: nil,
	}
}

func NewWebSocketMsg() *WebSocketMsg {
	return &WebSocketMsg{
		Data: NewWebSocketMsgData(),
	}
}

type Result struct {
	Code int 							`json:"code"`
	Msg string							`json:"msg"`
}

type LoginResult struct {
	Result
	Data map[string]string				`json:"data"`
}

type ListBoardResult struct {
	Result
	Data map[string][]map[string]string	`json:"data"`
}

type UploadFileResult struct {
	Result
	Data map[string]string				`json:"data"`
}

type SubmitTaskResult struct {
	Result
	Data map[string]string
}

type UploadPara struct {
	BoardName string 					`json:"boardname"`
}

type Task struct {
	BoardName string 					`json:"boardname"`
	DeviceID string						`json:"deviceid"`
	Runtime int							`json:"runtime"`
	FileHash string						`json:"filehash"`
	ClientID string						`json:"clientid"`
	TaskIndex int 						`json:"taskindex"`
}

type SubmitPara struct {
	Tasks []*Task 						`json:"tasks"`
}

func NewTask() *Task {
	return &Task{
		"",
		"",
		60,
		"",
		"",
		-1,
	}
}
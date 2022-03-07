package ldc

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	login()
}

func TestListBoard(t *testing.T) {
	token, err := login()
	if err != nil {
		log.Fatal("login failed!", err.Error())
	}
	fmt.Println("Login success!")
	listBoard(token)
}

func TestUploadFile(t *testing.T) {
	token, err := login()
	if err != nil {
		log.Fatal("login failed!", err.Error())
	}
	fmt.Println("Login success!")
	//fmt.Printf("token = %s\n", token)
	fmt.Println(flag.Args())
	var boardname, filepath string
	for _, arg := range flag.Args() {
		arglist := strings.Split(arg, "=")
		if len(arglist) != 2 {
			arglist = strings.Split(arg, " ")
		}
		if len(arglist) != 2 {
			fmt.Printf(" unknown arg:%v\n", arg)
			continue
		}
		arg0 := arglist[0]
		if arg0 == "f" || arg0 == "F" {
			filepath = arglist[1]
		} else if arg0 == "b" || arg0 == "B" {
			boardname = arglist[1]
		}
	}
	uploadFile(boardname, filepath, token)
}

func TestSubmitTasks(t *testing.T) {
	token, err := login()
	if err != nil {
		log.Fatal("login failed!", err.Error())
	}
	//fmt.Println("Login success!")
	var boardname, filepath string
	var runtime int
	for _, arg := range flag.Args() {
		arglist := strings.Split(arg, "=")
		if len(arglist) != 2 {
			arglist = strings.Split(arg, " ")
		}
		if len(arglist) != 2 {
			fmt.Printf(" unknown arg:%v\n", arg)
			continue
		}
		arg0 := arglist[0]
		if arg0 == "f" || arg0 == "F" {
			filepath = arglist[1]
		} else if arg0 == "b" || arg0 == "B" {
			boardname = arglist[1]
		} else if arg0 == "r" || arg0 == "R" {
			runtime, err = strconv.Atoi(arglist[1])
			if err != nil {
				log.Fatal("runtime not integer!", err.Error())
				return
			}
		}
	}
	filehash := uploadFile(boardname, filepath, token)
	groupid := submitTasks(2, runtime, boardname, filehash, token)
	fmt.Println("groupid = ", groupid)
}

func TestRunWebSocket(t *testing.T) {
	Init()
	token, err := login()
	if err != nil {
		log.Fatal("login failed!", err.Error())
	}
	//fmt.Println("Login success!")
	var boardname, filepath string
	var runtime int
	for _, arg := range flag.Args() {
		arglist := strings.Split(arg, "=")
		if len(arglist) != 2 {
			arglist = strings.Split(arg, " ")
		}
		if len(arglist) != 2 {
			fmt.Printf(" unknown arg:%v\n", arg)
			continue
		}
		arg0 := arglist[0]
		if arg0 == "f" || arg0 == "F" {
			filepath = arglist[1]
		} else if arg0 == "b" || arg0 == "B" {
			boardname = arglist[1]
		} else if arg0 == "r" || arg0 == "R" {
			runtime, err = strconv.Atoi(arglist[1])
			if err != nil {
				log.Fatal("runtime not integer!", err.Error())
				return
			}
		}
	}
	filehash := uploadFile(boardname, filepath, token)
	groupid := submitTasks(2, runtime, boardname, filehash, token)
	//groupid :=

	runWebSocket(groupid, 2, token)
	//runtime := 80
	//time.Sleep(time.Second * time.Duration(runtime))
	//stopWebSocket(done)
}

//func
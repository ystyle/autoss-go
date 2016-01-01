package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/opesun/goquery"
)

/**
 * 服务器-结构体
 */
type Server struct {
	Remarks     string `json:"remarks"`     // 备注
	Server      string `json:"server"`      // 服务器ip
	Server_port int    `json:"server_port"` // 服务器端口
	Password    string `json:"password"`    // 密码
	Timeout     int    `json:"timeout"`     // 超时时间
	Local_port  int    `json:"local_port"`  // 本地代理端口
	Method      string `json:"method"`      // 加密方式
}

/**
 * 程序配置-结构体
 */
type Configs struct {
	Cmd        string `json:"cmd"`        // ss客户端所有在的位置
	Json       string `json:"json"`       // json配置文件
	Timeout    int    `json:"timeout"`    // 超时时间
	Local_port int    `json:"local_port"` // 本地代理端口
	Args       string `json:"args"`       // 参数
}

// 程序配置信息
var env = Configs{}

func main() {
	var url = "http://www.ishadowsocks.com"
	p, err := goquery.ParseUrl(url)
	servers := []Server{}
	if err != nil {
		panic(err)
	} else {
		Setup()
		divs := p.Find("section#free div.container div.row:last-child div")
		for i := 0; i < divs.Length(); i++ {
			div := divs.Eq(i)
			h4 := div.Find("h4")
			if len(h4) == 0 {
				continue
			}
			server_port, _ := strconv.Atoi(strings.Split(h4.Eq(1).Text(), ":")[1])
			server := Server{
				Remarks:     strings.Split(h4.Eq(0).Text(), ":")[0],
				Server:      strings.Split(h4.Eq(0).Text(), ":")[1],
				Server_port: server_port,
				Password:    strings.Split(h4.Eq(2).Text(), ":")[1],
				Timeout:     env.Timeout,
				Local_port:  env.Local_port,
				Method:      strings.Split(h4.Eq(3).Text(), ":")[1],
			}
			servers = append(servers, server)
		}
		data := readJson(env.Json) // 读ss的json文件
		var setting string
		if "windows" != runtime.GOOS {
			s := rand.NewSource(time.Now().UnixNano())
			r := rand.New(s)
			index := r.Intn(3)
			body, _ := json.Marshal(servers[index])
			setting = string(body)
		} else {
			data["configs"] = servers
			data["localPort"] = env.Local_port
			body, _ := json.Marshal(data)
			setting = string(body)
		}
		save(setting) // 保存信息
		startSS()     // 启动ss代理
	}
}

/**
 * 初始化程序配置信息
 * [Setup description]
 */
func Setup() {
	fmt.Println("初始化程序配置...")
	data := readJson("./config.json")
	env.Cmd = data["cmd"].(string)
	env.Json = data["json"].(string)
	env.Local_port = int(data["local_port"].(float64))
	env.Timeout = int(data["timeout"].(float64))
	env.Args = data["args"].(string)
}

/**
 * 读json文件
 * @param  {[type]} path string        路径
 * @return {[type]}      map接口
 */
func readJson(path string) map[string]interface{} {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	str := string(fd)
	var data interface{}
	json.Unmarshal([]byte(str), &data)
	return data.(map[string]interface{})
}

/**
 * 保存json
 * @param  {[type]} json string        json字符串
 */
func save(json string) {
	fmt.Println("保存SS配置...")
	var str = []byte(json)
	_ = ioutil.WriteFile(env.Json, str, 0666)
}

/**
 * 启动ss代理
 */
func startSS() {
	cmd := exec.Command(env.Cmd, env.Args)
	err := cmd.Start()
	if err != nil {
		fmt.Println("SS代理启动失败:", err)
	}
	fmt.Println("SS代理正在启动...")
}

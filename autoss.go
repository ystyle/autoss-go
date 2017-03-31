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
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"bytes"
	"net"
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
var env = &Configs{}

func main() {
	var url = "http://abc.ishadow.online/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("创建请求时出错错误", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Arch Linux kernel 4.6.5) AppleWebKit/537.36 (KHTML, like Gecko) Maxthon/4.0 Chrome/39.0.2146.0 Safari/537.36")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("发起请求时出错错误", err)
	}
	p, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatal("解析html时出现错误", err)
	} else {
		Setup()
		servers := []Server{}
		divs := p.Find(".portfolio-items .portfolio-item")
		for i := range divs.Nodes {
			div := divs.Eq(i)
			h4 := div.Find("h4")
			server_port, _ := strconv.Atoi(strings.Split(h4.Eq(1).Text(), "：")[1])
			server_name := h4.Eq(0).Find("span[id]").First().Text()

			server_ip := getServerIP(server_name,strconv.Itoa(server_port))
			server := Server{
				Remarks:     server_name,
				Server:      server_ip,
				Server_port: server_port,
				Password:    h4.Eq(2).Find("span[id]").First().Text(),
				Timeout:     env.Timeout,
				Local_port:  env.Local_port,
				Method:      strings.Split(h4.Eq(3).Text(), ":")[1],
			}
			servers = append(servers, server)
		}
		if len(servers) > 0 {
			jsonstr := readJson(env.Json) // 读ss的json文件
			var jsonData interface{}
			json.Unmarshal([]byte(jsonstr), &jsonData)
			data := jsonData.(map[string]interface{})
			if "windows" != runtime.GOOS {
				s := rand.NewSource(time.Now().UnixNano())
				r := rand.New(s)
				index := r.Intn(3)
				save(servers[index]) // 保存信息
			} else {
				data["configs"] = servers
				data["localPort"] = env.Local_port
				save(data) // 保存信息
			}
		}
		startSS() // 启动ss代理
	}
}

func getServerIP(serverName string,serverPort string) string {
	conn, _ := net.Dial("udp", serverName+":"+serverPort)
	defer conn.Close()
	return strings.Split(conn.RemoteAddr().String(), ":")[0]
}

/**
 * 初始化程序配置信息
 * [Setup description]
 */
func Setup() {
	fmt.Println("初始化程序配置...")
	data := readJson("./config.json")
	json.Unmarshal([]byte(data), &env)
}

/**
 * 读json文件
 * @param  {[type]} path string        路径
 * @return {[type]}      map接口
 */
func readJson(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	str := string(fd)
	return str
}

/**
 * 保存json
 * @param  {[type]} json string        json字符串
 */
func save(data interface{}) {
	fmt.Println("保存SS配置...")
	body, _ := json.Marshal(data)
	var out bytes.Buffer
	json.Indent(&out, body, "", "  ")
	_ = ioutil.WriteFile(env.Json, out.Bytes(), 0666)
}

/**
 * 启动ss代理
 */
func startSS() {
	cmdstr := env.Cmd + " " + env.Args
	list := strings.Split(cmdstr, " ")
	cmd := exec.Command(list[0], list[1:]...)
	err := cmd.Start()
	if err != nil {
		fmt.Println("SS代理启动失败:", err)
	}
	fmt.Println("SS代理正在启动...")
}

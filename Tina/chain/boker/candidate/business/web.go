package business

import (
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/Tinachain/Tina/chain/boker/candidate/header"
	log4plus "github.com/Tinachain/Tina/chain/boker/common/log4go"
	"github.com/Tinachain/Tina/chain/common"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	StatusJson = 600 // 解析Json格式失败
	StatusIp   = 601 // 解析Json格式失败
)

type ResponseCommon struct {
	Result int    `json:"result"`
	Msg    string `json:"msg"`
}

type RequestRegisterCandidate struct {
	Addr        string `json:"addrCandidate"` //注册节点地址
	Description string `json:"description"`   //节点注释
	Team        string `json:"team"`          //节点团队
	Name        string `json:"name"`          //节点名称
	Tickets     int64  `json:"tickets"`       //节点注册票数
}

type RequestVoteCandidate struct {
	AddrVoter     string `json:"addrVoter"`     //投票者账号
	AddrCandidate string `json:"addrCandidate"` //投票给的账号
	Tokens        int64  `json:"tokens"`        //投票票数
}

type RequestGetCandidate struct {
	AddrCandidate string `json:"addrCandidate"` //投票给的账号
}

type ResponseCandidate struct {
	Addr   string `json:"addrCandidate"` //候选人地址
	Ticket int64  `json:"tickets"`       //候选人票数
}

type ResponseCandidates struct {
	Candidates []*ResponseCandidate `json:"candidates"` //候选节点
}

type CandidateInfo struct {
	Addr        string `json:"addrCandidate"` //注册节点地址
	Description string `json:"description"`   //节点注释
	Team        string `json:"team"`          //节点团队
	Name        string `json:"name"`          //节点名称
	Tickets     int64  `json:"tickets"`       //节点注册票数
}

var (
	logFilePath = "./"
	logFileName = "download_encryption.log"
)

type WebManager struct {
	Listen  string
	UserGin *gin.Engine
}

/****gin需要处理的固定信息****/

//解决跨域问题
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "PUT, DELETE, POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func logerMiddleware() gin.HandlerFunc {
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)

	// 写入文件
	src, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("err", err)
	}

	// 实例化
	logger := logrus.New()
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	//设置输出
	logger.Out = src
	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	logger.AddHook(lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	}))

	return func(c *gin.Context) {
		//开始时间
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//请求方式
		reqMethod := c.Request.Method
		//请求路由
		reqUrl := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		//请求ip
		clientIP := c.ClientIP()

		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUrl,
		}).Info()

	}
}

//获取到客户端ip地址；
func (wm *WebManager) GetClientIP(c *gin.Context) string {
	reqIP := c.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

// RemoteIP 通过 RemoteAddr 获取 IP 地址， 只是一个快速解析方法。
func (wm *WebManager) RemoteIP(r *http.Request) string {
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func (wm *WebManager) HttpRegisterCandidate(c *gin.Context) {

	//获取到客户端IP地址
	clientIP := wm.GetClientIP(c)
	log4plus.Info("---->>>>HttpRegisterCandidate Request clientIP=%s****", clientIP)

	//分解数据
	var request RequestRegisterCandidate
	if err := c.BindJSON(&request); err != nil {
		c.JSON(StatusJson, gin.H{
			"result":  clientIP,
			"message": "Unmarshal Set Word Boby Json Failed",
		})

		log4plus.Error("---->>>>HttpRegisterCandidate Error=%s", err.Error())
		return
	}

	log4plus.Info("HttpRegisterCandidate: addrCandidate=%s\n description=%s\n team=%s\n name=%s\n tickets=%d\n",
		request.Addr,
		request.Description,
		request.Team,
		request.Name,
		request.Tickets)

	gCandidate.cm.RegisterCandidate(common.HexToAddress(request.Addr), request.Description, request.Team, request.Name, new(big.Int).SetInt64(request.Tickets))

	c.JSON(http.StatusOK, gin.H{
		"result":  http.StatusOK,
		"message": "OK",
	})

	/*defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	var request RequestRegisterCandidate
	if err := json.Unmarshal(body, &request); err != nil {

		log4plus.Error("HttpService err=%s", err.Error())
		header.HttpErrorEx(w, header.NewError(header.ErrorJsonParseError, "Unmarshal Json Failed"))
		return
	}
	log4plus.Info("HttpRegisterCandidate: addrCandidate=%s\n description=%s\n team=%s\n name=%s\n tickets=%d\n",
		request.Addr,
		request.Description,
		request.Team,
		request.Name,
		request.Tickets)

	gCandidate.cm.RegisterCandidate(common.HexToAddress(request.Addr), request.Description, request.Team, request.Name, new(big.Int).SetInt64(request.Tickets))

	bytes, _ := json.Marshal(&ResponseCommon{0, ""})
	w.Write(bytes)*/
}

func (wm *WebManager) HttpVoteCandidate(c *gin.Context) {

	/*defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	var request RequestVoteCandidate
	if err := json.Unmarshal(body, &request); err != nil {

		log4plus.Error("HttpService err=%s", err.Error())
		header.HttpErrorEx(w, header.NewError(header.ErrorJsonParseError, "Unmarshal Json Failed"))
		return
	}
	log4plus.Info("HttpVoteCandidate: addrVoter=%s\n addrCandidate=%s\n tokens=%d\n",
		request.AddrVoter,
		request.AddrCandidate,
		request.Tokens)

	gCandidate.cm.VoteCandidate(common.HexToAddress(request.AddrVoter), common.HexToAddress(request.AddrCandidate), new(big.Int).SetInt64(request.Tokens))
	bytes, _ := json.Marshal(&ResponseCommon{0, ""})
	w.Write(bytes)*/
}

func (wm *WebManager) HttpFlushEpoch(c *gin.Context) {

	/*defer r.Body.Close()

	//判断发起刷新周期的IP地址信息
	remoteIp := wm.RemoteIP(r)
	log4plus.Info("HttpFlushEpoch: remoteIp=%s\n", remoteIp)
	gCandidate.cm.FlushEpoch()
	bytes, _ := json.Marshal(&ResponseCommon{0, ""})
	w.Write(bytes)*/
}

func (wm *WebManager) HttpCurCandidates(c *gin.Context) {

	/*defer r.Body.Close()

	//判断发起刷新周期的IP地址信息
	remoteIp := wm.RemoteIP(r)
	log4plus.Info("HttpCurCandidates: remoteIp=%s\n", remoteIp)
	err, addrs, tickets := gCandidate.cm.CurCandidates()
	if err != nil {

	}

	var response ResponseCandidates
	for i, v := range addrs {

		candidate := ResponseCandidate{

			Addr:   v.String(),
			Ticket: tickets[i].Int64(),
		}
		response.Candidates = append(response.Candidates, &candidate)
	}

	bytes, _ := json.Marshal(response)
	w.Write(bytes)*/
}

func (wm *WebManager) HttpGetCandidate(c *gin.Context) {

	/*defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	var request RequestGetCandidate
	if err := json.Unmarshal(body, &request); err != nil {

		log4plus.Error("HttpGetCandidate err=%s", err.Error())
		header.HttpErrorEx(w, header.NewError(header.ErrorJsonParseError, "Unmarshal Json Failed"))
		return
	}

	//判断发起刷新周期的IP地址信息
	remoteIp := wm.RemoteIP(r)
	log4plus.Info("HttpGetCandidate: remoteIp=%s\n", remoteIp)
	candidate, err := gCandidate.cm.GetCandidate(request.AddrCandidate)
	if err != nil {

	}

	var response CandidateInfo
	response.Addr = request.AddrCandidate
	response.Description = candidate.Description
	response.Name = candidate.Name
	response.Team = candidate.Team
	response.Tickets = candidate.Tickets.Int64()

	bytes, _ := json.Marshal(response)
	w.Write(bytes)*/
}

func (wm *WebManager) startWeb() {

	//分组处理
	userGroup := wm.UserGin.Group("/user")
	{
		userGroup.GET("/RegisterCandidate", wm.HttpRegisterCandidate)
		userGroup.GET("/VoteCandidate", wm.HttpVoteCandidate)
		userGroup.GET("/CurCandidates", wm.HttpCurCandidates)
		userGroup.GET("/GetCandidate", wm.HttpGetCandidate)
		userGroup.GET("/FlushEpoch", wm.HttpFlushEpoch)
	}
	wm.UserGin.Run(wm.Listen)
}

func NewWeb(listen string) *WebManager {

	web := &WebManager{
		Listen: listen,
	}

	//启动Web的gin
	log4plus.Info("Create User Web Manager")
	web.UserGin = gin.Default()
	gin.SetMode(gin.ReleaseMode)
	web.UserGin.Use(logerMiddleware())
	web.UserGin.Use(Cors())

	go web.startWeb()

	return web
}

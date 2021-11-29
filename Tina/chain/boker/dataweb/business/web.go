package business

import (
	"fmt"
	"io/ioutil"
	_ "math/big"
	"net/http"
	"os"
	"path"
	"time"

	log4plus "github.com/Tinachain/Tina/chain/boker/common/log4go"

	_ "github.com/Tinachain/Tina/chain/boker/dataweb/header"
	"github.com/Tinachain/Tina/chain/common"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	ErrorOK = iota + 600
	JsonParseError
	TxFailed
)

type RequestSetWord struct {
	Word string `json:"word"` //写入内容
}

type JsonGetWordFromTx struct {
	TxHash string `json:"txHash"` //交易hash
}

type ResponseFile struct {
	Name     string `json:"name"`     //文件名称
	SaveTime string `json:"savetime"` //保存时间
	TxHash   string `json:"txHash"`   //交易Hash
}

type ResponseFileList struct {
	FileList []*ResponseFile `json:"files"` //文件列表
}

type WebManager struct {
	userListen    string
	UserGin       *gin.Engine
	channelListen string
	ChannelGin    *gin.Engine
	dp            *DetailProcess
}

var (
	logFilePath = "./"
	logFileName = "dataweb.log"
)

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

func (wm *WebManager) RequestSetWord(c *gin.Context) {

	//获取到客户端IP地址
	clientIP := wm.GetClientIP(c)
	log4plus.Info("---->>>>RequestSetWord Request clientIP=%s****", clientIP)

	//分解数据
	var request RequestSetWord
	if err := c.BindJSON(&request); err != nil {
		c.JSON(JsonParseError, gin.H{
			"result":  JsonParseError,
			"message": "Unmarshal Set Word Boby Json Failed",
		})

		log4plus.Error("---->>>>RequestSetWord Error=%s", err.Error())
		return
	}

	//设置文字
	err, txHash := wm.dp.RequestSetWord(request)
	if err != nil {

		c.JSON(TxFailed, gin.H{
			"result":  TxFailed,
			"message": err.Error(),
		})
	} else {

		c.JSON(http.StatusOK, gin.H{
			"hash": txHash,
		})
	}
}

func (wm *WebManager) RequestGetWord(c *gin.Context) {

	//获取到客户端IP地址
	clientIP := wm.GetClientIP(c)
	log4plus.Info("---->>>>RequestGetWord Request clientIP=%s****", clientIP)

	//分解数据
	var request JsonGetWordFromTx
	if err := c.BindJSON(&request); err != nil {
		c.JSON(JsonParseError, gin.H{
			"result":  JsonParseError,
			"message": "Unmarshal Get Word Boby Json Failed",
		})

		log4plus.Error("---->>>>RequestGetWord Error=%s", err.Error())
		return
	}

	//获取文字
	wordErr, word := wm.dp.RequestGetWord(common.HexToHash(request.TxHash))
	if wordErr != nil {

		//返回错误;
		c.JSON(TxFailed, gin.H{
			"result":  TxFailed,
			"message": wordErr.Error(),
		})
	} else {

		//返回成功;
		c.JSON(http.StatusOK, gin.H{
			"word": word,
		})
	}
}

func (wm *WebManager) RequestSetPic(c *gin.Context) {

	//获取到客户端IP地址
	clientIP := wm.GetClientIP(c)
	log4plus.Info("---->>>>RequestSetPic Request clientIP=%s****", clientIP)

	picFile, _, err := c.Request.FormFile("pic")
	if err != nil {
		c.JSON(JsonParseError, gin.H{
			"result":  JsonParseError,
			"message": "Not Found pic File Context " + err.Error(),
		})
	}
	if picFile != nil { // 记得及时关闭文件，避免内存泄漏
		defer picFile.Close()
	}

	picContent, err := ioutil.ReadAll(picFile)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":  false,
			"message": "Read pic File Failed " + err.Error(),
		})
	}

	//设置图片
	err, hash := wm.dp.RequestSetPic(picContent)
	if err != nil {
		c.JSON(TxFailed, gin.H{
			"result":  TxFailed,
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"hash": hash,
		})
	}
}

func (wm *WebManager) RequestGetPic(c *gin.Context) {

	//获取到客户端IP地址
	clientIP := wm.GetClientIP(c)
	log4plus.Info("---->>>>RequestGetPic Request clientIP=%s****", clientIP)

	//分解数据
	var request JsonGetWordFromTx
	if err := c.BindJSON(&request); err != nil {
		c.JSON(JsonParseError, gin.H{
			"result":  JsonParseError,
			"message": "Unmarshal Get Word Boby Json Failed",
		})

		log4plus.Error("---->>>>RequestGetPic Error=%s", err.Error())
		return
	}

	//设置图片
	err, data := wm.dp.RequestGetPic(common.HexToHash(request.TxHash))
	if err != nil {
		c.JSON(TxFailed, gin.H{
			"result":  TxFailed,
			"message": err.Error(),
		})
	}
	c.Writer.WriteString(string(data))
}

func (wm *WebManager) RequestSetFile(c *gin.Context) {

	//获取到客户端IP地址
	clientIP := wm.GetClientIP(c)
	log4plus.Info("---->>>>RequestSetFile Request clientIP=%s****", clientIP)

	fileHandle, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(JsonParseError, gin.H{
			"result":  JsonParseError,
			"message": "Not Found file Context " + err.Error(),
		})
	}
	if fileHandle != nil { // 记得及时关闭文件，避免内存泄漏
		defer fileHandle.Close()
	}

	fileContent, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":  false,
			"message": "Read File Context Failed " + err.Error(),
		})
	}

	//设置图片
	err, hash := wm.dp.RequestSetFile(fileContent)
	if err != nil {
		c.JSON(TxFailed, gin.H{
			"result":  TxFailed,
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result":  http.StatusOK,
			"message": "OK",
			"hash":    hash,
		})
	}
}

func (wm *WebManager) RequestGetFile(c *gin.Context) {

	//获取到客户端IP地址
	clientIP := wm.GetClientIP(c)
	log4plus.Info("---->>>>RequestGetFile Request clientIP=%s****", clientIP)

	//分解数据
	var request JsonGetWordFromTx
	if err := c.BindJSON(&request); err != nil {
		c.JSON(JsonParseError, gin.H{
			"result":  JsonParseError,
			"message": "Unmarshal Get Word Boby Json Failed",
		})

		log4plus.Error("---->>>>RequestGetFile Error=%s", err.Error())
		return
	}

	//设置图片
	err, data := wm.dp.RequestGetFile(common.HexToHash(request.TxHash))
	if err != nil {
		c.JSON(TxFailed, gin.H{
			"result":  TxFailed,
			"message": err.Error(),
		})
	}

	c.Writer.WriteString(string(data))
}

func (wm *WebManager) RequestFileList(c *gin.Context) {

	//获取到客户端IP地址
	clientIP := wm.GetClientIP(c)
	log4plus.Info("---->>>>RequestFileList Request clientIP=%s****", clientIP)

	//分解数据
	var request JsonGetWordFromTx
	if err := c.BindJSON(&request); err != nil {
		c.JSON(JsonParseError, gin.H{
			"result":  JsonParseError,
			"message": "Unmarshal Get Word Boby Json Failed",
		})

		log4plus.Error("---->>>>RequestGetFile Error=%s", err.Error())
		return
	}

	//设置图片
	//	err, data := wm.dp.RequestGetFile(common.StringToHash(request.TxHash))
	//	if err != nil {
	//		c.JSON(header.ErrTxFailed, gin.H{
	//			"result":  header.ErrTxFailed,
	//			"message": err.Error(),
	//		})
	//	}
}

func (wm *WebManager) RequestGetChannel(c *gin.Context) {

	wm.RequestGetFile(c)
}

func (wm *WebManager) startUser() {

	//分组处理
	userGroup := wm.UserGin.Group("/user")
	{
		userGroup.POST("/setWord", wm.RequestSetWord)
		userGroup.POST("/getWord", wm.RequestGetWord)

		userGroup.POST("/setPic", wm.RequestSetPic)
		userGroup.POST("/getPic", wm.RequestGetPic)

		userGroup.POST("/setFile", wm.RequestSetFile)
		userGroup.POST("/getFile", wm.RequestGetFile)

		userGroup.POST("/getFileList", wm.RequestFileList)
	}
	wm.UserGin.Run(wm.userListen)
}

func (wm *WebManager) startChannel() {

	//分组处理
	channelGroup := wm.ChannelGin.Group("/channel")
	{
		channelGroup.GET("/getChannel", wm.RequestGetChannel)

	}
	wm.ChannelGin.Run(wm.channelListen)
}

func NewWeb(userlisten string) *WebManager {

	//创建对象
	web := &WebManager{
		userListen:    userlisten,
		channelListen: "",
		dp:            NewProcess(),
	}

	//启动Web的gin
	log4plus.Info("Create User Web Manager")
	web.UserGin = gin.Default()
	gin.SetMode(gin.DebugMode)
	web.UserGin.Use(logerMiddleware())
	web.UserGin.Use(Cors())

	//启动频道gin
	//	web.ChannelGin = gin.Default()
	//	web.ChannelGin.Use(logerMiddleware())
	//	web.ChannelGin.Use(Cors())

	//启动Web
	go web.startUser()
	//	go web.startChannel()

	return web
}

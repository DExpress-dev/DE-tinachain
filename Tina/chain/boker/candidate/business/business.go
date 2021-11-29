package business

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/Tinachain/Tina/chain/boker/candidate/header"
	log4plus "github.com/Tinachain/Tina/chain/boker/common/log4go"
	"github.com/Tinachain/Tina/chain/common"
)

type ConfigInfo struct {
	RPC               string //rpc的地址
	InterfaceAddr     string //
	InterfaceBaseAddr string //
	KeystoreFile      string //Keystore文件名
	Passwrod          string //Keystore的密码
	Listen            string //web监听信息
}

type CandidateBusiness struct {
	rpcAddress        string //rpc的地址
	keystoreFile      string //Keystore文件名
	passwrod          string //Keystore的密码
	listen            string //web监听信息
	interfaceAddr     common.Address
	interfaceBaseAddr common.Address
	web               *WebManager
	cm                *header.Client
}

var gCandidate CandidateBusiness

func (p *CandidateBusiness) loadConfig() bool {

	log4plus.Info("(p *CandidateBusiness) loadConfig")
	cfgFile, err := os.Open("config.json")
	if err != nil {
		log4plus.Error("(p *CandidateBusiness) loadConfig Failed Open File Error %s", err.Error())
		return false
	}
	defer cfgFile.Close()

	var config ConfigInfo
	cfgBytes, _ := ioutil.ReadAll(cfgFile)
	jsonErr := json.Unmarshal(cfgBytes, &config)
	if jsonErr != nil {
		log4plus.Error("(p *CandidateBusiness) loadConfig json.Unmarshal Failed %s", jsonErr.Error())
		return false
	}
	log4plus.Info("(p *CandidateBusiness) read config.json-> \n RPC=%s\n InterfaceAddr=%s\n InterfaceBaseAddr=%s\n KeystoreFile=%s\n Passwrod=%s\n Listen=%s\n",
		config.RPC,
		config.InterfaceAddr,
		config.InterfaceBaseAddr,
		config.KeystoreFile,
		config.Passwrod,
		config.Listen)

	p.rpcAddress = config.RPC
	p.interfaceAddr = common.HexToAddress(config.InterfaceAddr)
	p.interfaceBaseAddr = common.HexToAddress(config.InterfaceBaseAddr)
	p.passwrod = config.Passwrod
	p.keystoreFile = config.KeystoreFile
	p.listen = config.Listen

	log4plus.Info("(p *CandidateBusiness) loadConfig Success")
	return true
}

func Init() {

	//加载配置信息
	log4plus.Info("Start Load Config")
	if !gCandidate.loadConfig() {
		log4plus.Error("loadConfig Failed Exit Program")
		return
	}

	//启动rpc客户端
	var err error
	log4plus.Info("Create Client Manager")
	gCandidate.cm, err = header.NewClient(gCandidate.rpcAddress, gCandidate.interfaceAddr, gCandidate.interfaceBaseAddr, gCandidate.keystoreFile, gCandidate.passwrod)
	if err != nil {
		log4plus.Error("Create Client Failed Exit Program rpcAddress=%s interfaceAddr=%s interfaceBaseAddr=%s", gCandidate.rpcAddress,
			gCandidate.interfaceAddr.String(),
			gCandidate.interfaceBaseAddr.String())
		return
	}

	//启动Web
	log4plus.Info("Start Web Manager")
	gCandidate.web = NewWeb(gCandidate.listen)
}

func Run() {

	for {
		time.Sleep(time.Duration(10) * time.Second)
	}
}

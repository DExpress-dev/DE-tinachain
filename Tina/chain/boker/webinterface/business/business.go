package business

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	log4plus "github.com/Tinachain/Tina/chain/boker/common/log4go"
	"github.com/Tinachain/Tina/chain/boker/webinterface/header"
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

type InterfaceBusiness struct {
	rpcAddress        string //rpc的地址
	keystoreFile      string //Keystore文件名
	passwrod          string //Keystore的密码
	listenPort        string //web监听信息
	interfaceAddr     common.Address
	interfaceBaseAddr common.Address
	web               *WebManager
	cm                *header.Client
}

var gInterface InterfaceBusiness

func (i *InterfaceBusiness) loadConfig() bool {

	log4plus.Info("(i *InterfaceBusiness) loadConfig")
	cfgFile, err := os.Open("config.json")
	if err != nil {
		log4plus.Error("(i *InterfaceBusiness) loadConfig Failed Open File Error %s", err.Error())
		return false
	}
	defer cfgFile.Close()

	var config ConfigInfo
	cfgBytes, _ := ioutil.ReadAll(cfgFile)
	jsonErr := json.Unmarshal(cfgBytes, &config)
	if jsonErr != nil {
		log4plus.Error("(i *InterfaceBusiness) loadConfig json.Unmarshal Failed %s", jsonErr.Error())
		return false
	}
	log4plus.Info("(i *InterfaceBusiness) read config.json-> \n RPC=%s\n InterfaceAddr=%s\n InterfaceBaseAddr=%s\n KeystoreFile=%s\n Passwrod=%s\n Listen=%s\n",
		config.RPC,
		config.InterfaceAddr,
		config.InterfaceBaseAddr,
		config.KeystoreFile,
		config.Passwrod,
		config.Listen)

	i.rpcAddress = config.RPC
	i.interfaceAddr = common.HexToAddress(config.InterfaceAddr)
	i.interfaceBaseAddr = common.HexToAddress(config.InterfaceBaseAddr)
	i.passwrod = config.Passwrod
	i.keystoreFile = config.KeystoreFile
	i.listenPort = config.Listen

	log4plus.Info("(i *InterfaceBusiness) loadConfig Success")
	return true
}

func Init() {

	//加载配置信息
	log4plus.Info("Start Load Config")
	if !gInterface.loadConfig() {
		log4plus.Error("loadConfig Failed Exit Program")
		return
	}

	//启动rpc客户端
	var err error
	log4plus.Info("Create Client Manager")
	gInterface.cm, err = header.NewClient(gInterface.rpcAddress,
		gInterface.interfaceAddr,
		gInterface.interfaceBaseAddr,
		gInterface.keystoreFile,
		gInterface.passwrod)
	if err != nil {
		log4plus.Error("Create Client Failed Exit Program rpcAddress=%s interfaceAddr=%s interfaceBaseAddr=%s", gInterface.rpcAddress,
			gInterface.interfaceAddr.String(),
			gInterface.interfaceBaseAddr.String())
		return
	}

	//启动Web
	log4plus.Info("Start Web Manager")
	gInterface.web = NewWeb(gInterface.listenPort)
}

func Run() {

	for {
		time.Sleep(time.Duration(10) * time.Second)
	}
}

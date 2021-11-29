package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Tinachain/Tina/chain/boker/common"
	log4plus "github.com/Tinachain/Tina/chain/boker/common/log4go"
	"github.com/Tinachain/Tina/chain/boker/dataweb/business"
)

//版本号
var (
	ver     string = "1.0.2"
	exeName string = ""
	pidFile string = ""
)

type Flags struct {
	Help    bool
	Version bool
}

func (f *Flags) Init() {
	flag.BoolVar(&f.Help, "h", false, "help")
	flag.BoolVar(&f.Version, "V", false, "show version")
}

func (f *Flags) Check() (needReturn bool) {
	flag.Parse()

	if f.Help {
		flag.Usage()
		needReturn = true
	} else if f.Version {
		verString := exeName + " Version: " + ver + "\r\n"
		fmt.Println(verString)
		needReturn = true
	}

	return needReturn
}

var flags *Flags = &Flags{}

func init() {
	flags.Init()
	exeName = getExeName()
	pidFile = GetCurrentDirectory() + "/" + exeName + ".pid"
}

func getExeName() string {
	ret := ""
	ex, err := os.Executable()
	if err == nil {
		ret = filepath.Base(ex)
	}
	return ret
}

func setLog() {
	logJson := "log.json"
	set := false
	if bExist := common.PathExist(logJson); bExist {
		if err := log4plus.SetupLogWithConf(logJson); err == nil {
			set = true
		}
	}

	if !set {
		fileWriter := log4plus.NewFileWriter()
		exeName := getExeName()
		fileWriter.SetPathPattern("./log/" + exeName + "-%Y%M%D.log")
		log4plus.Register(fileWriter)
		log4plus.SetLevel(log4plus.DEBUG)
	}
}

func main() {
	needReturn := flags.Check()
	if needReturn {
		return
	}
	setLog()
	defer log4plus.Close()
	log4plus.Info("DataWeb Server Version=%s", ver)
	business.Init()
	business.Run()
}

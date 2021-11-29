package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "math/big"
	"os"
	"path/filepath"

	log4plus "github.com/Tinachain/Tina/chain/boker/common/log4go"
	ethcommon "github.com/Tinachain/Tina/chain/common"
	"github.com/Tinachain/Tina/chain/contracts/deploy/common"
	"github.com/Tinachain/Tina/chain/contracts/deploy/common/config"
	"github.com/Tinachain/Tina/chain/contracts/deploy/common/tinachain"
)

//版本号
var (
	Version = "1.0.6"
)

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

func GetInterfaceBaseContract() string {
	contractAddress := ""
	interfaceBaseFile := tinachain.ContractInterfaceBase + ".contract"
	file, err := os.Open(interfaceBaseFile)
	if err != nil {
		return ""
	}
	defer file.Close()

	addressBytes, _ := ioutil.ReadAll(file)
	contractAddress = string(addressBytes)
	return contractAddress
}

func SaveInterfaceBaseContract(address string) error {
	interfaceBaseFile := tinachain.ContractInterfaceBase + ".contract"
	file, err := os.Create(interfaceBaseFile)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(address)
	return nil
}

type DeployConfigEntry struct {
	File      string
	Contracts []string
}

type DeployConfig struct {
	Deploy []DeployConfigEntry
}

var interfaceBaseAddress ethcommon.Address

func LoadDeployConfig() *DeployConfig {

	cfg := &DeployConfig{}

	cfgFile, err := os.Open(config.GetInstance().BokerchainSolFolder + "/" + "deploy.json")
	if err != nil {
		return cfg
	}
	defer cfgFile.Close()

	cfgBytes, _ := ioutil.ReadAll(cfgFile)

	json.Unmarshal(cfgBytes, cfg)
	return cfg
}

func Deploy() {

	log4plus.Info("--->>>Deploy Contracts to Tinachain")

	client, err := tinachain.NewClient(config.GetInstance().BokerchainRpc)
	if err != nil {
		log4plus.Error("tinachain.NewClient Failed Err=%s", err.Error())
		return
	}

	err = client.Unlock(config.GetInstance().BokerchainAdminKeystore, config.GetInstance().BokerchainAdminPassword)
	if err != nil {
		log4plus.Error("client.Unlock Failed Err=%s", err.Error())
		return
	}

	log4plus.Info("Deploy Connect to Tinachain Account=%s", config.GetInstance().BokerchainAdminKeystore)

	fJs, err := os.Create("contract.js")
	if err != nil {
		log4plus.Error("Create contract.js Failed Err=%s", err.Error())
		return
	}
	defer fJs.Close()
	format := `
Ware.%s = web3.eth.contract(%s).at('%s', function(error, contract){
    if(!error){
        console.log('%s at : ' + contract.address);
        }
    else
        console.error('%s at error : ' + error);
})

				`
	fJs.WriteString("var Ware = new Object();\n")

	deployCfg := LoadDeployConfig()
	var managerAddress ethcommon.Address
	if config.GetInstance().BokerchainManagerAddress != "" {
		managerAddress = tinachain.HexToAddress(config.GetInstance().BokerchainManagerAddress)
	}

	for _, entry := range deployCfg.Deploy {

		filePath := config.GetInstance().BokerchainSolFolder + "/" + entry.File
		if !common.PathExist(filePath) {

			log4plus.Error("Not Found Contracts File %s", filePath)
			continue
		}

		compiledContracts, err := client.ContractCompile(filePath)
		if err != nil {

			log4plus.Error("client.ContractCompile Failed Err=%s", err.Error())
			return
		}

		for _, contractName := range entry.Contracts {

			log4plus.Info("Deploying Contract Name=%s", contractName)
			contract := compiledContracts[contractName]
			if nil == contract {

				log4plus.Error("Not Found Contract Name=%s", contractName)
				return
			}

			if contractName == tinachain.ContractManager {

				err = client.ContractDeploy(contract)
				if err != nil {

					log4plus.Error("Deploy Contract %s Fail Err=%s", contractName, err.Error())
					return
				}
				managerAddress = contract.Address
				fJs.WriteString(fmt.Sprintf(format, contract.Name, contract.Abi, contract.Address.String(), contract.Name, contract.Name))

				log4plus.Info("Deployed Contract %s OK Address %s", contractName, managerAddress.String())
				err = client.AtManager(managerAddress.String())
				if err != nil {

					log4plus.Error("Deploy Program Set Manager Contract Address %s Fail Err=%s", managerAddress.String(), err.Error())
					return
				}
				log4plus.Info("Set Manager Contract Address %s OK", managerAddress.String())

			} else if contractName == tinachain.ContractInterfaceBase {

				if managerAddress.String() == tinachain.ZeroAddressString {

					log4plus.Error("Manager Address is Empty")
					return
				}

				err = client.ContractDeploy(contract, managerAddress)
				if err != nil {

					log4plus.Error("Contract %s Deploy Fail Err=%s", contractName, err.Error())
					return
				}
				log4plus.Info("Contract %s Deploy Ok", contractName)

				fJs.WriteString(fmt.Sprintf(format, contract.Name, contract.Abi, contract.Address.String(), contract.Name, contract.Name))

				log4plus.Info("Contract %s Manager Address Setring...", contractName)
				err = client.SetContract(contract.Name, contract.Address)
				if err != nil {

					log4plus.Error("Contract %s Manager Address Set Fail", contractName, err.Error())
					return
				}
				log4plus.Info("Contract %s SetContract Ok", contract.Name)

				contractOld := GetInterfaceBaseContract()
				if "" != contractOld {

					log4plus.Info("CancelBaseContracts Old Contract %s", contractOld)
					/*err = client.CancelBaseContracts(tinachain.HexToAddress(contractOld))
					if err != nil {

						log4plus.Error("CancelBaseContracts Err=%s Address=%s", err.Error(), contractOld)
						return
					}*/
				}

				log4plus.Info("SetBaseContracts Address %s", contract.Address.String())
				err = client.SetSystemBaseContracts(contract.Address)
				if err != nil {

					log4plus.Error("SetBaseContracts Address %s Err=%s", contract.Address.String(), err.Error())
					return
				}

				log4plus.Info("SetBaseContracts Address %s Ok", contract.Address.String())
				SaveInterfaceBaseContract(contract.Address.String())

				interfaceBaseAddress = contract.Address

			} else {
				if managerAddress.String() == tinachain.ZeroAddressString {

					log4plus.Error("Manager Address Is Empty")
					return
				}

				err = client.ContractDeploy(contract, managerAddress)
				if err != nil {

					log4plus.Error("Deploy Contract %s Fail Err=%s", contractName, err.Error())
					return
				}
				log4plus.Info("Deploy Contract %s Ok", contract.Name)

				fJs.WriteString(fmt.Sprintf(format, contract.Name, contract.Abi, contract.Address.String(), contract.Name, contract.Name))

				log4plus.Info("SetContract %s Manager Address Setting...", contract.Name)
				err = client.SetContract(contract.Name, contract.Address)
				if err != nil {

					log4plus.Error("SetContract %s Manager Address Fail Err=%s", contract.Name, err.Error())
					return
				}
				log4plus.Info("SetContract %s Manager Address %s Ok \n", contract.Name, contract.Address.String())
			}

		}
	}

	return
}

func main() {

	setLog()
	defer log4plus.Close()
	log4plus.Info("Deploy Version=%s", Version)

	config.Initialize()
	Deploy()
}

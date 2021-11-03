package header

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Tinachain/Tina/chain/accounts/abi/bind"
	log4plus "github.com/Tinachain/Tina/chain/boker/common/log4go"
	"github.com/Tinachain/Tina/chain/boker/protocol"
	"github.com/Tinachain/Tina/chain/common"
	"github.com/Tinachain/Tina/chain/contracts/boker_interface"
	"github.com/Tinachain/Tina/chain/ethclient"
)

type Client struct {
	rpc                string
	keyJson            string
	password           string
	interfaceAddr      common.Address
	interfaceBaseAddr  common.Address
	from               common.Address
	signer             bind.SignerFn
	EthClient          *ethclient.Client
	bokerInterface     *boker_contract.BokerInterface
	bokerInterfaceBase *boker_contract.BokerInterfaceBase
}

func NewClient(rpc string, interfaceAddr common.Address, interfaceBaseAddr common.Address, keystoreFile string, password string) (c *Client, err error) {

	c = &Client{
		rpc:               rpc,
		interfaceAddr:     interfaceAddr,
		interfaceBaseAddr: interfaceBaseAddr,
	}
	log4plus.Info("header.NewClient c.interfaceAddr=%s c.interfaceBaseAddr=%s c.rpc=%s", c.interfaceAddr.String(), c.interfaceBaseAddr.String(), c.rpc)

	c.EthClient, err = ethclient.Dial(rpc)
	if err != nil {

		log4plus.Error("header.NewClient Dial Failed rpc=%s : %s", rpc, err.Error())
		return nil, err
	}

	fKeystore, _ := os.Open(keystoreFile)
	keystore, _ := ioutil.ReadAll(fKeystore)

	err = c.Unlock(string(keystore), password)
	if err != nil {

		log4plus.Error("header.NewClient Unlock Failed : %s", err.Error())
		return nil, err
	}
	log4plus.Info("header.NewClient Unlock Success")

	//得到Interface的类
	c.bokerInterface, err = boker_contract.NewBokerInterface(c.interfaceAddr, c.EthClient)
	if err != nil {
		log4plus.Error("header.NewClient NewBokerInterface err=%s", err.Error())
		return nil, err
	}
	log4plus.Info("header.NewClient NewBokerInterface Success")

	//得到InterfaceBase的类
	c.bokerInterfaceBase, err = boker_contract.NewBokerInterfaceBase(c.interfaceBaseAddr, c.EthClient)
	if err != nil {
		log4plus.Error("header.NewClient NewBokerInterfaceBase err=%s", err.Error())
		return nil, err
	}
	log4plus.Info("header.NewClient NewBokerInterfaceBase Success")

	return c, nil
}

func (c *Client) GetPendingNonce() (nonce uint64, err error) {
	return c.EthClient.PendingNonceAt(context.Background(), c.from)
}

func (c *Client) Unlock(keyJson, password string) (err error) {

	opts, err := bind.NewTransactor(strings.NewReader(keyJson), password)
	if err != nil {
		return err
	}

	c.signer = opts.Signer
	c.from = opts.From
	c.keyJson = keyJson
	c.password = password

	return nil
}

func (c *Client) NewOpts() *bind.TransactOpts {

	return &bind.TransactOpts{
		From:   c.from,
		Signer: c.signer,
	}
}

func (c *Client) NewCalls() *bind.CallOpts {

	return &bind.CallOpts{
		From: c.from,
	}
}

//设置文字
func (c *Client) SetWord(word string) (hash common.Hash, err error) {

	if len(word) <= 0 {
		return common.Hash{}, errors.New("Word Length is Zero")
	}

	if len(word) > int(protocol.MaxWordSize) {
		return common.Hash{}, errors.New("SetWord length too more than MaxWordSize(1MB)")
	}

	if nil == c.EthClient {
		return common.Hash{}, errors.New("interface not specified, call AtInterface first!")
	}

	hash, err = c.EthClient.SetWord(context.Background(), word)
	if err != nil {
		return common.Hash{}, err
	}

	return hash, nil
}

//获取文字
func (c *Client) GetWord(hash common.Hash) (string, error) {

	if nil == c.EthClient {
		return "", errors.New("interface not specified, call AtInterface first!")
	}

	word, err := c.EthClient.GetWord(context.Background(), hash)
	if err != nil {
		return "", err
	}

	return word, nil
}

//判断文件是否存在
func fileExist(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			return false
		}
	}
	return true
}

//设置图片
func (c *Client) SetPic(fileName string) (hash common.Hash, err error) {

	//判断文件是否存在
	if exist := fileExist(fileName); !exist {
		return common.Hash{}, errors.New("Picture Not Found")
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {

		errTxt := fmt.Sprintf("Picture File Reading error %s", err.Error())
		return common.Hash{}, errors.New(errTxt)
	}

	//判断eth对象是否存在
	if nil == c.EthClient {
		return common.Hash{}, errors.New("interface not specified, call AtInterface first!")
	}

	hash, err = c.EthClient.SetData(context.Background(), data)
	if err != nil {
		return common.Hash{}, err
	}

	return hash, nil
}

//设置图片
func (c *Client) SetPicFromData(data []byte) (hash common.Hash, err error) {

	if len(data) <= 0 {
		return common.Hash{}, errors.New("Picture File Length is Zero")
	}

	if len(data) > int(protocol.MaxWordSize) {
		return common.Hash{}, errors.New("SetPicFromData length too more than MaxWordSize(1MB)")
	}

	//判断eth对象是否存在
	if nil == c.EthClient {
		return common.Hash{}, errors.New("interface not specified, call AtInterface first!")
	}

	hash, err = c.EthClient.SetData(context.Background(), data)
	if err != nil {
		return common.Hash{}, err
	}

	return hash, nil
}

//获取图片
func (c *Client) GetPic(hash common.Hash) ([]byte, error) {

	//判断eth对象是否存在
	if nil == c.EthClient {
		return []byte{}, errors.New("interface not specified, call AtInterface first!")
	}

	data, err := c.EthClient.GetData(context.Background(), hash)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

//设置数据
func (c *Client) SetData(data []byte) (hash common.Hash, err error) {

	//判断eth对象是否存在
	if nil == c.EthClient {
		return common.Hash{}, errors.New("interface not specified, call AtInterface first!")
	}

	if len(data) >= protocol.MaxDataSize {
		return common.Hash{}, errors.New("SetData length too more than MaxDataSize(1MB)")
	}

	hash, err = c.EthClient.SetData(context.Background(), data)
	if err != nil {
		return common.Hash{}, err
	}

	return hash, nil
}

//获取数据
func (c *Client) GetData(hash common.Hash) ([]byte, error) {

	//判断eth对象是否存在
	if nil == c.EthClient {
		return []byte{}, errors.New("interface not specified, call AtInterface first!")
	}

	data, err := c.EthClient.GetData(context.Background(), hash)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

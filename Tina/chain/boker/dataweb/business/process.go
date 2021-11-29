package business

import (
	"math/big"

	"github.com/Tinachain/Tina/chain/common"
)

type DetailProcess struct {
}

func (dp *DetailProcess) RequestGetNumber(hash common.Hash) (err error, number *big.Int) {

	number, err = gInterface.cm.GetTx(hash)
	return err, number
}

func (dp *DetailProcess) RequestSetWord(request RequestSetWord) (error, common.Hash) {

	//设置文字
	hash, err := gInterface.cm.SetWord(request.Word)
	return err, hash
}

func (dp *DetailProcess) RequestGetWord(hash common.Hash) (error, string) {

	//获取文字
	word, err := gInterface.cm.GetWord(hash)
	return err, word
}

func (dp *DetailProcess) RequestSetPic(data []byte) (error, common.Hash) {

	//设置图片
	hash, err := gInterface.cm.SetPicFromData(data)
	return err, hash
}

func (dp *DetailProcess) RequestGetPic(hash common.Hash) (error, []byte) {

	//获取图片
	data, err := gInterface.cm.GetPic(hash)
	return err, data
}

func (dp *DetailProcess) RequestSetFile(data []byte) (error, common.Hash) {

	//设置文件
	hash, err := gInterface.cm.SetData(data)
	return err, hash
}

func (dp *DetailProcess) RequestGetFile(hash common.Hash) (error, []byte) {

	//获取文件
	data, err := gInterface.cm.GetData(hash)
	return err, data
}

func (dp *DetailProcess) RequestGetChannel(hash common.Hash) (error, []byte) {

	return dp.RequestGetFile(hash)
}

func NewProcess() *DetailProcess {

	//创建对象
	dp := &DetailProcess{}

	return dp
}

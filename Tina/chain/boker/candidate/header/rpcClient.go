package header

import (
	"context"
	"errors"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/Tinachain/Tina/chain/accounts/abi/bind"
	log4plus "github.com/Tinachain/Tina/chain/boker/common/log4go"
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

//注册候选人
func (c *Client) RegisterCandidate(Addr common.Address, Description string, Team string, Name string, Tickets *big.Int) error {

	/*log4plus.Info("(c *Client) RegisterCandidate-> Addr=%s Description=%s", Addr.String(), Description)

	if nil == c.bokerInterface {

		log4plus.Error("RegisterCandidate Failed bokerInterface is nil")
		return errors.New("RegisterCandidate Failed bokerInterface is nil")
	}

	tx, err := c.bokerInterface.RegisterOtherCandidate(c.NewOpts(), Addr, Description, Team, Name, Tickets)
	if err != nil {

		log4plus.Error("(c *Client) RegisterCandidate Failed")
		return protocol.ErrPair
	}

	ctx := context.Background()
	_, err = bind.WaitMined(ctx, c.EthClient, tx)
	if err != nil {

		log4plus.Error("bind.WaitMined :%s", err.Error())
		return err
	}*/

	return nil
}

//当前候选人
func (c *Client) CurCandidates() (error, []common.Address, []*big.Int) {

	log4plus.Info("(c *Client) CurCandidate->")

	if nil == c.bokerInterface {

		log4plus.Error("CurCandidate Failed bokerInterfaceBase is nil")
		return errors.New("CurCandidate Failed bokerInterfaceBase is nil"), nil, nil
	}

	candidates, err := c.bokerInterface.GetCandidates(c.NewCalls())
	if err != nil {

		log4plus.Error("(c *Client) RegisterCandidate Failed")
		return err, nil, nil
	}
	log4plus.Info("(c *Client) GetCandidates Result %d", len(candidates.Addresses))

	return nil, candidates.Addresses, candidates.Tickets
}

//当前候选人
func (c *Client) GetCandidate(Addr string) (struct {
	Description string
	Team        string
	Name        string
	Tickets     *big.Int
}, error) {

	log4plus.Info("(c *Client) GetCandidate->")

	ret := new(struct {
		Description string
		Team        string
		Name        string
		Tickets     *big.Int
	})

	if nil == c.bokerInterfaceBase {

		log4plus.Error("GetCandidate Failed bokerInterfaceBase is nil")
		return *ret, errors.New("GetCandidate Failed bokerInterfaceBase is nil")
	}

	var err error
	*ret, err = c.bokerInterface.GetCandidate(c.NewCalls(), common.HexToAddress(Addr))
	if err != nil {

		log4plus.Error("(c *Client) RegisterCandidate Failed")
		return *ret, err
	}
	log4plus.Info("(c *Client) GetCandidate->Description=%s Name=%s Team=%s Tickets=%d", ret.Description, ret.Name, ret.Team, ret.Tickets.Int64())

	return *ret, nil
}

//投票
func (c *Client) VoteCandidate(addrVoter common.Address, addrCandidate common.Address, tokens *big.Int) error {

	log4plus.Info("(c *Client) Vote->")

	if nil == c.bokerInterface {

		log4plus.Error("Vote Failed bokerInterface is nil")
		return errors.New("Vote Failed bokerInterface is nil")
	}

	/*tx, err := c.bokerInterface.Vote(c.NewOpts(), addrVoter, addrCandidate, tokens)
	if err != nil {

		log4plus.Error("(c *Client) Vote Failed")
		return err
	}

	ctx := context.Background()
	_, err = bind.WaitMined(ctx, c.EthClient, tx)
	if err != nil {

		log4plus.Error("bind.WaitMined :%s", err.Error())
		return err
	}*/

	return nil
}

//刷新周期数据，让投票马上起作用
func (c *Client) FlushEpoch() error {

	log4plus.Info("(c *Client) FlushEpoch->")

	if nil == c.bokerInterfaceBase {

		log4plus.Error("FlushEpoch Failed bokerInterfaceBase is nil")
		return errors.New("FlushEpoch Failed bokerInterfaceBase is nil")
	}

	/*tx, err := c.bokerInterfaceBase.BokerInterfaceBaseTransactor.RotateOtherVote(c.NewOpts())
	if err != nil {

		log4plus.Error("(c *Client) FlushEpoch Failed %s", err.Error())
		return err
	}

	ctx := context.Background()
	_, err = bind.WaitMined(ctx, c.EthClient, tx)
	if err != nil {

		log4plus.Error("bind.WaitMined :%s", err.Error())
		return err
	}
	log4plus.Info("(c *Client) FlushEpoch tx=%s", tx.Hash().String())*/

	return nil
}

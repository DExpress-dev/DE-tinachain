运行流程

准备流程：
1：同步时间：
	/usr/sbin/ntpdate te cn.pool.ntp.org > / > /dev/null 2>&1

2：设置可运行路径
	vi /etc/profile
	PATH=$PATH:/projects/tinachain
	source /etc/profile

第一步：初始化创世文件
	geth --datadir "/projects/tinachain/node" init genesis.json


第二步：启动geth
	nohup geth --nodiscover  \
	--maxpeers 3 \
	--identity "tinachain" \
	--rpc \
	--rpcaddr 0.0.0.0 \
	--rpccorsdomain "*" \
	--rpccorsdomain "http://localhost:8000" \
	--datadir "/projects/tinachain/node" \
	--port 30304 \
	--rpcapi "db,eth,net,web3" \
	--networkid 98765 &


第三步：进入geth控制台
	geth attach ipc:/projects/tinachain/node/geth.ipc

第四步：创建账号(0x4c386449ddee1ca0eb0f71a2daa5543cedb9ac83)
	personal.newAccount()

第五步：设置帐号解锁（这里使用假定账号、密码）
	personal.unlockAccount(eth.accounts[0], "123456", 0)

第六步：设置自己为验证人
	miner.setLocalValidator()

第七步：设置验证人（这里使用假定账号、票数）
	eth.addValidator(eth.accounts[0], 10000)

第八步：启动挖矿
	miner.start()

第九步：设置链所属
	eth.setStockManager(eth.accounts[0])

第十步：获得链所属
	eth.getStockManager()

创建新用户
	personal.newAccount()

第十一步：设置股权
	eth.stockSet(eth.accounts[0], 10000)

	eth.stockGet(eth.accounts[0])

第十二步：查看当前区块高度
	eth.blockNumber

第九步：终止挖矿
miner.stop()

eth.sendTransaction({from: eth.accounts[0], to: "1ceed12b103d9e76fea7de5410f08684eb5ef113", value: web3.toWei(1000,"ether")})


新增特殊指令：

刻字：
eth.setWord("测试的刻字内容")
eth.getWord("0x8901173dc329a98311c96786b4e42f9bc43f68f38a227a3d509acc5facf5926d")

图片：
eth.setPicture("/projects/tinachain/1.jpg")
eth.getPicture("0x26635445ae6e1f20bc2a7ed5be45c3a0b7e847e1c79167c9b1564fe77ef72094", "/projects/tina")


第十一步：设置合约为基础合约
eth.setBaseContracts("0xff2e5867f89e7be22e8c4a3cd9fb59bfd31ce681", 1, "[{\"constant\":true,\"inputs\":[],\"name\":\"cfoAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ceoAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getLogSize\",\"outputs\":[{\"name\":\"size\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"logKeyDefault\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"assgineTokenPeriodDefault\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCEO\",\"type\":\"address\"}],\"name\":\"setCEO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"index\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCOO\",\"type\":\"address\"}],\"name\":\"setCOO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"enable\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getLog\",\"outputs\":[{\"name\":\"level\",\"type\":\"uint8\"},{\"name\":\"time\",\"type\":\"uint256\"},{\"name\":\"key\",\"type\":\"string\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"v1\",\"type\":\"uint256\"},{\"name\":\"v2\",\"type\":\"uint256\"},{\"name\":\"v3\",\"type\":\"uint256\"},{\"name\":\"remarks\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"eventType\",\"type\":\"uint256\"},{\"name\":\"addrFrom\",\"type\":\"address\"},{\"name\":\"addrTo\",\"type\":\"address\"},{\"name\":\"eventValue1\",\"type\":\"uint256\"},{\"name\":\"eventValue2\",\"type\":\"uint256\"}],\"name\":\"fireUserEvent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"checkAssignToken\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"assginedTokensPerPeriodDefault\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCFO\",\"type\":\"address\"}],\"name\":\"setCFO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"clearLog\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"implAddress\",\"type\":\"address\"}],\"name\":\"setImpl\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"logLevel\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"impl\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"assignToken\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cooAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"level\",\"type\":\"uint8\"}],\"name\":\"setLevel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"logEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"}]")


运行接口JS（交易类型，不能直接使用控制台执行，用控制台执行，有可能无法执行出结果，最好使用remix执行）
Ware.BokerInterface.registerCandidate.sendTransaction("test", "chainware", "demo", 10000,{from:eth.accounts[0]})

Ware.BokerInterfaceBase.setCandidates.sendTransaction(eth.accounts[0], 10000, {from:eth.accounts[0]})

Ware.BokerInterfaceBase.getCandidates()

Ware.BokerInterface.getCandidate(eth.accounts[0])


搭建子节点

第一步：查询创世服务器信息
admin.nodeInfo

返回
{
  enode: "enode://97730b179a8c1af4c21613c78fcc9c607f2ce9614ec927ab2f788e0c1faab2c3a4f532fa326220a589b59af421ea30c0e7224c6c0837b66a2ff4e391314da679@[::]:30304?discport=0",
  id: "97730b179a8c1af4c21613c78fcc9c607f2ce9614ec927ab2f788e0c1faab2c3a4f532fa326220a589b59af421ea30c0e7224c6c0837b66a2ff4e391314da679",
  ip: "::",
  listenAddr: "[::]:30304",
  name: "Geth/bokerchain--rpc/v1.7.4-stable/linux-amd64/go1.9.2",
  ports: {
    discovery: 0,
    listener: 30304
  },
  protocols: {
    eth: {
      difficulty: 28684,
      genesis: "0xb0a53b645ea52f7971d06ada36e004a9961105d4acb013f7e2d337a4ce90e280",
      head: "0x73d7519a2edc78be4b12c8137d92266f234631ebf976a08bb79d55614b0b0c20",
      network: 66666
    }
  }
}


第二步：在子节点中添加创世服务器信息

admin.addPeer("enode://b767e6d26287eccd8eec72a63e1e116d31aad781aa79a5dfa9b92424317ca98bb714f4f9772141657a6fc9dc0300a7cf749e6c7ad3b269313d66f4ca661d82e3@192.168.22.135:30303")
此处注意IP信息和端口信息


在创世服务器中做相应的操作，并等待同步完成。

查看节点同步信息
admin.peers

查看块同步信息
eth.blockNumber


























assigntoken.checkAssignToken()


查询信息

1：查询区块数量
eth.blockNumber

2：查询区块内容
eth.getBlock(1) 

3：查询交易内容（这里使用假定交易的Hash）
eth.getTransaction("0x63dfdfa9e14f187c9b35a5c1f5bcc5e4401bcf81c2682dd832a0f09da8bea17c")

eth.getTransaction("0x82a9326290d246c59405e47154feeac0cbbe14ec1ed70b5f3f7f047438cad4d5")

4：发布合约
loadScript("assigntoken.js")












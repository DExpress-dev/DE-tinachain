# 子节点升级为出块节点工具

## 代码路径
```python
    Dposchain_V2\boker\candidate
```

## 普通节点升级为超级节点步骤

- [x]  假设1：创世节点：`0x1CEEd12B103D9e76fEa7De5410f08684Eb5eF113`
- [x]  假设2：普通节点：`0x8978BfdD91fD9c52D60737b5359ee544836Fe65b`
- [x]  假设3：创世节点已经运行并能正常出块。

#### 启动子节点

1：进入子节点同步区块：
```python
>admin.addPeer("enode://e7293e1a1d30606d2393c093d31e7a6bd85612842194ee5cd74fefd08781e73a1b20fe0468aa8f2ad21abec2a4ba65657ffaae2f9df3a825bd03e58a86258441@192.168.3.81:30303")
```
2：解锁Coinbase（***为密码）
```python
>personal.unlockAccount(eth.accounts[0], ***, 0)
```
3：启动子节点挖矿
```python
>miner.start()
```
注意：当前由于子节点不是出块节点，因此子节点不会进行出块。

#### 启动Candidate工具
```python
1：设置Candidate配置信息，例如：
{
    "RPC":"http://192.168.3.81:8545",
    "InterfaceAddr":"0x518d1f05A5E67eEAb6c9914a2042F0F53c6E3583",
    "InterfaceBaseAddr":"0xbe691343cF3054a950189D4838BbB3423E5D5d1A",
    "KeystoreFile":"keystore.json",
    "Passwrod":"123456",
    "Listen":"0.0.0.0:80"
}
```
其中：
"InterfaceAddr" 是合约Interface的地址
"InterfaceBaseAddr": 是合约InterfaceBase的地址
这两个地址可以在创世节点发布合约后在contract.js中查找到

"Listen" 是Web Service监听的IP和端口。
```python
2：启动Candidate
./candidate
```
启动Candidate后，Candidate会创建一个Web Service，并监听 Listen设置的端口。

#### 注册候选人
```python
1：通过Postman或者其他web客户端发送注册创世节点候选人（***为Candidate的IP和监听端口）：
http://***/RegisterCandidate
```
POST的Body为（其中addrCandidate为创世节点的Coinbase）：
```python
{
    "addrCandidate":"0x1ceed12b103d9e76fea7de5410f08684eb5ef113",
    "description":"First Blockchain Node",
    "team":"fxh7622",
    "name":"张超",
    "tickets":5000
}
```
```python
2：相同方法，注册子节点Coinbase为候选人（其中addrCandidate为创世节点的Coinbase）
http://***/RegisterCandidate
```
POST的Body为
```python
{
    "addrCandidate":"8978bfdd91fd9c52d60737b5359ee544836fe65b",
    "description":"Second Blockchain Node",
    "team":"zgc7622",
    "name":"张轶钦",
    "tickets":2000
}
```
```python
3：设置好后可以通过指令进行查询

http://***/CurCandidates
```
```python
4：启动周期转换
http://***/FlushEpoch
```
周期转换好后，子节点账号`"8978bfdd91fd9c52d60737b5359ee544836fe65b"`将作为第二个出块节点。

#### 日志解释

当日志中出现以下内容：
```python
INFO [03-14|16:06:50] Mint Block Consume Timer                 second=0
INFO [03-14|16:06:55] (d *Dpos) CheckProducer                  now=1584173215 firstTimer=0 producer=0x8978BfdD91fD9c52D60737b5359ee544836Fe65b d.signer=0x1CEEd12B103D9e76fEa7De5410f08684Eb5eF113
ERROR[03-14|16:06:55] CheckProducer Failed CoinBase is`t Current Producer err="invalid current producer"
INFO [03-14|16:06:55] (d *Dpos) Finalize                       Number=201 Time=1584173215 Validator=0x8978BfdD91fD9c52D60737b5359ee544836Fe65b ProducerRewards=8183172125000000000
INFO [03-14|16:06:55] (d *Dpos) verifySeal                     header.Time=1584173215 producer=0x8978BfdD91fD9c52D60737b5359ee544836Fe65b firstTimer=0
INFO [03-14|16:07:00] (d *Dpos) CheckProducer                  now=1584173220 firstTimer=0 producer=0x1CEEd12B103D9e76fEa7De5410f08684Eb5eF113 d.signer=0x1CEEd12B103D9e76fEa7De5410f08684Eb5eF113
INFO [03-14|16:07:00] (d *Dpos) Finalize                       Number=202 Time=1584173220 Validator=0x1CEEd12B103D9e76fEa7De5410f08684Eb5eF113 ProducerRewards=8183172125000000000
INFO [03-14|16:07:00] Mint Block Consume Timer                 second=0
```
上面日志中出现：
```python
INFO [03-14|16:06:55] (d *Dpos) CheckProducer                  now=1584173215 firstTimer=0 producer=0x8978BfdD91fD9c52D60737b5359ee544836Fe65b d.signer=0x1CEEd12B103D9e76fEa7De5410f08684Eb5eF113
```
说明，当前应该出块的节点为 `[0x8978BfdD91fD9c52D60737b5359ee544836Fe65b]`，但是当前的Coinbase为 `[0x1CEEd12B103D9e76fEa7De5410f08684Eb5eF113]`，因此这个时间点不能由 `[0x1CEEd12B103D9e76fEa7De5410f08684Eb5eF113]` 来进行出块。
因此下一句错误提示为正常情况。
```python
ERROR[03-14|16:06:55] CheckProducer Failed CoinBase is`t Current Producer err="invalid current producer"
```
```python
INFO [03-14|16:07:00] (d *Dpos) CheckProducer                  now=1584173220 firstTimer=0 producer=0x1CEEd12B103D9e76fEa7De5410f08684Eb5eF113 d.signer=0x1CEEd12B103D9e76fEa7De5410f08684Eb5eF113
```
说明：当前出块的节点为 `[0x1CEEd12B103D9e76fEa7De5410f08684Eb5eF113]` 和当前的Coinbase一致，可以出块。
至此将子节点 `[0x8978BfdD91fD9c52D60737b5359ee544836Fe65b]` 升级为出块节点完成。
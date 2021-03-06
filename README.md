
### Tina链（Tinachain）
Tina链是一个专门为服务于文字、图片、文章以及文件保存而开发的垂直型区块链平台，用户可以使用Tina链提供的接口将自己的相关信息永久保存到Tina链上。同时为了保证信息的安全和保密，Tina链对于信息的保存提供了可选加密功能。从而保证了用户信息的安全性。

Tinachain is a vertical blockchain platform specifically designed for preservation of text, images, articles and files. With rpc interfaces provided by Tinachain, it is convenient for users to save personal information permanently. Meanwhile, Tinachain provides optional encryption function for users to ensure the security and privacy of their information.


### Tina链特色（System characteristic）

#### 数据切片（Data slice）：
Tina链在对需求用户需要存储数据时，为了防止数据块过大，Tina链首先采取了数据切片，将大的存储数据切片成多个更小的数据块（一个数据切片大小为32MB）。

When Tina chain needs users to store data, in order to prevent the data block from being too large, the Tina chain first adopts data slicing, slicing the large stored data into multiple smaller data blocks (a data slice size is 32MB).

#### 数据冗余（Data redundancy）：
Tina链在对数据进行切片后，为了防止存储供应商在存储数据切片是出现个别数据切片丢失、损坏，造成需求用户无法最终获取完整的数据。因此Tina链使用了FEC方式，对于数据切片加入了一定数量的冗余数据切片。当数据切片丢失后，利用冗余切片可以对丢失的数据切片进行恢复。

After the Tina chain slices the data, in order to prevent the loss or damage of individual data slices when the storage supplier stores the data slices, the demanding users cannot finally obtain the complete data. Therefore, the Tina chain uses the FEC method and adds a certain number of redundant data slices to the data slices. When data slices are lost, redundant slices can be used to recover the lost data slices.

#### 数据加密（Data encryption）：
为了保证用户数据的安全性，Tina链对每一个数据切片丢进行了AES加密，从而保证了用户数据在存储提供商处也无法进行数据泄露。

In order to ensure the security of user data, the Tina chain performs AES encryption on each data slice, thereby ensuring that user data cannot be leaked at the storage provider.

#### 切片多向分发（Multi-directional distribution of slices）：
为了保障不因为存储供应商成为单点故障节点（存储设备损坏），造成用户数据损坏无法恢复。Tina链采取切片多向分发方式，会将一个用户的多个切片数据分发给不同的存储供应商（由撮合机来进行撮合匹配选择）。

In order to ensure that the storage provider does not become a single point of failure (damage of the storage device), causing user data to be damaged and cannot be recovered. The Tina chain adopts a slice multi-directional distribution method, which distributes multiple slice data of a user to different storage providers (matching machine for matching and matching selection).

#### 多副本存储（Multi-copy storage）：
为了避免用户数据只保存一个副本而带来的存储损坏隐患，Tina链提供了多副本存储功能，用户可以根据自己的实际需求选择存储的副本数量。Tina链将按照指定副本数量选择合适的存储供应商（同一存储供应商不会有相同切片的多重副本），从而达到“不将鸡蛋放在一个篮子里”的效果。

In order to avoid the hidden danger of storage damage caused by only one copy of user data, the Tina chain provides a multi-copy storage function, and users can choose the number of copies to store according to their actual needs. The Tina chain will select the appropriate storage supplier according to the specified number of copies (the same storage supplier will not have multiple copies of the same slice), so as to achieve the effect of "not putting eggs in one basket".

#### 多地切片获取（Multi-place slice acquisition）：
用户在从Tina链中获取完整的存储数据时，会首先获取到包含冗余切片的所有切片列表，获取器会根据这个切片列表中的信息从不同的存储供应商处获取相应的切片数据。最终由获取器在本地合并出完整的用户数据。

When users obtain complete storage data from the Tina chain, they will first obtain a list of all the slices containing redundant slices, and the obtainer will obtain corresponding slice data from different storage providers based on the information in this slice list. Finally, the complete user data is merged locally by the getter.


### Tina链系统架构（System architecture）
![Image text](https://github.com/DExpress-dev/DE-tinachain/blob/main/Tina/image/Architecture.png)

### Tina链文件存储流程图（System flow chart）
![Image text](https://github.com/DExpress-dev/DE-tinachain/blob/main/Tina/image/process.png)

### Tina链组织图（System Organization chart）
![Image text](https://github.com/DExpress-dev/DE-tinachain/blob/main/Tina/image/combination.png)

### Tina链微信小程序截图（Screenshot of Tina chain wechat applet）
![Image text](https://github.com/DExpress-dev/DE-tinachain/blob/main/Tina/image/Mini_Programs1.jpg)
![Image text](https://github.com/DExpress-dev/DE-tinachain/blob/main/Tina/image/Mini_Programs2.jpg)

### 目录（Folders）

#### [chain](https://github.com/DExpress-dev/DE-tinachain/tree/master/chain)
    采用DPOS共识实现的基础链代码（基于ethereum 1.7.4版本）
    Main chain code, implementing DPOS.

#### [contracts](https://github.com/DExpress-dev/DE-tinachain/tree/master/contracts)
    采用Solidity编写的基础合约代码
    Basic contract code in solidity.

#### [explorer](https://github.com/DExpress-dev/DE-tinachain/tree/master/explorer)
    区块链浏览器以及文章上传页面代码
    Basic contract code in solidity.

#### [NFT](https://github.com/DExpress-dev/DE-tinachain/tree/main/NFT)
    NFT相关的合约代码
    NFT related contract code.

### 公众号
![logo](https://github.com/DExpress-dev/DE-tinachain/blob/main/Tina/image/wechat.png)

### 小程序
![logo](https://github.com/DExpress-dev/DE-tinachain/blob/main/Tina/image/Mini_Programs.jpg)

# codechain

​	codechain基于 Tendermint共识机制的代码链，采用go语言实现，并且依赖于Tendermint环境，建议在ubuntu16.04下开发，

- ​	为了便于后续的开发，首先建立起git、golang、tendermint代码库、vscode的开发环境，参考  [建立开发环境](https://github.com/little51/codechain/blob/master/getting-start.md)  。
- 为了验证Tendermint的四个节点的运行情况，参考 [运行环境组网](https://github.com/little51/codechain/blob/master/making-testnet.md) 。
- 为了验证Tendermint更多的节点加入，参考 [Seed模式组网](https://github.com/little51/codechain/blob/master/making-seeds.md) 。
- Tendermint一些难点，可参考 [FAQ](https://github.com/little51/codechain/blob/master/faq.md) 。

## Codechain核心

### 安装Mongodb

​	资产数据采用Mongodb保存，所以请先安装Mongodb。

```shell
sudo apt-get install mongodb
```

### 二进制文件安装

codechain采用go语言开发，所以无依赖库，直接下载可执行文件即可运行，预览功能，安装脚本见

https://github.com/little51/codechain/blob/master/release/v0.0.1/codechain_v0.0.1_linux_amd64.sh

### 程序编译

如果已经配置好go开发环境，可以编译程序源码。

```shell
export GO111MODULE=on
export GOPROXY=https://goproxy.io
cd core
go build
```

### 运行

#### 链初始化

```shell
tendermint init
```

链初始化只做一次，执行完成后，会在~/.tendermint目录下生成配置文件和数据文件。

#### 运行应用
进入codechain/core目录下执行下面操作：
```shell
cd codechain/core
go build
./core
```

core程序会在26658端口监听tendermint进程发过来的交易。

#### 运行单节点

​	再开一个命令行，执行：

```shell
tendermint node
```

tendermint会在26657端口监听。

### 测试

```shell
curl -s 'localhost:26657/broadcast_tx_commit?tx="key1:value1"'
curl -s 'localhost:26657/abci_query?data="key=key1"'
```

## Web应用

​	Web应用目前提供了账户管理和简单的资产管理，可实现账户密钥对的创建，资产的签名、登记。

### 编译运行

```shell
cd webserver
go build
./webserver
```

### 测试

​	新建账户，返回公私钥对，对资产字符串进行私钥签名，用公钥验证。

#### 新建账户

```shell
curl -X POST http://localhost:3000/account/new
```

结果如下：

```json
{
	"address":"155703E01281055E6C198040FB53FB203157B903",
	"error":"",
	"privateKey":"78fde71c52eb009b862a9041071f4cfe6721ea639d0a671dd1dd19292bd6799dbab233aa5573f4d06ec8376585357e7eee4026adabb6fed7ccb081df5ce9ddd3",
	"publicKey":"BAB233AA5573F4D06EC8376585357E7EEE4026ADABB6FED7CCB081DF5CE9DDD3"}
```

#### 资产签名
由于msg中的value可以为一个string（字符串），也可以为格式为{ token: 'xxx', from: 'xxx', to: 'xxx', amount: 'xxx'}的json格式，所以我们分两种情况演示：

<font>**① msg_value为string**</font>

```shell
curl  -H "Content-Type: application/json" -d '{"privatekey":"78fde71c52eb009b862a9041071f4cfe6721ea639d0a671dd1dd19292bd6799dbab233aa5573f4d06ec8376585357e7eee4026adabb6fed7ccb081df5ce9ddd3","msg":"myasset:string"}' -X POST http://localhost:3000/account/sign
```

结果如下：

```json
{"error":"","sign":"3ae68c53970c000e39d195647bd50488e17c147d7757dcebbca433f66d6279a822c5763e79ae1c53de58add275c0e02785dd69c45b816c8eab97f9f77dd1220d"}
```

<font>**② msg_value为json**</font>

由于格式为{ token: "token", from: "from", to: "to", amount: 10}的json格式,我们需要进行base64加密然后进行签名，所以我们测试如下
```shell
curl  -H "Content-Type: application/json" -d '{"privatekey":"78fde71c52eb009b862a9041071f4cfe6721ea639d0a671dd1dd19292bd6799dbab233aa5573f4d06ec8376585357e7eee4026adabb6fed7ccb081df5ce9ddd3","msg":"myasset:eyJ0b2tlbiI6InRva2VuIiwiZnJvbSI6ImZyb20iLCJ0byI6InRvIiwiYW1vdW50IjoiMTAifQ=="}' -X POST http://localhost:3000/account/sign
```

结果如下：

```json
{"error":"","sign":"70474bf8a326677d92299cafcd33b5f27a7f2fe0de5694d775bcbdd126caf90804bae0d6a88a5be4b97d13e605898f575e982d73d2854bff29017c81b216da0b"}
```

#### 资产登记
由于msg中的value可以为一个string（字符串），也可以为格式为{ token: 'xxx', from: 'xxx', to: 'xxx', amount: 'xxx'}的json格式，所以我们和资产签名一样分两种情况演示：

<font>**① msg_value为string**</font>

```shell
curl  -H "Content-Type: application/json" -d '{"publickey":"BAB233AA5573F4D06EC8376585357E7EEE4026ADABB6FED7CCB081DF5CE9DDD3","sign":"3ae68c53970c000e39d195647bd50488e17c147d7757dcebbca433f66d6279a822c5763e79ae1c53de58add275c0e02785dd69c45b816c8eab97f9f77dd1220d","msg":"myasset:string"}' -X POST http://localhost:3000/assets/new
```

结果如下：

```json
{
  "error": "",
  "info": "{
  "jsonrpc": "2.0",
  "id": "",
  "result": {
    "check_tx": {
      "code": 0,
      "data": null,
      "log": "",
      "info": "",
      "gasWanted": "1",
      "gasUsed": "0",
      "events": [],
      "codespace": ""
    },
    "deliver_tx": {
      "code": 0,
      "data": null,
      "log": "",
      "info": "",
      "gasWanted": "0",
      "gasUsed": "0",
      "events": [],
      "codespace": ""
    },
    "hash": "9D4F905666FB9BB194D9D9C6CEDB4E921A0597B7DA89B27F4419CB25A0B21A77",
    "height": "54"
  }
}",
  "result": true
}
```

<font>**② msg_value为json**</font>

```shell
curl  -H "Content-Type: application/json" -d '{"publickey":"BAB233AA5573F4D06EC8376585357E7EEE4026ADABB6FED7CCB081DF5CE9DDD3","sign":"3ae68c53970c000e39d195647bd50488e17c147d7757dcebbca433f66d6279a822c5763e79ae1c53de58add275c0e02785dd69c45b816c8eab97f9f77dd1220d","msg":"myasset:eyJ0b2tlbiI6InRva2VuIiwiZnJvbSI6ImZyb20iLCJ0byI6InRvIiwiYW1vdW50IjoiMTAifQ=="}' -X POST http://localhost:3000/assets/new
```

结果如下：

```json
{
  "error": "",
  "info": "{
  "jsonrpc": "2.0",
  "id": "",
  "result": {
    "check_tx": {
      "code": 0,
      "data": null,
      "log": "",
      "info": "",
      "gasWanted": "1",
      "gasUsed": "0",
      "events": [],
      "codespace": ""
    },
    "deliver_tx": {
      "code": 0,
      "data": null,
      "log": "",
      "info": "",
      "gasWanted": "0",
      "gasUsed": "0",
      "events": [],
      "codespace": ""
    },
    "hash": "48AFC166273C0A8968F29323C914B9830A35616BF9A368FA45E674E9E085CCEF",
    "height": "55"
  }
}",
  "result": true
}
```

### 数据库查询

​	资产会保存到Mongodb中，通过Mongodb客户端可以查询chain数据库的assets集合。

## 链浏览器

参见 [链浏览器](https://github.com/little51/codechain/blob/master/browser/README.md) 的说明。

## 路线图

| 序号 | 类别  | 任务                                         | 完成情况 |
| :--: | ----- | -------------------------------------------- | -------- |
|  1   | 基础  | 开发环境                                     | 完成     |
|  2   |       | 运行环境                                     | 完成     |
|  3   |       | 一键自动安装应用                             | 完成     |
|  4   | 资产  | 基本key-value(mongledb)                      | 完成     |
|  5   |       | 账户account与签名                            | 完成     |
|  6   |       | 可分隔资产                                   | 开发中   |
|  7   | API   | API规划                                      | 完成     |
|  8   | WEBUI | 块浏览器                                     | 开发中   |
|  9   |       | 资产管理                                     | 开发中   |
|  10  |       | 账户管理                                     | 完成     |
|  11  | 网络  | 协调节点                                     | 完成     |
|  12  |       | 节点管理（增减节点）                         | 开发中   |
|  13  | 应用  | 生态系统（代码仓库、分布式应用、资产、价值） | 规划完成 |


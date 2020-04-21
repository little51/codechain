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
curl -s 'localhost:26657/broadcast_tx_commit?tx="eyJwdWJsaWNrZXkiOiJBNUZFQTY0NUMwOTlGNjlFRTVGQ0Q4QjY3RkU4M0VBQ0QwQ0JBQzQxMEFCQzVCOEZBNUNCNEU0NTFENUY1Q0VDIiwic2lnbiI6IjNiOTM2ODc5MzkyMDJiNTJhZWJlZDA2YjU1ZjY2M2ZiN2M0ZjYyNzU2MGY3OTA3NDM5ZDg4ZjYzMWM5NWI4ZjNjYzdkODBkNzEzMTgwZjRhMDg1MDBkZDU5YTgzNjlhODZmODlhODQ4NDI1NzZjMDJkMTFkYTI3NTdhNzBiZDBkIiwibXNnIjoiZXlKMGIydGxiaUk2SWtWUFJpSXNJbVp5YjIwaU9pSkJOVVpGUVRZME5VTXdPVGxHTmpsRlJUVkdRMFE0UWpZM1JrVTRNMFZCUTBRd1EwSkJRelF4TUVGQ1F6VkNPRVpCTlVOQ05FVTBOVEZFTlVZMVEwVkRJaXdpZEc4aU9pSkVOemM1UkRZelJqZzNRa0kyT1RreE4wTTBSRE5FUkRKQ056WXlNekUxT1RVeE9URTVOVVExUVRWQlJEYzBSakV3UWpjNFJEbEJNamRCT1RZMk5UVkRJaXdpWVcxdmRXNTBJam9pTlRBaWZRPT0ifQ=="'
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
	"address":"ECDCE3D5B6164768C0D4DFB74BF6B9E2C4D1682B",
	"error":"",
	"privateKey":"80d7f6a76d6bbc9893efff3721a49432f5ffdba39e0c98256d7e3fdd53e6b807d4e4cf08ec338970efc0e6b8e1dde3d4129a6bb5588ba90f72f9a0d6c0ed62ce",
  "publicKey":"D4E4CF08EC338970EFC0E6B8E1DDE3D4129A6BB5588BA90F72F9A0D6C0ED62CE"
}
```

#### 资产签名
资产交易需要先进行资产签名，验证发起资产交易账户的公私钥对

测试如下
```shell
curl  -H "Content-Type: application/json" -d '{"privatekey":"80d7f6a76d6bbc9893efff3721a49432f5ffdba39e0c98256d7e3fdd53e6b807d4e4cf08ec338970efc0e6b8e1dde3d4129a6bb5588ba90f72f9a0d6c0ed62ce","msg":"eyJ0b2tlbiI6IktLSyIsImZyb20iOiJENEU0Q0YwOEVDMzM4OTcwRUZDMEU2QjhFMURERTNENDEyOUE2QkI1NTg4QkE5MEY3MkY5QTBENkMwRUQ2MkNFIiwidG8iOiIiLCJhbW91bnQiOiI2MDAifQ=="}' -X POST http://localhost:3000/account/sign
```

结果如下：

```json
{"error":"","sign":"25a76ff76196706d63c3190cfb97edd45f322f046cb9d78592eda1368b42c15c63a6780a777a40de626b9b9530e9433f51d8ee300b881ea55941cc53b5b27507"}
```

#### 资产创世
在资产交易之前，应该确认发起者拥有足够的token，可以通过资产创世的方式来创建一种新的token

测试如下：
```shell
curl  -H "Content-Type: application/json" -d '{"publickey":"D4E4CF08EC338970EFC0E6B8E1DDE3D4129A6BB5588BA90F72F9A0D6C0ED62CE","sign":"25a76ff76196706d63c3190cfb97edd45f322f046cb9d78592eda1368b42c15c63a6780a777a40de626b9b9530e9433f51d8ee300b881ea55941cc53b5b27507","msg":"eyJ0b2tlbiI6IktLSyIsImZyb20iOiJENEU0Q0YwOEVDMzM4OTcwRUZDMEU2QjhFMURERTNENDEyOUE2QkI1NTg4QkE5MEY3MkY5QTBENkMwRUQ2MkNFIiwidG8iOiIiLCJhbW91bnQiOiI2MDAifQ=="}' -X POST http://localhost:3000/assets/new
```

结果如下：

```json
{
  "error": "",
  "info": "{ \n \"jsonrpc\": "2.0",\n \"id\": \"\",\n \"result\": {\n \"check_tx\": {\n \"code\": 0,\n \"data\": null,\n \"log\": \"\",\n \"info\": \"CheckTx successfully carried out the signVerify\",\n \"gasWanted\": \"1\" \n \"gasUsed\": "0",\n \"events\": [], \n \"codespace\": \"\"},\n \"deliver_tx\": {\n \"code\": 0,\n \"data\": null, \n \"log\": \"\", \n \"info\": \"\",\n \"gasWanted\": \"0\",\n \"gasUsed\": \"0\",\n \"events\": [],\n \"codespace\": \"\"},\n \"hash\": \"91C2DA3A01E9093FB7CF1FD4A30C93A5164A8F20FFCE28020BD7AB1497FB3AE8\",\n \"height\": \"82\"}}\n",
  "result": true
}
```

#### 资产交易
两个不同账户之间可以进行相同token的交易

```shell
curl  -H "Content-Type: application/json" -d '{"publickey":"D4E4CF08EC338970EFC0E6B8E1DDE3D4129A6BB5588BA90F72F9A0D6C0ED62CE","sign":"df2d0652c8994ddc14ec70f9532aff7fb9b4ccd3077bf9c88c936a65a18a500e380f639f410d4058deddc2d4047bee070e229b0b5e86aaa97448374ac854770d","msg":"eyJ0b2tlbiI6IktLSyIsImZyb20iOiJENEU0Q0YwOEVDMzM4OTcwRUZDMEU2QjhFMURERTNENDEyOUE2QkI1NTg4QkE5MEY3MkY5QTBENkMwRUQ2MkNFIiwidG8iOiJBNUZFQTY0NUMwOTlGNjlFRTVGQ0Q4QjY3RkU4M0VBQ0QwQ0JBQzQxMEFCQzVCOEZBNUNCNEU0NTFENUY1Q0VDIiwiYW1vdW50IjoiMzAwIn0="}' -X POST http://localhost:3000/assets/new
```

结果如下：

```json
{
  "error": "",
  "info": "{ \n \"jsonrpc\": "2.0",\n \"id\": \"\",\n \"result\": {\n \"check_tx\": {\n \"code\": 0,\n \"data\": null,\n \"log\": \"\",\n \"info\": \"CheckTx successfully carried out the signVerify\",\n \"gasWanted\": \"1\" \n \"gasUsed\": "0",\n \"events\": [], \n \"codespace\": \"\"},\n \"deliver_tx\": {\n \"code\": 0,\n \"data\": null, \n \"log\": \"\", \n \"info\": \"\",\n \"gasWanted\": \"0\",\n \"gasUsed\": \"0\",\n \"events\": [],\n \"codespace\": \"\"},\n \"hash\": \"B80F6B79569989513E5F23552CDBB53BA854F468FC2A4BFF1696C4C897660734\",\n \"height\": \"83\"}}\n",
  "result": true
}
```

#### 资产查询
通过公钥可以来查询到当前所有不同token的资产信息

测试如下：
```shell
curl  -H "Content-Type: application/json" -d '{"key": "D4E4CF08EC338970EFC0E6B8E1DDE3D4129A6BB5588BA90F72F9A0D6C0ED62CE"}' -X POST http://localhost:3000/assets/query
```

结果如下：
```json
{
  "error":"",
  "info":"{\n  \"jsonrpc\": \"2.0\",\n  \"id\": -1,\n  \"result\": {\n    \"response\": {\n      \"code\": 0,\n      \"log\": \"\",\n      \"info\": \"D4E4CF08EC338970EFC0E6B8E1DDE3D4129A6BB5588BA90F72F9A0D6C0ED62CE\",\n      \"index\": \"0\",\n      \"key\": null,\n      \"value\": \"eyJhcnJheSI6IFt7InB1YmxpY2tleSI6IkQ0RTRDRjA4RUMzMzg5NzBFRkMwRTZCOEUxRERFM0Q0MTI5QTZCQjU1ODhCQTkwRjcyRjlBMEQ2QzBFRDYyQ0UiLCJ0b2tlbiI6IktLSyIsImFtb3VudCI6MzAwfV19\",\n      \"proof\": null,\n      \"height\": \"0\",\n      \"codespace\": \"\"\n    }\n  }\n}",
  "result":true
}
```

#### msg说明
在资产签名，资产创世，资产查询测试中，POST请求数据中包含msg字段，该字段是处理后的一段加密Base64的字符串。其原始信息为包含token、from、to和amount四个属性的javascript对象，在前端经JSON.stringify()和Base64.encode()先后处理后形成最终的msg值。

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


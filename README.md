# codechain

​	codechain基于 Tendermint共识机制的代码链，采用go语言实现，并且依赖于Tendermint环境，建议在ubuntu16.04下开发，首先 [建立开发环境](https://github.com/little51/codechain/blob/master/getting-start.md) ，为了验证Tendermint的四个节点的运行情况，可以参考 [运行环境组网](https://github.com/little51/codechain/blob/master/making-testnet.md) 。

## 资产管理

### 安装Mongodb

​	资产数据采用Mongodb保存，所以请先安装Mongodb。

```shell
sudo apt-get install mongodb
```

### 编译

```shell
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

```shell
./core
```

core程序会在26658端口监听tendermint进程发过来的交易。

#### 运行NODE

​	再开一个命令行，执行：

```shell
tendermint node
```

tendermint会在26657端口监听。

### 测试

```shell
curl -s 'localhost:26657/broadcast_tx_commit?data="key1=value1"'
curl -s 'localhost:26657/abci_query?data="key=key1"'
```

## 路线图

| 序号 | 类别  | 任务                                         | 完成情况 |
| :--: | ----- | -------------------------------------------- | -------- |
|  1   | 基础  | 开发环境                                     | 完成     |
|  2   |       | 运行环境                                     | 完成     |
|  3   |       | CURL一键自动安装应用                         |          |
|  4   | 资产  | 基本key-value(mongledb)                      | 完成     |
|  5   |       | 账户account与签名                            | 开发中   |
|  6   |       | 可分隔资产                                   |          |
|  7   | API   | API规划                                      |          |
|  8   | WEBUI | 块浏览器                                     | 开发中   |
|  9   |       | 资产管理                                     |          |
|  10  |       | 账户管理                                     |          |
|  11  | 网络  | 协调节点                                     |          |
|  12  |       | 节点管理（增减节点）                         |          |
|  13  | 应用  | 生态系统（代码仓库、分布式应用、资产、价值） | 规划中   |


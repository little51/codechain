# codechain

​	codechain基于 Tendermint共识机制的代码链，采用go语言实现，并且依赖于Tendermint环境，建议在ubuntu16.04下开发，首先 [建立开发环境](https://github.com/little51/codechain/blob/master/getting-start.md) ，为了验证Tendermint的四个节点的运行情况，可以参考 [运行环境组网](https://github.com/little51/codechain/blob/master/making-testnet.md) 。

## 资产管理

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


# Tendermint组网

## 机器准备

- 4台ubuntu16.04
- IP必须连续，如果不连续，修改的配置文件比较多

## 安装Tendermint

在四台机器上，分别安装Tendermint环境，安装脚本参见 https://github.com/little51/codechain/blob/master/inst-tm.sh 

## 创建节点配置文件

在第一台机器上，执行以下命令：

```shell
tendermint testnet --v 4 --o ./mytestnet --populate-persistent-peers --starting-ip-address 第一个IP(x.x.x.x)
```

参数解释如下：

- --v：指定节点数
- --o：输出目录
- --populate-persistent-peers：自动生成持久节点参数
- --starting-ip-address：起始IP

这样会在mytestnet目录下，生成4个文件夹，分别是node0、node1、node2、node3，每个文件夹下都有config和data文件夹。

- priv_validator_key.json：当前节点的公私钥对
- genesis.json：初始块配置文件，保存了链的名称、链的参数、四个验证节点的公钥。四个结点中这个文件是一样的
- config.toml：链的运行参数，如persistent_peers会按4个IP自动生成，如果IP不连续，四个节点都得修改这个参数

- 链的数据存放在data

## 部署运行

将mytestnet复制到4个节点的home下，其实每个节点只要本节点那个文件夹即可

在4个节点上分别运行：

```shell
tendermint node --home ./mytestnet/node0 --proxy_app=kvstore
tendermint node --home ./mytestnet/node1 --proxy_app=kvstore
tendermint node --home ./mytestnet/node2 --proxy_app=kvstore
tendermint node --home ./mytestnet/node3 --proxy_app=kvstore 
```

可以观察到，起到第三个节点，链已能正常接受交易了。4个节点的height是同步的

## 测试

在任意一节点提交交易，查询交易，得到的都是同样的结果

```shell
curl -s 'localhost:26657/broadcast_tx_commit?tx="abcd"'

curl -s 'localhost:26657/abci_query?data="abcd"'
```

## 节点的同步

如果一个节点异常中断，然后再启动后丢失的块如何同步？Tendermint提供了一种快速同步（fast-sync）的机制

可以用以下方法测试一下：

- 中断掉一个节点的tendermint（ctrl +c）,然后发交易，这样可以看到3个节点的交易是正常的（> 3/4个节点）
- 中断的结点起动后，会很快追平数据
- 中断掉两个节点，交易就不正常了，但提交的交易不会丢，节点正常后，会自动应有交易，这个机制被称为预写日志 (WAL)，再加上内存池mempool.wal，极其类似oracle的在线归档日志和离线归档日志。


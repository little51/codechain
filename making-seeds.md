# Tendermint Seed模式

​	使用tendermint testnet命令可以生成4个或更多的节点，每个节点的ID和IP地址都设置在各节点的config.toml文件的persistent_peers参数中。节点启动时，会按这个参数不断尝试与其他节点建立连接。那么要加入第5个节点如何操作？

## 种子模式

​		persistent_peers是持久连接节点，适用于较固定的节点，如果动态新增节点，用这种模式灵活性差一些，所以得用P2P的方式——种子模式(seed_mode)。在四个节点中，只需要其中一个设置成种子节点即可，其他节点都连接到这个节点来，种子节点的addressbook中，只要保留另外几个的其中一个节点信息即可，随着各节点与种子节点连接，各节点的ID和IP信息会通过种子节点扩散开来，一旦其他节点互相联通，就不再需要种子节点。在这个Tendermint中被形象地描述成乒乓（ping-pong）,从实践来看，**种子节点只要能打出第一个球就可以了**。

## 配置

- 种子节点：将seed_mode设为 true（注意，其他节点的seed_mode要设成false）
- 其他节点：将seeds设成种子节点的id@ip:port,类似于：

```shell
seeds = "0e58be64d8da79457778bca73c48d7f2b8c9694e@172.16.62.43:26656"
```

## 第五节点加入

1. 初始化节点：tendermint init

2. 同步genesis.json：从原有节点复制genesis.json，覆盖本地文件

3. 设置seeds参数：连接到种子节点

   完成以上三步，就可以实现第五个节点的接入，新增的节点以此类推。

## 验证器节点加入

​	按以上三步加入的新节点，只参与链块的同步，不参与共识的验证，如果新加入验证器节点，则有两种方式：

1. 在genesis.json的validators加入新节点，所有节点都更新genesis.json
2. 修改链码的endBlock方法，可在线更改验证器节点列表


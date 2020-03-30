# FAQ

Q、如何动态增删验证结点

```
A、在EndBlock中响应客户端的增、删的请求，需要注意的是：validator列中是在块的Header中，所以必须记住从新加的验证节点是在Height在多少时开始生效的，否则replay会报错。
```

Q、如何在程序重启时不从第0块重新replay

```
A、在Info中给客户端返回最后的高度和apphash，这两个值在Commit时保存下来，apphash可以给空数组，如果给定值，则会造成是否创建空块参数失效，因为每个Commit的apphash不同
```

Q、如何重置所有块数据

```
A、执行tendermint unsafe_reset_all，同时，要将MongodbMongodb中的chain.state删除掉，否则statestate中记录的Height之前的块将不会replay
```


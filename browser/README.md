# vue-tendermint-explorer
`vue-tendermint-explorder`区块链浏览器是以`tendermint-explorder`为依据，在`vue-admin-template`基础上进行的框架二次开发，该项目启动我们需要两个步骤，分别是:
+ 配置`tendermint`跨域
+ 启动`vue-admin-template`

## 配置跨域
### 1. 链初始化
```shell
tendermint init
```
链初始化只做一次，如果之前已经执行了这一步可以忽略了，执行完成后，会在`~/.tendermint`目录下面生成配置文件`config`和数据文件`data`。

### 2. 配置跨域
由于`tendermint`会监听26657端口，前端和后端项目都会通过`rpc`和`http`的方式去访问，所以我们需要进入`tendermint`的配置文件：
```shell
cd ~/.tendermint/config
sudo vim config.toml
```
然后我们我们修改`config.toml`文件其中的`rpc`配置如下：
```shell
##### rpc server configuration options #####
[rpc]

# TCP or UNIX socket address for the RPC server to listen on
laddr = "tcp://0.0.0.0:26657"

# A list of origins a cross-domain request can be executed from
# Default value '[]' disables cors support
# Use '["*"]' to allow any origin
cors_allowed_origins = ["*"]
```
然后我们就可以跨域请求`127.0.0.1:26657`或者`localhost:26657`了，从而实现和`tendermint`的交互。

## 启动项目
### 1. 启动后端服务
进入后端文件开始编译代码：
```shell
cd ~/codechain/webserver/
go build
```
启动后端服务，启动在`4000`端口:
```shell
./webserver
```
### 2. 启动前端服务
```shell
cd ~/codechain/browser
npm install --registry=https://registry.npm.taobao.org
npm run dev
```
访问`localhost:9528`
# 开发环境的建立

​	codechain基于Tendermint(以拜占庭容错BFT实现多台机器安全共识的应用)采用golang开发，操作系统选用了ubuntu16.04,为了便于开发，在ubuntu上安装vscode用于go开发的IDE。在开发环境建立之初，初步先建立Tendermint单节点(Single node)，基本上能够满足开发的要求，随着项目的进展再扩展到4个节点组成一个最基本的BFT网络。

## Tendermint安装

​	Tendermint有两种安装方式，用压缩包直接解压和最新源码编译安装，由于Tendermint在快速开发中，且项目主要要用go实现，所以建议安装go环境和用源码方式安装。安装过程比较简单，但不能按照Tendermint官方的文档安装，那个文档跟不上Tendermint开发的进度，对go的依赖和编译命令都存在问题，我已在Tendermint的github中提了issues(https://github.com/tendermint/tendermint/issues/4499),给出了目前为止正确的安装脚本。

​	执行脚本注意以下问题：

- 不要用sudo执行脚本，用当前用户
- 由于网络原因，从github下载源码经常会中断，所以一般可能会多次执行，除第一次执行以外，后续的执行都要把echo那三行注释掉，不然会重复往~/.profile文件中增加环境变量。
- 另外一种方法是一句一句执行，保证每一步都成功即可

脚本如下：

```shell
#!/usr/bin/env bash

## install tools
sudo apt-get install curl
sudo apt-get install git
sudo apt install make

## download golang
curl -O https://dl.google.com/go/go1.14.linux-amd64.tar.gz
tar -xvf go1.14.linux-amd64.tar.gz

## install golang
sudo rm -fr /usr/local/go
sudo mv go /usr/local
mkdir goApps

## init environment variable
## execute next time, comment 3 lines "echo"
echo "export GOPATH=~/goApps" >> ~/.profile
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
echo "export PATH=\$PATH:\$GOPATH/bin" >> ~/.profile

## apply variable
source ~/.profile

## clone tendermint
REPO=github.com/tendermint/tendermint
go get $REPO
cd $GOPATH/src/$REPO

## build tendermint
git checkout master
make tools
make install
make build
```

## 测试Tendermint

参考https://github.com/tendermint/tendermint/blob/master/docs/introduction/quick-start.md

```shell
# 初始tendermint节点，会在~/.tendermint目录下生成配置和密钥文件
tendermint init
# 启动内部应用kvstore
tendermint node --proxy_app=kvstore
# 检查应用状态
curl -s localhost:26657/status
# 发送交易
curl -s 'localhost:26657/broadcast_tx_commit?tx="abcd"'
# 查询交易状态
curl -s 'localhost:26657/abci_query?data="abcd"'
```

## vncserver安装

​	要在ubuntu16.04上安装vscode，必须要用linux桌面系统，如果是远程操作，需要先安装vncserver。安装过程如下：

```shell
sudo apt-getinstallgnome-session-flashback
sudo apt-get install gnome-session-flashback
sudo apt install ubuntu-desktop gnome-panel gnome-settings-daemon metacity nautilus gnome-terminal -y
vncserver
cd ~/.vnc
vi xstartup 
```

xstartup内容如下：

```shell
#!/bin/sh
export XKL_XMODMAP_DISABLE=1
unset SESSION_MANAGER
unset DBUS_SESSION_BUS_ADDRESS
[ -x /etc/vnc/xstartup ] && exec /etc/vnc/xstartup
[ -r $HOME/.Xresources ] && xrdb $HOME/.Xresources
xsetroot -solid grey
vncconfig -iconic &
gnome-session &
gnome-panel &
gnome-settings-daemon &
metacity &
nautilus &
gnome-terminal &
```

然后重启vncserver，指定分辨率，不然分辨率是800*600，且没办法修改：

```shell
vncserver -kill :1
vncserver -geometry 1440×900
```

## VSCODE安装

​	ubuntu下安装vscode，有多种方法，建议采用deb包安装，https://code.visualstudio.com/Download#

```shell
sudo dpkg -i code_1.31.1-1549938243_amd64.deb
安装后，快捷方式会出现在ubuntu主菜单的开发工具下，但点不开，需要改以下文件：
sudo cp /usr/lib/x86_64-linux-gnu/libxcb.so.1 /usr/share/code/
cd /usr/share/code
sudo sed -i 's/BIG-REQUESTS/_IG-REQUESTS/' libxcb.so.1
```

## 第一个应用
## vue-tendermint-explorer
`vue-tendermint-explorder`区块链浏览器是以`tendermint-explorder`为依据，在`vue-admin-template`基础上进行的框架二次开发，该项目启动我们需要两个步骤，分别是:
+ 使用`nginx`作为t`endermint`节点机器的反向代理服务器
+ 使用npm包管理器下载依赖并且启动项目

### nginx配置
以`172.16.62.48`节点（`ubuntu16.04`）为测试实例，在下载好`nginx`之后,首先使用下面的命令来获取`nginx`配置文件的路径：
```shell
sudo nginx -t
```
并且进入`nginx`文件夹,并且查看当前文件下所有文件
```shell
cd /etc/nginx/
ls
```
其中`nginx.conf`是主要的配置文件，`conf.d`文件下所有以`.conf`结尾的文件都会包含在`nginx.conf`的配置当中，所以我们在`conf.d`下面创建新的文件，并配置`tendermint`的反向代理：
```shell
cd conf.d
sudo vim tendermint.conf
```
然后我们编辑`tendermint.conf`文件如下:
```shell
map $http_upgrade $connection_upgrade {
  default upgrade;
  ''      close;
}

upstream websocket {
  #ip_hash;
  server localhost:26657;
}

server {
  listen       26659;
  server_name localhost;

  location / {
    add_header 'Access-Control-Allow-Origin' $http_origin;
    add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
    add_header 'Access-Control-Allow-Headers' '*';

    if ($request_method = 'OPTIONS') {
        return 204;
    }
    proxy_pass http://websocket;
    proxy_read_timeout 300s;
    proxy_send_timeout 300s;

    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection $connection_upgrade;
  }
}
```
这样设置后同时对所需要的`http`和`websocket`同时做了代理和跨域，保存文件之后，我们使用下面的命令来验证配置是否正确：
```shell
sudo nginx -t
```
配置正确就启动`nginx`:
```shell
// 第一次启动
sudo nginx
// 不是第一次启动
sudo nginx -s reload
```

### 启动项目
```shell
cd browser
npm install --registry=https://registry.npm.taobao.org
npm run dev
```
访问`localhost:9528`
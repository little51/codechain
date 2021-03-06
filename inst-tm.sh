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
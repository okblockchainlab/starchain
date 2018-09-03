### 编译

```shell
yum install glide # setup glide, whick is a package management for golang.
export GOPATH=/your/go/path/directory  #设置GOPATH路径
cd $GOPATH/src
git clone https://github.com/okblockchainlab/starchain.git ./starchain
cd ./starchain
./build.sh #run this script only if you first time build the project
./runbuild.sh
ls *.so
ls *.dylib
```

### 其它注意项
- 调用"createrawtransaction"时，需要传入一些tx信息，这些信息可以从 __"/api/v1/asset/utxos/:addr"__ 和 __"/api/v1/transaction/:hash"__ 这些api中获取。
- 使用 __"/api/v1/transaction"__ 发送已签名的tx.

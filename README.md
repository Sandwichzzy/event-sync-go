# event-sync 代码调式
## 一.部署合约
- 生成钱包地址 
- cast wallet n
```
Successfully created new keypair.
Address:     0x5d008115fEB33aE09F135A724b3343112C162B6c
Private key: 0x5fb48f2298b35d82e632f7385ab918a3896a9c6ed774bd0787d9461bb3a3d2e4
```
- 环境变量配置
```
export PRIVATE_KEY=0x5fb48f2298b35d82e632f7385ab918a3896a9c6ed774bd0787d9461bb3a3d2e4
export RPC_URL=https://rpc-testnet.roothashpay.com
```
- 部署合约
`forge script ./script/TreasureManagerScript.s.sol:TreasureManagerScript --rpc-url $RPC_URL --private-key $PRIVATE_KEY --broadcast`
- 部署完成之后的合约地址
```treasureManagerImplementation = 0x09bc3071DD385DFe5A10c09F747Ac9037D66499f
treasureManager = 0x388fF618Ca5c1b8F28D4E845B431Ca3D4200140e
treasureManagerProxyAdmin = 0x7bC3b56AE67698632Bb25DbedDB86D00f81AF0F7
```
- 调用合约抛出合约事件
```
cast send --rpc-url $ROOTHASH_RPC_URL --private-key $PRIVATE_KEY 0x388fF618Ca5c1b8F28D4E845B431Ca3D4200140e --value 100000000000 "depositETH()"
cast send --rpc-url $ROOTHASH_RPC_URL --private-key $PRIVATE_KEY 0x388fF618Ca5c1b8F28D4E845B431Ca3D4200140e "setWithdrawManager(address)"  0x5d008115fEB33aE09F135A724b3343112C162B6c
```

## 二.启动项目扫链
- 编译代码
* `go mod tidy`
* `make`
- 配置环境变量
```
export EVENT_SYNC_MIGRATIONS_DIR="./migrations"

export EVENT_SYNC_CHAIN_ID=1
export EVENT_SYNC_CHAIN_RPC="https://rpc-testnet.roothashpay.com"
export EVENT_SYNC_STARTING_HEIGHT=1140200 区块的配置在创建合约那个高度就行
export EVENT_SYNC_CONFIRMATIONS=10
export EVENT_SYNC_LOOP_INTERVAL=1s
export EVENT_SYNC_BLOCKS_STEP=10

export EVENT_SYNC_HTTP_PORT=8989
export EVENT_SYNC_HTTP_HOST="127.0.0.1"

export EVENT_SYNC_SLAVE_DB_ENABLE=false

export EVENT_SYNC_MASTER_DB_HOST="127.0.0.1"
export EVENT_SYNC_MASTER_DB_PORT=5432
export EVENT_SYNC_MASTER_DB_USER="sandwichzzy"
export EVENT_SYNC_MASTER_DB_PASSWORD=""
export EVENT_SYNC_MASTER_DB_NAME="event_sync"

export EVENT_SYNC_SLAVE_DB_HOST="127.0.0.1"
export EVENT_SYNC_SLAVE_DB_PORT=5432
export EVENT_SYNC_SLAVE_DB_USER="sandwichzzy"
export EVENT_SYNC_SLAVE_DB_PASSWORD=""
export EVENT_SYNC_SLAVE_DB_NAME="event_sync"
区块的配置在创建合约那个高度就行
```

- 让环境变量生效
`source .env`
- migrate 数据库
`./event-sync migrate`
- 启动扫链服务
`./event-sync migrate`
接下来不断发合约执行命令

## 三. GRPC 和 HTTP Server
- 启动 grpc server
`./event-sync grpc`
- 使用 grpcui 测试 grpc 接口
`grpcui -plaintext ip:porr`
`grpcui -plaintext 127.0.0.1:8987`
- 启动 http server
`./event-sync api`
- 测试 http api
`http://127.0.0.1:8989/api/v1/deposit/tokens?page=1&pageSize=10`
## 四.RootHash Chain 附属资料
- 测试网 RPC 与浏览器
* https://rpc-testnet.roothashpay.com
* https://wss-testnet.roothashpay.com
* https://explorer-testnet.roothashpay.com
- 主网 RPC 与浏览器
* https://rpc.roothashpay.com
* https://wss.roothashpay.com
* https://explorer.roothashpay.com
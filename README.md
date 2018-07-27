# fakegeth

fake [geth](https://geth.ethereum.org/) or [parity](https://github.com/paritytech/parity-ethereum) node.

## why?

I want to use this for testing web3 interface ([Web3.py](https://github.com/ethereum/web3.py), [rust-web3](https://github.com/tomusdrw/rust-web3), etc...).

## Installation

```
$ go build  # and fakegeth copy to your PATH's direcotry
```

## Usage

for example of config.toml:
```toml
[http]
host = "localhost"
port = 8545
```

run:
```
$ fakegeth config.toml
2018/07/27 14:31:12 config: &{localhost 8545}, <nil>, <nil>
2018/07/27 14:31:12 listen http, endpoint: localhost:8545
```

access from web3.py:
```python
# accounts.py
from web3 import Web3
from web3 import HTTPProvider

w3 = Web3(HTTPProvider())
print(w3.personal.listAccounts())
```

```
$ python accounts.py   # error...
```

fakegeth log:
```
2018/07/27 14:37:28 http: user-agent=[Web3.py/4.4.1/<class 'web3.providers.rpc.HTTPProvider'>], id=0, method=personal_listAccounts, params=[]
```

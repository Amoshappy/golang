# Compose Grand Tour - Go - etcdv3

## Build notes

Libraries are included in the vendor directory. Assembled using the glide.sh tool for vendoring.

## Run notes

Three environment variables need to be set:

* `COMPOSE_ETCD_ENDPOINTS` = set to a comma delimited list of HTTPS URLs for etcd v3 endpoints
* `COMPOSE_ETCD_USER` = User name for etcd user
* `COMPOSE_ETCD_PASS` = Password for etcd user

## Example:

```
export COMPOSE_ETCD_ENDPOINTS=https://portal219-5.threetcd.compose-3.composedb.com:18279,https://portal227-0.threetcd.compose-3.composedb.com:18279
export COMPOSE_ETCD_USER=root
export COMPOSE_ETCD_PASS=YOURPASSWORD
```



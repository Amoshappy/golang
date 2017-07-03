# etcd v3 example

Three environment variables need to be set:

* `COMPOSEETCDENDPOINTS` = set to a comma delimited list of HTTPS URLs for etcd v3 endpoints
* `COMPOSEETCDUSER` = User name for etcd user
* `COMPOSEETCSPASS` = Password for etcd user

## Example:

```
export COMPOSEETCDENDPOINTS=https://portal219-5.threetcd.compose-3.composedb.com:18279,https://portal227-0.threetcd.compose-3.composedb.com:18279
export COMPOSEETCDUSER=root
export COMPOSEETCDPASS=YOURPASSWORD
```



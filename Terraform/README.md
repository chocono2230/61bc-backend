## コマンド他

```
function tfi() {
    if [[ $# != 1 ]]; then
        echo 引数エラー: $# args
        echo \<env\>
        return
    fi
    terraform init -backend-config="./envs/${1}/backend.conf" -reconfigure
}

function tfp() {
    if [[ $# != 1 ]]; then
        echo 引数エラー: $# args
        echo \<env\>
        return
    fi
    terraform plan -var-file="./envs/${1}/default.tfvars"
}

function tfxa() {
    if [[ $# != 1 ]]; then
        echo 引数エラー: $# args
        echo \<env\>
        return
    fi
    terraform apply -var-file="./envs/${1}/default.tfvars"
}

function tfxd() {
    if [[ $# != 1 ]]; then
        echo 引数エラー: $# args
        echo \<env\>
        return
    fi
    terraform destroy -var-file="./envs/${1}/default.tfvars"
}
```

## 注意点（git 上に存在しないもの）

- /terraform.tfvars

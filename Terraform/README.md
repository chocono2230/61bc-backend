## コマンド他

alias tfi='terraform init -backend-config="./envs/dev/backend.conf" -reconfigure'

alias tfp='terraform plan -var-file="./envs/dev/default.tfvars"'

alias tfxa='terraform apply -var-file="./envs/dev/default.tfvars"'

alias tfxd='terraform destroy -var-file="./envs/dev/default.tfvars"'

## 注意点（git 上に存在しないもの）

- /terraform.tfvars

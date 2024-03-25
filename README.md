# taskapp
Monorepo of simple task management application

## Prerequisites
- Container execution environment (Docker Desktop, Rancher Desktop, WSL2, etc.)
- make
- asdf (https://asdf-vm.com/)

## Setup

### Create Secrets
create secrets file
```shell
$ mkdir secrets
$ touch secrets/mysql_root_password
$ touch secrets/mysql_user_password
```

generate passwords
```shell
$ make make-mysql-passwords
```

### Setup asdf

Please read the setup guide of asdf. (https://asdf-vm.com/guide/getting-started.html)

### Install CLI tools with asdf

`hack/install-tools.sh` installs necessary CLI tools with asdf.

```shell
$ sh hack/install-tools.sh
```

## Migrate
Using golang migrate
```shell
# チェックポイント以降のすべてのup.sqlを実行する
$ migrate -path [マイグレーションSQLディレクトリ] -database [データベース接続文字列] up

# チェックポイント以降のup.sqlを1つだけ実行する
$ migrate -path [マイグレーションSQLディレクトリ] -database [データベース接続文字列] up 1

# チェックポイント以前のすべてのdown.sqlを実行する(初期化)
$ migrate -path [マイグレーションSQLディレクトリ] -database [データベース接続文字列] down

# チェックポイント以降のup.sqlを1つだけ実行する
$ migrate -path [マイグレーションSQLディレクトリ] -database [データベース接続文字列] down 1
```

Check CheckPoint
```
mysql> SELECT * FROM schema_migrations;
+---------+-------+
| version | dirty |
+---------+-------+
|    1003 |     0 |
+---------+-------+
```

## Connect to MySQL
```shell
$ docker compose exec mysql mysql -u taskapp_user -p taskapp
```
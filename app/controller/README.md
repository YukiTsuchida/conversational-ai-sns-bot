# controller

## パッケージ構成

```
.
├── ai # 対話型AIのモデルに依存するロジックを実装
│   ├── ai.go
│   ├── chatgpt_3_5_turbo
│   │   └── ai.go
│   └── chatgpt_4_0_turbo
│       └── ai.go
├── cli # main系
│   └── main.go
├── cmd # cmdに依存するロジックを実装
│   ├── cmd.go
│   ├── v0_1
│   │   └── cmd.go
│   └── v0_2
│       └── cmd.go
├── go.mod
├── model # 共通で利用するstructなど
│   ├── ai 
│   │   └── message_role.go
│   ├── cmd
│   │   ├── command.go
│   │   └── type.go
│   └── sns
│       ├── account.go
│       └── response.go
├── server # handlerとか入れる想定
├── service # コアとなるユースケース群
│   ├── abort_conversation.go
│   ├── create_account.go
│   ├── reply_conversation.go
│   └── start_conversation.go
└── sns # SNSに依存するロジックを実装
    ├── misskey
    │   └── sns.go
    ├── sns.go
    └── twitter
        └── sns.go
```

## ローカルで起動する

普通に起動する。

```sh
docker-compose up -d
curl localhost:8080/health
```

DBを初期化して起動する。

```sh
docker-compose down && docker-compose up -d --renew-anon-volumes
goose -dir ./ent/migrate/migrations postgres "host=localhost port=5432 user=admin password=admin dbname=db sslmode=disable" up # tableを初期化
```

### twitter bot動かし方

http://localhost:8080/accounts/twitter_login にアクセスするとtwitterの認可画面が表示されるので許可する。

以下のcurlコマンドでbotが動き出す。

```
curl -X POST localhost:8080/conversations/twitter -d '{"twitter_id":"hoge","ai_model":"gpt-3.5-turbo","cmd_version":"v0.1"}'
```

以下のcurlコマンドでbotを停止できる。

```
curl -X DELETE http://localhost:8080/conversations/twitter -d '{"twitter_id":"hoge"}'
```



## DBのマイグレーションについて

DBのスキーマ管理にentを利用しています。

スキーマを変更・追加する場合はentのスキーマを書き換えてpostgresqlのマイグレーションファイルを生成してください。

例えばUserというテーブルを追加する場合は以下のようにコマンドを打ちます。
```sh
go run -mod=mod entgo.io/ent/cmd/ent new User
```

マイグレーションファイルは`ent/migrate/migrations`に溜まっていて、以下のコマンドで生成します。
```sh
go generate ./ent
go run -mod=mod ent/migrate/main.go hogehoge(任意のmigration名)

# 適用
$ goose -dir ./ent/migrate/migrations postgres "host=localhost port=5432 user=admin password=admin dbname=db sslmode=disable" up
```

DBの中身を確認する場合は以下のコマンドを活用してください。
```sh
PGPASSWORD=admin docker-compose exec postgresql psql -d db -U admin -c "\dt"
```

詳しくは[entの公式ドキュメント](https://entgo.io/ja/docs/getting-started)も参考にしてください。

## DBへの接続

DBコンテナに接続するには以下のコマンドを実行してください。

```
docker-compose exec postgresql psql -Uadmin -ddb
```

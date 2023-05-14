# controller

## パッケージ構成

```sh
.
├── cmd
│   ├── ai_requestor
│   │   └── main.go # requestorのエントリポイント
│   └── controller
│       └── main.go # controllerのエントリポイント
├── config # 環境変数から値を読み出したりなど
├── ent # DBのスキーマ、自動生成されたORM、goose用のmigrationファイルなど
├── http
│   └── handlers # httpリクエストのパース・バリデーション、DI、usecaseの呼び出しを行う
├── models # ドメイン全体で引き回す構造体
├── repositories # いわゆるrepository
│   └── conversation.go
├── services
│   ├── ai # AIに依存するロジック、例えばAPIに投げるリクエストを組み立てたり、会話履歴をDBに保存するなど
│   │   ├── ai.go
│   │   ├── chatgpt_3_5_turbo
│   │   │   └── chatgpt_3_5_turbo.go
│   │   └── chatgpt_4_0_turbo
│   │       └── chatgpt_4_0_turbo.go
│   ├── prompt # prompt生成とパースに関連するロジック
│   │   ├── prompt.go
│   │   ├── v0_1
│   │   │   ├── v0_1.go
│   │   │   └── v0_1_test.go
│   │   └── v0_2
│   │       └── v0_2.go
│   └── sns # SNSに依存するロジック
│       ├── misskey
│       │   └── misskey.go
│       ├── sns.go
│       └── twitter
│           └── twitter.go
└── usecases # handlerから呼び出されるusecaseたち、DIして使うのでDBやAI,prompt,sns実装についての知識を一切持たない
    ├── abort_conversation.go
    ├── register_account.go
    ├── reply_conversation.go
    └── start_conversation.go
```

## ローカルで起動する

普通に起動する。hotreloadに[air](https://github.com/cosmtrek/air)を使っているので初回起動に時間がかかります。

```sh
docker-compose up -d
curl localhost:8080/health
```

DBを初期化して起動する。

```sh
docker-compose down && docker-compose up -d --renew-anon-volumes
goose -dir ./ent/migrate/migrations postgres "host=localhost port=5432 user=admin password=admin dbname=db sslmode=disable" up # tableを初期化
curl localhost:8080/health
```

### twitter bot動かし方

http://localhost:8080/accounts/twitter_login にアクセスするとtwitterの認可画面が表示されるので許可する。

以下のcurlコマンドでbotが動き出す。

```sh
curl -X POST localhost:8080/conversations/twitter -d '{"twitter_id":"hoge","ai_model":"gpt-3.5-turbo","cmd_version":"v0.1"}'
```

以下のcurlコマンドでbotを停止できる。

```sh
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

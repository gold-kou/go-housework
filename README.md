# はじめに
こちらのソフトウェアは私のバックエンドスキルを証明するためのポートフォリオの1つです。

現時点では、フロントエンドやネイティブアプリなどのクライアントツールは存在せず、サービスの一般公開もしておりません。
したがって、ローカル動作までとなっております。あらかじめご了承ください。

# サービス概要
サービス名：SmartChores

家事を管理するソフトウェアです。

「誰が、いつ、何の家事」を担当しているかを管理することで、家庭内の家事の負担を可視化します。

# 提供機能
下記に関する16個のAPIを実装しています。

- ユーザ管理機能
- 世帯管理機能
- 家事タスク管理機能
- サーバのヘルスチェック

# 技術スタック
Go/OpenAPI(Swagger)/Nginx/PostgreSQL/Docker/CircleCI/Gitなど

# このポートフォリオで証明できるスキル
- 要件定義(機能要件のみ)
- REST-APIの設計(OpenAPI)
- RDBの論理設計、物理設計
- Goによる実装(UT含む)
- Dockerを使ったローカル環境の構築
- CI環境の整備(CircleCI)
- Git操作
- 認証・認可(JWT)

# 動かしてみる
基本的な実行方法を記載します。

## サーバ側準備
前提として、Dockerがインストールされている必要があります。

コンテナ群起動：

```sh
$ ./serverrun.sh
```

DBマイグレーション：

```sh
 # make migrate
1/u initialize_schema (150.3ms)
```

アプリケーション起動：

```sh
# go run app/cmd/run-server/main.go
INFO[0000] Server started
```

## リクエスト実行
ユーザ作成：

```sh
curl --location --request POST 'http://localhost:8080/user' \
--header 'Content-Type: application/json' \
--data-raw '{
  "email": "test-user1@example.com",
  "user_name": "test-user1",
  "password": "123456"
}'
```
```json
{
    "user": {
        "user_id": 1,
        "user_name": "test-user1"
    }
}
```

ログイン：

```sh
curl --location --request POST 'http://localhost:8080/login?user_name=test-user1&password=123456'
```
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODczNjUzNDgsImlhdCI6IjIwMjAtMDQtMTlUMTU6NDk6MDguMzI1NjIyMyswOTowMCIsIm5hbWUiOiJ0ZXN0LXVzZXIxIn0.OMwFqWwixoZi9RElquMgfRENH3-l6x9_9P6-QJfxKjc"
}
```

世帯登録（世帯主として）：

```sh
curl --location --request POST 'http://localhost:8080/family' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODczNjUzNDgsImlhdCI6IjIwMjAtMDQtMTlUMTU6NDk6MDguMzI1NjIyMyswOTowMCIsIm5hbWUiOiJ0ZXN0LXVzZXIxIn0.OMwFqWwixoZi9RElquMgfRENH3-l6x9_9P6-QJfxKjc' \
--header 'Content-Type: application/json' \
--data-raw '{
  "family_name": "TestFamily1"
}'
```
```json
{
    "family": {
        "family_id": 1,
        "family_name": "TestFamily1"
    }
}
```

世帯メンバ追加：

```sh
curl --location --request POST 'http://localhost:8080/user' \
--header 'Content-Type: application/json' \
--data-raw '{
  "email": "test-user2@example.com",
  "user_name": "test-user2",
  "password": "123456"
}'
```
```json
{
    "user": {
        "user_id": 2,
        "user_name": "test-user2"
    }
}
```

```sh
curl --location --request POST 'http://localhost:8080/family/member' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODczNjUzNDgsImlhdCI6IjIwMjAtMDQtMTlUMTU6NDk6MDguMzI1NjIyMyswOTowMCIsIm5hbWUiOiJ0ZXN0LXVzZXIxIn0.OMwFqWwixoZi9RElquMgfRENH3-l6x9_9P6-QJfxKjc' \
--header 'Content-Type: application/json' \
--data-raw '{
  "member_name": "test-user2"
}'
```
```json
{
    "family": {
        "family_id": 1,
        "family_name": "TestFamily1"
    },
    "member": {
        "member_id": 2,
        "member_name": "test-user2"
    }
}
```

世帯メンバ一覧取得：

```sh
curl --location --request GET 'http://localhost:8080/family/members' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODczNjUzNDgsImlhdCI6IjIwMjAtMDQtMTlUMTU6NDk6MDguMzI1NjIyMyswOTowMCIsIm5hbWUiOiJ0ZXN0LXVzZXIxIn0.OMwFqWwixoZi9RElquMgfRENH3-l6x9_9P6-QJfxKjc'
```
```json
{
    "family": {
        "family_id": 1,
        "family_name": "TestFamily1"
    },
    "members": [
        {
            "member_id": 1,
            "member_name": "test-user1"
        },
        {
            "member_id": 2,
            "member_name": "test-user2"
        }
    ]
}
```

家事登録：

```sh
curl --location --request POST 'http://localhost:8080/task' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODczNjUzNDgsImlhdCI6IjIwMjAtMDQtMTlUMTU6NDk6MDguMzI1NjIyMyswOTowMCIsIm5hbWUiOiJ0ZXN0LXVzZXIxIn0.OMwFqWwixoZi9RElquMgfRENH3-l6x9_9P6-QJfxKjc' \
--header 'Content-Type: application/json' \
--data-raw '{
  "task_name": "make lunch",
  "member_name": "test-user1",
  "status": "done",
  "date": "2020-05-01"
}'
```
```json
{
    "family": {
        "family_id": 1,
        "family_name": "TestFamily1"
    },
    "task": {
        "task_id": 1,
        "task_name": "make lunch",
        "member_name": "test-user1",
        "status": "done",
        "date": "2020-05-01"
    }
}
```

家事一覧取得：

```sh
curl --location --request GET 'http://localhost:8080/tasks?date=2020-05-01' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODczNjUzNDgsImlhdCI6IjIwMjAtMDQtMTlUMTU6NDk6MDguMzI1NjIyMyswOTowMCIsIm5hbWUiOiJ0ZXN0LXVzZXIxIn0.OMwFqWwixoZi9RElquMgfRENH3-l6x9_9P6-QJfxKjc'
```
```json
{
    "family": {
        "family_id": 1,
        "family_name": "TestFamily1"
    },
    "tasks": [
        {
            "task_id": 1,
            "task_name": "make lunch",
            "member_name": "test-user1",
            "status": "done",
            "date": "2020-05-01"
        }
    ]
}
```

# 今後追加予定の機能
- React+TypeScriptを使ったSPAフロントエンド開発
- ゲストユーザでのログイン機能（各種操作権限無し）
- パスワードリセット機能
- スマートスピーカーを使ったUIの追加
- IaCを使ったインフラ構築
- デプロイパイプラインの構築
- サービスの一般リリース

# for developers
## IF
openapi/openapi.yamlをSwagger Editorなどでご確認ください。

## mockgenの実行例
```
$ cd app/server/repository
$ mockgen -source user.go -destination mock_user.go -package repository
```

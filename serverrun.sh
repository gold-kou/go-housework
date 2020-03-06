#!/bin/bash
base=${0##*/}

# dockerコンテナのビルドと起動
docker-compose up -d --build

# appコンテナにログイン
docker-compose exec app sh

# appコンテナから抜けたら自動的に停止とする
docker-compose down
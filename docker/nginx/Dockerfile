FROM nginx:1.17.4-alpine

# timezoneをJSTに変更
ENV TZ='Asia/Tokyo'

# ログファイルに関してホストOSとシンボリックリンク設定
RUN ln -sf /dev/stdout /var/log/nginx/access.log
RUN ln -sf /dev/stderr /var/log/nginx/error.log

# ローカルからconfigファイルコピー
ADD ./nginx/common/nginx.conf /etc/nginx/nginx.conf
ADD ./nginx/prd/default.conf /etc/nginx/conf.d/default.conf

# ホストOSに外部ボリュームマウント
VOLUME ["/var/cache/nginx"]

# ポート設定
EXPOSE 80

# フォアグラウンドで実行
CMD ["nginx", "-g", "daemon off;"]
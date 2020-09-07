# 概要

構造を整理する

# ファイル構造と役割

参考にしたサイト  
[クリーンアーキテクチャ概論](https://qiita.com/nrslib/items/a5f902c4defc83bd46b8#%E5%90%8C%E5%BF%83%E5%86%86%E3%81%AE%E5%9B%B3%E3%81%AE%E8%AA%AC%E6%98%8E)  
[レイヤードアーキテクチャと golang](https://yyh-gl.github.io/tech-blog/blog/go_web_api/)

```
pkg
|- constant: 定数を保存
|- db: DBの設定
|- dontext: contextとの接続（userIDの取得処理）
|- http
  |- middleware: 認証系
  |- response: レスポンスの形を定義
|
|- server
  |- domain/
    |- model/ : 扱うデータをまとめたもの
      |- user.go: ユーザ関連
      |- collection.go: APIのレスポンス（/collection/list）
      |- gacha.go: APIのレスポンス（/gacha/draw）
      |- gacha_probability.go: table:gacha_probabilityのコラムの形
      |- user.go: APIのレスポンス（/user/~）
      |- user_collection_item.go: table:user_collection_itemのコラムの形
    |
    |- repository/ : infraとusecaseをつなぐ役割
      |- collection_item.go
      |- gacha_probability.go
      |- user.go
      |- user_collection_item.go
  |
  |- infra/
    |- persistence/ : DBとの接続（MySQL）
      |- collection_item.go
      |- gacha_probability.go
      |- user.go
      |- user_collection_item.go
  |
  |- interface/
    |- handler/ : ハンドラ部分（リクエストとレスポンスを担当）
      |- collection/
      |- gacha/
      |- game/
      |- ranking/
      |- setting/
      |- user/
    |- middleware/ : userIDの認証部分
    |- myerror/ : エラーハンドリング部分
    |- response/ : レスポンスの形を定義
  |
  |- usecase/ : ハンドラ層から受け取った情報を使って処理をする．API処理の実質的な部分
    |- collection/
    |- gacha/
    |- game/
    |- ranking/
    |- user/
  |
  |- server.go
```

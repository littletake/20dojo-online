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
      |- error.go: エラー関連
    |
    |- repository/ : infraとusecaseをつなぐ役割
      |- user.go: ユーザ関連
  |
  |- infra/
    |- persistence/ : DBとの接続（MySQL）
      |- user.go: ユーザ関連
  |
  |- interface/
    |- handler/ : ハンドラ部分（リクエストとレスポンスを担当）
    |- error/ : エラーハンドリング
  |
  |- usecase/
    |- user.go
    |- game.go
  |
  |- server.go
```

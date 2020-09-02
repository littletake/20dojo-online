package server

import (
	"log"
	"net/http"

	"20dojo-online/pkg/server/infra/persistence"
	"20dojo-online/pkg/server/usecase"

	"20dojo-online/pkg/http/middleware"
	"20dojo-online/pkg/server/handler"
	"20dojo-online/pkg/server/initializer"
)

// Serve HTTPサーバを起動する
func Serve(addr string) {
	// レイヤードアーキテクチャ
	userPersistence := persistence.NewUserPersistence()
	userUseCase := usecase.NewUserUseCase(userPersistence)
	userHandler := handler.NewUserHandler(userUseCase)
	// ユーザ情報取得
	http.HandleFunc("/userL/get",
		get(middleware.Authenticate(userHandler.HandleUserLGet)))
	// ユーザ作成
	http.HandleFunc("/userL/create", post(userHandler.HandleUserLCreate))
	// // ユーザ情報更新
	// http.HandleFunc("/userL/update",
	// 	post(middleware.Authenticate(userHandler.HandleUserLUpdate)))
	// // インゲーム終了
	// http.HandleFunc("/game/finish",
	// 	post(middleware.Authenticate(userHandler.HandleGameFinish)))

	/* ===== URLマッピングを行う ===== */
	http.HandleFunc("/setting/get", get(handler.HandleSettingGet()))
	// http.HandleFunc("/user/create", post(handler.HandleUserCreate()))

	// 認証を行うmiddlewareを追加する
	// http.HandleFunc("/user/get",
	// 	get(middleware.Authenticate(handler.HandleUserGet())))
	// http.HandleFunc("/user/update",
	// 	post(middleware.Authenticate(handler.HandleUserUpdate())))

	// // インゲーム終了処理
	// http.HandleFunc("/game/finish",
	// 	post(middleware.Authenticate(handler.HandleGameFinish())))

	// ランキング情報取得
	http.HandleFunc("/ranking/list",
		get(middleware.Authenticate(handler.HandleRankingList())))

	// ガチャ実行
	http.HandleFunc("/gacha/draw",
		post(middleware.Authenticate(handler.HandleGachaDraw())))

	// コレクションアイテム一覧情報取得
	http.HandleFunc("/collection/list",
		get(middleware.Authenticate(handler.HandleCollectionList())))

	// 初期化
	// アイテム対応表の作成
	if err := initializer.CreateItemRatioSliceOnce(); err != nil {
		log.Println(err)
	}

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}

// get GETリクエストを処理する
func get(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodGet)
}

// post POSTリクエストを処理する
func post(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodPost)
}

// httpMethod 指定したHTTPメソッドでAPIの処理を実行する
func httpMethod(apiFunc http.HandlerFunc, method string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// CORS対応
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

		// プリフライトリクエストは処理を通さない
		if request.Method == http.MethodOptions {
			return
		}
		// 指定のHTTPメソッドでない場合はエラー
		if request.Method != method {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method Not Allowed"))
			return
		}

		// 共通のレスポンスヘッダを設定
		writer.Header().Add("Content-Type", "application/json")
		apiFunc(writer, request)
	}
}

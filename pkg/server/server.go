package server

import (
	"log"
	"net/http"

	"20dojo-online/pkg/server/infra/persistence"
	"20dojo-online/pkg/server/interface/handler"
	"20dojo-online/pkg/server/interface/middleware"
	"20dojo-online/pkg/server/usecase"
)

// Serve HTTPサーバを起動する
func Serve(addr string) {
	// レイヤードアーキテクチャ
	userPersistence := persistence.NewUserPersistence()
	cItemPersistence := persistence.NewCItemPersistence()
	ucItemPersistence := persistence.NewUCItemPersistence()
	gachaProbPersistence := persistence.NewGachaProbPersistence()

	userUseCase := usecase.NewUserUseCase(userPersistence)
	gameUseCase := usecase.NewGameUseCase(userPersistence)
	rankingUseCase := usecase.NewRankingUseCase(userPersistence)
	gachaUseCase := usecase.NewGachaUseCase(userPersistence, cItemPersistence, ucItemPersistence, gachaProbPersistence)
	collectionUseCase := usecase.NewCollectionUseCase(userPersistence, cItemPersistence, ucItemPersistence)

	settingHandler := handler.NewSettingHandler()
	userHandler := handler.NewUserHandler(userUseCase)
	gameHandler := handler.NewGameHandler(gameUseCase)
	rankingHandler := handler.NewRankingHandler(rankingUseCase)
	gachaHandler := handler.NewGachaHandler(gachaUseCase)
	collectionHandler := handler.NewCollectionHandler(collectionUseCase)

	middleware := middleware.NewMiddleware(userUseCase)

	/* ===== URLマッピングを行う ===== */
	// // 設定情報取得
	http.HandleFunc("/setting/get", get(settingHandler.HandleSettingGet))
	// ユーザ情報取得
	http.HandleFunc("/user/get",
		get(middleware.Authenticate(userHandler.HandleUserGet)))
	// ユーザ作成
	http.HandleFunc("/user/create", post(userHandler.HandleUserCreate))
	// ユーザ情報更新
	http.HandleFunc("/user/update",
		post(middleware.Authenticate(userHandler.HandleUserUpdate)))
	// インゲーム終了
	http.HandleFunc("/game/finish",
		post(middleware.Authenticate(gameHandler.HandleGameFinish)))
	// ランキング情報取得
	http.HandleFunc("/ranking/list",
		get(middleware.Authenticate(rankingHandler.HandleRankingList)))
	// ガチャ実行
	http.HandleFunc("/gacha/draw",
		post(middleware.Authenticate(gachaHandler.HandleGachaDraw)))
	// コレクションアイテム一覧情報取得
	http.HandleFunc("/collection/list",
		get(middleware.Authenticate(collectionHandler.HandleCollectionList)))

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

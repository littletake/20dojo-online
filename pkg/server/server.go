package server

import (
	"log"
	"net/http"
	"time"

	cp "20dojo-online/pkg/server/infra/persistence/collectionitem"
	gp "20dojo-online/pkg/server/infra/persistence/gachaprobability"
	up "20dojo-online/pkg/server/infra/persistence/user"
	ucp "20dojo-online/pkg/server/infra/persistence/usercollectionitem"
	ch "20dojo-online/pkg/server/interface/handler/collection"
	gch "20dojo-online/pkg/server/interface/handler/gacha"
	gh "20dojo-online/pkg/server/interface/handler/game"
	rh "20dojo-online/pkg/server/interface/handler/ranking"
	sh "20dojo-online/pkg/server/interface/handler/setting"
	uh "20dojo-online/pkg/server/interface/handler/user"
	"20dojo-online/pkg/server/interface/middleware"
	cu "20dojo-online/pkg/server/usecase/collection"
	gcu "20dojo-online/pkg/server/usecase/gacha"
	gu "20dojo-online/pkg/server/usecase/game"
	ru "20dojo-online/pkg/server/usecase/ranking"
	uu "20dojo-online/pkg/server/usecase/user"
)

// Serve HTTPサーバを起動する
func Serve(addr string) {
	// レイヤードアーキテクチャ
	userPersistence := up.NewUserPersistence()
	cItemPersistence := cp.NewCItemPersistence()
	ucItemPersistence := ucp.NewUCItemPersistence()
	gachaProbPersistence := gp.NewGachaProbPersistence()

	userUseCase := uu.NewUserUseCase(userPersistence)
	gameUseCase := gu.NewGameUseCase(userPersistence)
	rankingUseCase := ru.NewRankingUseCase(userPersistence)
	gachaUseCase := gcu.NewGachaUseCase(userPersistence, cItemPersistence, ucItemPersistence, gachaProbPersistence, time.Now().UnixNano())
	collectionUseCase := cu.NewCollectionUseCase(userPersistence, cItemPersistence, ucItemPersistence)

	settingHandler := sh.NewSettingHandler()
	userHandler := uh.NewUserHandler(userUseCase)
	gameHandler := gh.NewGameHandler(gameUseCase)
	rankingHandler := rh.NewRankingHandler(rankingUseCase)
	gachaHandler := gch.NewGachaHandler(gachaUseCase)
	collectionHandler := ch.NewCollectionHandler(collectionUseCase)
	// TODO: httpとauthを一つにしたい
	m := middleware.NewMiddleware(userUseCase)

	/* ===== URLマッピングを行う ===== */
	// // 設定情報取得
	http.HandleFunc("/setting/get", m.Get(settingHandler.HandleSettingGet()))
	// ユーザ情報取得
	http.HandleFunc("/user/get",
		m.Get(m.Authenticate(userHandler.HandleUserGet())))
	// ユーザ作成
	http.HandleFunc("/user/create", m.Post(userHandler.HandleUserCreate()))
	// ユーザ情報更新
	http.HandleFunc("/user/update",
		m.Post(m.Authenticate(userHandler.HandleUserUpdate())))
	// インゲーム終了
	http.HandleFunc("/game/finish",
		m.Post(m.Authenticate(gameHandler.HandleGameFinish())))
	// ランキング情報取得
	http.HandleFunc("/ranking/list",
		m.Get(m.Authenticate(rankingHandler.HandleRankingList())))
	// ガチャ実行
	http.HandleFunc("/gacha/draw",
		m.Post(m.Authenticate(gachaHandler.HandleGachaDraw())))
	// コレクションアイテム一覧情報取得
	http.HandleFunc("/collection/list",
		m.Get(m.Authenticate(collectionHandler.HandleCollectionList())))

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}

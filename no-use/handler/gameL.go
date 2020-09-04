// package handler

// // HandleGameFinish インゲーム終了処理
// func (uh userHandler) HandleGameFinish(writer http.ResponseWriter, request *http.Request) {
// 	// GameFinishRequest ゲーム終了request
// 	type GameFinishRequest struct {
// 		Score int32 `json:"score"`
// 	}
// 	// GameFinishResponse ゲーム終了response
// 	type GameFinishResponse struct {
// 		Coin int32 `json:"coin"`
// 	}

// 	// コンテキストからuserID取得
// 	ctx := request.Context()
// 	userID := dcontext.GetUserIDFromContext(ctx)
// 	if userID == "" {
// 		myErr := uh.userUseCase.CreateMyErr(
// 			fmt.Errorf("userID is empty"),
// 			500,
// 		)
// 		DoErr(writer, myErr)
// 		return
// 	}
// 	// リクエストBodyから更新後情報を取得
// 	var requestBody GameFinishRequest
// 	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
// 		myErr := uh.userUseCase.CreateMyErr(
// 			err,
// 			500,
// 		)
// 		DoErr(writer, myErr)
// 		return
// 	}
// 	// coinとscoreを更新
// 	coin, myErr := uh.userUseCase.UpdateCoinAndHighScore(userID, requestBody.Score)
// 	if myErr != nil {
// 		DoErr(writer, myErr)
// 		return
// 	}
// 	// レスポンスに必要な情報を詰めて返却
// 	response.Success(writer, &GameFinishResponse{
// 		Coin: coin,
// 	})
// }
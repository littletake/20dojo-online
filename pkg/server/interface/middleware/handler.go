package middleware

// import (
// 	"net/http"

// 	usecase "20dojo-online/pkg/server/usecase/user"
// )

// // Middleware middlewareのインターフェース
// type Middleware interface {
// 	Authenticate(http.HandlerFunc) http.HandlerFunc
// 	Get(http.Handler) http.Handler
// 	Post(http.HandlerFunc) http.HandlerFunc
// }

// type middleware struct {
// 	userUseCase usecase.UserUseCase
// }

// // NewMiddleware userUseCaseと疎通
// func NewMiddleware(uu usecase.UserUseCase) Middleware {
// 	return &middleware{
// 		userUseCase: uu,
// 	}
// }

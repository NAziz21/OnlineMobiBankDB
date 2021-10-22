package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

const UserClaimGet = "User"


func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		tokenString = strings.TrimSpace(tokenString)

		// Parse the token
		claim := jwt.RegisteredClaims{}
		_, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
			return []byte("MySecretCode"), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "auth error",err)
			return
		}
		r = r.WithContext(WriteTokenToContext(r.Context(), claim))
		next.ServeHTTP(w,r)
	})
}

func WriteTokenToContext(ctx context.Context, claim jwt.RegisteredClaims) context.Context {
	return context.WithValue(ctx, UserClaimGet, claim)
}

func ReadTokenFromContext(ctx context.Context) jwt.RegisteredClaims {
	claim := ctx.Value(UserClaimGet).(jwt.RegisteredClaims)
	return claim
}



package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"market/internal/logic"
	"market/internal/svc"
)

func getCategoriesHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetCategoriesLogic(r.Context(), ctx)
		resp, err := l.GetCategories()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

package handler

import (
	"github.com/tal-tech/go-zero/rest/httpx"
	"market/internal/logic"
	"market/internal/svc"
	"market/internal/types"
	"net/http"
)

func getMyNftHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMyNftReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetMyNftLogic(r.Context(), ctx)
		resp, err := l.GetMyNft(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

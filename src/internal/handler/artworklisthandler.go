package handler

import (
	"net/http"

	"market/internal/logic"
	"market/internal/svc"
	"market/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func ArtworkListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ArtworkListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewArtworkListLogic(r.Context(), ctx)
		resp, err := l.ArtworkList(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

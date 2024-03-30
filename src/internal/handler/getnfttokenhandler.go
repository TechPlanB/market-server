package handler

import (
	"net/http"

	"market/internal/logic"
	"market/internal/svc"
	"market/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func getNFTTokenHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NFTTokenIdsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetNFTTokenLogic(r.Context(), ctx)
		resp, err := l.GetNFTTokenIds(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func getNFTTokenByNameHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NFTTokenIdsByNameReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetNFTTokenLogic(r.Context(), ctx)
		resp, err := l.GetNFTTokenIdsByName(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

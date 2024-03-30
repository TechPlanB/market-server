package handler

import (
	"market/internal/types"
	"net/http"

	"market/internal/logic"
	"market/internal/svc"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func getSyncNFTListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetNFTListenersLogic(r.Context(), ctx)
		resp, err := l.GetSyncNftList()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func syncNFTHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SyncNFTReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetNFTListenersLogic(r.Context(), ctx)
		resp, err := l.SyncNFT(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func sync721WithoutTokenIdHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Sync721WithoutTokenIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetNFTListenersLogic(r.Context(), ctx)
		resp, err := l.Sync721WithoutTokenId(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

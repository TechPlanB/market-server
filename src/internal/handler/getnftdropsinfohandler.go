package handler

import (
	"net/http"

	"market/internal/logic"
	"market/internal/svc"
	"market/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func getNFTDropsInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NFTDropsInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetNFTDropsInfoLogic(r.Context(), ctx)
		resp, err := l.GetNFTDropsInfo(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

package handler

import (
	"net/http"

	"market/internal/logic"
	"market/internal/svc"
	"market/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func NFTSalesListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NFTSalesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewNFTSalesListLogic(r.Context(), ctx)
		resp, err := l.NFTSalesList(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

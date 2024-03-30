package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"market/internal/logic"
	"market/internal/svc"
	"market/internal/types"
)

func getSaleContractHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SaleContractReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetSaleContractLogic(r.Context(), ctx)
		resp, err := l.GetSaleContract(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

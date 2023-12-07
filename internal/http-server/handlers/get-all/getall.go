package getall

import (
	"context"
	"github.com/go-chi/render"
	"github.com/rvecwxqz/apod/internal/core"
	"github.com/rvecwxqz/apod/internal/lib/api/response"
	"log"
	"net/http"
)

type Resp struct {
	response.Response
	Info []core.APODInfo `json:"info"`
}

type AllInfoProvider interface {
	GetAllInfo(ctx context.Context) ([]core.APODInfo, error)
}

type AllImagesProvider interface {
	SetUrls(names []core.APODInfo)
}

func New(provider AllInfoProvider, iProvider AllImagesProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := provider.GetAllInfo(r.Context())
		if err != nil {
			log.Println("getall handler error:", err)
			render.JSON(w, r, response.Error("storage error"))
			return
		}

		iProvider.SetUrls(info)

		render.JSON(w, r, Resp{
			Response: response.OK(),
			Info:     info,
		})
	}
}

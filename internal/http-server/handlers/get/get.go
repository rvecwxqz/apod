package get

import (
	"context"
	"github.com/go-chi/render"
	"github.com/rvecwxqz/apod/internal/core"
	"github.com/rvecwxqz/apod/internal/lib/api/response"
	"log"
	"net/http"
)

type Request struct {
	Date core.Date `json:"date"`
}

type Resp struct {
	response.Response
	Info core.APODInfo `json:"info"`
}

type InfoProvider interface {
	GetInfo(ctx context.Context, date core.Date) (core.APODInfo, error)
}

type ImageProvider interface {
	GetUrl(oName string) string
}

func New(provider InfoProvider, iProvider ImageProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Println("get handler error:", err)
			render.JSON(w, r, response.Error("error reading body"))
			return
		}

		info, err := provider.GetInfo(r.Context(), req.Date)
		if err != nil {
			log.Println("get handler, storage error:", err)
			render.JSON(w, r, response.Error("storage error"))
			return
		}

		info.Image = iProvider.GetUrl(info.Title)

		render.JSON(w, r, Resp{
			Response: response.OK(),
			Info:     info,
		})

	}
}

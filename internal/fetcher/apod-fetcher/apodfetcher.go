package apodfetcher

import (
	"context"
	"fmt"
	"github.com/go-chi/render"
	"github.com/rvecwxqz/apod/internal/core"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	host            = "api.nasa.gov"
	path            = "/planetary/apod"
	defaultWaitTime = 1 * time.Second
)

type Response struct {
	Date        core.Date `json:"date"`
	Title       string    `json:"title"`
	Explanation string    `json:"explanation"`
	URL         string    `json:"hdurl"`
}

type APODFetcher struct {
	ctx          context.Context
	u            url.URL
	errCh        chan error
	ticker       *time.Ticker
	client       *http.Client
	retriesCount int
}

type InfoSaver interface {
	SaveInfo(
		ctx context.Context,
		title string,
		explanation string,
		date core.Date,
	) error
}

type ImageUploader interface {
	UploadImage(ctx context.Context, reader io.Reader, oName string, oSize int64) error
}

func New(
	ctx context.Context,
	apiKey string,
	workerInterval time.Duration,
	retriesCount int,
	saver InfoSaver,
	uploader ImageUploader,
) {
	u := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path,
	}
	q := u.Query()
	q.Set("api_key", apiKey)
	u.RawQuery = q.Encode()

	fetcher := APODFetcher{
		ctx:          ctx,
		u:            u,
		errCh:        make(chan error, 10),
		ticker:       time.NewTicker(workerInterval),
		client:       &http.Client{},
		retriesCount: retriesCount,
	}

	go func() {
		for {
			go fetcher.Fetch(saver, uploader)

			select {
			case <-fetcher.ctx.Done():
				fetcher.ticker.Stop()
				return
			case <-fetcher.ticker.C:
				continue
			}

		}
	}()

	go func() {
		for {
			select {
			case <-fetcher.ctx.Done():
				close(fetcher.errCh)
				return
			case v := <-fetcher.errCh:
				log.Println(v)
			}
		}
	}()

}

func (f *APODFetcher) Fetch(saver InfoSaver, uploader ImageUploader) {
	resp, err := f.getWithRetry(f.u.String())
	if err != nil {
		f.errCh <- NewRequestError(err)
		return
	}

	var r Response
	err = render.DecodeJSON(resp.Body, &r)
	if err != nil {
		f.errCh <- NewDecodeJSONError(err)
		return
	}

	err = f.fetchAndUploadImg(r.URL, r.Title, uploader)
	if err != nil {
		f.errCh <- NewRequestError(err)
		return
	}

	err = saver.SaveInfo(f.ctx, r.Title, r.Explanation, r.Date)
	if err != nil {
		f.errCh <- err
		return
	}
}

func (f *APODFetcher) fetchAndUploadImg(url, oName string, u ImageUploader) error {

	resp, err := f.getWithRetry(url)
	if err != nil {
		return fmt.Errorf("get image error: %w", err)
	}
	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return fmt.Errorf("parse int error: %w", err)
	}
	defer resp.Body.Close()

	log.Println("image fetched len", size)
	err = u.UploadImage(f.ctx, resp.Body, oName, size)

	if err != nil {
		return fmt.Errorf("fetcher upload error: %w", err)
	}

	return nil
}

func (f *APODFetcher) getWithRetry(url string) (*http.Response, error) {
	var (
		retries = f.retriesCount
		err     error
		resp    *http.Response
	)

	for retries > 0 {
		resp, err = f.client.Get(url)
		if err == nil {
			break
		} else if err != nil && retries > 1 {
			retries--
			time.Sleep(defaultWaitTime)
		} else {
			return nil, err
		}
	}
	if resp == nil {
		return nil, fmt.Errorf("nil response")
	}

	return resp, nil
}

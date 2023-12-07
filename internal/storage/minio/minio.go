package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rvecwxqz/apod/internal/core"
	"io"
	"net/url"
)

type Provider struct {
	user          string
	password      string
	url           string
	bName         string
	port          string
	serverAddress string
	client        *minio.Client
}

func NewProvider(user, password, url, bName, serverAddress, port string) (*Provider, error) {
	m := Provider{
		user:          user,
		password:      password,
		url:           url,
		bName:         bName,
		port:          port,
		serverAddress: serverAddress,
	}
	err := m.Connect()

	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Provider) Connect() error {
	var err error
	m.client, err = minio.New(m.url, &minio.Options{
		Creds: credentials.NewStaticV4(m.user, m.password, ""),
	})
	if err != nil {
		return fmt.Errorf("minio connect error: %w", err)
	}

	return nil
}

func (m *Provider) UploadImage(ctx context.Context, reader io.Reader, oName string, oSize int64) error {
	_, err := m.client.PutObject(
		ctx,
		m.bName,
		oName,
		reader,
		oSize,
		minio.PutObjectOptions{ContentType: "image/png"},
	)
	if err != nil {
		return fmt.Errorf("error upload image: %w", err)
	}
	return nil
}

func (m *Provider) GetUrl(oName string) string {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%v:%v", m.serverAddress, m.port),
		Path:   fmt.Sprintf("%v/%v", m.bName, oName),
	}

	return u.String()
}

func (m *Provider) SetUrls(info []core.APODInfo) {
	for i := 0; i < len(info); i++ {
		info[i].Image = m.GetUrl(info[i].Title)
	}
}

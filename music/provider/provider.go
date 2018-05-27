package provider

import "time"

type Song struct {
	Name      string
	URL       string
	Duration  time.Duration
	Thumbnail string
}

type Provider interface {
	URL(url string) ([]Song, error)
	Search(keyword string, maxResult int) ([]Song, error)
}

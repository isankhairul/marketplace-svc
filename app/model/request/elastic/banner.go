package requestelastic

import "fmt"

type BannerRequest struct {
	Query        string `json:"q" schema:"q" binding:"omitempty"`
	CategorySlug string `json:"category_slug" schema:"category_slug" binding:"omitempty"`
	Slug         string `json:"slug" schema:"slug" binding:"omitempty"`
	Fields       string `json:"fields" schema:"fields" binding:"omitempty"`
	ID           int    `json:"id" schema:"id" binding:"omitempty"`
	ChannelID    int    `json:"channel_id" schema:"channel_id" binding:"omitempty"`
	Page         int    `json:"page" schema:"page" binding:"omitempty"`
	Limit        int    `json:"limit" schema:"limit" binding:"omitempty"`
}

func (b BannerRequest) ToString() string {
	return fmt.Sprintf("%s-%s-%s-%s-%d-%d-%d-%d", //nolint:govet
		b.Query, b.Fields, b.Slug, b.CategorySlug, b.ID, b.ChannelID, b.Page, b.Limit)
}

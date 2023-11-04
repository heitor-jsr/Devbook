package models

import "time"

type Puyblication struct {
	Id uint64 `json:"id,omitempty"`
	Title string `json:"titulo,omitempty"`
	Content string `json:"conteúdo,omitempty"`
	AuthorId uint64 `json:"autorId,omitempty"`
	AuthorNick string `json:"autorNick,omitempty"`
	Likes uint64 `json:"likes"`
	CriadaEm time.Time `json:"criadaEm,imitempty"`
}
module github.com/bryce4651/pkg/es

go 1.21.4

replace github.com/bryce4651/pkg/log => ../log

require (
	github.com/bryce4651/pkg/log v0.0.0-00010101000000-000000000000
	github.com/olivere/elastic/v7 v7.0.32
)

require (
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)

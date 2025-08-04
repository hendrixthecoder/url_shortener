module github.com/hendrixthecoder/url_shortener

go 1.23.0

toolchain go1.23.11

require (
	github.com/boj/redistore v1.4.1
	github.com/go-chi/chi/v5 v5.2.2
	github.com/go-chi/cors v1.2.2
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
)

require github.com/lib/pq v1.10.9

require github.com/gorilla/csrf v1.7.3 // indirect

require (
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	golang.org/x/crypto v0.40.0
)

package main

//go:generate go run github.com/99designs/gqlgen generate

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/af-afk/stylusup.sh/cmd/popcon.stylusup.sh/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

const (
	// EnvBackendType to use to listen the server with, (http|lambda).
	EnvBackendType = "STYLUSUP_LISTEN_BACKEND"

	// EnvListenAddr to listen on when using the HTTP option.
	EnvListenAddr = "STYLUSUP_LISTEN_ADDR"

	// EnvDatabaseUri to chat with to query info about the database.
	EnvDatabaseUri = "STYLUSUP_DATABASE_URI"
)

// ChangelogLen to send to the user at max on request for the changelog endpoint.
const ChangelogLen = 10

type corsMiddleware struct {
	srv *handler.Server
}

func (m corsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	ipAddr := strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
	ctx := context.WithValue(r.Context(), "ip", ipAddr)
	m.srv.ServeHTTP(w, r.WithContext(ctx))
}

func main() {
	db, err := sql.Open("postgres", os.Getenv(EnvDatabaseUri))
	if err != nil {
		log.Fatal("database open", err)
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			Db: db,
		},
	}))
	http.Handle("/", corsMiddleware{srv})
	http.Handle("/playground", playground.Handler("9lives.so playground", "/"))
	switch typ := os.Getenv(EnvBackendType); typ {
	case "lambda":
		lambda.Start(httpadapter.NewV2(http.DefaultServeMux).ProxyWithContext)
	case "http":
		err := http.ListenAndServe(os.Getenv(EnvListenAddr), nil)
		log.Fatalf( // This should only return if there's an error.
			"err listening, %#v not set?: %v",
			EnvListenAddr,
			err,
		)
	default:
		log.Fatalf(
			"unexpected listen type: %#v, use either (lambda|http) for STYLUSUP_LISTEN_BACKEND",
			typ,
		)
	}
}

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.73

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/af-afk/stylusup.sh/cmd/popcon.stylusup.sh/graph/model"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, arch string, lang string, os string) (*bool, error) {
	// We intentionally don't record the same user's information more than
	// once if we can. At this point, in the client side code here, we hash
	// their information as well, using the recorded language as the seed for
	// a HMAC interaction. We don't allow doubling up with the IP info
	// that they gave us.
	ipAddr := ctx.Value("ip").(string)
	if ipAddr == "" {
		return nil, fmt.Errorf("X-Forwarded-For not set")
	}
	ipHash := hashIpAddr(ipAddr, arch, lang, os)
	s := `
INSERT INTO stylusup_popularity_1 (ip_hash, lang, arch, os)
VALUES ($1, $2, $3, $4)`
	if _, err := r.Db.Exec(s, ipHash, lang, arch, os); err != nil {
		slog.Error("failed to insert", "err", err)
		return nil, fmt.Errorf("insertion")
	}
	v := true
	return &v, nil
}

// Popularity is the resolver for the popularity field.
func (r *queryResolver) Popularity(ctx context.Context) ([]*model.Popularity, error) {
	q := `
SELECT lang, arch, os, COUNT(*) AS n
FROM stylusup_popularity_1
GROUP BY lang, arch, os
ORDER BY n DESC
LIMIT 1`
	rows, err := r.Db.Query(q)
	if err != nil {
		slog.Error("failed to query row", "err", err)
		return nil, fmt.Errorf("query row")
	}
	defer rows.Close()
	pop := make([]*model.Popularity, 0)
	for rows.Next() {
		p := new(model.Popularity)
		if err := rows.Scan(&p.Lang, &p.Arch, &p.Os, &p.Instances); err != nil {
			slog.Error("scan row", "err", err)
			return nil, fmt.Errorf("scan row")
		}
		pop = append(pop, p)
	}
	return pop, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

package asymmetric

import (
	"context"
	"encoding/json"
	"net/http"
)

func (r *Config) getJwks(ctx context.Context, cfg *Config) (err error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.JwksURL, nil)
	if err != nil {
		return
	}

	httpClient := http.Client{}
	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		return
	}
	defer httpRes.Body.Close()

	json.NewDecoder(httpRes.Body).Decode(&cfg.jwks)

	return
}

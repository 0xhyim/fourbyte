package fourbyte

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type link string

func (l *link) Follow(ctx context.Context) (sr *SignaturesResponse, err error) {
	if *l == "" {
		return nil, errors.New("no link to follow")
	}
	url, err := url.Parse(fmt.Sprintf("%v", *l))
	if err != nil {
		return nil, err
	}
	_, data, err := newRequest(*url).Get(ctx)
	if err != nil {
		return nil, err
	}
	return sr, json.Unmarshal(data, &sr)
}

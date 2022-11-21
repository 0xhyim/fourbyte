package fourbyte

import (
	"strconv"
)

type filterOption interface {
	apply(*request)
}

type funcFilterOption struct {
	f func(*request)
}

func (ffo *funcFilterOption) apply(to *request) {
	ffo.f(to)
}

func newFuncFilterOption(f func(*request)) *funcFilterOption {
	return &funcFilterOption{
		f: f,
	}
}

func WithPageNumber(num int) filterOption {
	return newFuncFilterOption(func(r *request) {
		r.setQueryValue("page", strconv.Itoa(num))
	})
}

func WithHexSignature(hex string) filterOption {
	return newFuncFilterOption(func(r *request) {
		r.addQueryValue("hex_signature", hex)
	})
}

func WithTextSignature(text string) filterOption {
	return newFuncFilterOption(func(r *request) {
		r.addQueryValue("text_signature", text)
	})
}

func withId(id int) filterOption {
	return newFuncFilterOption(func(r *request) {
		r.url = *r.url.JoinPath(strconv.Itoa(id))
	})
}

package fourbyte

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

type Signature struct {
	Id            int    `json:"id,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	TextSignature string `json:"text_signature,omitempty"`
	ByteSignature string `json:"bytes_signature,omitempty"`
	HexSignature  string `json:"hex_signature,omitempty"`
}

type SignaturesResponse struct {
	Next       *link       `json:"next"`
	Previous   *link       `json:"previous"`
	Count      int         `json:"count"`
	Signatures []Signature `json:"results"`
}

type sourceImport struct {
	SourceCode string `json:"source_code"`
}

type abiImport struct {
	ContractAbi string `json:"contract_abi"`
}

type ImportResponse struct {
	Processed  int `json:"num_processed"`
	Imported   int `json:"num_imported"`
	Duplicated int `json:"num_duplicates"`
	Ignored    int `json:"num_ignored"`
}

func BaseUrl() *url.URL {
	return &url.URL{Scheme: "https", Host: "www.4byte.directory"}
}

func GetFunctionSignatures(ctx context.Context, opts ...filterOption) (*SignaturesResponse, error) {
	return getSignatures(ctx, *BaseUrl().JoinPath("/api/v1/signatures/"), opts...)
}

func GetFunctionSignatureById(ctx context.Context, id int) (*Signature, error) {
	return getSignatureById(ctx, *BaseUrl().JoinPath("/api/v1/signatures/"), id)
}

func CreateFunctionSignature(ctx context.Context, textSignature string) (*Signature, error) {
	return postSignature(ctx, *BaseUrl().JoinPath("/api/v1/signatures/"), textSignature)
}

func GetEventSignatures(ctx context.Context, opts ...filterOption) (*SignaturesResponse, error) {
	return getSignatures(ctx, *BaseUrl().JoinPath("/api/v1/event-signatures/"), opts...)
}

func GetEventSignatureById(ctx context.Context, id int) (*Signature, error) {
	return getSignatureById(ctx, *BaseUrl().JoinPath("/api/v1/event-signatures/"), id)
}

func CreateEventSignature(ctx context.Context, textSignature string) (*Signature, error) {
	return postSignature(ctx, *BaseUrl().JoinPath("/api/v1/event-signatures/"), textSignature)
}

func ImportFromSourceCode(ctx context.Context, sourceCode string) (resp *ImportResponse, err error) {
	body, _ := json.Marshal(sourceImport{SourceCode: sourceCode})
	return postSolidity(ctx, *BaseUrl().JoinPath("/api/v1/import-solidity/"), body)
}

func ImportFromABI(ctx context.Context, contractAbi string) (resp *ImportResponse, err error) {
	body, _ := json.Marshal(abiImport{ContractAbi: contractAbi})
	return postSolidity(ctx, *BaseUrl().JoinPath("/api/v1/import-abi/"), body)
}

func getSignatures(ctx context.Context, url url.URL, opts ...filterOption) (resp *SignaturesResponse, err error) {
	_, data, err := newRequest(url, opts...).Get(ctx)
	if err != nil {
		return nil, err
	}
	return resp, json.Unmarshal(data, &resp)
}

func getSignatureById(ctx context.Context, url url.URL, id int) (sig *Signature, err error) {
	_, data, err := newRequest(url, withId(id)).Get(ctx)
	if err != nil {
		return nil, err
	}
	return sig, json.Unmarshal(data, &sig)
}

func postSignature(ctx context.Context, url url.URL, textSignature string) (sig *Signature, err error) {
	body, _ := json.Marshal(Signature{TextSignature: textSignature})
	status, data, err := newRequest(url).Post(ctx, "application/json", body)
	if err != nil {
		return nil, err
	}
	if status >= http.StatusBadRequest {
		return nil, apiError(data)
	}
	return sig, json.Unmarshal(data, &sig)
}

func postSolidity(ctx context.Context, url url.URL, body []byte) (resp *ImportResponse, err error) {
	status, data, err := newRequest(url).Post(ctx, "application/json", body)
	if err != nil {
		return nil, err
	}
	if status >= http.StatusBadRequest {
		return nil, apiError(data)
	}
	return resp, json.Unmarshal(data, &resp)
}

func apiError(data []byte) error {
	errResponse := make(map[string][]string)
	json.Unmarshal(data, &errResponse)
	for k, v := range errResponse {
		if len(v[0]) == 0 {
			return errors.New(string(data))
		}
		msg := nonAlphanumericRegex.ReplaceAllString(v[0], "")
		return fmt.Errorf("*%v* %s", k, strings.ToLower(msg))
	}
	return errors.New(string(data))
}

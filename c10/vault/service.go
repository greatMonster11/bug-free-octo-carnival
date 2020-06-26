package vault

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/crypto/bcrypt"
)

// Service provides password hashing capabilities
type Service interface {
	Hash(ctx context.Context, password string) (string, error)
	Validate(ctx context.Context, password, hash string) (bool, error)
}

type vaultService struct{}

// NewService makes a new Service
func NewService() Service {
	return vaultService{}
}

// Hashing and validiting password with bcrypt
func (vaultService) Hash(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (vaultService) Validate(ctx context.Context, password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, nil
	}

	return true, nil
}

type hashRequest struct {
	Password string `json:"password"`
}

type hashResponse struct {
	Hash string `json:"hash"`
	Err  string `json:"err,omitempty"`
}

type validateRequest struct {
	Password string `json:"password"`
	Hash     string `json:"hash"`
}

type validateResponse struct {
	Valid bool   `json:"valid"`
	Err   string `json:"err,omitempty"`
}

// helper method that will able to decode the http json's body
func decodeHashRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req hashRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	return req, err
}

func decodeValidateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req validateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	return req, err
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

func MakeHashEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(hashRequest)
		v, err := srv.Hash(ctx, req.Password)
		if err != nil {
			return hashResponse{v, err.Error()}, nil
		}
		return hashResponse{v, ""}, nil
	}
}

func MakeValidateEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(validateRequest)
		v, err := srv.Validate(ctx, req.Password, req.Hash)
		if err != nil {
			return validateResponse{false, err.Error()}, nil
		}

		return validateResponse{v, ""}, nil
	}
}

// Endpoint represents all endpoint for the vault Service.
type Endpoints struct {
	HashEndpoint     endpoint.Endpoint
	ValidateEndpoint endpoint.Endpoint
}

// Hash uses the HashEndpoint to hash a password.
func (e Endpoints) Hash(ctx context.Context, password string) (string, error) {
	req := hashRequest{Password: password}
	resp, err := e.HashEndpoint(ctx, req)
	if err != nil {
		return "", err
	}

	hashResp := resp.(hashResponse)
	if hashResp.Err != "" {
		return "", errors.New(hashResp.Err)
	}
	return hashResp.Hash, nil
}

// Validate uses the HashValidate to validate a hash string
func (e Endpoints) Validate(ctx context.Context, hash, password string) (bool, error) {
	req := validateRequest{Hash: hash, Password: password}
	resp, err := e.ValidateEndpoint(ctx, req)
	if err != nil {
		return false, err
	}

	validateResp := resp.(validateResponse)
	if validateResp.Err != "" {
		return false, errors.New(validateResp.Err)
	}
	return validateResp.Valid, nil
}

package testutils

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"certitrack/internal/services"
)

type RegisterRequestOptions struct {
	Email     *string
	Password  *string
	FirstName *string
	LastName  *string
	Phone     *string
}

type RequestBuilder struct {
	services.RegisterRequest
}

func NewRegisterRequest(opts ...func(*RegisterRequestOptions)) *RequestBuilder {
	options := RegisterRequestOptions{
		Email:     stringPtr("test-" + uuid.New().String() + "@example.com"),
		Password:  stringPtr("ValidPass123!"),
		FirstName: stringPtr("Test"),
		LastName:  stringPtr("User"),
		Phone:     stringPtr("+1234567890"),
	}

	for _, opt := range opts {
		opt(&options)
	}

	return &RequestBuilder{
		RegisterRequest: services.RegisterRequest{
			Email:     *options.Email,
			Password:  *options.Password,
			FirstName: *options.FirstName,
			LastName:  *options.LastName,
			Phone:     *options.Phone,
		},
	}
}

func WithEmail(email string) func(*RegisterRequestOptions) {
	return func(o *RegisterRequestOptions) { o.Email = &email }
}

func WithPassword(password string) func(*RegisterRequestOptions) {
	return func(o *RegisterRequestOptions) { o.Password = &password }
}

func WithFirstName(firstName string) func(*RegisterRequestOptions) {
	return func(o *RegisterRequestOptions) { o.FirstName = &firstName }
}

func WithLastName(lastName string) func(*RegisterRequestOptions) {
	return func(o *RegisterRequestOptions) { o.LastName = &lastName }
}

func WithPhone(phone string) func(*RegisterRequestOptions) {
	return func(o *RegisterRequestOptions) { o.Phone = &phone }
}

func stringPtr(s string) *string {
	return &s
}

func (r *RequestBuilder) ToJSON() string {
	jsonData, err := json.Marshal(r.RegisterRequest)
	if err != nil {
		panic(fmt.Sprintf("error al convertir a JSON: %v", err))
	}
	return string(jsonData)
}

func (r *RequestBuilder) ToRegisterRequest() *services.RegisterRequest {
	return &r.RegisterRequest
}

package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Filter struct {
	Search         string `json:"search" query:"search"`
	Page           int64  `json:"page" query:"page"`
	Limit          int64  `json:"limit" query:"limit"`
	WithPagination bool   `json:"with_pagination" query:"with_pagination"`
	ShowDeleted    bool   `json:"show_deleted" query:"show_deleted"`
	SortBy         string `json:"sort_by" query:"sort_by"`
	OrderBy        string `json:"order_by" query:"order_by"`
	UserRole       string `json:"user_role" query:"user_role"`
}

type MetaResponse struct {
	TotalRecords    int64 `json:"total_records"`
	FilteredRecords int64 `json:"filtered_records"`
	Page            int64 `json:"page"`
	PerPage         int64 `json:"per_page"`
	TotalPages      int64 `json:"total_pages"`
}

type JwtCustomClaims struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.StandardClaims
}

type JwtCustomRefreshClaims struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.StandardClaims
}

type JsonResponse struct {
	Data     interface{} `json:"data,omitempty"`
	Resource string      `json:"resource,omitempty"`
	Meta     interface{} `json:"meta,omitempty"`
	Message  string      `json:"message"`
	Success  bool        `json:"success"`
	Detail   interface{} `json:"detail,omitempty"`
}

type TokenData struct {
	Name      string    `json:"name"`
	UserID    uuid.UUID `json:"user_id"`
	Role      string    `json:"role"`
	TokenType string    `json:"token_type"`
	Token     string    `json:"token"`
}

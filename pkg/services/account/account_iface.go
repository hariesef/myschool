package account

import "context"

type AccountServiceIface interface {
	Create(context.Context, string, string) error
	Login(context.Context, string, string) (*TokenInfo, error)
	Logout(context.Context, string) error
}

type TokenInfo struct {
	Token  string
	Expiry int
}

package services

import (
    "golang.org/x/oauth2"
)

type OAuthUserInfo struct {
    Email     string
    FirstName string
    LastName  string
    ID        string
    Provider  string
}

type OAuthProvider interface {
    GetConfig() *oauth2.Config
    GetUserInfo(token *oauth2.Token) (*OAuthUserInfo, error)
    GetProviderName() string
}

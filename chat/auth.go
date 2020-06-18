package main

import (
	"crypto/md5"
	"fmt"
	"github.com/stretchr/gomniauth"
	gomniauthcommon "github.com/stretchr/gomniauth/common"
	"github.com/stretchr/objx"
	"io"
	"net/http"
	"strings"
)

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}

type authHandler struct {
	next http.Handler
}

func (auth *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//cek apakah cookie menyimpan data user login
	cookie, err := r.Cookie("auth")
	if err == http.ErrNoCookie || cookie.Value == "" {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	auth.next.ServeHTTP(w,r)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandle(lenroom int) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		segs := strings.Split(r.URL.Path, "/")
		if len(segs) < 4 {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Parameter doesn't exist, length = %d", len(segs))
			return
		}
		action := segs[2]
		provider := segs[3]
		switch action {
		case "login" :
			provider, err := gomniauth.Provider(provider)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s",provider, err), http.StatusBadRequest)
				return
			}
			loginUrl, err := provider.GetBeginAuthURL(nil, nil)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s:%s", provider, err), http. StatusInternalServerError)
				return
			}
			w.Header().Set("Location", loginUrl)
			w.WriteHeader(http.StatusTemporaryRedirect)
		case "callback" :
			provider, err := gomniauth.Provider(provider)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s",provider, err), http.StatusBadRequest)
				return
			}
			//fmt.Println(r.URL.RawQuery)
			//parse RawQuery dari request dan dimasukkan ke CompleteAuth
			//Method CompleteAuth menggunakan nilai tsb untuk menyelesaikan handshake oauth2
			creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
			if err != nil {
				http.Error(w, fmt.Sprintf("Error when trying to complete auth for%s: %s", provider, err), http.StatusInternalServerError)
				return
			}
			//jika credential sudah didapatkan, maka data user bisa diakses
			user, err := provider.GetUser(creds)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error when trying to get user from %s: %s", provider, err), http.StatusInternalServerError)
				return
			}
			//log.Println(user)

			//dapetin cookie value
			//field Name dari data user diubah ke encode-base64 dan disimpan dalam object json
			chatUser := &chatUser{User: user}
			m := md5.New()
			fmt.Println(user.Email(), lenroom)
			//fmt.Println(len())
			io.WriteString(m, strings.ToLower(user.Email()))
			chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
			avatarURL, err := avatars.GetAvatarURL(chatUser)
			authCookieValue := objx.New(map[string]interface{}{
				"userid": chatUser.uniqueID,
				"name": user.Name(),
				"avatar_url": avatarURL,
				"email": user.Email(),
			}).MustBase64()



			//authCookieValue := objx.New(map[string]interface{}{
			//	"name": user.Name(),
			//	"avatar_url": user.AvatarURL(),
			//}).MustBase64()


			//set cookie dengan key "auth"
			http.SetCookie(w, &http.Cookie{
				Name: "auth",
				Value: authCookieValue,
				Path: "/"})
			w.Header().Set("Location", "/chat")
			w.WriteHeader(http.StatusTemporaryRedirect)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Auth action %s not supported", action)
		}
	}
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	if len(segs) < 4 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Parameter doesn't exist, length = %d", len(segs))
		return
	}
	action := segs[2]
	provider := segs[3]
	switch action {
		case "login" :
			provider, err := gomniauth.Provider(provider)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s",provider, err), http.StatusBadRequest)
				return
			}
			loginUrl, err := provider.GetBeginAuthURL(nil, nil)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s:%s", provider, err), http. StatusInternalServerError)
				return
			}
			w.Header().Set("Location", loginUrl)
			w.WriteHeader(http.StatusTemporaryRedirect)
		case "callback" :
			provider, err := gomniauth.Provider(provider)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s",provider, err), http.StatusBadRequest)
				return
			}
			//fmt.Println(r.URL.RawQuery)
			//parse RawQuery dari request dan dimasukkan ke CompleteAuth
			//Method CompleteAuth menggunakan nilai tsb untuk menyelesaikan handshake oauth2
			creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
			if err != nil {
				http.Error(w, fmt.Sprintf("Error when trying to complete auth for%s: %s", provider, err), http.StatusInternalServerError)
				return
			}
			//jika credential sudah didapatkan, maka data user bisa diakses
			user, err := provider.GetUser(creds)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error when trying to get user from %s: %s", provider, err), http.StatusInternalServerError)
				return
			}
			//log.Println(user)

			//dapetin cookie value
			//field Name dari data user diubah ke encode-base64 dan disimpan dalam object json
			chatUser := &chatUser{User: user}
			m := md5.New()
			fmt.Println(user.Email())
			//fmt.Println(len())
			io.WriteString(m, strings.ToLower(user.Email()))
			chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
			avatarURL, err := avatars.GetAvatarURL(chatUser)
			authCookieValue := objx.New(map[string]interface{}{
				"userid": chatUser.uniqueID,
				"name": user.Name(),
				"avatar_url": avatarURL,
				"email": user.Email(),
			}).MustBase64()



			//authCookieValue := objx.New(map[string]interface{}{
			//	"name": user.Name(),
			//	"avatar_url": user.AvatarURL(),
			//}).MustBase64()


			//set cookie dengan key "auth"
			http.SetCookie(w, &http.Cookie{
				Name: "auth",
				Value: authCookieValue,
				Path: "/"})
			w.Header().Set("Location", "/chat")
			w.WriteHeader(http.StatusTemporaryRedirect)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
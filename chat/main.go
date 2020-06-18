package main

import (
	"BlueprintChatApp/trace"
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar}

// template struct represent a single template
type templateHandler struct {
	once     sync.Once //guarantee the function can only be executed once
	filename string
	templ    *template.Template
}

// servehttp handles the http request
func (t *templateHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		// template/t.filename
		// template/chat.html
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
	// w?
	// r *http.Request mengandung data port
	//t.templ.Execute(w, r)
}

func main() {
	// 8080 adalah nilai default port,
	var addr = flag.String("addr", //flag name
		":8080",  //default value
		"The addr of the application")

	var traceswitch = flag.String("traceswitch",
		"off",
		"Tracer Switch")
	flag.Parse()

	gomniauth.SetSecurityKey("cobacoba")
	gomniauth.WithProviders(
		//auth dengan akun github
		github.New("f83f0b8270791ae58691",
			"a3cfa76f66c6370e4118336e2cb6ad572e857b44",
			"http://localhost:8080/auth/callback/github"),
		facebook.New("678574552923450",
			"6e66b55db7e7622ce9842c9f6b58affd",
			"http://localhost:8080/auth/callback/facebook"),
		)

	//UseGravatar
	newroom := newRoom()
	if *traceswitch != "off" {
		newroom.tracer = trace.New(os.Stdout)
	}

	// HOMEPAGE
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/chat", 302)
	})
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"})) //pass templateHandler object address
	http.Handle("/room", newroom)
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandle(len(newroom.clients)))
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	//logout
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: "",
			Path: "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	//serve images
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))

	//run room
	go newroom.run() //room running di thread yang lain,
	//agar main goroutine menjalankan web server
	// INITIALIZE WEBSERVER
	log.Println("Starting web server on", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}




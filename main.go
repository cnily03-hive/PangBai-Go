package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"text/template"
)

type Token struct {
	Stringer
	Name string
}

type Config struct {
	Stringer
	Name          string
	JwtKey        string
	SignaturePath string
}

type Helper struct {
	Stringer
	User   string
	Config Config
}

var config = Config{
	Name:          "PangBai 过家家 (4)",
	JwtKey:        RandString(64),
	SignaturePath: "./sign.txt",
}

func (c Helper) Curl(url string) string {
	fmt.Println("Curl:", url)
	cmd := exec.Command("curl", "-fsSL", "--", url)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error: curl:", err)
		return "error"
	}
	return "ok"
}

func routeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/index.html")
}

func routeEye(w http.ResponseWriter, r *http.Request) {

	input := r.URL.Query().Get("input")
	if input == "" {
		input = "{{ .User }}"
	}

	// get template
	content, err := ioutil.ReadFile("views/eye.html")
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	tmplStr := strings.Replace(string(content), "%s", input, -1)
	tmpl, err := template.New("eye").Parse(tmplStr)
	if err != nil {
		input := "[error]"
		tmplStr = strings.Replace(string(content), "%s", input, -1)
		tmpl, err = template.New("eye").Parse(tmplStr)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
	}

	// get user from cookie
	user := "PangBai"
	token, err := r.Cookie("token")
	if err != nil {
		token = &http.Cookie{Name: "token", Value: ""}
	}
	o, err := validateJwt(token.Value)
	if err == nil {
		user = o.Name
	}

	// renew token
	newToken, err := genJwt(Token{Name: user})
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: newToken,
	})

	// render template
	helper := Helper{User: user, Config: config}
	err = tmpl.Execute(w, helper)
	if err != nil {
		http.Error(w, "[error]", http.StatusInternalServerError)
		return
	}
}

func routeFavorite(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPut {

		// ensure only localhost can access
		requestIP := r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
		fmt.Println("Request IP:", requestIP)
		if requestIP != "127.0.0.1" && requestIP != "[::1]" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Only localhost can access"))
			return
		}

		token, _ := r.Cookie("token")

		o, err := validateJwt(token.Value)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		if o.Name == "PangBai" {
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("Hello, PangBai!"))
			return
		}

		if o.Name != "Papa" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("You cannot access!"))
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
		}
		config.SignaturePath = string(body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}

	// render

	tmpl, err := template.ParseFiles("views/favorite.html")
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	sig, err := ioutil.ReadFile(config.SignaturePath)
	if err != nil {
		http.Error(w, "Failed to read signature files: "+config.SignaturePath, http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, string(sig))

	if err != nil {
		http.Error(w, "[error]", http.StatusInternalServerError)
		return
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", routeIndex)
	r.HandleFunc("/eye", routeEye)
	r.HandleFunc("/favorite", routeFavorite)
	r.PathPrefix("/assets").Handler(http.StripPrefix("/assets", noDirList(http.FileServer(http.Dir("./assets")))))

	fmt.Println("Starting server on :8000")
	http.ListenAndServe(":8000", r)
}

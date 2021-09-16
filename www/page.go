package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"log"
	"net/http"
	"regexp"

	"entry_task/rpcclient"
)


var templates = template.Must(template.ParseFiles("edit.html", "view.html", "index.html"))
var validPath = regexp.MustCompile("^/(edit|upload|view|showpic)/([a-zA-Z0-9]+)$")

type profile struct {
	Username string
	Nickname string
	Picture []byte
	PicName string
}

func (p *profile) save() error {
	filename := p.Username + ".jpg"
	return ioutil.WriteFile(filename, p.Picture, 0600)
}

func loadProfile(username string) (*profile, error) {
	// filename := username + ".jpg"
	// picture, _ := ioutil.ReadFile(filename)

	_, u, n, p, err := rpcclient.GetAuth(username, "123456")
	if err != nil {
		fmt.Println("loadProfile failed, ", err)
	}
	fmt.Println("load profile")
	return &profile{Username: u, Nickname: n, PicName: p}, nil
}


func indexHandler(w http.ResponseWriter, r *http.Request)  {
    t, _ := template.ParseFiles("index.html")
    t.Execute(w, nil)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("auth handler")
	username := r.FormValue("username")
	password := r.FormValue("password")

	res, u, n, p, err := rpcclient.GetAuth(username, password)
	if res != true {
		http.Error(w, "User not exit!", http.StatusForbidden)
		return
	}
	fmt.Println("rediret to view page after auth ", res, u, n, p, err)

	pfile := &profile{Username: u, Nickname: n, PicName: p}
	renderTemplate(w, "view", pfile)
}

func viewHandler(w http.ResponseWriter, r *http.Request, username string) {
	p, err := loadProfile(username)

	if err != nil {
		http.Redirect(w, r, "/index", http.StatusFound)
		fmt.Println("in view page, profile empty")
		return
	}

	fmt.Println("render view page")
	renderTemplate(w, "view", p)
}

func showPicHandle( w http.ResponseWriter, r *http.Request) {
	filename := "./imgs/" + r.URL.Path[len("/showpic/"):];
    file, err := os.Open(filename)
	fmt.Println("show picture handler ", r.URL.Path, " ", filename)
    errorHandle(err, w);
 
    defer file.Close()
    buff, err := ioutil.ReadAll(file)
    errorHandle(err, w);
    w.Write(buff)
}


func uploadHandle(w http.ResponseWriter, r *http.Request, username string) {
    w.Header().Set("Content-Type", "text/html")
 
    r.ParseForm()
    if r.Method != "POST" {
        http.Redirect(w, r, "/view/" + username, http.StatusOK);
    } else {
        uploadFile, handle, err := r.FormFile("image")
        errorHandle(err, w)
 
        ext := strings.ToLower(path.Ext(handle.Filename))
        if ext != ".jpg" && ext != ".jpeg" &&ext != ".png" {
            errorHandle(errors.New("只支持jpg/png图片上传"), w);
            return
        }
 
        saveFile, err := os.OpenFile("./imgs/" + handle.Filename, os.O_WRONLY|os.O_CREATE, 0666);
        errorHandle(err, w)
        io.Copy(saveFile, uploadFile);
		defer uploadFile.Close()
        defer saveFile.Close()

		u, n, p, err := rpcclient.UpdatePicture(username, handle.Filename)
		if err != nil {
			fmt.Println("uploadHandle ", err)
			return
		}

		pf := &profile{Username: u, Nickname: n, PicName: p}

        renderTemplate(w, "view", pf)
    }
}


func updateNicknameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[len("/updatenickname/"):]
	new_nickname := r.FormValue("nickname")
	u, n, p, err := rpcclient.UpdateNickname(username, new_nickname)
	if err != nil {
		fmt.Println("updateNicknameHandler ", err)
		http.Error(w, "Update Nickname failed", http.StatusInternalServerError)
		return
	}

	pf := &profile{Username: u, Nickname: n, PicName: p}
	fmt.Println("updateNicknameHandler ", pf)
	renderTemplate(w, "view", pf)
}


func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if  err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r * http.Request) {
		m := validPath.FindStringSubmatch((r.URL.Path))
		fmt.Println("Url ", r.URL.Path, "not found")
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}


func errorHandle(err error, w http.ResponseWriter) {
    if  err != nil {
        w.Write([]byte(err.Error()))
    }
}


func main() {
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/auth", authHandler)

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/updatenickname/", updateNicknameHandler)
	http.HandleFunc("/upload/", makeHandler(uploadHandle))
	http.HandleFunc("/showpic/", showPicHandle)

	log.Fatal(http.ListenAndServe(":5000", nil))
}
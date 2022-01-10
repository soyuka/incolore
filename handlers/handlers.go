package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
    "strconv"
	"log"
	"net/http"
	"strings"
	"path/filepath"
	"io/ioutil"
	"io"
	"os"
	"bytes"
	"crypto/sha256"

	"github.com/matoous/go-nanoid"
	"github.com/h2non/filetype"
)

func Exists(name string) (bool, error) {
    _, err := os.Stat(name)
    if err == nil {
        return true, nil
    }
    if errors.Is(err, os.ErrNotExist) {
        return false, nil
    }
    return false, err
}

const cookieName = "incolore"

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func makeStatusError(code int) StatusError {
	return StatusError{
		Code: code,
		Err:  errors.New(http.StatusText(code)),
	}
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}

type Handler struct {
	Env *Env
	Handler func(e *Env, w http.ResponseWriter, r *http.Request) error
}

// Adapted from https://blog.questionable.services/article/http-handler-error-handling-revisited/
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.Handler(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

/// Single endpoint /
/// When there's a query we're using it's value to save the link in the database
/// If there's a code (eg: hostname.com/EnYQkRXzK30d) we redirect to the given value
func GetIndex(env *Env, w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		cookie := &http.Cookie{Name: cookieName, SameSite: http.SameSiteStrictMode, Secure: true, HttpOnly: true}
		http.SetCookie(w, cookie)
		return CreateLink(env, w, r)
	}

	key := strings.Replace(r.URL.Path, "/", "", 1)

	if key != "" {
		_, err := r.Cookie(cookieName)
		source, err := env.Transport.Get(key)

		if source == "" {
			return makeStatusError(http.StatusNotFound)
		}

		if err == nil {
			cookie := &http.Cookie{Name: cookieName, MaxAge: -1, SameSite: http.SameSiteStrictMode, Secure: true, HttpOnly: true}
			http.SetCookie(w, cookie)
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		buf, err := ioutil.ReadFile(source)

		if err != nil {
			return makeStatusError(http.StatusNotFound)
		}

		kind, _ := filetype.Match(buf)

		w.Header().Add("Content-Type", kind.MIME.Value) 
		w.Write(buf)
		return nil
	}

	return Index(env, w, r)
}

// GET ?http://link creates the link and redirect to the link
func CreateLink(env *Env, w http.ResponseWriter, r *http.Request) error {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("f")

	if err != nil {
		log.Println(err)
		return makeStatusError(http.StatusBadRequest)
	}

	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Println(err)
		return makeStatusError(http.StatusBadRequest)
	}

	hash := sha256.Sum256(buf.Bytes())
	hashStr := string(hash[:])
	existingId, _ := env.Transport.Get(hashStr)
	if existingId != "" {
		http.Redirect(w, r, fmt.Sprintf("%s/%s", env.Config.ShortenerHostname, existingId), 302)
		return nil
	}

	if !filetype.IsImage(buf.Bytes()) {
		return makeStatusError(http.StatusUnsupportedMediaType)
	}

	kind, _ := filetype.Match(buf.Bytes())
	if kind == filetype.Unknown {
		log.Println("File type is unknown")
		return makeStatusError(http.StatusBadRequest)
	}

	id, err := gonanoid.Generate(env.Config.IdAlphabet, env.Config.IdLength)
	if err != nil {
		return StatusError{http.StatusInternalServerError, err}
	}

	size := int64(buf.Len())
	if size <= 0 || size > env.Config.MaxSize {
		return makeStatusError(http.StatusRequestEntityTooLarge)
	}
	
	destination := filepath.Join(env.Config.Directory, handler.Filename)
	exists, _ := Exists(destination)
	if exists {
		destination = filepath.Join(env.Config.Directory, id + "-" + handler.Filename)
	}

    err = ioutil.WriteFile(destination, buf.Bytes(), 0644)
	if err != nil {
		log.Println(err)
		return StatusError{http.StatusInternalServerError, err}
	}

	err = env.Transport.Put(hashStr, id)
	if err != nil {
		return StatusError{http.StatusInternalServerError, err}
	}

	err = env.Transport.Put(id, destination)
	if err != nil {
		return StatusError{http.StatusInternalServerError, err}
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", env.Config.ShortenerHostname, id), 302)
	return nil
}

/// Favicon just for fun
func Favicon(env *Env, w http.ResponseWriter, r *http.Request) error {
	decoded, _ := base64.StdEncoding.DecodeString("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")
	w.Header().Set("Content-Type", "image/x-icon")
	w.Write(decoded)
	return nil
}

/// Favicon just for fun
func Index(env *Env, w http.ResponseWriter, r *http.Request) error {
	count, _ := env.Transport.Count()

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	index := `<!DOCTYPE html>
<html lang="en">
<head>
  <title>Incolore 🎨</title>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <meta name="description" content="Image and picture host" />
  <style>
	body {margin: 5% auto; background: #f2f2f2; color: #444444; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; font-size: 16px; line-height: 1.8; text-shadow: 0 1px 0 #ffffff; max-width: 73%;}
	code {background: white;}
	a {border-bottom: 1px solid #444444; color: #444444; text-decoration: none;}
	a:hover {border-bottom: 0;}
  </style>
</head>
<body>
  <h1>Incolore 🎨</h1>
  <h2>Upload an image</h2>
  <form enctype="multipart/form-data" method="POST" action="/">
	<input type="file" name="f" autofocus />
	<input type="submit" value="Upload"/>
	<p><small>Data has no warranty and can be removed at any time.</small></p>
  </form>
  <h2>API</h2>
  <p>POST <code>`+env.Config.ShortenerHostname+`</code> with multipart/form-data with f.</p>
  <p><a href="https://github.com/soyuka/incolore">Code on github</a></p>
  <h2>Statistics</h2>
  <p>`+strconv.FormatInt(count, 10)+` images online</p>
  <script type="text/javascript">
	var fileUpload = document.querySelector('input[type="file"]')
	var submit = document.querySelector('input[type="submit"]')
	fileUpload.addEventListener('change', function(event) {
		submit.click()
	})

	document.onpaste = function(event){
		var items = (event.clipboardData || event.originalEvent.clipboardData).items;
		for (index in items) {
			var item = items[index];
			if (item.kind === 'file') {
				var file = item.getAsFile();
				var container = new DataTransfer();
				container.items.add(file);
				fileUpload.files = container.files;
				submit.click()
				return;
			}
		}
	}
  </script>
</body>
</html>`
	w.Write([]byte(index))
	return nil
}

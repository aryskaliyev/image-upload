package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"lincoln.boris/forum/pkg/forms"
	"lincoln.boris/forum/pkg/models"

	"github.com/gofrs/uuid/v5"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	p, err := app.posts.AllPosts(0, 0)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	c, err := app.categories.All()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	user_id := 0
	cookies := r.Cookies()
	for _, c := range cookies {
		if c.Name == "session_cookie" {
			token := c.Value
			user_id = app.sessions.GetUser(token)
			break
		}
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Posts: p,
		Categories: c,
		AuthenticatedUser: user_id,
	})
}

func (app *application) showPost(w http.ResponseWriter, r *http.Request) {
	post_id, err := strconv.Atoi(r.URL.Query().Get(":post_id"))
	if err != nil || post_id < 1 {
		app.notFound(w, r)
		return
	}

	p, err := app.posts.Get(post_id)
	if err == models.ErrNoRecord {
		app.notFound(w, r)
		return
	} else if err != nil {
		app.serverError(w, r, err)
		return
	}


	pc, err := app.post_categories.Get(post_id)
	if err == models.ErrNoRecord {
		app.notFound(w, r)
		return
	} else if err != nil {
		app.serverError(w, r, err)
		return
	}

	c, err := app.comments.Get(post_id)
	if err == models.ErrNoRecord {
		app.notFound(w, r)
		return
	} else if err != nil {
		app.serverError(w, r, err)
		return
	}

	for i := range c {
		cv, err := app.comment_votes.Sum(c[i].ID)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		c[i].Votes = cv
	}

	pv, err := app.post_votes.Sum(post_id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	p.Votes = pv


	app.render(w, r, "show.page.tmpl", &templateData{
		Post: p,
		PostCategories: pc,
		Comments: c,
	})

}

func (app *application) createPostForm(w http.ResponseWriter, r *http.Request) {
	c, err := app.categories.All()
	if err == models.ErrNoRecord {
		app.notFound(w, r)
		return
	} else if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.render(w, r, "create.page.tmpl", &templateData{
		Categories: c,
		Form: forms.New(nil),
	})
}

func getFileExtension(fileHeader *multipart.FileHeader) string {
	parts := strings.Split(fileHeader.Filename, ".")
	extension := parts[len(parts)-1]
	return extension
}

func generateImageID(file multipart.File, ext string) (string, error) {
	imageUUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", imageUUID.String(), ext), nil
}

func saveImageToFile(file multipart.File, destination string) error {
	if err := os.MkdirAll(filepath.Dir(destination), os.ModePerm); err != nil {
		return err
	}

	newFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(30 << 20)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	user_id := app.authenticatedUser(r)

	form := forms.New(r.PostForm)
	form.Required("title", "body", "category")
	form.MaxLength("title", 75)
	form.MaxLength("body", 500)

	var categories []int
	for _, category := range r.PostForm["category"] {
		c, err := strconv.Atoi(string(category))
		if err != nil {
			app.clientError(w, r, http.StatusBadRequest)
			return
		}
		categories = append(categories, c)
	}

	file, handler, err := r.FormFile("image")

	var imageID string
	var imageSizeInBytes int64	    
	var fileTooLargeError string

	if err == nil {
		defer file.Close()

		imageSizeInBytes = handler.Size
		maxAllowedSize := int64(20 * 1024 * 1024)
		if imageSizeInBytes > maxAllowedSize {
			fileTooLargeError = "Image size exceeds the allowed limit (20MB)"
		} else {
			extension := getFileExtension(handler)
			imageID, err = generateImageID(file, extension)
			if err != nil {
                        	app.serverError(w, r, err)
                        	return
                	}

			err = saveImageToFile(file, filepath.Join("./uploads", imageID))
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		}
	}

	if !form.Valid() || fileTooLargeError != "" {
		c, err := app.categories.All()
		if err == models.ErrNoRecord {
			app.notFound(w, r)
			return
		} else if err != nil {
			app.serverError(w, r, err)
			return
		}

		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
			Categories: c,
			FileTooLargeErr: fileTooLargeError,
		})
		return
	}

	id, err := app.posts.Insert(user_id, form.Get("title"), form.Get("body"), categories, imageID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts/%d", id), http.StatusSeeOther)
}

func (app *application) showCategoryPosts (w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":category_id"))
	if err != nil || id < 1 {
		app.notFound(w, r)
		return
	}

	p, err := app.posts.ByCategory(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	c, err := app.categories.Get(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Posts: p,
		Category: c,
	})
}

func (app *application) createComment(w http.ResponseWriter, r *http.Request) {
	post_id, err := strconv.Atoi(r.URL.Query().Get(":post_id"))
	if err != nil || post_id < 1 {
		app.notFound(w, r)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	user_id := app.authenticatedUser(r)
	body := r.PostForm.Get("body")

	form := forms.New(r.PostForm)
	form.Required("body")

	if !form.Valid() {
		p, err := app.posts.Get(post_id)
		if err == models.ErrNoRecord {
			app.notFound(w, r)
			return
		} else if err != nil {
			app.serverError(w, r, err)
			return
		}


		pc, err := app.post_categories.Get(post_id)
		if err == models.ErrNoRecord {
			app.notFound(w, r)
			return
		} else if err != nil {
			app.serverError(w, r, err)
			return
		}

		c, err := app.comments.Get(post_id)
		if err == models.ErrNoRecord {
			app.notFound(w, r)
			return
		} else if err != nil {
			app.serverError(w, r, err)
			return
		}

		for i := range c {
			cv, err := app.comment_votes.Sum(c[i].ID)
			if err != nil {
				app.serverError(w, r, err)
				return
			}

			c[i].Votes = cv
		}

		pv, err := app.post_votes.Sum(post_id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		p.Votes = pv

		app.render(w, r, "show.page.tmpl", &templateData{
			Post: p,
			Comments: c,
			PostCategories: pc,
			Form: form,
		})
		return
	}

	_, err = app.comments.Insert(body, post_id, user_id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts/%d", post_id), http.StatusSeeOther)
}

func (app *application) votePost(w http.ResponseWriter, r *http.Request) {
	user_id := app.authenticatedUser(r)

	post_id, err := strconv.Atoi(r.URL.Query().Get(":post_id"))
	if err != nil || post_id < 1 {
		app.notFound(w, r)
		return
	}

	vote := r.URL.Query().Get(":vote")
	if !(vote == "up" || vote == "down") {
		app.notFound(w, r)
		return
	}

	vm := map[string]int {
		"up": 1,
		"down": -1,
	}

	pv, err := app.post_votes.Get(post_id, user_id)
	if err == models.ErrNoRecord {
		_, err = app.post_votes.Insert(post_id, user_id, vm[vote])
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	} else if pv.Vote == vm[vote] {
		_, err = app.post_votes.Delete(post_id, user_id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	} else if pv.Vote != vm[vote] {
		_, err = app.post_votes.Delete(post_id, user_id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		_, err = app.post_votes.Insert(post_id, user_id, vm[vote])
		if err != nil {
			app.serverError(w, r, err)
			return
		}

	} else if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts/%d", post_id), http.StatusSeeOther)
}

func (app *application) voteComment(w http.ResponseWriter, r *http.Request) {
	user_id := app.authenticatedUser(r)

	post_id, err := strconv.Atoi(r.URL.Query().Get(":post_id"))
	if err != nil || post_id < 1 {
		app.notFound(w, r)
		return
	}
	comment_id, err := strconv.Atoi(r.URL.Query().Get(":comment_id"))
	if err != nil || comment_id < 1 {
		app.notFound(w, r)
		return
	}

	vote := r.URL.Query().Get(":vote")
	if !(vote == "up" || vote == "down") {
		app.notFound(w, r)
		return
	}

	vm := map[string]int {
		"up": 1,
		"down": -1,
	}

	cv, err := app.comment_votes.Get(comment_id, user_id)
	if err == models.ErrNoRecord {
		_, err = app.comment_votes.Insert(comment_id, user_id, vm[vote])
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	} else if cv.Vote == vm[vote] {
		_, err = app.comment_votes.Delete(comment_id, user_id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	} else if cv.Vote != vm[vote] {
		_, err = app.comment_votes.Delete(comment_id, user_id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		_, err = app.comment_votes.Insert(comment_id, user_id, vm[vote])
		if err != nil {
			app.serverError(w, r, err)
			return
		}

	} else if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts/%d", post_id), http.StatusSeeOther)
}

func (app *application) showUserCreatedPosts(w http.ResponseWriter, r *http.Request) {
	user_id := app.authenticatedUser(r)

	p, err := app.posts.AllPosts(user_id, 0)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Posts: p,
	})
}

func (app *application) showUserUpvotedPosts(w http.ResponseWriter, r *http.Request) {
	user_id := app.authenticatedUser(r)

	p, err := app.posts.AllPosts(user_id, 1)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Posts: p,
	})
}

func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("username", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	err = app.users.Insert(form.Get("username"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}



func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	user_id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, r, err)
		return
	}

	token := app.sessions.GetToken(user_id)
	if token != "" {
		err = app.sessions.Delete(user_id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	u, err := uuid.NewV4()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	name := "session_cookie"
	token = u.String()
	expires := time.Now().Add(5 * time.Minute)

	http.SetCookie(w, &http.Cookie{
		Name: name,
		Value: token,
		Expires: expires,
		HttpOnly: true,
		Path: "/",
	})

	err = app.sessions.Insert(user_id, token, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*
var (
    githubOauthConfig = &oauth2.Config{
        ClientID:     "8674939c0230d52d5fdf",
        ClientSecret: "9f9699969a9bd4fda311112e3bfa2bbaa4d37aa7",
        Endpoint:     github.Endpoint,
    }
)

func (app *application) githubLogin(w http.ResponseWriter, r *http.Request) {
	authURL := githubOauthConfig.AuthCodeURL("state")
	http.Redirect(w, r, authURL, http.StatusSeeOther)
}

func (app *application) githubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, err := githubOauthConfig.Exchange(r.Context(), code)
  	  if err != nil {
        	http.Error(w, "Authentication failed", http.StatusInternalServerError)
     	   return
    	}
}
*/
func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	var token string
	cookies := r.Cookies()

	if len(cookies) > 0 {
		for _, c := range cookies {
			if c.Name == "session_cookie" {
				token = c.Value
				break
			}
		}
	}

	user_id := app.sessions.GetUser(token)

	err := app.sessions.Delete(user_id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", 303)
}

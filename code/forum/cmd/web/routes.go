package main

import (
	"net/http"

	"lincoln.boris/forum/pkg/alice"
	"lincoln.boris/forum/pkg/pat"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/posts/create", app.requireAuthenticatedUser(http.HandlerFunc(app.createPostForm)))
	mux.Post("/posts/create", app.requireAuthenticatedUser(http.HandlerFunc(app.createPost)))
	mux.Get("/posts/:post_id", http.HandlerFunc(app.showPost))
	mux.Get("/posts/:post_id/vote/:vote", app.requireAuthenticatedUser(http.HandlerFunc(app.votePost)))
	mux.Get("/posts/:post_id/comments/:comment_id/vote/:vote", app.requireAuthenticatedUser(http.HandlerFunc(app.voteComment)))

	mux.Get("/categories/:category_id", http.HandlerFunc(app.showCategoryPosts))
	mux.Post("/comments/create/:post_id", app.requireAuthenticatedUser(http.HandlerFunc(app.createComment)))

	mux.Get("/user/posts", app.requireAuthenticatedUser(http.HandlerFunc(app.showUserCreatedPosts)))
	mux.Get("/user/upvotes", app.requireAuthenticatedUser(http.HandlerFunc(app.showUserUpvotedPosts)))

	mux.Get("/user/signup", http.HandlerFunc(app.signupForm))
	mux.Post("/user/signup", http.HandlerFunc(app.signup))
	mux.Get("/user/login", http.HandlerFunc(app.loginForm))
	mux.Post("/user/login", http.HandlerFunc(app.login))
	mux.Post("/user/logout", app.requireAuthenticatedUser(http.HandlerFunc(app.logout)))

	mux.Get("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static/"))))
    	mux.Get("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	return standardMiddleware.Then(mux)
}

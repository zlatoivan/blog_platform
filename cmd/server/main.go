package main

import (
	"blog-platform/internal/api"
	"blog-platform/internal/dbPkg"
	"github.com/go-chi/chi/v5/middleware"

	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	fmt.Print("\nStart\n\n")

	db := dbPkg.DB{}
	db.InitDB()
	db.CreateTables()
	//db.DropTables()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//r.Use(middleware.RequestID)
	//r.Use(middleware.Recoverer)
	//r.Use(middleware.URLFormat)

	a := api.Api{
		Db: &db,
	}

	r.Get("/register", a.GetRegister)
	r.Post("/register", a.PostRegister)

	r.Get("/login", a.GetLogin)
	r.Post("/login", a.PostLogin)

	r.Group(func(r chi.Router) {
		r.Use(a.CheckAuth) // Если пользователь не залогинен, то редирект на страницу логина
		// А есди залогинен, то выз. след. обработчик, в котором лежит объект пользователя

		r.Get("/", a.GetSubsArticles)

		r.Route("/bloggers", func(r chi.Router) {
			r.Get("/", a.GetBloggers)
			r.Route("/{bloggerId}", func(r chi.Router) {
				//r.Use(a.BloggerViewCtx)
				r.Get("/", a.GetBlogger)
				r.Post("/", a.PostBlogger)
				r.Get("/subscribed", a.GetBloggerSubscribed)
				r.Post("/subscribed", a.PostBloggerSubscribed)
			})
		})

		r.Get("/logout", a.GetLogout)

		r.Get("/profile", a.GetProfile)

		r.Get("/insert", a.GetInsertArticle)
		r.Post("/insert", a.PostInsertArticle)
	})

	fmt.Println("Listening on http://localhost:3333/...")
	err := http.ListenAndServe(":3333", r)
	log.Fatal(err)
}

//
//r.Get("/bloggers", a.GetBloggers)
//r.Get("/bloggers/{bloggerId}", a.GetBlogger)
//r.Post("/bloggers/{bloggerId}", a.PostBlogger)
//r.Get("/bloggers/{bloggerId}/subscribed", a.GetBloggerSubscribed)
//r.Post("/bloggers/{bloggerId}/subscribed", a.PostBloggerSubscribed)

//r.Route("/articles", func(r chi.Router) {
//	r.Get("/", a.GetInsertArticle)     // +
//	r.Post("/", a.PostInsertArticle) // +
//	r.Route("/{articleID}", func(r chi.Router) {
//		r.Use(ArticleCtx)
//		r.Get("/", GetArticle)       // GET /articles/1234
//		r.Put("/", UpdateArticle)    // PUT /articles/1234
//		r.Delete("/", DeleteArticle) // DELETE /articles/1234
//		r.Get("/edit", EditArticle)  // GET /articles/1234/edit
//	})
//})

//// RESTy routes for "articles" resource
//r.Route("/articles", func(r chi.Router) {
//	r.With(paginate).Get("/", listArticles)                           // GET /articles
//	r.With(paginate).Get("/{month}-{day}-{year}", listArticlesByDate) // GET /articles/01-16-2017
//
//	r.Article("/", createArticle)       // POST /articles
//	r.Get("/search", searchArticles) // GET /articles/search
//
//	// Regexp url parameters:
//	r.Get("/{articleSlug:[a-z-]+}", getArticleBySlug) // GET /articles/home-is-toronto
//
//	// Subrouters:
//	r.Route("/{articleID}", func(r chi.Router) {
//		r.Use(ArticleCtx)
//		r.Get("/", getArticle)       // GET /articles/123
//		r.Put("/", updateArticle)    // PUT /articles/123
//		r.Delete("/", deleteArticle) // DELETE /articles/123
//	})
//})
//
//// Mount the admin sub-router
//r.Mount("/admin", adminRouter())

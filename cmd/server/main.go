package main

import (
	"blog-platform/internal/api"
	"blog-platform/internal/dbPkg"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FileServer is serving static files. | https://github.com/go-chi/chi/issues/403
func FileServer(router *chi.Mux) {
	root := "./"
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer2(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/internal/static" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

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

	////////////////////////////////////////////////////////////////////
	//FileServer(r)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "/"))
	FileServer2(r, "/", filesDir)

	//_ = mime.AddExtensionType(".js", "text/javascript")

	//fs := http.FileServer(http.Dir("static"))
	//r.Handle("/static/*", http.StripPrefix("/static/", fs))
	////////////////////////////////////////////////////////////////////

	r.Get("/register", a.GetRegister)
	r.Post("/register", a.PostRegister)

	r.Get("/login", a.GetLogin)
	r.Post("/login", a.PostLogin)

	r.Group(func(r chi.Router) {
		r.Use(a.CheckAuth) // Если пользователь не залогинен, то редирект на страницу логина
		// А есди залогинен, то выз. след. обработчик, в котором лежит объект пользователя

		r.Get("/", a.GetSubsArticles)

		//r.Route("", func(r chi.Router) {
		r.Get("/bloggers", a.GetBloggers)
		//r.Use(a.BloggerViewCtx)
		r.Get("/bloggers/{bloggerId}", a.GetBlogger)
		r.Get("/bloggers/{bloggerId}/{articleId}", a.GetBloggerArticle)
		r.Post("/bloggers/{bloggerId}/{articleId}", a.PostBloggerArticle)
		//})

		r.Get("/logout", a.GetLogout)

		r.Get("/profile", a.GetProfile)

		r.Get("/insert", a.GetInsertArticle)
		r.Post("/insert", a.PostInsertArticle)

		r.Post("/someoneIsLiked", a.SomeoneIsLiked)
		r.Post("/showLikes", a.ShowLikes)
		r.Post("/someoneIsSubscribed", a.SomeoneIsSubscribed)
		r.Post("/showSubscriptions", a.ShowSubscriptions)
	})

	fmt.Println("Listening on http://localhost:3333/...")
	err := http.ListenAndServe(":3333", r)
	log.Fatal(err)
}

//

//r.Route("/articles", func(r chi.Router) {
//	r.Get("/", a.GetInsertArticle)     // +
//	r.Post("/", a.PostInsertArticle) // +
//	r.Route("/{articleID}", func(r chi.Router) {
//		r.Use(ArticleCtx)
//		r.Get("/", GetBloggerArticle)       // GET /articles/1234
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

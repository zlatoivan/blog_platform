package api

import (
	"blog-platform/internal/dbPkg"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Server. Смысл - все зависимости, которые нужны в обработчиках
type Api struct {
	Db     *dbPkg.DB
	Logger *log.Logger
}

const headerHtmlPath = "internal/static/templates/header.html"
const footerHtmlPath = "internal/static/templates/footer.html"
const navbarHtmlPath = "internal/static/templates/navbar.html"
const insertArticleHtmlPath = "internal/static/templates/insertArticle.html"
const subsArticlesHtmlPath = "internal/static/templates/subsArticles.html"
const registerHtmlPath = "internal/static/templates/register.html"
const loginHtmlPath = "internal/static/templates/login.html"
const profileHtmlPath = "internal/static/templates/profile.html"
const bloggersHtmlPath = "internal/static/templates/bloggers.html"
const bloggerProfileHtmlPath = "internal/static/templates/bloggerProfile.html"
const bloggerProfileSubscribedHtmlPath = "internal/static/templates/bloggerProfileSubscribed.html"

func (a *Api) GetSubsArticles(w http.ResponseWriter, r *http.Request) {
	bloggerId := getBloggerFromCtx(r.Context()).BloggerId
	subs, _ := a.Db.GetSubsAndNotSubs(bloggerId)

	data := make([]struct {
		Blogger  *dbPkg.Blogger
		Articles []*dbPkg.Article
	}, 0)
	for _, blogger := range subs {
		articles := a.Db.GetArticlesByBloggerId(blogger.BloggerId)
		combo := struct {
			Blogger  *dbPkg.Blogger
			Articles []*dbPkg.Article
		}{
			Blogger:  blogger,
			Articles: articles,
		}
		data = append(data, combo)
	}

	t, err := template.ParseFiles(subsArticlesHtmlPath, headerHtmlPath, footerHtmlPath, navbarHtmlPath)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *Api) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("PassToken")
		if err != nil {
			log.Fatal()
		}
		fmt.Printf("\n[Check auth]: %s\n", cookie.Value)

		blogger := a.Db.GetBloggerByCookie(cookie)

		if blogger == nil { // Не происходит, так как /logout сразу переводит на /login
			http.Redirect(w, r, "/login", http.StatusFound)
			next.ServeHTTP(w, r)
		} else {
			ctx := context.WithValue(r.Context(), "blogger", blogger)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

	})
}

func (a *Api) GetRegister(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(registerHtmlPath, headerHtmlPath, footerHtmlPath, navbarHtmlPath)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *Api) PostRegister(w http.ResponseWriter, r *http.Request) {
	blogger := &dbPkg.Blogger{
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		AuthToken: uuid.New().String(),
		Login:     r.FormValue("login"),
		Name:      r.FormValue("name"),
		Surname:   r.FormValue("surname"),
		Country:   r.FormValue("country"),
	}

	a.Db.InsertBlogger(blogger)

	cookie := http.Cookie{
		Name:  "PassToken",
		Value: blogger.AuthToken,
	}

	http.SetCookie(w, &cookie)

	//w.Header().Add("Set-Cookie", cookie.String())

	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *Api) GetLogin(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(loginHtmlPath, headerHtmlPath, footerHtmlPath, navbarHtmlPath)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *Api) PostLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	blogger := a.Db.GetBloggerByEmailPassword(email, password)
	if blogger != nil {
		cookie := http.Cookie{
			Name:  "PassToken",
			Value: blogger.AuthToken,
		}
		http.SetCookie(w, &cookie)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *Api) GetLogout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:  "PassToken",
		Value: "NoCookie",
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/login", http.StatusFound)
}

func getBloggerFromCtx(ctx context.Context) *dbPkg.Blogger {
	return ctx.Value("blogger").(*dbPkg.Blogger)
}

func (a *Api) GetProfile(w http.ResponseWriter, r *http.Request) {
	blogger := getBloggerFromCtx(r.Context())
	articles := a.Db.GetArticlesByBloggerId(blogger.BloggerId)
	data := struct {
		Blogger  *dbPkg.Blogger
		Articles []*dbPkg.Article
	}{
		Blogger:  blogger,
		Articles: articles,
	}

	t, err := template.ParseFiles(profileHtmlPath, headerHtmlPath, footerHtmlPath, navbarHtmlPath)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
}

func (a *Api) GetInsertArticle(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(insertArticleHtmlPath, headerHtmlPath, footerHtmlPath, navbarHtmlPath)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *Api) PostInsertArticle(w http.ResponseWriter, r *http.Request) {
	blogger := getBloggerFromCtx(r.Context())

	article := &dbPkg.Article{
		BloggerId:      blogger.BloggerId,
		Title:          r.FormValue("title"),
		ArticleMessage: template.HTML(r.FormValue("articleMessage")),
		Date:           r.FormValue("date"),
	}

	a.Db.InsertArticle(article)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *Api) GetBloggers(w http.ResponseWriter, r *http.Request) {
	exceptBloggerId := getBloggerFromCtx(r.Context()).BloggerId
	bloggersSub, bloggersNotSub := a.Db.GetSubsAndNotSubs(exceptBloggerId)

	data := struct {
		BloggersSub    []*dbPkg.Blogger
		BloggersNotSub []*dbPkg.Blogger
	}{
		BloggersSub:    bloggersSub,
		BloggersNotSub: bloggersNotSub,
	}

	t, err := template.ParseFiles(bloggersHtmlPath, headerHtmlPath, footerHtmlPath, navbarHtmlPath)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

//func (a *Api) BloggerViewCtx(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		bloggerIdFromURL := chi.URLParam(r, "bloggerId")
//		bloggerId, _ := strconv.Atoi(bloggerIdFromURL)
//		blogger := a.Db.GetBloggerByBloggerId(bloggerId)
//		fmt.Println(bloggerId, "\n", blogger)
//		ctx := context.WithValue(r.Context(), "bloggerView", blogger)
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}
//
//func getBloggerViewFromCtx(ctx context.Context) *dbPkg.Blogger {
//	return ctx.Value("bloggerView").(*dbPkg.Blogger)
//}

func (a *Api) GetBlogger(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("IN GET Blogger")
	////bloggerId := getBloggerViewFromCtx(r.Context()).BloggerId
	//bloggerId := r.Context().Value("bloggerView").(*dbPkg.Blogger).BloggerId
	//blogger := a.Db.GetBloggerByBloggerId(bloggerId)
	//fmt.Println(bloggerId, "\n", blogger)
	//
	//bloggerId1 := getBloggerFromCtx(r.Context()).BloggerId
	//blogger1 := a.Db.GetBloggerByBloggerId(bloggerId1)
	//fmt.Println(bloggerId1, "\n", blogger1)

	bloggerIdFromURL := chi.URLParam(r, "bloggerId")
	bloggerId, _ := strconv.Atoi(bloggerIdFromURL)
	blogger := a.Db.GetBloggerByBloggerId(bloggerId)

	articles := a.Db.GetArticlesByBloggerId(bloggerId)

	data := struct {
		Blogger  *dbPkg.Blogger
		Articles []*dbPkg.Article
	}{
		Blogger:  blogger,
		Articles: articles,
	}

	fmt.Println("Before template")
	t, err := template.ParseFiles(bloggerProfileHtmlPath, headerHtmlPath, footerHtmlPath, navbarHtmlPath)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
}

func (a *Api) PostBlogger(w http.ResponseWriter, r *http.Request) {
	bloggerId := getBloggerFromCtx(r.Context()).BloggerId

	//bloggerViewId := getBloggerViewFromCtx(r.Context()).BloggerId

	bloggerViewIdFromURL := chi.URLParam(r, "bloggerId")
	bloggerViewId, _ := strconv.Atoi(bloggerViewIdFromURL)

	fmt.Println("TWO ID's:", bloggerId, bloggerViewId)
	a.Db.InsertSubscription(bloggerId, bloggerViewId)

	http.Redirect(w, r, "/bloggers/"+strconv.Itoa(bloggerViewId)+"/subscribed", http.StatusFound)
}

func (a *Api) GetBloggerSubscribed(w http.ResponseWriter, r *http.Request) {
	//bloggerId := getBloggerViewFromCtx(r.Context()).BloggerId

	bloggerIdFromURL := chi.URLParam(r, "bloggerId")
	bloggerId, _ := strconv.Atoi(bloggerIdFromURL)

	blogger := a.Db.GetBloggerByBloggerId(bloggerId)
	articles := a.Db.GetArticlesByBloggerId(bloggerId)

	data := struct {
		Blogger  *dbPkg.Blogger
		Articles []*dbPkg.Article
	}{
		Blogger:  blogger,
		Articles: articles,
	}

	t, err := template.ParseFiles(bloggerProfileSubscribedHtmlPath, headerHtmlPath, footerHtmlPath, navbarHtmlPath)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
}

func (a *Api) PostBloggerSubscribed(w http.ResponseWriter, r *http.Request) {
	bloggerId := getBloggerFromCtx(r.Context()).BloggerId

	//bloggerViewId := getBloggerViewFromCtx(r.Context()).BloggerId

	bloggerViewIdFromURL := chi.URLParam(r, "bloggerId")
	bloggerViewId, _ := strconv.Atoi(bloggerViewIdFromURL)

	a.Db.DeleteSubscription(bloggerId, bloggerViewId)

	http.Redirect(w, r, "/bloggers/"+strconv.Itoa(bloggerViewId), http.StatusFound)
}

//func Logger(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		log.Println(r.URL.Path)
//		next.ServeHTTP(w, r)
//	})
//}

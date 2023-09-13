package api

import (
	"blog-platform/internal/dbPkg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
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
const bloggerArticleHtmlPath = "internal/static/templates/bloggerArticle.html"

func (a *Api) GetSubsArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Print("-----------------------------------\n\n\n\n\n")
	bloggerId := getBloggerFromCtx(r.Context()).BloggerId
	subs, _ := a.Db.GetSubsAndNotSubs(bloggerId)

	articles := make([]*dbPkg.Article, 0)
	for _, blogger := range subs {
		articlesOfBlogger := a.Db.GetArticlesByBloggerId(blogger.BloggerId)
		for _, ar := range articlesOfBlogger {
			articles = append(articles, ar)
		}
	}
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Date > articles[j].Date
	})

	data := make([]struct {
		Blogger     *dbPkg.Blogger
		Article     *dbPkg.Article
		Likes       int
		CommentsCnt int
	}, 0)
	for _, article := range articles {
		blogger := a.Db.GetBloggerByBloggerId(article.BloggerId)
		likes := a.Db.GetLikesCntByArticleId(article.ArticleId)
		commentsCnt := len(a.Db.GetCommentsByArticleId(article.ArticleId))
		articleForm := struct {
			Blogger     *dbPkg.Blogger
			Article     *dbPkg.Article
			Likes       int
			CommentsCnt int
		}{
			Blogger:     blogger,
			Article:     article,
			Likes:       likes,
			CommentsCnt: commentsCnt,
		}
		data = append(data, articleForm)
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
		Date:           time.Now().Format("02.01.2006, 15:04:05"),
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

	articleForms := make([]struct {
		Article *dbPkg.Article
		Likes   int
		//Comments *dbPkg.Comment
	}, 0)
	for _, article := range articles {
		likes := a.Db.GetLikesCntByArticleId(article.ArticleId)
		//comme
		dataElem := struct {
			Article *dbPkg.Article
			Likes   int
			//Comments *dbPkg.Comment
		}{
			Article: article,
			Likes:   likes,
			//Comments: comments,
		}
		articleForms = append(articleForms, dataElem)
	}

	data := struct {
		Blogger      *dbPkg.Blogger
		ArticleForms []struct {
			Article *dbPkg.Article
			Likes   int
			//Comments *dbPkg.Comment
		}
	}{
		Blogger:      blogger,
		ArticleForms: articleForms,
	}

	t, err := template.ParseFiles(bloggerProfileHtmlPath, headerHtmlPath, footerHtmlPath, navbarHtmlPath)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
}

func (a *Api) GetBloggerArticle(w http.ResponseWriter, r *http.Request) {
	articleIdFromURL := chi.URLParam(r, "articleId")
	articleId, _ := strconv.Atoi(articleIdFromURL)
	article := a.Db.GetArticleByArticleId(articleId)
	likesCnt := a.Db.GetLikesCntByArticleId(articleId)
	comments := a.Db.GetCommentsByArticleId(articleId)
	commentsCnt := len(comments)
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].Date > comments[j].Date
	})
	fmt.Println(likesCnt, commentsCnt)

	commentsForms := make([]struct {
		Blogger *dbPkg.Blogger
		Comment *dbPkg.Comment
	}, 0)
	for _, comment := range comments {
		dataElem := struct {
			Blogger *dbPkg.Blogger
			Comment *dbPkg.Comment
		}{
			Blogger: a.Db.GetBloggerByBloggerId(comment.BloggerId),
			Comment: comment,
		}
		commentsForms = append(commentsForms, dataElem)
	}

	data := struct {
		Article       *dbPkg.Article
		LikesCnt      int
		CommentsCnt   int
		CommentsForms []struct {
			Blogger *dbPkg.Blogger
			Comment *dbPkg.Comment
		}
	}{
		Article:       article,
		LikesCnt:      likesCnt,
		CommentsForms: commentsForms,
		CommentsCnt:   commentsCnt,
	}

	t, err := template.ParseFiles(bloggerArticleHtmlPath, headerHtmlPath, footerHtmlPath, navbarHtmlPath)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err) //a.Db.Logger.Fatal(err)
	}
}

func (a *Api) PostBloggerArticle(w http.ResponseWriter, r *http.Request) {
	bloggerId := getBloggerFromCtx(r.Context()).BloggerId
	articleIdFromURL := chi.URLParam(r, "articleId")
	articleId, _ := strconv.Atoi(articleIdFromURL)

	comment := &dbPkg.Comment{
		BloggerId:      bloggerId,
		ArticleId:      articleId,
		CommentMessage: template.HTML(r.FormValue("commentMessage")),
		Date:           time.Now().Format("02.01.2006, 15:04:05"),
	}

	a.Db.InsertComment(comment)

	http.Redirect(w, r, "/bloggers/"+chi.URLParam(r, "bloggerId")+"/"+articleIdFromURL+"#", http.StatusFound)
}

type ArticleDataInput struct {
	ArticleId int
}

type ArticleDataOutput struct {
	LikesCnt int
	IsLiked  bool
}

func (a *Api) SomeoneIsLiked(w http.ResponseWriter, r *http.Request) {
	bloggerId := getBloggerFromCtx(r.Context()).BloggerId
	var data ArticleDataInput
	_ = json.NewDecoder(r.Body).Decode(&data)
	articleId := data.ArticleId

	isLiked := a.Db.IsLiked(bloggerId, articleId)
	if isLiked {
		a.Db.DeleteLike(bloggerId, articleId)
	} else {
		a.Db.InsertLike(bloggerId, articleId)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ArticleDataOutput{
		LikesCnt: a.Db.GetLikesCntByArticleId(articleId),
		IsLiked:  !isLiked,
	})
}

type IsLikedOutput struct {
	IsLiked bool
}

func (a *Api) ShowLikes(w http.ResponseWriter, r *http.Request) {
	bloggerId := getBloggerFromCtx(r.Context()).BloggerId
	var data ArticleDataInput
	_ = json.NewDecoder(r.Body).Decode(&data)
	articleId := data.ArticleId

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(IsLikedOutput{
		IsLiked: a.Db.IsLiked(bloggerId, articleId),
	})
}

type SubDataInput struct {
	BloggerId int
}

type SubDataOutput struct {
	IsSubscribed bool
}

func (a *Api) SomeoneIsSubscribed(w http.ResponseWriter, r *http.Request) {
	bloggerIdCur := getBloggerFromCtx(r.Context()).BloggerId
	var data SubDataInput
	_ = json.NewDecoder(r.Body).Decode(&data)
	bloggerId := data.BloggerId

	isSubscribed := a.Db.IsSubscribed(bloggerIdCur, bloggerId)
	if isSubscribed {
		a.Db.DeleteSubscription(bloggerIdCur, bloggerId)
	} else {
		a.Db.InsertSubscription(bloggerIdCur, bloggerId)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(SubDataOutput{
		IsSubscribed: !isSubscribed,
	})
}

func (a *Api) ShowSubscriptions(w http.ResponseWriter, r *http.Request) {
	bloggerIdCur := getBloggerFromCtx(r.Context()).BloggerId
	var data SubDataInput
	_ = json.NewDecoder(r.Body).Decode(&data)
	bloggerId := data.BloggerId

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(SubDataOutput{
		IsSubscribed: a.Db.IsSubscribed(bloggerIdCur, bloggerId),
	})
}

//func Logger(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		log.Println(r.URL.Path)
//		next.ServeHTTP(w, r)
//	})
//}

//
//func (a *Api) PostBlogger(w http.ResponseWriter, r *http.Request) {
//	bloggerId := getBloggerFromCtx(r.Context()).BloggerId
//
//	//bloggerViewId := getBloggerViewFromCtx(r.Context()).BloggerId
//
//	bloggerViewIdFromURL := chi.URLParam(r, "bloggerId")
//	bloggerViewId, _ := strconv.Atoi(bloggerViewIdFromURL)
//
//	fmt.Println("TWO ID's:", bloggerId, bloggerViewId)
//	//a.Db.InsertSubscription(bloggerId, bloggerViewId)
//
//	//http.Redirect(w, r, "/bloggers/"+strconv.Itoa(bloggerViewId)+"/subscribed", http.StatusFound)
//}
//
//func (a *Api) GetBloggerSubscribed(w http.ResponseWriter, r *http.Request) {
//	//bloggerId := getBloggerViewFromCtx(r.Context()).BloggerId
//
//	bloggerIdFromURL := chi.URLParam(r, "bloggerId")
//	bloggerId, _ := strconv.Atoi(bloggerIdFromURL)
//	blogger := a.Db.GetBloggerByBloggerId(bloggerId)
//
//	bloggerIdCur := getBloggerFromCtx(r.Context()).BloggerId
//	articlesLiked, articlesNotLiked := a.Db.GetLikedAndNotLiked(bloggerId, bloggerIdCur)
//
//	data := struct {
//		Blogger          *dbPkg.Blogger
//		ArticlesLiked    []*dbPkg.Article
//		ArticlesNotLiked []*dbPkg.Article
//	}{
//		Blogger:          blogger,
//		ArticlesLiked:    articlesLiked,
//		ArticlesNotLiked: articlesNotLiked,
//	}
//
//	t, err := template.ParseFiles(bloggerProfileSubscribedHtmlPath, headerHtmlPath, footerHtmlPath, navbarHtmlPath)
//	if err != nil {
//		log.Fatal(err) //a.Db.Logger.Fatal(err)
//	}
//	err = t.Execute(w, data)
//	if err != nil {
//		log.Fatal(err) //a.Db.Logger.Fatal(err)
//	}
//}
//
//func (a *Api) PostBloggerSubscribed(w http.ResponseWriter, r *http.Request) {
//	bloggerId := getBloggerFromCtx(r.Context()).BloggerId
//
//	//bloggerViewId := getBloggerViewFromCtx(r.Context()).BloggerId
//
//	bloggerViewIdFromURL := chi.URLParam(r, "bloggerId")
//	bloggerViewId, _ := strconv.Atoi(bloggerViewIdFromURL)
//
//	a.Db.DeleteSubscription(bloggerId, bloggerViewId)
//
//	http.Redirect(w, r, "/bloggers/"+strconv.Itoa(bloggerViewId), http.StatusFound)
//}

//func (a *Api) PostBloggerArticle(w http.ResponseWriter, r *http.Request) {
//	bloggerId := getBloggerFromCtx(r.Context()).BloggerId
//	articleIdFromURL := chi.URLParam(r, "articleId")
//	articleId, _ := strconv.Atoi(articleIdFromURL)
//
//	a.Db.InsertLike(bloggerId, articleId)
//
//	bloggerIdFromURL := chi.URLParam(r, "bloggerId")
//
//	http.Redirect(w, r, "/bloggers/"+bloggerIdFromURL+"/"+articleIdFromURL+"/liked", http.StatusFound)
//}
//
//func (a *Api) GetBloggerArticleLiked(w http.ResponseWriter, r *http.Request) {
//	articleIdFromURL := chi.URLParam(r, "articleId")
//	articleId, _ := strconv.Atoi(articleIdFromURL)
//	article := a.Db.GetArticleByArticleId(articleId)
//
//	//bloggerId := getBloggerFromCtx(r.Context()).BloggerId
//	//isLiked := a.Db.IsLiked(bloggerId, articleId)
//	//if isLiked {
//	//	http.Redirect(w, r, "/bloggers/"+chi.URLParam(r, "bloggerId")+"/"+articleIdFromURL, http.StatusFound)
//	//}
//
//	likesCnt := a.Db.GetLikesCntByArticleId(articleId)
//
//	data := struct {
//		Article  *dbPkg.Article
//		LikesCnt int
//	}{
//		Article:  article,
//		LikesCnt: likesCnt,
//	}
//
//	t, err := template.ParseFiles(bloggerArticleLikedHtmlPath, headerHtmlPath, footerHtmlPath, navbarHtmlPath)
//	if err != nil {
//		log.Fatal(err) //a.Db.Logger.Fatal(err)
//	}
//	err = t.Execute(w, data)
//	if err != nil {
//		log.Fatal(err) //a.Db.Logger.Fatal(err)
//	}
//}
//
//func (a *Api) PostBloggerArticleLiked(w http.ResponseWriter, r *http.Request) {
//	bloggerId := getBloggerFromCtx(r.Context()).BloggerId
//	articleIdFromURL := chi.URLParam(r, "articleId")
//	articleId, _ := strconv.Atoi(articleIdFromURL)
//
//	a.Db.DeleteLike(bloggerId, articleId)
//
//	bloggerIdFromURL := chi.URLParam(r, "bloggerId")
//
//	http.Redirect(w, r, "/bloggers/"+bloggerIdFromURL+"/"+articleIdFromURL, http.StatusFound)
//}

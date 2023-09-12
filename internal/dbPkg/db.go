package dbPkg

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"slices"
	_ "slices"
)

type DB struct {
	Database *sql.DB
	Logger   *log.Logger
}

type Blogger struct {
	BloggerId int
	Email     string
	Password  string
	AuthToken string
	Login     string
	Name      string
	Surname   string
	Country   string
}

type Article struct {
	ArticleId      int
	BloggerId      int
	Title          string
	ArticleMessage template.HTML
	Date           string
}

type Comment struct {
	CommentId      int
	BloggerId      int
	ArticleId      int
	CommentMessage string
	Date           string
}

type Like struct {
	BloggerId int
	ArticleId int
}

type Subscriptions struct {
	BloggerId    int
	BloggerIdSub int
}

func (d *DB) InitDB() {
	db, err := sql.Open("sqlite3", "internal/dbPkg/data.sqlite")
	if err != nil {
		d.Logger.Fatal(err) // Logger можно настроить и насройки изменятся во всей программе (инъекция зависимости!)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		d.Logger.Fatal(err)
	}
	d.Database = db
}

const dropTablesConst = `
	DROP TABLE IF EXISTS Blogger;
 	DROP TABLE IF EXISTS Article;
 	DROP TABLE IF EXISTS Comment;
 	DROP TABLE IF EXISTS Like;
`

const createTablesConst string = `
	CREATE TABLE IF NOT EXISTS Blogger (
		BloggerId INTEGER NOT NULL PRIMARY KEY,
		Email NVARCHAR NOT NULL,
		Password NVARCHAR NOT NULL,
		AuthToken NVARCHAR,
		Login NVARCHAR NOT NULL,
		Name NVARCHAR NOT NULL,
		Surname NVARCHAR NOT NULL,
		Country NVARCHAR
	);

	CREATE TABLE IF NOT EXISTS Article (
		ArticleId INTEGER NOT NULL PRIMARY KEY,
		BloggerId INTEGER NOT NULL,
		Title NVARCHAR NOT NULL,
		ArticleMessage NVARCHAR NOT NULL,
		Date NVARCHAR NOT NULL,
		FOREIGN KEY(BloggerId) REFERENCES Blogger(BloggerId) ON DELETE CASCADE ON UPDATE CASCADE
	);

	CREATE TABLE IF NOT EXISTS Subscriptions (
		BloggerId INTEGER NOT NULL,
		BloggerIdSub INTEGER NOT NULL,
		FOREIGN KEY(BloggerId) REFERENCES Blogger(BloggerId) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS Comment (
	    CommentId INTEGER NOT NULL PRIMARY KEY,
		BloggerId INTEGER NOT NULL,
		ArticleId INTEGER NOT NULL,
		CommentMessage NVARCHAR NOT NULL,
		Date DATETIME NOT NULL,
		FOREIGN KEY(BloggerId) REFERENCES Blogger(BloggerId) ON DELETE CASCADE,
		FOREIGN KEY(ArticleId) REFERENCES Article(ArticleId) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS Like (
		BloggerId INTEGER NOT NULL,
		ArticleId INTEGER NOT NULL,
		FOREIGN KEY(BloggerId) REFERENCES Blogger(BloggerId) ON DELETE CASCADE,
		FOREIGN KEY(ArticleId) REFERENCES Article(ArticleId) ON DELETE CASCADE,
	    PRIMARY KEY (BloggerId, ArticleId)
	);
`

// Create tables
func (d *DB) CreateTables() {
	_, err := d.Database.Exec(createTablesConst)
	if err != nil {
		log.Fatal(err)
	}
}

// Drop tables
func (d *DB) DropTables() {
	_, err := d.Database.Exec(dropTablesConst)
	if err != nil {
		log.Fatal(err)
	}
}

// Insert Blogger
const insertBloggerConst = `INSERT INTO Blogger VALUES(NULL, ?, ?, ?, ?, ?, ?, ?);`

func (d *DB) InsertBlogger(b *Blogger) {
	_, err := d.Database.Exec(insertBloggerConst, b.Email, b.Password, b.AuthToken, b.Login, b.Name, b.Surname, b.Country)
	if err != nil {
		log.Fatal(err)
	}
}

// Insert Article
const insertArticleConst = `INSERT INTO Article VALUES(NULL, ?, ?, ?, ?);`

func (d *DB) InsertArticle(a *Article) int64 {
	fmt.Println(a.Date)
	lastPost, err := d.Database.Exec(insertArticleConst, a.BloggerId, a.Title, a.ArticleMessage, a.Date)
	if err != nil {
		d.Logger.Fatal(err)
	}

	lastId, _ := lastPost.LastInsertId()
	return lastId
}

// Insert Comment
const insertCommentConst = `INSERT INTO Comment VALUES(NULL, ?, ?, ?, ?);`

func (d *DB) insertComment(bloggerId int, articleId int, commentMessage string, date string) {
	_, err := d.Database.Exec(insertCommentConst, bloggerId, articleId, commentMessage, date)
	if err != nil {
		d.Logger.Fatal(err)
	}
}

// Insert Like
const insertLikeConst = `INSERT INTO Like VALUES(?, ?);`

func (d *DB) InsertLike(bloggerId int, articleId int) {
	_, err := d.Database.Exec(insertLikeConst, bloggerId, articleId)
	if err != nil {
		d.Logger.Fatal(err)
	}
}

// Delete Blogger
const deleteBloggerConst = `DELETE FROM Blogger WHERE BloggerId = ?;`

func (d *DB) deleteBlogger(bloggerId int) {
	_, err := d.Database.Exec(deleteBloggerConst, bloggerId)
	if err != nil {
		d.Logger.Fatal(err)
	}
}

// Delete Article
const deletePostConst = `DELETE FROM Article WHERE ArticleId = ?;`

func (d *DB) deletePost(articleId int) {
	_, err := d.Database.Exec(deletePostConst, articleId)
	if err != nil {
		d.Logger.Fatal(err)
	}
}

// Delete Comment
const deleteCommentConst = `DELETE FROM Comment WHERE CommentId = ?;`

func (d *DB) deleteComment(commentId int) {
	_, err := d.Database.Exec(deleteCommentConst, commentId)
	if err != nil {
		d.Logger.Fatal(err)
	}
}

// Delete Comment
const deleteLikeConst = `DELETE FROM Like WHERE BloggerId = ? AND articleId = ?;`

func (d *DB) DeleteLike(bloggerId int, articleId int) {
	_, err := d.Database.Exec(deleteLikeConst, bloggerId, articleId)
	if err != nil {
		d.Logger.Fatal(err)
	}
}

const insertSubscriptionConst = `INSERT INTO Subscriptions VALUES(?, ?);`

func (d *DB) InsertSubscription(bloggerId int, bloggerViewId int) {
	_, err := d.Database.Exec(insertSubscriptionConst, bloggerId, bloggerViewId)
	if err != nil {
		log.Fatal(err)
	}
}

const deleteSubscriptionConst = `DELETE FROM Subscriptions WHERE BloggerId = ? AND BloggerIdSub = ?;`

func (d *DB) DeleteSubscription(bloggerId int, bloggerViewId int) {
	_, err := d.Database.Exec(deleteSubscriptionConst, bloggerId, bloggerViewId)
	if err != nil {
		log.Fatal(err)
	}
}

// Update Blogger
const updateBloggerConst = `
	UPDATE Blogger SET Email = ?, Login = ?, Name = ?, Surname = ?, Country = ? WHERE BloggerId = ?;
`

func (d *DB) updateBlogger(email, login, name, surname, country string, bloggerId int) {
	_, err := d.Database.Exec(updateBloggerConst, email, login, name, surname, country, bloggerId)
	if err != nil {
		d.Logger.Fatal(err)
	}
}

const selectBlogger = `SELECT * FROM Blogger;`

func (d *DB) GetBloggerByCookie(cookie *http.Cookie) *Blogger {
	rows, err := d.Database.Query(selectBlogger)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		blogger := new(Blogger)
		err = rows.Scan(
			&blogger.BloggerId,
			&blogger.Email,
			&blogger.Password,
			&blogger.AuthToken,
			&blogger.Login,
			&blogger.Name,
			&blogger.Surname,
			&blogger.Country,
		)
		if err != nil {
			log.Fatal(err)
		}
		if blogger.AuthToken == cookie.Value {
			return blogger
		}
	}
	return nil
}

func (d *DB) GetBloggerByEmailPassword(emailForm string, passwordForm string) *Blogger {
	rows, err := d.Database.Query(selectBlogger)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		blogger := new(Blogger)
		err = rows.Scan(
			&blogger.BloggerId,
			&blogger.Email,
			&blogger.Password,
			&blogger.AuthToken,
			&blogger.Login,
			&blogger.Name,
			&blogger.Surname,
			&blogger.Country,
		)
		if err != nil {
			log.Fatal(err)
		}
		if blogger.Email == emailForm && blogger.Password == passwordForm {
			return blogger
		}
	}
	return nil
}

const selectArticles = `SELECT ArticleId, BloggerId, Title, ArticleMessage, Date FROM Article WHERE BloggerId = ?;`

func (d *DB) GetArticlesByBloggerId(bloggerId int) []*Article {
	rows, err := d.Database.Query(selectArticles, bloggerId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	articles := make([]*Article, 0)
	for rows.Next() {
		data := new(Article)
		err = rows.Scan(
			&data.ArticleId,
			&data.BloggerId,
			&data.Title,
			&data.ArticleMessage,
			&data.Date,
		)
		if err != nil {
			log.Fatal(err)
		}
		articles = append(articles, data)
	}
	return articles
}

const selectArticle = `SELECT ArticleId, BloggerId, Title, ArticleMessage, Date FROM Article WHERE ArticleId = ?;`

func (d *DB) GetArticleByArticleId(articleId int) *Article {
	rows, err := d.Database.Query(selectArticle, articleId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		article := new(Article)
		err = rows.Scan(
			&article.ArticleId,
			&article.BloggerId,
			&article.Title,
			&article.ArticleMessage,
			&article.Date,
		)
		if err != nil {
			log.Fatal(err)
		}
		return article
	}
	return nil
}

const selectLike = `SELECT COUNT(*) FROM Like WHERE ArticleId = ?;`

func (d *DB) GetLikesCntByArticleId(articleId int) int {
	rows, err := d.Database.Query(selectLike, articleId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var likesCnt int
	for rows.Next() {
		err = rows.Scan(&likesCnt)
		if err != nil {
			d.Logger.Fatal(err)
		}
	}
	return likesCnt
}

const selectIsLiked = `SELECT COUNT(*) FROM Like WHERE BloggerId = ? AND ArticleId = ?;`

func (d *DB) IsLiked(bloggerId int, articleId int) bool {
	rows, err := d.Database.Query(selectIsLiked, bloggerId, articleId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var isLiked int
	for rows.Next() {
		err = rows.Scan(&isLiked)
		if err != nil {
			d.Logger.Fatal(err)
		}
	}
	fmt.Println("isLiked", isLiked)
	if isLiked == 0 {
		return false
	} else {
		return true
	}
}

const selectIsSubscribed = `SELECT COUNT(*) FROM Subscriptions WHERE BloggerId = ? AND BloggerIdSub = ?;`

func (d *DB) IsSubscribed(bloggerId int, articleId int) bool {
	rows, err := d.Database.Query(selectIsSubscribed, bloggerId, articleId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var isSubscribed int
	for rows.Next() {
		err = rows.Scan(&isSubscribed)
		if err != nil {
			d.Logger.Fatal(err)
		}
	}
	if isSubscribed == 0 {
		return false
	} else {
		return true
	}
}

const selectSubsByBloggerId = `SELECT BloggerIdSub FROM Subscriptions WHERE BloggerId = ?`

func (d *DB) GetSubsIdByBloggerId(bloggerId int) []int {
	rows, err := d.Database.Query(selectSubsByBloggerId, bloggerId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var subs []int
	for rows.Next() {
		var bloggerIdSub int
		err = rows.Scan(&bloggerIdSub)
		if err != nil {
			d.Logger.Fatal(err)
		}
		subs = append(subs, bloggerIdSub)
	}
	return subs
}

const selectForBloggers = `SELECT BloggerId, Email, Login, Name, Surname, Country FROM Blogger`

func (d *DB) GetSubsAndNotSubs(exceptBloggerId int) ([]*Blogger, []*Blogger) {
	rows, err := d.Database.Query(selectForBloggers)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	subs := d.GetSubsIdByBloggerId(exceptBloggerId)
	bloggersNotSub := make([]*Blogger, 0)
	bloggersSub := make([]*Blogger, 0)
	for rows.Next() {
		data := new(Blogger)
		err = rows.Scan(
			&data.BloggerId,
			&data.Email,
			&data.Login,
			&data.Name,
			&data.Surname,
			&data.Country)
		if err != nil {
			log.Fatal(err)
		}
		if exceptBloggerId == data.BloggerId {
			continue
		}
		if slices.Contains(subs, data.BloggerId) {
			bloggersSub = append(bloggersSub, data)
		} else {
			bloggersNotSub = append(bloggersNotSub, data)
		}
	}

	return bloggersSub, bloggersNotSub
}

const selectLikedByBloggerId = `SELECT ArticleId FROM Like WHERE BloggerId = ?`

func (d *DB) GetLikedByBloggerId(bloggerId int) []int {
	rows, err := d.Database.Query(selectLikedByBloggerId, bloggerId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var liked []int
	for rows.Next() {
		var articleIdLiked int
		err = rows.Scan(&articleIdLiked)
		if err != nil {
			log.Fatal(err)
		}
		liked = append(liked, articleIdLiked)
	}
	return liked
}

const selectLiked = `SELECT ArticleId, BloggerId, Title, ArticleMessage, Date FROM Article WHERE BloggerId = ?;`

func (d *DB) GetLikedAndNotLiked(bloggerId int, bloggerIdCur int) ([]*Article, []*Article) {
	rows, err := d.Database.Query(selectLiked, bloggerId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	liked := d.GetLikedByBloggerId(bloggerIdCur)
	fmt.Println("liked id:", liked)
	articlesNotLiked := make([]*Article, 0)
	articlesLiked := make([]*Article, 0)
	for rows.Next() {
		data := new(Article)
		err = rows.Scan(
			&data.ArticleId,
			&data.BloggerId,
			&data.Title,
			&data.ArticleMessage,
			&data.Date,
		)
		if err != nil {
			log.Fatal(err)
		}
		if slices.Contains(liked, data.ArticleId) {
			articlesLiked = append(articlesLiked, data)
		} else {
			articlesNotLiked = append(articlesNotLiked, data)
		}
	}
	return articlesLiked, articlesNotLiked
}

const selectBloggerByBloggerId = `SELECT BloggerId, Login, Name, Surname, Country FROM Blogger WHERE BloggerId = ?`

func (d *DB) GetBloggerByBloggerId(bloggerId int) *Blogger {
	rows, err := d.Database.Query(selectBloggerByBloggerId, bloggerId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		blogger := new(Blogger)
		err = rows.Scan(
			&blogger.BloggerId,
			&blogger.Login,
			&blogger.Name,
			&blogger.Surname,
			&blogger.Country,
		)
		if err != nil {
			log.Fatal(err)
		}
		return blogger
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// ---------------------------------------------------------------------------------------------------

// Print table Blogger
func (d *DB) printTableBlogger() {
	fmt.Println("Blogger:")
	rows, err := d.Database.Query("SELECT * FROM Blogger")
	if err != nil {
		d.Logger.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var bloggerId int
		var email, login, name, surname, country string
		err = rows.Scan(&bloggerId, &email, &login, &name, &surname, &country)
		if err != nil {
			d.Logger.Fatal(err)
		}
		fmt.Printf("\t%d) %s %s %s %s %s\n", bloggerId, email, login, name, surname, country)
	}
}

// Print table Blogger
func (d *DB) printTablePost() {
	fmt.Println("Article:")
	rows, err := d.Database.Query("SELECT * FROM Article")
	if err != nil {
		d.Logger.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var bloggerId, articleId int
		var title, postMessage, date string
		err = rows.Scan(&articleId, &bloggerId, &title, &postMessage, &date)
		if err != nil {
			d.Logger.Fatal(err)
		}
		fmt.Printf("\t%d) %d-BloggerId %s %s %s\n", articleId, bloggerId, title, postMessage, date)
	}
}

// Print table Comment
func (d *DB) printTableComment() {
	fmt.Println("Comment:")
	rows, err := d.Database.Query("SELECT * FROM Comment")
	if err != nil {
		d.Logger.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var commentId, bloggerId, articleId int
		var commentMessage, date string
		err = rows.Scan(&commentId, &bloggerId, &articleId, &commentMessage, &date)
		if err != nil {
			d.Logger.Fatal(err)
		}
		fmt.Printf("\t%d) %d-BloggerId %d-articleId %s %s\n", commentId, bloggerId, articleId, commentMessage, date)
	}
}

// Print table Like
func (d *DB) printTableLike() {
	fmt.Println("Like:")
	rows, err := d.Database.Query("SELECT * FROM Like")
	if err != nil {
		d.Logger.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var bloggerId, articleId int
		err = rows.Scan(&bloggerId, &articleId)
		if err != nil {
			d.Logger.Fatal(err)
		}
		fmt.Printf("\t%d-BloggerId %d-articleId\n", bloggerId, articleId)
	}
}

func (d *DB) getAllBloggers() {
	rows, err := d.Database.Query("SELECT * FROM Blogger")
	if err != nil {
		d.Logger.Fatal(err)
	}
	defer rows.Close()

	bloggers := make([]*Blogger, 0)
	for rows.Next() {
		data := new(Blogger)
		err = rows.Scan(&data.BloggerId, &data.Email, &data.Login, &data.Name, &data.Surname, &data.Country)
		if err != nil {
			d.Logger.Fatal(err)
		}
		bloggers = append(bloggers, data)
	}

	fmt.Println(bloggers[0].Name)
	fmt.Println(bloggers)
}

func DbWork() {
	//blogger := &Blogger{
	//	Email:    "alexthecage@gmail.com",
	//	Password: "123",
	//	Login:    "molodoy",
	//	Name:     "Alex",
	//	Surname:  "Cage",
	//	Country:  "USA",
	//}

	d := DB{}

	d.InitDB()

	d.CreateTables()

	//d.InsertBlogger(blogger)

	d.printTableBlogger()

	fmt.Println("----------------------")

	d.getAllBloggers()

	//d.dropTables()
}

//dbPkg.insertBlogger("zlatoivan4@gmail.com", "dovolniy", "Ivan", "Zlat", "Russia")
//dbPkg.insertBlogger("alexthecage@gmail.com", "molodoy", "Alex", "Cage", "USA")
//dbPkg.printTableBlogger()
//fmt.Println()

//dbPkg.insertPost(1, "Sport", "Sport is heath", "2023-08-29 20:37")
//dbPkg.insertPost(2, "Theatre", "I visited Gogol Theatre", "2023-08-29 22:08")
//dbPkg.insertPost(2, "Music", "Rap is the best", "2023-08-29 22:11")
//lastarticleId := dbPkg.insertPost(2, "Dancing", "My love is Hip-hop", "2023-08-29 23:23")
//dbPkg.printTablePost()
//fmt.Println()
//fmt.Printf("\tlastarticleId = %d\n\n", lastarticleId)

//dbPkg.insertComment(2, 2, "Good post!", "2023-08-29 20:56")
//dbPkg.insertComment(1, 3, "WTF?", "2023-08-29 19:59")

//dbPkg.printTableComment()
//fmt.Println()

//dbPkg.insertLike(1, 2)
//dbPkg.insertLike(2, 3)
//dbPkg.printTableLike()

//deleteBlogger(1)
//deletePost(4)
//deleteComment(1)
//deleteLike(2, 3)

//dbPkg.updateBlogger("alexthecage@gmail.com", "stariy", "Alex", "Cage", "USA", 2)

//fmt.Println("\n---------------------------------------------------------------\n")

//dbPkg.printTableBlogger()
//fmt.Println()
//dbPkg.printTablePost()
//fmt.Println()
//dbPkg.printTableComment()
//fmt.Println()
//dbPkg.printTableLike()
//fmt.Println()

//dbPkg.dropTables()

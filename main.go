package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

const createTablesConst string = `
	DROP TABLE IF EXISTS Blogger;
	DROP TABLE IF EXISTS Post;
	DROP TABLE IF EXISTS Comment;
	DROP TABLE IF EXISTS Like;

	CREATE TABLE IF NOT EXISTS Blogger (
		BloggerId INTEGER NOT NULL PRIMARY KEY,
		Email NVARCHAR NOT NULL,
		Login NVARCHAR NOT NULL,
		Name NVARCHAR NOT NULL,
		Surname NVARCHAR NOT NULL,
		Country NVARCHAR
	);

	CREATE TABLE IF NOT EXISTS Post (
		PostId INTEGER NOT NULL PRIMARY KEY,
		BloggerId INTEGER NOT NULL,
		Title NVARCHAR NOT NULL,
		PostMessage NVARCHAR NOT NULL,
		Date DATETIME NOT NULL,
		FOREIGN KEY(BloggerId) REFERENCES Blogger(BloggerId) ON DELETE CASCADE ON UPDATE CASCADE
	);

	CREATE TABLE IF NOT EXISTS Comment (
	    CommentId INTEGER NOT NULL PRIMARY KEY,
		BloggerId INTEGER NOT NULL,
		PostId INTEGER NOT NULL,
		CommentMessage NVARCHAR NOT NULL,
		Date DATETIME NOT NULL,
		FOREIGN KEY(BloggerId) REFERENCES Blogger(BloggerId) ON DELETE CASCADE,
		FOREIGN KEY(PostId) REFERENCES Post(PostId) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS Like (
		BloggerId INTEGER NOT NULL,
		PostId INTEGER NOT NULL,
		FOREIGN KEY(BloggerId) REFERENCES Blogger(BloggerId) ON DELETE CASCADE,
		FOREIGN KEY(PostId) REFERENCES Post(PostId) ON DELETE CASCADE
	    PRIMARY KEY (BloggerId, PostId)
	);
`

// Create tables
func createTables() {
	db, err := sql.Open("sqlite3", "platform.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(createTablesConst)
	if err != nil {
		log.Fatal(err)
	}
}

// Insert Blogger
const insertBloggerConst = `INSERT INTO Blogger VALUES(NULL, ?, ?, ?, ?, ?);`

func insertBlogger(email string, login string, name string, surname string, country string) {
	db, err := sql.Open("sqlite3", "platform.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(insertBloggerConst, email, login, name, surname, country)
	if err != nil {
		log.Fatal(err)
	}
}

// Insert Post
const insertPostConst = `INSERT INTO Post VALUES(NULL, ?, ?, ?, ?);`

func insertPost(bloggerId int, title string, postMessage string, date string) int64 {
	db, err := sql.Open("sqlite3", "platform.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	lastPost, err := db.Exec(insertPostConst, bloggerId, title, postMessage, date)
	if err != nil {
		log.Fatal(err)
	}

	lastId, _ := lastPost.LastInsertId()
	return lastId
}

// Insert Comment
const insertCommentConst = `INSERT INTO Comment VALUES(NULL, ?, ?, ?, ?);`

func insertComment(bloggerId int, postId int, commentMessage string, date string) {
	db, err := sql.Open("sqlite3", "platform.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(insertCommentConst, bloggerId, postId, commentMessage, date)
	if err != nil {
		log.Fatal(err)
	}
}

// Insert Like
const insertLikeConst = `INSERT INTO Like VALUES(?, ?);`

func insertLike(bloggerId int, postId int) {
	db, err := sql.Open("sqlite3", "platform.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(insertLikeConst, bloggerId, postId)
	if err != nil {
		log.Fatal(err)
	}
}

// Delete Blogger
const deleteBloggerConst = `DELETE FROM Blogger WHERE BloggerId = ?;`

func deleteBlogger(bloggerId int) {
	db, err := sql.Open("sqlite3", "platform.db")
	db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(deleteBloggerConst, bloggerId)
	if err != nil {
		log.Fatal(err)
	}
}

// Delete Post
const deletePostConst = `DELETE FROM Post WHERE PostId = ?;`

func deletePost(postId int) {
	db, err := sql.Open("sqlite3", "platform.db")
	db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(deletePostConst, postId)
	if err != nil {
		log.Fatal(err)
	}
}

// Delete Comment
const deleteCommentConst = `DELETE FROM Comment WHERE CommentId = ?;`

func deleteComment(commentId int) {
	db, err := sql.Open("sqlite3", "platform.db")
	db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(deleteCommentConst, commentId)
	if err != nil {
		log.Fatal(err)
	}
}

// Delete Comment
const deleteLikeConst = `DELETE FROM Like WHERE BloggerId = ? AND PostId = ?;`

func deleteLike(bloggerId int, postId int) {
	db, err := sql.Open("sqlite3", "platform.db")
	db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(deleteLikeConst, bloggerId, postId)
	if err != nil {
		log.Fatal(err)
	}
}

// Update Blogger
const updateBloggerConst = `UPDATE Blogger SET Login = ? WHERE BloggerId = ?;`

func updateBlogger(login string, bloggerId int) {
	db, err := sql.Open("sqlite3", "platform.db")
	db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(updateBloggerConst, login, bloggerId)
	if err != nil {
		log.Fatal(err)
	}
}

// Print table Blogger
func printTableBlogger() {
	fmt.Println("Blogger:")
	db, err := sql.Open("sqlite3", "platform.db")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * FROM Blogger")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var bloggerId int
		var email, login, name, surname, country string
		err = rows.Scan(&bloggerId, &email, &login, &name, &surname, &country)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\t%d) %s %s %s %s %s\n", bloggerId, email, login, name, surname, country)
	}
}

// Print table Blogger
func printTablePost() {
	fmt.Println("Post:")
	db, err := sql.Open("sqlite3", "platform.db")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * FROM Post")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var bloggerId, postId int
		var title, postMessage, date string
		err = rows.Scan(&postId, &bloggerId, &title, &postMessage, &date)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\t%d) %d-bloggerId %s %s %s\n", postId, bloggerId, title, postMessage, date)
	}
}

// Print table Comment
func printTableComment() {
	fmt.Println("Comment:")
	db, err := sql.Open("sqlite3", "platform.db")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * FROM Comment")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var commentId, bloggerId, postId int
		var commentMessage, date string
		err = rows.Scan(&commentId, &bloggerId, &postId, &commentMessage, &date)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\t%d) %d-bloggerId %d-postId %s %s\n", commentId, bloggerId, postId, commentMessage, date)
	}
}

// Print table Like
func printTableLike() {
	fmt.Println("Like:")
	db, err := sql.Open("sqlite3", "platform.db")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * FROM Like")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var bloggerId, postId int
		err = rows.Scan(&bloggerId, &postId)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\t%d-bloggerId %d-postId\n", bloggerId, postId)
	}
}

func main() {
	fmt.Println()

	createTables()

	insertBlogger("zlatoivan4@gmail.com", "dovolniy", "Ivan", "Zlat", "Russia")
	insertBlogger("alexthecage@gmail.com", "molodoy", "Alex", "Cage", "USA")
	printTableBlogger()
	fmt.Println()

	insertPost(1, "Sport", "Sport is heath", "2023-08-29 20:37")
	insertPost(2, "Theatre", "I visited Gogol Theatre", "2023-08-29 22:08")
	insertPost(2, "Music", "Rap is the best", "2023-08-29 22:11")
	lastPostId := insertPost(2, "Dancing", "My love is Hip-hop", "2023-08-29 23:23")
	printTablePost()
	fmt.Println()
	fmt.Printf("\tlastPostId = %d\n\n", lastPostId)

	insertComment(2, 2, "Good post!", "2023-08-29 20:56")
	insertComment(1, 3, "WTF?", "2023-08-29 19:59")

	printTableComment()
	fmt.Println()

	insertLike(1, 2)
	insertLike(2, 3)
	printTableLike()

	//deleteBlogger(1)
	//deletePost(4)
	//deleteComment(1)
	//deleteLike(2, 3)

	updateBlogger("stariy", 2)

	fmt.Println("\n---------------------------------------------------------------\n")

	printTableBlogger()
	fmt.Println()
	printTablePost()
	fmt.Println()
	printTableComment()
	fmt.Println()
	printTableLike()
	fmt.Println()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":3333", nil)

}

/*
CREATE Blogger
CREATE Post
CREATE Comment
CREATE Like

INSERT Blogger
INSERT Post
INSERT Comment
INSERT Like

DELETE Blogger
DELETE Post
DELETE Comment
DELETE Like
*/

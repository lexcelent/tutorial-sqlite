package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

/*
TODO: Сделать структуру для БД, определить для нее методы
*/

// Книга
type Book struct {
	id       int
	title    string
	author   string
	numPages int
	rating   float64
}

// Создать таблицу
func CreateTable(db *sql.DB) error {
	query := `
		drop table if exists books;
		create table if not exists books(
			id integer primary key,
			title text,
			author text,
			num_pages integer,
			rating real
		);
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	fmt.Println("create books table")
	return nil
}

// Добавление новых книг в таблицу
func InsertBooks(db *sql.DB) error {
	query := `
		insert into books(title, author, num_pages, rating)
		values (?, ?, ?, ?)
	`

	data := [][]any{
		{"Code Complete", "Steve McConnell", 500, 5},
		{"Белый Клык", "Джек Лондон", 300, 4.5},
	}

	for _, vals := range data {
		res, err := db.Exec(query, vals...)
		if err != nil {
			return err
		}
		bookId, err := res.LastInsertId()
		fmt.Printf("added new book: id=%d, error=%v\n", bookId, err)
	}
	return nil
}

// Добавить одну книгу
func InsertBook(db *sql.DB, title, author string, num_pages int) error {
	query := `
		insert into books(title, author, num_pages, rating)
		values (?, ?, ?, ?) 
	`

	res, err := db.Exec(query, title, author, num_pages, 0)
	if err != nil {
		return err
	}
	bookId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Printf("added new book: id=%d\n", bookId)

	return nil
}

// Обновить информацию о книгах
func UpdateBooks(db *sql.DB) error {
	query := "update books set author = ? where author = ?"
	res, err := db.Exec(query, "Steve McConnell", "Стив Макконнел")
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	fmt.Printf("updated %d books, error=%v\n", count, err)
	return nil
}

// Удалить книги
func DeleteBooks(db *sql.DB) error {
	query := "delete from books where rating < 4"
	res, err := db.Exec(query)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	fmt.Printf("deleted %d books, error=%v\n", count, err)
	return nil
}

// Вывести на экран информацию о книгах
func SelectAll(db *sql.DB) error {
	query := "select id, title, author, num_pages, rating from books"
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.id, &book.title, &book.author, &book.numPages, &book.rating)
		if err != nil {
			return err
		}
		fmt.Println(book)
	}

	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func main() {
	// Инициализация
	db, err := sql.Open("sqlite3", "books.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("connected to books.db")

	// Работа с БД
	CreateTable(db)

	err = InsertBooks(db)
	if err != nil {
		panic(err)
	}

	err = InsertBook(db, "Зов предков", "Джек Лондон", 400)
	if err != nil {
		panic(err)
	}

	err = InsertBook(db, "Морской волк", "Джек Лондон", 350)
	if err != nil {
		panic(err)
	}
	SelectAll(db)
}

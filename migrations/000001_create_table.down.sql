DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS book_rental;
DROP TABLE IF EXISTS unique_book_rental;

DROP INDEX IF EXISTS idx_books_author_id;
DROP INDEX IF EXISTS idx_book_rental_book_id;
DROP INDEX IF EXISTS idx_book_rental_user_id;
DROP INDEX IF EXISTS idx_book_rental_return_date;

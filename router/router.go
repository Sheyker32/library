package router

import (
	"library/internal/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewApiRouter(authorController handler.Authorer, bookController handler.Booker, rentController handler.Rentaler, userController handler.Userer) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/author", authorController.CreateAuthor)
		r.Get("/author/{authorId}", authorController.GetAuthor)
		r.Get("/author/top", authorController.GetTopAuthors)
		r.Get("/author/all", authorController.GetAllAuthors)
		r.Delete("/author/{authorId}", authorController.DeleteAuthor)
		r.Get("/author/books/{authorId}", authorController.GetByBooksAuthor)

	})

	r.Group(func(r chi.Router) {
		r.Post("/book", bookController.AddBook)
		r.Get("/book/{bookId}", bookController.GetBook)
		r.Delete("/book/author/{authorId}", bookController.DeleteBook)

	})

	r.Group(func(r chi.Router) {
		r.Post("/rental/{bookId}/{userId}", rentController.RentBook)
		r.Delete("/rental/{bookId}", rentController.ReturnBook)
	})

	r.Group(func(r chi.Router) {
		r.Post("/user", userController.Create)
		r.Get("/user/{userId}", userController.GetByID)
		r.Delete("/user/{userId}", userController.DeleteUser)
		r.Get("/user/all", userController.GetAll)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	return r
}

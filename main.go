package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v2"

	"atorgold/auth"
	"atorgold/database"
	"atorgold/product"
	"atorgold/web"
)

const (
	viewPath   = "./client/templates"
	publicPath = "./client/public"
	renderType = ".html"
)

func add(x, y int) int {
	return x + y
}

func main() {
	database.InitiliazeDatabaseConnection()
	engine := html.New(viewPath, renderType)
	engine.AddFunc("unescape", func(s string) template.HTML {
		return template.HTML(s)
	})

	engine.AddFunc("safe", func(s string) template.HTML {
		return template.HTML(s) // HTML olarak işaretler, güvenli kabul eder
	})

	engine.AddFunc("attr", func(s string) template.HTMLAttr {
		return template.HTMLAttr(s) // Attribute olarak işaretler
	})
	engine.AddFunc("safeHTML", func(s string) template.HTML {
		return template.HTML(s) // HTML olarak işaretle
	})
	engine.AddFunc("raw", func(s string) template.HTML {
		return template.HTML(s) // Mark string as raw HTML
	})
	engine.AddFunc("add", add)
	app := fiber.New(fiber.Config{
		ReadTimeout:   time.Minute * time.Duration(5),
		StrictRouting: false,
		CaseSensitive: true,
		BodyLimit:     4 * 1024 * 1024,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		AppName:       "atorgold",
		Immutable:     true,
		Views:         engine,
		//ViewsLayout: "layouts/main",
		ErrorHandler: func(c fiber.Ctx, err error) error {
			var e *fiber.Error
			if errors.As(err, &e) {
				if e.Code == fiber.StatusNotFound {
					return c.Render("404", fiber.Map{
						"Title": "Page Not Found",
					})
				}
				return c.Status(e.Code).Render("error", fiber.Map{
					"Title":   "Error",
					"Message": e.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).Render("error", fiber.Map{
				"Title":   "Internal Server Error",
				"Message": "An unexpected error occured.",
			})
		},
	})

	app.Use(static.New(publicPath))


	

	app.Get("/", web.IndexPage)
	app.Get("/about-us", web.AboutPage)
	app.Get("/contacts", web.ContactsPage)
	app.Get("/products", web.ProductsListPage)
	app.Get("/product-single/:main_id", web.ProductSinglePage)
	app.Get("/add-product", web.AddProductPage, auth.IsAuthorized, auth.RateLimiter(120, time.Minute))
	app.Get("/login", web.LoginPage, auth.RateLimiter(5, time.Minute))

	route := app.Group("/auth")
	route.Post("/login", auth.Login, auth.RateLimiter(5, time.Minute), auth.LoginValidation)
	route.Post("/register", auth.Register, auth.RateLimiter(5, time.Minute), auth.RegisterValidation)


	document := app.Group("/documenter")
	document.Post("/main", product.CreateMainDocument)
	document.Post("/sub", product.CreateSubDocument)
	document.Post("/content", product.CreateContentDocument)
	document.Get("/all", product.GetAllDocuments)
	document.Get("/all-join", product.GetAllDocumentsByJoin)

	
	
	app.Use(web.NotFoundPage)
	
	//s.app.Get("/dashboard", s.DashboardWeb, s.authMiddleware)

	log.Fatal(app.Listen("0.0.0.0:3000"))
}

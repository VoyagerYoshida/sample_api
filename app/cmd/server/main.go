package main

import (
    "database/sql"
    // "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    // "reflect"
    // "strings"
    // "time"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    _ "github.com/lib/pq"
    "github.com/coopernurse/gorp"
    // "github.com/go-playground/locales/ja_JP"
    // "gopkg.in/go-playground/validator.v9"
    // ja "gopkg.in/go-playground/validator.v9/translations/ja"
    // ut "github.com/go-playground/universal-translator"
)

var dbDriver = "postgres"

type Comment struct {
    Id      int64  `json:"id" db:"id,primarykey,autoincrement"`
    Name    string `json:"name" db:"name,notnull"`
    Content string `json:"content" db:"content,notnull"`
}

func setupDB() (*gorp.DbMap, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("CONTAINER_NAME_DB"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
        os.Getenv("POSTGRES_DB"), os.Getenv("PORT_DB"))
    db, err := sql.Open(dbDriver, dsn)
    if err != nil {
        return nil, err
    }

    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
    dbmap.AddTableWithName(Comment{}, "tab_comments").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}

    return dbmap, nil
}

type Error struct {
	Error string `json:"error"`
}

type Controller struct {
	dbmap *gorp.DbMap
}

// func (controller *Controller) GetComment(c echo.Context) error {
// 	var comment Comment
// 	err := controller.dbmap.SelectOne(&comment,
// 		"SELECT * FROM comments WHERE id = $1", c.Param("id"))
// 	if err != nil {
// 		if err != sql.ErrNoRows {
// 			c.Logger().Error("SelectOne: ", err)
// 			return c.String(http.StatusBadRequest, "SelectOne: "+err.Error())
// 		}
// 		return c.String(http.StatusNotFound, "Not Found")
// 	}
// 	return c.JSON(http.StatusOK, comment)
// }

func (controller *Controller) ListComments(c echo.Context) error {
	var comments []Comment
	_, err := controller.dbmap.Select(&comments, "SELECT id, name, content FROM tab_comments")
		// "SELECT * FROM comments ORDER BY created desc LIMIT 10")
	if err != nil {
		c.Logger().Error("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}
	return c.JSON(http.StatusOK, comments)
}

func (controller *Controller) InsertComment(c echo.Context) error {
	var comment Comment
	if err := c.Bind(&comment); err != nil {
		c.Logger().Error("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}
	if err := c.Validate(&comment); err != nil {
		c.Logger().Error("Validate: ", err)
		return c.JSON(http.StatusBadRequest, &Error{Error: err.Error()})
	}
	if err := controller.dbmap.Insert(&comment); err != nil {
		c.Logger().Error("Insert: ", err)
		return c.String(http.StatusBadRequest, "Insert: "+err.Error())
	}
	c.Logger().Infof("inserted comment: %v", comment.Id)
	return c.NoContent(http.StatusCreated)
}

func main() {
    dbmap, err := setupDB()
    if err != nil {
        log.Fatal(err)
    }

    controller := &Controller{dbmap: dbmap}

    e := echo.New()

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    e.GET("/", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]string{"hello": "world"})
    })
    // e.GET("/api/comments/:id", controller.GetComment)
	e.GET("/api/comments", controller.ListComments)
	e.POST("/api/comments", controller.InsertComment)

    e.Logger.Fatal(e.Start(":8080"))
}

// type Validator struct {
// 	trans     ut.Translator
// 	validator *validator.Validate
// }

// func (v *Validator) Validate(i interface{}) error {
// 	err := v.validator.Struct(i)
// 	if err == nil {
// 		return nil
// 	}
// 	errs := err.(validator.ValidationErrors)
// 	msg := ""
// 	for _, v := range errs.Translate(v.trans) {
// 		if msg != "" {
// 			msg += ", "
// 		}
// 		msg += v
// 	}
// 	return errors.New(msg)
// }

// func setupEcho() *echo.Echo {
// 	e := echo.New()
// 	e.Debug = true
// 	e.Logger.SetOutput(os.Stderr)

// 	// setup japanese translation
// 	japanese := ja_JP.New()
// 	uni := ut.New(japanese, japanese)
// 	trans, _ := uni.GetTranslator("ja")
// 	validate := validator.New()
// 	err := ja.RegisterDefaultTranslations(validate, trans)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// register japanese translation for input field
// 	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
// 		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
// 		switch name {
// 		case "name":
// 			return "お名前"
// 		case "text":
// 			return "コメント"
// 		case "-":
// 			return ""
// 		}
// 		return name
// 	})

// 	e.Validator = &Validator{validator: validate, trans: trans}
// 	return e
// }

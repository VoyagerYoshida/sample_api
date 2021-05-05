package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    _ "github.com/lib/pq"
    "github.com/coopernurse/gorp"
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

func (controller *Controller) ListComments(c echo.Context) error {
	var comments []Comment
	_, err := controller.dbmap.Select(&comments, "SELECT * FROM tab_comments")
	if err != nil {
		c.Logger().Error("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}
	return c.JSON(http.StatusOK, comments)
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
	e.GET("/api/comments", controller.ListComments)

    e.Logger.Fatal(e.Start(":8080"))
}

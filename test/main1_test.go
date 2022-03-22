package main

import (
    "bytes"
    v1 "dcard-project/controller/api/v1"
    "dcard-project/database"
    "dcard-project/models"
    "fmt"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestSetupRouter(t *testing.T) {
    if err := initDBMock(); err != nil {
        println(err)
        return
    }
    router := initRoutesMock()
    w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
    var jsonStr = []byte(`{"org_url":"https://www.google.com"}`)
    req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "miles")
}

func initRoutesMock() *gin.Engine {
    r := gin.Default()

    v1Group := r.Group("/")
    {
        v1Group.POST("", v1.CreateUrl)
        v1Group.GET("/:shortUrl", v1.ToOrgPage)
    }

    return r
}

type Result struct {
    ID   int
    Name string
}

func initDBMock() error {
    db, _, err := sqlmock.New()
    if err != nil {
        return err
    }
    //測試時不需要連接資料庫
    gdb, err := gorm.Open(postgres.New(postgres.Config{
        DriverName:           "postgres",
        PreferSimpleProtocol: true,
        Conn:                 db,
    }), &gorm.Config{})
    if err != nil {
        return err
    }

    //rows := sqlmock.NewRows([]string{"id", "name"}).
    //    AddRow(1, 111).AddRow(2, 222)
    //
    //mock.ExpectQuery("Select * from users").WillReturnRows(rows)
    //var result Result
    //res := gdb.Raw("Select * from users").Scan(&result)
    //fmt.Printf("%#v", res)

    database.Db = gdb
    err = gdb.Migrator().CreateTable(&models.Url{})
    if err != nil {
        return err
    }
    gdb.Table("urls").Create(models.Url{
        ID:     1,
        OrgUrl: "123123",
    })
    gdb.Table("urls").First(&models.Url{})
    fmt.Printf("%v", models.Url{})

    return nil
}

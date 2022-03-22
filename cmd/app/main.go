package main

import (
    "dcard-project/database"
    _ "dcard-project/docs"
    "dcard-project/routers"
    "fmt"
)

func main() {
    _, err := database.NewPgClient()
    if err != nil {
        fmt.Printf("err open databases : %v", err)
        return
    }

    database.NewRedisClient()

    r := routers.SetRouter()
    r.Run()
}

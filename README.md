# 做此專案的動機
在實務開發上，時常遇到專案ORM架構不優雅的案例。
比如說簡單的All, Find功能，這在所有的模型都是一樣的邏輯，只是換了張表或是collection。
我認為golang的ORM應能設計成常態通用的方法就不必重複刻輪子。

鑒於golang不是OOP，沒有繼承、抽象的概念，實作上也是思考了許久。

本專案乃個人嘗試做一個簡潔優雅mongodb 的 orm

目標為開發出可重複利用、可擴展、容易維護的ORM
# todo
- eloquent
    - create InsertMultiple
    - create DeleteMultiple
    - create UpdateMultiple
    - create FindMultiple
    - create Count
    - create Paginate
- repo
    - create GetUnderage
# 開始前的作業
- cp .env.example .env
# 建立屬於個別collection的struct
- 範例如下
```go
type User struct {
	ID        string `bson:"_id"`
	Name      string `bson:"name"`
	Age       uint16 `bson:"age"`
	CreatedAt uint64 `bson:"created_at"`
	UpdatedAt uint64 `bson:"updated_at"`
}
```

# usage example
```go
package main

import (
	"fmt"
	"log"

	"github.com/LIOU2021/go-eloquent-mongodb/core"
	"github.com/LIOU2021/go-eloquent-mongodb/orm"
)

type User struct {
	ID        string `bson:"_id"`
	Name      string `bson:"name"`
	Age       uint16 `bson:"age"`
	CreatedAt uint64 `bson:"created_at"`
	UpdatedAt uint64 `bson:"updated_at"`
}

func main() {
	core.Setup()
	defer core.Cleanup()

	userOrm := orm.NewEloquent[User]("users")
	id := "6380cf16742f1bd2061f28b8"
	user, ok := userOrm.Find(id)

	if !ok {
		log.Fatal("fail !")
	}

	fmt.Printf("id : %s, name : %s, age : %d, created_time : %d, updated_time : %d\n", user.ID, user.Name, user.Age, user.CreatedAt, user.UpdatedAt)
}

```

# Ref
- https://www.mongodb.com/docs/drivers/go/current/quick-start/
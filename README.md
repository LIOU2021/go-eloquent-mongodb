# 做此專案的動機
在實務開發上，時常遇到專案ORM架構不優雅的案例。
比如說簡單的All, Find功能，這在所有的模型都是一樣的邏輯，只是換了張表或是collection。
我認為golang的ORM應能設計成常態通用的方法就不必重複刻輪子。

鑒於golang不是OOP，沒有繼承、抽象的概念，實作上也是思考了許久。

本專案乃個人嘗試做一個簡潔優雅mongodb 的 orm

目標為開發出可重複利用、可擴展、容易維護的ORM

開發的後期，因為model與ORM本身的依賴與責任設計的不良，也時常導致出現一堆model混亂的場景，本ORM將會克服此情境。

# todo
- eloquent
    - create DeleteMultiple
    - create UpdateMultiple
    - create FindMultiple
    - create Paginate
- repo
    - create GetUnderage
# 開始前的作業
- cp .env.example .env
# usage example
- 更多範例請直接參考 tests\test

```go
package main

import (
	"fmt"
	"log"

	"github.com/LIOU2021/go-eloquent-mongodb/core"
	"github.com/LIOU2021/go-eloquent-mongodb/orm"
)

type User struct {
	ID        *string `bson:"_id,omitempty"`
	Name      *string `bson:"name,omitempty"`
	Age       *uint16 `bson:"age,omitempty"`
	CreatedAt *uint64 `bson:"created_at,omitempty"`
	UpdatedAt *uint64 `bson:"updated_at,omitempty"`
}

func main() {
	core.Setup()
	defer core.Cleanup()

	userOrm := orm.NewEloquent[User]("users")
	id := "6380cf16742f1bd2061f28b8"
	user, ok := userOrm.Find(id)

	if !ok {
		log.Fatal("user id not found !")
	}

	fmt.Printf("id : %s, name : %s, age : %d, created_time : %d, updated_time : %d\n", *user.ID, *user.Name, *user.Age, *user.CreatedAt, *user.UpdatedAt)
}

```

# Ref
- https://www.mongodb.com/docs/drivers/go/current/quick-start/
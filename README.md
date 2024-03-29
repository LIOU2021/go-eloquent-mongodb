[![Go Reference](https://pkg.go.dev/badge/golang.org/x/pkgsite.svg)](https://pkg.go.dev/github.com/LIOU2021/go-eloquent-mongodb)

# 做此專案的動機
在實務開發上，時常遇到專案ORM架構不優雅的案例。
比如說簡單的All, Find功能，這在所有的模型都是一樣的邏輯，只是換了張表或是collection。
我認為golang的ORM應能設計成常態通用的方法就不必重複刻輪子。

鑒於golang不是OOP，沒有繼承、抽象的概念，實作上也是思考了許久。

本專案乃個人嘗試做一個簡潔優雅mongodb 的 orm

目標為開發出可重複利用、可擴展、容易維護的ORM

開發的後期，因為model與ORM本身的依賴與責任設計的不良，也時常導致出現一堆model混亂的場景，本ORM將會克服此情境。

# todo
- createIndex and ttl
	- ref
		- https://christiangiacomi.com/posts/mongodb-index-using-go/
		- https://www.mongodb.com/docs/drivers/go/v1.8/fundamentals/indexes/
		- https://www.mongodb.com/docs/manual/core/index-ttl/
		- https://fmabid.medium.com/creating-ttl-index-in-mongodb-with-go-mongo-driver-48a51d899241

# usage example
- more sample see tests\test

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/LIOU2021/go-eloquent-mongodb/orm"
)

type User struct {
	ID        *string `bson:"_id,omitempty"`
	Name      *string `bson:"name,omitempty"`
	Age       *int    `bson:"age,omitempty"`
	CreatedAt *int64  `bson:"created_at,omitempty"`
	UpdatedAt *int64  `bson:"updated_at,omitempty"`
}

func main() {
	orm.Setup("go-eloquent-mongo", "127.0.0.1", "27017", "")
	ctx := context.Background()
	orm.Connect(ctx)
	defer orm.Disconnect(ctx)

	userOrm := orm.NewEloquent[User]("users")
	id := "642d5b2298ba2bb73c55e5c4"
	user, err := userOrm.Find(context.Background(), id)

	if err != nil {
		log.Fatal("user id not found !")
	}

	fmt.Printf("id : %s, name : %s, age : %d, created_at : %d, updated_at : %d\n", *user.ID, *user.Name, *user.Age, *user.CreatedAt, *user.UpdatedAt)
}

```

# Ref
- https://www.mongodb.com/docs/drivers/go/current/quick-start/

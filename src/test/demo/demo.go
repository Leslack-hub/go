package main

import (
	"fmt"
	"os"
)

// 1. 基础结构体 Base
type Base struct {
	ID int
}

// 在 Base 上定义一个方法，我们希望嵌入它的结构体也能访问
func (b *Base) GetBaseID() int {
	return b.ID
}

// (可选) 可以给 Base 添加更多方法
func (b *Base) Identify() string {
	return fmt.Sprintf("Base ID: %d", b.ID)
}

// 2. 定义一个接口，要求类型具有 Base 的关键方法
type HasBase interface {
	GetBaseID() int
	Identify() string // 包含 Base 提供的 Identify 方法
}

// 3. "子类" 结构体，嵌入了 Base
type DerivedA struct {
	Base  // 嵌入 Base
	DataA string
}

// DerivedA 自动获得了 GetBaseID() 和 Identify() 方法，因此满足 HasBase 接口

type DerivedB struct {
	*Base // 也可以嵌入指针类型
	DataB bool
}

// DerivedB 也自动获得了 GetBaseID() 和 Identify() 方法，满足 HasBase 接口

// 其他不相关的结构体
type Other struct {
	Name string
}

// Other 没有嵌入 Base，也没有 GetBaseID() 或 Identify() 方法，不满足 HasBase

// 4. 定义泛型函数，使用 HasBase 接口作为约束
// T 必须是实现了 HasBase 接口的类型
// 由于 DerivedA 和 DerivedB 嵌入了 Base (并且 Base 有 GetBaseID/Identify 方法),
// 它们隐式地满足了 HasBase 接口
func ProcessItemWithBase[T HasBase](item T) {
	fmt.Printf("Processing item with Base ID: %d\n", item.GetBaseID())
	fmt.Printf("Identification: %s\n", item.Identify())

	// 注意：这里不能直接访问 item.ID，因为 T 的约束是 HasBase 接口，
	// 接口本身不知道有 ID 字段。我们只能访问接口定义的方法。
	// 如果需要访问嵌入的 Base 实例本身，可以稍作修改，见下面的进阶技巧。
}

type MyConstraint interface {
	*User | *Product | ~*Order
}

type User struct {
	ID   int
	Name string
}

type Product struct {
	ID    int
	Title string
}

type Order struct {
	ID     int
	Amount float64
}

type Order2 struct {
	*Order
	Name string
}

func ProcessData[T MyConstraint](data T) {
	var anyData any = data

	switch v := anyData.(type) {
	case *User:
		fmt.Printf("处理用户: %s (ID: %d)\n", v.Name, v.ID)
	case *Product:
		fmt.Printf("处理产品: %s (ID: %d)\n", v.Title, v.ID)
	case *Order:
		fmt.Printf("处理订单: ID %d, 金额 %.2f\n", v.ID, v.Amount)
	default:
	}
	// 这里只能传入 User、Product 或 Order 类型
	fmt.Printf("Processing: %+v\n", data)
}

type BaseEntity struct {
	ID  int
	Age uint32
}

type EntityConstraint interface {
	~struct {
		BaseEntity
		Age uint32
	}
}

func P[T EntityConstraint](data T) {
	fmt.Println(data)
}

type User1 struct {
	Age uint32
}

type Product1 struct {
	BaseEntity
	Title string
}

type HasAge interface {
	~struct{ Age uint32 }
}

func hasAge[T HasAge](data T) {
	fmt.Println(data)
}

func main() {
	user := User1{
		Age: 1,
	}

	// 需要传递 AgeEntity 部分
	hasAge(user) //ProcessData(&User{})
	//ProcessData(&Order2{})
	//ProcessData(&Product{})
	os.Exit(1)
	da := DerivedA{Base: Base{ID: 101}, DataA: "Hello"}
	db := DerivedB{Base: &Base{ID: 102}, DataB: true}
	// other := Other{Name: "unrelated"}

	//ProcessItemWithBase(da)  // 正确，DerivedA 满足 HasBase
	ProcessItemWithBase(&da) // 如果方法接收者是指针，传递指针也可以
	ProcessItemWithBase(db)  // 正确，DerivedB 满足 HasBase (通过嵌入的 *Base)
	ProcessItemWithBase(&db) // 指针类型也满足

	// 下面这行会编译错误，因为 Other 没有实现 HasBase 接口
	// ProcessItemWithBase(other)

	// 下面这行也会编译错误，因为 Base 本身虽然有 GetBaseID 和 Identify 方法，
	// 但如果我们传递的是值类型 Base{}，而方法的接收者是指针 (*Base)，则值类型不满足接口。
	// 如果接收者是值类型 (b Base)，则 Base{} 会满足。
	// bValue := Base{ID: 0}
	// ProcessItemWithBase(bValue) // 可能编译错误，取决于接收者类型

	// 如果传递指针，通常满足接口（假设接收者是指针或值）
	bPtr := &Base{ID: 0}
	ProcessItemWithBase(bPtr) // 正确
}

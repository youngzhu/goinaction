# 学习笔记

## ch02 
### defer
并不需要放在函数的最后一行。它会保证后面的语句在函数返回时执行。
可以在文件打开后立即跟着`defer`，例如：
```Go
file, err := os.Open("filePaht")
if err != nil {
    return nil, err
}
defer file.Close()

// 其他对file的操作
```

### 声明函数的接收者用值类型还是指针类型？
先说结论：因为大部分方法在被调用后都需要维护接收者的值的状态，所以，**最佳实践是将接收者声明为指针**。

|接收者\调用者|值|指针|接口类型的值|接口类型的指针|
|-|-|-|-|-|
|值类型|√|√|√|√|
|指针类型|√|√|x|√|

详见 reciver.go

## ch03 打包和工具链
### 给相同的包重新命名
```Go
import (
    "fmt"
    myfmt "mylib/fmt"
)
```

## ch04 数组/切片/映射
### 数组
```Go
// 声明一个包含5个元素的整型数组
var array [5]int

// 声明一个包含5个元素的整型数组
// 用具体值初始化每个元素
array := [5]int{10, 20, 30, 40, 50}

// 如果不知道具体长度（懒得数）
// 以下的数组容量由初始化的数量决定
array := [...]int{10, 20}

// 给指定元素特定的值
// 其他的还是0
array := [5]int{1:10, 2:20}

// 指针数组
array := [5]*int{0:new(int), 1:new(int)}
// 为指针元素赋值
*array[0] = 10
*array[1] = 20

// 将数组赋值给另一个数组
// 赋值后，两个数组的值完全一样
// 只有元素类型和长度相同的数组才能互相赋值
var arr1 [5]string
arr2 := [5]string{"Red", "Blue"}
arr1 = arr2

// 复制指针数组，只会复制指针，而不会复制指针所指向的值
var arr1 [3]*string
arr2 := [3]*string(new (string), new (string), new (string))
*arr2[0] = "Red"
*arr2[1] = "Blue"
*arr2[2] = "Green"
// 复制后，两个数组指向同一组字符串
arr1 = arr2

// 申明二维数组
var arr [3][2]int
arr := [2][2]int{{0, 0}, {1, 1}}
// 初始化外层元素
arr := [2][2]int{1: {1, 1}}
// 初始化内层元素
arr := [2][2]int{1: {0: 10}}
```

```Go
/*在函数间传递数组*/
// 声明一个需要8M内存的数组
var arr [1e6]int
// 将数组传递给函数
foo(arr)
func foo(arr [1e6]int)

// 每次foo被调用时，必须在栈上分配8M的内存。
// 之后，整个数组的值被复制到刚分配的内存里。

// 如果复制的是指针而不是值，只需要复制8byte的数据到栈内存上
foo(&arr)
func foo(arr *[1e6]int)
// 但这样也有一个缺点：会改变指针指向的值
// 于是就有了切片
```

### 切片
切片是一种数据结构，这种数据结构便于使用和管理数据集合。切片是围绕动态数组的概念构建的，可以按需自动增长或缩小。切片的动态增长是通过内置函数`append`来实现的。这个函数可以快速且高效地增长切片。还可以通过对切片再次切片来缩小一个切片的大小。

```Go
// 创建和初始化
// 1. make
// 长度和容量都是5
slice := make([]string, 5)
// 长度3，容量5
slice := make([]int, 3, 5)
// 编译错误
slice := make([]int, 5, 3)

// 通过字面量
// 长度和容量都是3
// 跟数组的区别就是 [] 里没数字！！！！
slice := []string{"Red", "Green", "Blue"}
// 创建了长度和容量都是100的切片
slice := []string{99: ""}

// nil切片
var slice []int

// 空切片
slice := make([]int, 0)
slice := []int{}

// 切片的赋值和数组一样
slice[1] = 25

// 切片（动词）
slice := []int{1, 2, 3, 4, 5}
// 创建一个新的切片
// 长度为2，容量为4
newSlice := slice[1:3]
// 索引为3的元素不存在
// Runtime Exception:
// panic: runtime error: index out of range
newSlice[3] = 45

/*
容量为k的切片slice
slice[i:j]
长度：j - i
容量：k - i
*/

// 使用3个索引创建切片
// 长度为1，容量为2
newSlice := slice[2:3:4]

/*
slice[i:j:k]
长度：j - i
容量：k - i
*/

// 将一个切片追加到另一个切片
s1 := []int{1, 2}
s2 := []int{3, 4}

// 使用 ...
s1 = append(s1, s2...)

```

#### 长度和容量相等的好处
append方法会首先使用可用容量。只有在容量不够的情况下，才会分配一个新的底层数组。
这样，在很多时候，就不知道新的切片和老的切片是否在共用一个数组。
如果在创建切片时设置容量和长度一样，就可以强制让新切片的第一个append操作创建新的底层数组，与原有的底层数组分离。

详见 slice.go

#### 迭代切片
```Go
// range创建了元素的副本，而不是对元素的直接引用
// 详见 sliceRange.go
for index, value := range slice {
    fmt.Printf("Index: %d Value: %d\n", index, value)
}

// range 总是从头开始
// 如果需要更多的控制，可以使用传统的for循环
for index := 2; index < len(slice); index++ {
    fmt.Printf("Index: %d Value: %d\n", index, slice[index])
}
```

#### 多维切片
```Go
slice := [][]int{{10}, {100, 200}}

slice[0] = append(slice[0], 20)
```

### Map
#### 创建和初始化
**由于切片/函数以及包含切片的结构类型具有引用语义，不能作为键类型**
```Go
// 使用make
// 创建一个键类型是string，值类型是int的映射
dict := make(map[string]int)

// 使用键值对
dict := map[string]string{"Read":"#da1337", "Orange":"#e95a22"}
```

#### 使用
```Go
colors := map[string]string{}
colors["Red"] = "#da1337"

// nil 映射（未初始化）不能赋值
// 通过声明创建一个nil映射
var colors map[string]string
colors["Red"] = "#da1337" // 编译错误

// 从映射获取值并判断是否存在
value, exists := colors["Blue"]
if exists {
    fmt.Println(value)
}

// 从映射中取值，并通过该值判断其键是否存在
// Go语言里，通过键来索引映射时，即便这个键不存在也会返回一个值
// 返回该类型的“零值”，如string就是""，int就是0
value := colors["Blue"]
if value != "" {
    fmt.Println(value)
}

// 使用range迭代映射
// 返回的是键值对
for key, value := range colors {
    fmt.Println("Key: %s Value: %s\n", key, value)
}

// 从映射中删除一项
// 使用内置函数 delete
delete(colors, "Coral")
```

#### 在函数间的传递
和切片类似，不会创建副本。修改会反应到所有的引用上。详见，map.go

### 小结
- 数组是构造切片和映射的基石
- 切片用来处理数据的集合，映射用来处理键值对结构的数据
- 内置函数make可以创建切片和映射，并指定长度和容量
- 切片有容量限制，不过可以使用内置函数append扩展容量
- 映射的增长没有容量或其他任何限制
- 内置函数len可以用来获取切片或映射的长度
- 内置函数cap只能用于切片
- 切片不能用作映射的键
- 将切片或映射传递给函数成本很小，并且不会复制底层的数据结构

## ch05 类型系统
### 5.1 自定义类型
```Go
// 自定义一个结构类型
type user struct {
    name string
    email string
    ext int
    privileged bool
}

// 声明一个类型变量，其初始值都为“零值”
var bill user

// 初始化所有值
// 顺序没要求
lisa := user{
    name:"Lisa",
    email:"lisa.email.com"
    ext:123,
    privileged:true,
}
// 也可以这样
// 顺序必须与声明中的一样
lisa := user{"Lisa", "lisa@mail.com", 123, true}

// 使用其他结构类型声明字段
// admin 需要一个user类型作为管理者，并附加权限
type admin struct {
    person user
    level string
}

fred := admin {
    person: user {
        name: "Lisa",
        email: "lisa@mail.com",
        ext: 123, 
        privileged: true,
    }, 
    level: "super"
}

// 基于已有的类型定义一个类型
type Duration int64
// 虽然这样定义，但Duration和int64不是同一种类型
// 以下代码会有编译错误
var dur Duration
dur = int64(10000)
```

### 5.2 方法
如果使用值接收者声明方法，调用时会使用这个值得一个副本来执行。
```Go
// list11.go
// 指针变量调用值接收者方法
// 可以认为Go编译器做了如下操作
(*lisa).notify() // 被转换为指针所指向的值

// 值变量调用指针接收者方法
// 编译器背后的动作
(&bill).changeEmail("")
```

### 5.3 类型的本质
#### 5.3.1 内置类型
#### 5.3.2 引用类型
Go语言里有以下几个引用类型：切片/映射/通道/接口和函数类型。
#### 5.3.3 结构类型
```Go
// golang.org/src/os/file.go
// 即使没有修改接收者的值，依然是用指针接收者来声明的。
// 因为File类型的值具备非原始的本质，所以总是应该被共享，而不是被复制

// 是使用值接收者还是指针接收者，不应该由该方法是否修改了入参值来决定。
// 应该基于该类型的本质
func (f *File) Chdir() error {
    if f == nil {
        return ErrInvalid
    }
    if e := syscall.Fchdir(f.fd); e != nil {
        return &PathError("chdir", f.name, e)
    }
    return nil
}
```
### 5.4 接口
### 5.5 嵌入类型
感觉嵌入类型实现了继承功能。

```Go
type user struct {
    name string
    email string
}

type admin struct {
    user // 嵌入类型。注意声明方式的不同
    level string
}

// 内部类型user，外部类型admin
// 可以直接访问内部类型的方法
ad.user.notify()

// 内部类型的方法也可以通过外部类型直接访问
ad.notify()
```

### 5.6 公开或未公开的标识符
- 可以通过工厂函数来创建一个未公开的类型。
- 永远不能显式创建一个未公开的类型的变量，不过`:=`可以这么做

### 5.7 小结
- 使用关键字struct或者通过指定已经存在的类型，可以声明用户自定义的类型
- 方法提供了一个给用户自定义类型增加行为的方式
- 设计类型时需要确认类型的本质是原始的，还是非原始的
- 接口是声明了一组行为并支持多态的类型
- 嵌入类型提供了扩展类型的能力，而无需使用继承
- 标识符要么是从包里公开的，要么是未公开的

## ch06 并发
并发（Concurrency），同时管理很多事情，这些事情可能只做了一半就被暂停去做别的事了。
并行（Parallelism），是让不同的代码片段同时在不同的物理处理器上执行。

### 6.5 通道 channel
```Go
// 无缓冲的整型通道
unbuffered := make(chan int)

// 有缓冲的字符串通道
buffered := make(chan string, 10)

// 向通道发送值
buffered <- "Gopher"

// 从通道接收值
value := <-buffered
```

#### 无缓冲的通道
无缓冲通道（unbuffered channel）是指在接收前没有能力保存任何值的通道。这种类型的通道要求发送者和接收者同时准备好，才能完成发送和接收操作。
如果两个goroutine没有同时准备好，通道会导致先执行发送或者接收操作的goroutine阻塞等待。
这种对通道进行发送和接收的交互行为本身就是同步的。其中任意一个操作都无法离开另一个操作单独存在。
两个有趣的例子：打网球和4x100米接力，这里球和接力棒就是通道。

#### 有缓冲的通道
有缓冲通道（buffered channel）是一种在被接收前能存储一个或多个值的通道。这种类型的通道并不强制要求goroutine之间必须同时完成发送和接收。
只有在通道中没有待接收的值时，接收才会被阻塞。
只有在通道缓冲已满时，发送才会阻塞。

### 6.6 小结
- 竞态是指两个或多个goroutine试图访问同一个资源
- 原子函数和互斥锁提供了一种防止出现竞态的方法
- 通道提供了一种在两个goroutine之间共享数据的简单方法
- 无缓冲通道保证同时交换数据，而有缓冲的没有这种保证

## ch09 测试和性能
### 9.1 单元测试
- 测试文件以`_test.go`结尾。
- 函数签名必须是`func TestXXX(t *testing.T)`
- 执行`go test -v`
- 以`_test`结尾的包，测试代码只能访问包里公开的标识符。即使测试代码文件和被测试的代码放在同一个文件夹里

### 9.2 示例
- `godoc`生成文档时有用
- 函数名必须以 `Example`开头
- 注释 `Output:` 用于测试时的比较
- 执行`go test -v -run="ExampleXXX"`

### 9.3 基准测试
- 签名`func BenchmarkXXX(b *testing.B)`
- 执行`go test -v -run="none" -bench="BechmarkXXX"`
- 默认情况下，基准测试最小运行时间是1秒。可以使用`-benchtime="3s"`
- 有时候，增加基准测试的时间，会得到更加精确的性能结果。对大多数测试来说，超过3s的基准测试并不会改变测试的精确度
- 执行所有的基准测试 `-bench=.`
- `-benchmem`，提供每次操作分配内存的次数，以及总共分配的字节数


```Go
// go test -v -run="none" -bench=. -benchtime="3s" -benchmem
func BenchmarkSprintf(b *testing.B) {
    number := 10

    // 重置计时器
    b.ResetTimer()

    // 为了让基准测试框架能准确测试性能
    // 它必须在一段时间内反复运行要测试的代码
    // 所以，这里使用了for循环
    // b.N 也是固定的
    for i := 0; i < b.N; i++ {
        fmt.Sprintf("%d", number)
    }
}
```

### 9.4 小结
- 示例代码，既能用于测试，也能用于文档

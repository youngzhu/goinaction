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
# KPack 粘包解决方案

### 基准测试

|       goos       |                 windows                  |            |           |              |
| :--------------: | :--------------------------------------: | :--------: | :-------: | :----------: |
|      goarch      |                  amd64                   |            |           |              |
|       pkg        |        github.com/hxoreyer/kpack         |            |           |              |
|       cpu        | Intel(R) Core(TM) i5-9400F CPU @ 2.90GHz |            |           |              |
| BenchmarkKpack-6 |                  232525                  | 4889 ns/op | 6445 B/op | 60 allocs/op |
|       PASS       |                                          |            |           |              |
|        ok        |        github.com/hxoreyer/kpack         |   1.400s   |           |              |

### 使用

```shell
go get -u github.com/hxoreyer/kpack
```



### 例子

首先我们模拟粘包

包:

```go
type Package struct {
	kpack.Kpack
	MsgLen int16 `ksize:"true"` 
	Msg    []byte
	SexLen int16 `ksize:"true"`
	Sex    []byte
	Age    int16
}
```

> 包中如果有slice类型的变量，就得在其上加入其长度变量，并加上tag后在pack阶段会自动填写其对应的slice大小

接下来填包

```go
pack := Package{
    Msg: []byte("现在时间是:" + time.Now().Format("2006-01-02 15:04:05")),
    Sex: []byte("男"),
    Age: 23,
}
pack.Length = int16(kpack.SizeOf(pack))
```

> SizeOf()函数可以计算包的大小

连续写入四次，模拟粘包

```go
buf := new(bytes.Buffer)
for i := 0; i < 4; i++ {
    kpack.Pack(buf, &pack)
}
```

> Pack()函数将包内的各个slice大小填入其大小变量，并将其以大端写入buf中
>
> 这里循环四次，并全部写入同一个buf中实现粘包模拟

开始解包

```go
scanner, size := kpack.ScanPack(buf)
res, _ := kpack.UnPack(scanner, size, func(r io.Reader) interface{} {
    //这里是解包函数
    p := &Package{}
    binary.Read(r, binary.BigEndian, &p.Length)
    binary.Read(r, binary.BigEndian, &p.MsgLen)
    p.Msg = make([]byte, p.MsgLen)
    binary.Read(r, binary.BigEndian, &p.Msg)
    binary.Read(r, binary.BigEndian, &p.SexLen)
    p.Sex = make([]byte, p.SexLen)
    binary.Read(r, binary.BigEndian, &p.Sex)
    binary.Read(r, binary.BigEndian, &p.Age)
    return p
})
```

> ScanPack()将会返回一个scanner和UnPack()用到的缓冲区大小
>
> UnPack()返回解包后的包组和数量
>
> 解包函数是让分割好的数据以你心想的方式写入你的包中

#### 完整代码

```go
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/hxoreyer/kpack"
)

type Package struct {
	kpack.Kpack
	MsgLen int16 `ksize:"true"`
	Msg    []byte
	SexLen int16 `ksize:"true"`
	Sex    []byte
	Age    int16
}

func (p *Package) String() string {
	return fmt.Sprintf("length:%d msglen:%d msg:%s sexlen:%d sex:%s age:%d", p.Length, p.MsgLen, p.Msg, p.SexLen, p.Sex, p.Age)
}

func main() {
    //模拟粘包
	pack := Package{
		Msg: []byte("现在时间是:" + time.Now().Format("2006-01-02 15:04:05")),
		Sex: []byte("男"),
		Age: 23,
	}
	pack.Length = int16(kpack.SizeOf(pack))
	buf := new(bytes.Buffer)
	// 写入四次，模拟TCP粘包效果
	for i := 0; i < 4; i++ {
		kpack.Pack(buf, &pack)
	}
    
    //解包
	scanner, size := kpack.ScanPack(buf)
	res, _ := kpack.UnPack(scanner, size, func(r io.Reader) interface{} {
		p := &Package{}
		binary.Read(r, binary.BigEndian, &p.Length)
		binary.Read(r, binary.BigEndian, &p.MsgLen)
		p.Msg = make([]byte, p.MsgLen)
		binary.Read(r, binary.BigEndian, &p.Msg)
		binary.Read(r, binary.BigEndian, &p.SexLen)
		p.Sex = make([]byte, p.SexLen)
		binary.Read(r, binary.BigEndian, &p.Sex)
		binary.Read(r, binary.BigEndian, &p.Age)
		return p
	})
	
    //输出解包结果
	for _, v := range res {
		fmt.Println(v)
	}
}

```


## 使用GDB调试Go语言##


用Go语言已经有一段时间了,总结一下如何用GDB来调试它!

ps:网上有很多文章都有描述,但是都不是很全面,这里将那些方法汇总一下

### GDB简介

	GDB是GNU开源组织发布的一个强大的UNIX下的程序调试工具。或许，各位比较喜欢那种图形界面方式的，像VC、BCB等IDE的调试，但如果你是在UNIX平台下做软件，你会发现GDB这个调试工具有比VC、BCB的图形化调试器更强大的功能。所谓“寸有所长，尺有所短”就是这个道理。
	
目前支持的语言 (GNU gdb (GDB) 7.8)

进入 gdb之后输入 set language 可以查看支持的语言列表


```
$ gdb
(gdb) set language
Requires an argument. Valid arguments are auto, local, unknown, ada, c, c++, asm, minimal, d, fortran, objective-c, go, java, modula-2, opencl, pascal.
(gdb)
```

### 准备工作

首先看下已经编写好的一个简单的go语言程序

```
➜  go-debug-example  tree
.
├── lib
│   └── calc.go
└── main.go

1 directory, 2 files
```

main.go

```
package main

import (
        "fmt"
        "github.com/beyondblog/go-debug-example/lib"
        "os"
        "runtime"
        //"runtime/debug"
)

func main() {
        var modify string
        argsLen := len(os.Args)
        if argsLen < 2 {
                fmt.Printf("Usage go-debug-example [username] \r\n")
                os.Exit(-1)
        }
        username := os.Args[1]
        var password string
        fmt.Printf("%s welcome!\r\nplease input password:", username)
        fmt.Scanf("%s", &password)
        fmt.Printf("%s  password: %s\r\n", username, password)

        sum := 0
        for i := 0; i < 10; i++ {
                sum += i
                if i == 5 {
                        modify = "modify!"
                }
        }

        fmt.Println(lib.Add(sum, 10))

        runtime.Breakpoint()
        //debug.PrintStack()
        fmt.Println(sum)
        fmt.Println(modify)
}
```

calc.go

```
package lib

func Add(a int, b int) int {
        c := 10
        a = c + b
        return a + b
}
```

程序很简单,就是从命令行获取一个值然后做了写简单的计算,最后输出一下

那么我们编译一下 ^.^

go build 
然后生成了我们要的文件,然后我们执行 gdb go-debug-example

输入 run(简写r)命令运行程序,这个时候可能会有一个提示

```
(gdb) run
Starting program: /Users/****/gopath/src/github.com/beyondblog/go-debug-example/go-debug-example
Unable to find Mach task port for process-id 40056: (os/kern) failure (0x5).
 (please check gdb is codesigned - see taskgated(8))
```
如果你用的是OSX应该就能看到这个,这个提示签名错误,
Darwin kernel出于安全考虑，在没有特殊授权的情况下不允许gdb调试任何程序，因为可以调试就掌握了进程的控制权。不过如果是root用户就没有这个问题，不过谁愿意用root来调试程序呢
解决办法可以通过[这篇文章](http://blog.csdn.net/powerlly/article/details/30323015)来查看
或者启动gdb的时候 加个sudo呗,那么重新来过

sudo gdb go-debug-example
r
然后提示了

```
(gdb) r
Starting program: /Users/****/gopath/src/github.com/beyondblog/go-debug-example/go-debug-example
Usage go-debug-example [username]
[Inferior 1 (process 40111) exited with code 0377]
```

程序成功执行,但是我们那个Go程序提示需要加一个参数才能够继续下去
可以直接 r [参数] 获取在 使用 set args [参数]

ps:可以用 r > file 或者 >> file 支持结果重定向到文件

set args 这个命令是设置参数信息
show args 查看启动的参数信息

既然调试代码就能看到源码啊,使用list(简写l) 命令查看代码执行位置附近10行

```
(gdb) list
1	package main
2
3	import (
4		"fmt"
5		"github.com/beyondblog/go-debug-example/lib"
6		"os"
7		"runtime"
8		//"runtime/debug"
9	)
10
(gdb)
```
默认显示了10行 查看更多的话可以在输入list 或者敲回车(gdb 会默认记住上一个指令然后回车就能继续执行了 -,- )
查看list帮助

```
(gdb) help list
List specified function or line.
With no argument, lists ten more lines after or around previous listing.
"list -" lists the ten lines before a previous ten-line listing.
One argument specifies a line, and ten lines are listed around that line.
Two arguments with comma between specify starting and ending lines to list.
Lines can be specified in these ways:
  LINENUM, to list around that line in current file,
  FILE:LINENUM, to list around that line in that file,
  FUNCTION, to list around beginning of that function,
  FILE:FUNCTION, to distinguish among like-named static functions.
  *ADDRESS, to list around the line containing that address.
With two args if one is empty it stands for ten lines away from the other arg.
(gdb)
```
大概的意思就是

	list 20 	//查看第20行周围的10行

	list - 		//查看上一个list代码之前的10行
	
	list 1,100	//查看1到100行 如果不足100行就显示末尾
	
	list main 	//查看main函数
	执行这个命令的时候会发现显示的不是main函数的信息而是一段汇编代码
	
	(gdb) list main
	9		MOVQ	0(SP), DI // argc
	10		MOVQ	$main(SB), AX
	11		JMP	AX
	12
	13	TEXT main(SB),NOSPLIT,$-8
	14		MOVQ	$runtime·rt0_go(SB), AX
	15		JMP	AX
	(gdb)
	
	大概的意思是go程序的真正入口点应该是这玩意 (ps: $main(SB) SB  = ,=)
	这个时候用 list main.main 即可 意思是查看main.go文件里面main函数附近的10行
	
	list main.main,20 	//查看从main文件中main函数中从函数开始到第20行
	
	l main.go:0 		//以 :方式查看指定文件源码
	l calc.go:0 		//查看calc.go的源码
	l github.com/beyondblog/go-debug-example/lib/calc.go:0 //绝对文件路径查看
	
	还可以搜索代码用 search text 	//可显示在当前文件中包含text串的下一行
	reverse-search text 		//显示包含text 的前一行
	forward-search text 		//不解释
	
现在设置一个断点看看,命令是 break (简写b) 后面的参数和list命令后面的参数大致一样,例如

	b main.main //在main下设置断点
	b 11 		//在11行设置断点
	
	//查看断点
	info breakpoints
	
	(gdb) info breakpoints
	Num     Type           Disp Enb Address            What
	5       breakpoint     keep y   0x0000000000002000 in main.main 	at /Users/****/github.com/	beyondblog/go-debug-example/main.go:11
	6       breakpoint     keep y   0x000000000000203b in main.main 	at /Users/****/github.com/	beyondblog/go-debug-example/main.go:14
	
	能够显示详细的信息,注意这上面有个End 是是否启用的意思可以使用
	disable Num或者 enable Num 来设置断点是否有效
	
	//删除断点
	
	delete(简写d) //不带参数清空所有断点
	d	1		  //删除编号为1的断点
	
	//条件断点 非常实用!
	b Num if [表达式] //在第Num设置断点当满足表达式的条件是触发
	例如
	b 26 if i=5
	
下面让程序运行起来

```
Argument list to give program being debugged when it is started is "".
(gdb) set args beyond
(gdb) show args
Argument list to give program being debugged when it is started is "beyond".
(gdb) b main.main
Breakpoint 1 at 0x2000: file /Users/****/gopath/src/github.com/beyondblog/go-debug-example/main.go, line 11.
(gdb) r
Starting program: /Users/****/gopath/src/github.com/beyondblog/go-debug-example/go-debug-example beyond
[New Thread 0x1617 of process 41584]
[New Thread 0x1803 of process 41584]

Breakpoint 1, main.main () at /Users/****/github.com/beyondblog/go-debug-example/main.go:11
11	func main() {
(gdb) n
12		var modify string
(gdb)
```

	next(简写n) 		//下一步
	step(简写s) 		//单步执行,例如跳进函数内部
	finish			//退出该函数返回到它的调用函数中
	until(简写u)		//直接执行到下一行,如果遇到循环语句,会执行完当前循环
	u Num 			//指哪打哪,继续执行直到Num行时触发断点
	continue(简写c) //从断点开始继续执行
	
	例如
	(gdb) u 22
	beyond welcome!
	please input password:123456
	main.main () at /Users/****/gopath/src/github.com/beyondblog/go-debug-example/main.go:22
	22		fmt.Printf("%s  password: %s\r\n", username, password)
	(gdb)
	
	frame(简写f)	//查看当前命令帧,也就是看当前程序执行到那一行了
	
	info locals //查看当前变量信息
	(gdb) info locals
	sum = 1005232
	&password = 0x2081b4210
	username = 0x7fff5fbffc58 "beyond"
	modify = 0x0 ""
	这个时候会发现,怎么有些变量没显示出来,官方说默认的编译会给调试带来一些不变的优化,可以使用
	go build -gcflags "-N -l" 来关闭这个优化从而方便调试
	那么重新编译后运行在 info local 下
	(gdb) info locals
	sum = 180144
	i = 53967
	argsLen = 2
	&password = 0x2081b4210
	username = 0x7fff5fbffc58 "beyond"
	modify = 0x0 ""
	
	会发现 i 和 sum的值不是一个预期的0,这个是我们程序还没执行到初始化那一块,默认取的一个随机数吧
	ps: 虽然从代码上看好像程序在这个时候还没有声明 i 和 sum,但是最终生成的go语言程序应该是有自己的优化自动声明了
	
	类似的还有很多种命令就不一一介绍了,下面来个逼格高的
	layout src 或者 ctrl x + ctrl a
	
	启动 tui 界面
	或者启动gdb 的时候加上 tui参数
	例如
	gdb -tui
	
![image](https://raw.githubusercontent.com/beyondblog/go-debug-example/master/img/tui.png)

逼格满满有木有哇!

tui有4中窗口模式分别是
	
	command 命令窗口. 可以键入调试命令
	source 源代码窗口. 显示当前行,断点等信息
	assembly 汇编代码窗口
	register 寄存器窗口
	
详细的说明:https://sourceware.org/gdb/current/onlinedocs/gdb/TUI.html#TUI

这个时候想查看源代码可以用list 也可以用方向键

首先先将窗口的焦点设置到源代码窗口上

	focus src
然后就可以用方向键来查看源代码了,这个时候如果想回到command窗口同理只需要

	focus cmd
	
	这个时候你执行命令n 会有一个文本的图形界面还显示

gdb 还有一个特别有用的jump命令,它允许强制的跳转，不会改变栈的结构.
意思就是 如果我的程序运行到第10行的时候我现在想回到第5行调试看看又不想重新运行一边
举个栗子
现在我们执行到第32行,直接输入命令

	u 32

	B+>│32              fmt.Println(lib.Add(sum, 10))
	
	s //单步进入
	
	3	func Add(a int, b int) int {
	(gdb) list
	1	package lib
	2
	3	func Add(a int, b int) int {
	4		c := 10
	5		a = c + b
	6		return a + b
	7	}
	(gdb) n
	4		c := 10
	(gdb) n
	5		a = c + b
	(gdb) f
	#0  github.com/beyondblog/go-debug-example/lib.Add (a=45, b=10, ~r2=0) at /Users/****/github.com/beyondblog/go-debug-example/lib/calc.go:5
	5		a = c + b
	
这个时候已经到了第5行 我现在要回到第4行,jump 4 发现没有用因为jump是跳转到第4行开始执行,但不触发断点所以一般先在要跳转的行号哪儿设置个断点

	github.com/beyondblog/go-debug-example/lib.Add (a=45, b=10, ~r2=8725978992) at /Users/****/github.com/beyondblog/go-debug-example/lib/calc.go:3
	3	func Add(a int, b int) int {
	(gdb) b 4
	Breakpoint 2 at 0x50fed: file /Users/****/github.com/beyondblog/go-debug-example/lib/calc.go, line 4.
	(gdb) n

	Breakpoint 2, github.com/beyondblog/go-debug-example/lib.Add (a=45, b=10, ~r2=0) at /Users/****/github.com/beyondblog/go-debug-example/lib/calc.go:4
	4		c := 10
	(gdb) n
	5		a = c + b
	(gdb) n
	6		return a + b
	(gdb) jump 4
	Continuing at 0x50fed.

	Breakpoint 2, github.com/beyondblog/go-debug-example/lib.Add (a=20, b=10, ~r2=0) at /Users/****/github.com/beyondblog/go-debug-example/lib/calc.go:4
	4		c := 10
	(gdb) n
	5		a = c + b
	(gdb)
	6		return a + b
	(gdb) jump 4
	Continuing at 0x50fed.

	Breakpoint 2, github.com/beyondblog/go-debug-example/lib.Add (a=20, b=10, ~r2=0) at /Users/****/github.com/beyondblog/go-debug-example/lib/calc.go:4
	4		c := 10
	(gdb)
	
	//还可以使用set 命令来设置变量的值 例如
	set b = 10
	
最后附上gdb的官方文档	
https://sourceware.org/gdb/current/onlinedocs/gdb/

这个example的链接
https://github.com/beyondblog/go-debug-example.git
	
----
参考文章:

[0] http://blog.csdn.net/haoel/article/details/2879

[1] http://blog.studygolang.com/2012/12/gdb%E8%B0%83%E8%AF%95go%E7%A8%8B%E5%BA%8F/

[2] http://laokaddk.blog.51cto.com/368606/945057/

[3] http://blog.csdn.net/lwbeyond/article/details/7839225
	
	





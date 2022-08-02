## 目录
```
server.go   一些配置项
user.go   用户的一些操作
```


## v0.5 在线用户查询
检测用户输入 如果是who 则只为当前用户展示在线用户数




## v0.4用户业务层封装
封装一些函数
把之前Server里面user的业务 封装到user里面

**这节课才比较了解面向对象思想了**

```go run .\main.go .\server.go .\user.go```




## v0.3 用户消息广播功能
server.go中有一个go程用来接收用户发送的消息 向用户读
![](https://cdn.jsdelivr.net/gh/Loveyless/img-clouding/img/20220802033312.png)
user.go中有一个go程用来写向用户写
![](https://cdn.jsdelivr.net/gh/Loveyless/img-clouding/img/20220802033540.png)

所以是一个读写分离的模型 没太懂 还需要理解


启动步骤
1. go run .\main.go .\server.go .\user.go 或者 build
2. 开多个命令行来测试上线 telnet localhost 8888
3. 随便输入数据

问题 我会给我自己也广播出来数据 不过这个应该好解决
```msg := string(buf[:n-1])```这个和```msg := string(buf[:len(buf)-1])```是不是没区别？
不过好像n不行

win只能输入一个字符 尝试检测用户按回车 但是失败了 写不来


## v0.2 用户上线以及广播功能
![](https://cdn.jsdelivr.net/gh/Loveyless/img-clouding/img/20220802014908.png)

增加了根据conn创建用户
增加了上线提示
还有锁的概念
具体看代码

我觉得面向对象也不难就是太绕了 但是一个一个理解还是可以的

启动步骤
1. go run .\main.go .\server.go .\user.go 或者 build
2. 开多个命令行来测试上线 telnet localhost 8888


## v0.1 基础server构建
启动步骤
1. go build -o server.exe main.go server.go 或者 go run main.go server.go
2. 如果build了就 ./server.exe打开文件 run就无需操作
3. 新终端telnet localhost 8888

基本上就是创建一个server.go里面封装了一个类
里面有个方法会导出类创建的实例 然后实例里有start方法 里面会自动调用handler方法



命令行编译
```
//直接编译当前目录所有文件
go build -o server.exe .



go build [-o 输出名] [-i] [编译标记] [包名]

-o output 指定编译输出的名称，代替默认的包名。
-i install 安装作为目标的依赖关系的包(用于增量编译提速)。


如果参数为***.go文件或文件列表，则编译为一个个单独的包。
当编译单个main包（文件），则生成可执行文件。
当编译单个或多个包非主包时，只构建编译包，但丢弃生成的对象（.a），仅用作检查包可以构建。
当编译包时，会自动忽略’_test.go’的测试文件。
```
其他命令
以下 build 参数可用在 build, clean, get, install, list, run, test
```
-a
    完全编译，不理会-i产生的.a文件(文件会比不带-a的编译出来要大？)
-n
    仅打印输出build需要的命令，不执行build动作（少用）。
-p n
    开多少核cpu来并行编译，默认为本机CPU核数（少用）。
-race
    同时检测数据竞争状态，只支持 linux/amd64, freebsd/amd64, darwin/amd64 和 windows/amd64.
-msan
    启用与内存消毒器的互操作。仅支持linux / amd64，并且只用Clang / LLVM作为主机C编译器（少用）。
-v
    打印出被编译的包名（少用）.
-work
    打印临时工作目录的名称，并在退出时不删除它（少用）。
-x
    同时打印输出执行的命令名（-n）（少用）.
-asmflags 'flag list'
    传递每个go工具asm调用的参数（少用）
-buildmode mode
    编译模式（少用）
    'go help buildmode'
-compiler name
    使用的编译器 == runtime.Compiler
    (gccgo or gc)（少用）.
-gccgoflags 'arg list'
    gccgo 编译/链接器参数（少用）
-gcflags 'arg list'
    垃圾回收参数（少用）.
-installsuffix suffix
    ？？？？？？不明白
    a suffix to use in the name of the package installation directory,
    in order to keep output separate from default builds.
    If using the -race flag, the install suffix is automatically set to race
    or, if set explicitly, has _race appended to it.  Likewise for the -msan
    flag.  Using a -buildmode option that requires non-default compile flags
    has a similar effect.
-ldflags 'flag list'
    '-s -w': 压缩编译后的体积
    -s: 去掉符号表
    -w: 去掉调试信息，不能gdb调试了
-linkshared
    链接到以前使用创建的共享库
    -buildmode=shared.
-pkgdir dir
    从指定位置，而不是通常的位置安装和加载所有软件包。例如，当使用非标准配置构建时，使用-pkgdir将生成的包保留在单独的位置。
-tags 'tag list'
    构建出带tag的版本.
-toolexec 'cmd args'
    ？？？？？？不明白
    a program to use to invoke toolchain programs like vet and asm.
    For example, instead of running asm, the go command will run
    'cmd args /path/to/asm <arguments for asm>'.
```
以上命令，单引号/双引号均可。

对包的操作'go help packages'
对路径的描述'go help gopath'
对 C/C++ 的互操作'go help c'
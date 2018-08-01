package reexec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var registeredInitializers = make(map[string]func()) //维护Initializers的包内全局map

// Register adds an initialization func under the specified name
func Register(name string, initializer func()) {
	if _, exists := registeredInitializers[name]; exists {
		panic(fmt.Sprintf("reexec func already registred under name %q", name))
	}

	registeredInitializers[name] = initializer
}

// Init is called as the first part of the exec process and returns true if an
// initialization function was called.

//docker启动时如果调用初始化函数，则会导致该函数返回true,从而导致主进程直接返回，仅仅执行初始化工作
//否则，返回false,继续执行docker主进程
func Init() bool {
	initializer, exists := registeredInitializers[os.Args[0]]
	if exists {
		initializer()

		return true
	}

	return false
}

// Self returns the path to the current processes binary
func Self() string {
	name := os.Args[0]
	//filepath.Dir(),filepath.Base()
	//可以将一个路径分解为目录和文件名两部分
	
	// 获取 path 中最后一个分隔符之前的部分（不包含分隔符）
	//Dir(path string) string
	// 获取 path 中最后一个分隔符之后的部分（不包含分隔符）
	//Base(path string) string
	
	
	//func LookPath(file string) (string, error) 
	//LookPath在环境变量中查找可执行二进制文件，如果file中包含一个斜杠，
	//则直接根据绝对路径或者相对本目录的相对路径去查找
	if filepath.Base(name) == name {
		if lp, err := exec.LookPath(name); err == nil {
			name = lp//这里的lp是找到的可执行文件的绝对路径
		}
	}

	return name
}

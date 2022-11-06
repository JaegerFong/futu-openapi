package module

import "sync"

var (
	logicInit = new(sync.Once)
)

func DoInit() {
	logicInit.Do(func() {
		// 初始化futu连接
		initFutuConnect()
	})
}

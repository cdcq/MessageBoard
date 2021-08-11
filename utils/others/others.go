package others

var PowerName = [4]string{
	"普通用户",
	"管理员",
	"超级管理员",
	"创始人",
}

func Min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

// パッケージ名はディレクトリと同じものにする必要があります。
package greeting

import (
	"fmt"
)

// 他のディレクトリ(パッケージ)から呼び出される関数は大文字から始まる必要があります。
// 意外と間違えやすいポイントなので注意しましょう。
func SayHello() {
	fmt.Println("Hello Golang!!")
	fmt.Println("V6だよ!!")
}

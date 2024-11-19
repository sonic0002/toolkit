package example

import (
	"fmt"
	"os"

	"github.com/ozgio/strutil"

	"importtidier/example/base"
)

func thirdWithLocal() {
	fmt.Println("Hello world")

	strutil.WordWrap("Lorem ipsum dolor sit amet", 15, false)

	fmt.Printf("min: %d\n", base.Min(3, 7))

	os.Exit(0)
}

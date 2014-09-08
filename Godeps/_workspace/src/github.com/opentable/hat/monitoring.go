package hat

import (
	"fmt"
	"time"
)

func now() time.Time {
	return time.Now()
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	println(fmt.Sprintf("%s took %s", name, elapsed))
}

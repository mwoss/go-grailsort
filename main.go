package wikisort

import (
	"log"
	"math/rand"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func generateNumbers() []float64{
	var values []float64
	for i:=0; i < 32000000; i++ {
		values = append(values, rand.Float64())
	}
	return values
}

func loopOver(values []float64) {
	defer timeTrack(time.Now(), "loop")
	sum := 0.0
	for _,v := range values[:] {
		sum = sum +  v*v
	}
}

func loopOver2(values []float64){
	defer timeTrack(time.Now(), "loop")
	sum := 0.0
	for i := 0; i < len(values); i++ {
		x := values[i]
		sum = sum +  x*x
	}
}

func main() {
	values := generateNumbers()
	//loopOver(values)
	loopOver2(values)
}

Timed queue
---

Done by some article, i've made some minor improvements.

the basic idea is to plan execution of arbitrary closures in future.

```go
import(
	"github.com/thefish/scheduleq"
	"time"
)

func main() {
	//queue init
	q := scheduleq.Newqueue()
	//plan a phrase printed to stdout in 1400 milliseconds
	q.Plan(func(){
		fmt.Println("trololo")
	}, 1400)
	
	//then some time passes
	time.Sleep(time.Second)
	//but nothing happens
	q.Advance(time.Now())
	//then some mmore time passes
	time.Sleep(time.Second)
	//and we get our message printed! 
	q.Advance(time.Now())
}

```

Useful to process events in time quants, allowing us to fall behind to some degree if payload in functions took some 
more time than was planned. There are some gotchas with scopes, **see test for details** and take note. 

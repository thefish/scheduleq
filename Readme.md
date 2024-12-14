Timed queue
---

Done by some article, i've made some minor improvements.

The basic idea is to plan execution of arbitrary closures in future, and have zero dependencies.

You can throttle planning by limiting number of functions being executed per period and set the period.

### Example:

```go
import(
	"github.com/thefish/scheduleq"
	"time"
)

func main() {
	//queue init
	q := scheduleq.Newqueue(0,0,0)
	task := NewTask(func() error {
            fmt.Println("trololo")
            return nil
	})
	//plan a phrase printed to stdout in 1400 milliseconds
	q.Schedule(task, time.Now().Add(1400*time.Millisecond))
	
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
more time than was planned. **More details in** [tests](scheduleq_test.go). 

It is not a distributed queue. Use proper message broker or database if you need one.


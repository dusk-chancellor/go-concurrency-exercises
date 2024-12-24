//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated) - done
//

package main

import (
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	ch := make(chan bool)
	startTime := time.Now()
	go func() {
		process()
		elapsed := time.Since(startTime).Seconds()
		u.TimeUsed += int64(elapsed)
		
		if !u.IsPremium && u.TimeUsed > 10 {
			ch <- false
		} else {
			ch <- true
		}
		close(ch)
	}()
	return <-ch
}

func main() {
	RunMockServer()
}

/*`$ go run .`
UserID: 0       Process 1 started.
UserID: 1       Process 2 started.
UserID: 0       Process 3 started.
UserID: 0       Process 4 started.
UserID: 1       Process 5 started.
UserID: 0       Process 1 done.
UserID: 0       Process 3 killed. (No quota left)
UserID: 1       Process 5 done.
UserID: 1       Process 2 done.
UserID: 0       Process 4 killed. (No quota left)
*/

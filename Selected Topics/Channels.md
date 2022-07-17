# Introduction 
A channel is a mechanism to communicate between different go routines. To create a channel, we can use the following statement:
```go
ch := make(chan int)
```
Note that a channel is a **reference** to the data structure created by `make()`, which is `ch` in this statement.
Two channels can be compared if they are in the same type, the expression returns true if they have the same one underlying data structure(i.e. refer to the same data structure).
# Unbuffered Channels
Like **Producer-Consumer** problem in OS, channels can be divided into two categories based on capacity. If the capacity is non-zero, `make` creates Buffered channel.
A send operation on an unbuffered channel blocks the sending goroutine until another go routine executes a corresponding receive on the same channel, which means that we can use this property as a **synchronize** approach. Thus, unbuffered channel is also called Synchronized channel. 
# Buffered Channel 
Buffered channel is also called Asynchronized channel, which is the classical producer-consumer problem.
Adding a second argument to `make` gives us a Buffered channel with capacity.
```go
ch := make(chan int, 3)    //Buffered channel with capacity 3
```
In go source code, the core structure that implements go channel is the struct `hchan` :
```go
type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters
	lock     mutex
}
```
- `buf` points to an circular queue.
- `recvq` and `sendq` are both lists that stores blocked go routines trying to receive/send data from/to channel due to limited capacity.
- `lock` guarantees receiving and sending are atomic operations.

If the channel is neither full nor empty, either send or receive operation could proceed without blocking.
- When we execute a receive operation on an empty channel, the goroutine will be put into blocking state and added to `recvx` linked list. When a new go routine performs a send operation on the channel, the data will be directly transferred to the first waiting goroutine. Then it will be removed from `recvq` and Go scheduler will set the blocked routine runnable.
- Almost the same operation when sending data to a full channel.
```go
type waitq struct {
	first *sudog
	last  *sudog
}
```
`waitq`is a struct which contains two pointers point to the head and the tail of linked list respectively. `sudog` contains information of goroutine, and all of these are necessary for the blocked goroutine to be awaken later.
***
Note that goroutine is not OS thread, **what blocks here is goroutine not the OS thread!**.Go uses its own scheduler that is analogous to kernel scheduler, but go scheduler is not rely on hardware timer interrupting the running thread to make context switch. It implements  scheduling by some go structs. No underlying context switch could reduce the cost of scheduling, and this is much cheaper than rescheduling a thread. 
***
```go
type sudog struct {
	// The following fields are protected by the hchan.lock of the
	// channel this sudog is blocking on. shrinkstack depends on
	// this for sudogs involved in channel ops.

	g *g

	next *sudog
	prev *sudog
	elem unsafe.Pointer // data element (may point to stack)

	// The following fields are never accessed concurrently.
	// For channels, waitlink is only accessed by g.
	// For semaphores, all fields (including the ones above)
	// are only accessed when holding a semaRoot lock.

	acquiretime int64
	releasetime int64
	ticket      uint32

	// isSelect indicates g is participating in a select, so
	// g.selectDone must be CAS'd to win the wake-up race.
	isSelect bool

	// success indicates whether communication over channel c
	// succeeded. It is true if the goroutine was awoken because a
	// value was delivered over channel c, and false if awoken
	// because c was closed.
	success bool

	parent   *sudog // semaRoot binary tree
	waitlink *sudog // g.waiting list or semaRoot
	waittail *sudog // semaRoot
	c        *hchan // channel
}
```
*** 
# How can we exploit the CPU performance to increase the speed? 
A typical answer would be using parallelism on problems that can be divided into independent smaller subproblems. Consider the following code trying to use gotoutine to speed up image translation from normal size to thumbnail : 
```go
func makeThumbnails(filenames []string) {
    for _, f := range filenames {
    go thumbnail.ImageFile(f) // NOTE: ignoring errors
    }
}
```
This program runs pretty fast, actually too fast...
`makeThumbnail` returns before all gotoutines finish. If you know something about Linux programming, you would probably know that a simple `pthread_join` syscall would solve the problem. However, there's no such thing that is similar to `wait` in go, and go routine is not OS thread. So, we can change the inner go routine to notify its completion by sending an event to outer goroutine on a shared channel. 
**This leads to the third usage of channel in go** i.e. channel as semaphore
# channel as semaphore
channel type is an anonymous struct.
```go
    sema := make(chan struct{})
```
Therefore, the revised version of `makeThumbnail` is :
```go
func makeThumbnails(filenames []string) {
    ch := make(chan struct{})
    for _, f := range filenames {
    go func(f string) {
        thumbnail.ImageFile(f) // NOTE: ignoring errors
        ch <- struct{}{}
        }(f)
    }
    for range filenames {
        <-ch
    }

}
```
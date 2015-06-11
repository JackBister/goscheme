#goscheme
This is a Lisp-y language interpreter implemented in Go. It also contains some Go concepts, like channels and goroutines for concurrent execution.

##Examples
    (define c (go (+ 1 1)))

Creates a new goroutine which calculates the value of 1+1. The result is sent on the channel c.

    (define r (-> c))

Receives the value from the earlier calculation and stores the value in r.

    (define c (chan))
    (define ex (lambda (x) (begin (<- c x) (ex (+ x 1)))))
    (go (ex 0))

Defines and executes a recursive function that runs forever and repeatedly sends numbers on the channel c. Channels are blocking, so the function will wait until another routine is ready to receive before continuing with the next iteration.

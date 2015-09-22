(define abs (lambda (x) (if (< x 0) (- x) x)))

(define expt (lambda (x y) (if (= y 0) 1 (* x (expt x (- y 1))))))


(define even? (lambda (x) (= (remainder x 2) 0)))

(define negative? (lambda (x) (< x 0)))

(define odd? (lambda (x) (not (even? x))))

(define positive? (lambda (x) (> x 0)))

(define zero? (lambda (x) (= x 0)))

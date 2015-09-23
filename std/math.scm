(define abs (lambda (x) (if (< x 0) (- x) x)))

(define expt (lambda (x y) (if (= y 0) 1 (* x (expt x (- y 1))))))

(define max (lambda (x y . z)
	(begin
	  (define li (cons y z))
	  (fold-left (lambda (v w) (if (> v w) v w)) x li))))

(define min (lambda (x y . z)
	(begin
	  (define li (cons y z))
	  (fold-left (lambda (v w) (if (< v w) v w)) x li))))

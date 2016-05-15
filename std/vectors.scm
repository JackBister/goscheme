(define vector (lambda li
  (begin
    (define ret (make-vector (length li)))
    (define i 0)
    (map (lambda (x)
	      (begin
		(vector-set! ret i x)
		(set! i (+ i 1))
		)) li)
    ret)))

(define vector-fill! (lambda (v f)
  (if (not (vector? v)) (error "vector-fill: Argument 1 is not a vector.")
    (begin
      (define fillVec (lambda (k)
			(if (= k (vector-length v)) v (begin (vector-set! v k f) (fillVec (+ k 1))))))
      (fillVec 0)))))

(define vector->list (lambda (v)
  (if (not (vector? v)) (error "vector->list: Argument 1 is not a vector.")
    (begin
      (define makelist (lambda (k)
			 (if (= k (vector-length v)) '() (apply list (vector-ref v k) (makelist (+ k 1))))))
      (makelist 0)))))


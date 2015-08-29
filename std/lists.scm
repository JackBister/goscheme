;Because map runs in order in this implementation, map and for-each are equivalent.
(define for-each map)

(define member (lambda (obj li) (if (eqv? li (list)) #f (if (equal? obj (car li)) li (member obj (cdr li))))))

(define memq (lambda (obj li) (if (eqv? li (list)) #f (if (eq? obj (car li)) li (memq obj (cdr li))))))

(define memv (lambda (obj li) (if (eqv? li (list)) #f (if (eqv? obj (car li)) li (memv obj (cdr li))))))


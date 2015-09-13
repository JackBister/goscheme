(define concatenate (lambda (li) (apply append li)))

(define empty? (lambda (li) (eqv? li (list))))

(define filter (lambda (pred li) (
	if (empty? li) li (
		if (pred (car li)) (append (list (car li)) (filter pred (cdr li)))
			(filter pred (cdr li))))))

(define filter-map (lambda (f li)
	(filter (lambda (x) (not (not x))) (map f li))))

(define fold-left (lambda (f z li)
	(if (empty? li) z (fold-left f (f (car li) z) (cdr li)))))

(define fold-right (lambda (f z li)
	(if (empty? li) z (f (car li) (fold-right f z (cdr li))))))

(define last (lambda (li)
	(if (= (length li) 1) (car li) (last (cdr li)))))

(define list-tail (lambda (li k) (if (zero? k) li (list-tail (cdr li) (- k 1)))))

(define list-ref (lambda (li k) (car (list-tail li k))))

(define map (lambda (f . lis)
	(begin
	  ;smap is a non-variadic map function, it is necessary for the implementation of the variadic version,
	  ;but it is not necessary outside this context.
	  (define smap (lambda (f li) (if (empty? li) (list) (append (list (f (car li))) (smap f (cdr li))))))
	  ;if any of the lists passed as arguments is empty, return the empty list.
	  (if (some? (lambda (x) (eqv? #t x)) (smap empty? lis)) (list)
	  ;otherwise, apply f to the head of the lists, and append the result to the result of mapping on the tails.
	  (append (list (apply f (smap car lis))) (apply map f (smap cdr lis)))))))

(define member (lambda (obj li) (if (empty? li) (list) #f (if (equal? obj (car li)) li (member obj (cdr li))))))

(define memq (lambda (obj li) (if (empty? li) #f (if (eq? obj (car li)) li (memq obj (cdr li))))))

(define memv (lambda (obj li) (if (empty? li) #f (if (eqv? obj (car li)) li (memv obj (cdr li))))))

(define remove (lambda (pred li) (
	if (empty? li) li (
		if (pred (car li)) (remove pred (cdr li)) (
			append (list (car li)) (remove pred (cdr li)))))))

(define reverse (lambda (li) (if (empty? li) (list) (append (reverse (cdr li)) (list (car li))))))

(define some? (lambda (pred li) (if (empty? li) #f (if (pred (car li)) #t (some? pred (cdr li))))))

(define zip (lambda lis (apply map list lis)))

;Because map runs in order in this implementation, map and for-each are equivalent.
(define for-each map)
(define concatenate (lambda (li) (apply append li)))

(define empty? (lambda (li) (eqv? li (list))))

(define filter (lambda (pred li) (
	if (empty? li) li (
		if (pred (car li)) (cons (car li) (filter pred (cdr li)))
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
	  (define smap (lambda (f li) (if (empty? li) (list) (cons (f (car li)) (smap f (cdr li))))))
	  ;if any of the lists passed as arguments is empty, return the empty list.
	  (if (some? (lambda (x) (eqv? #t x)) (smap empty? lis)) (list)
	  ;otherwise, apply f to the head of the lists, and append the result to the result of mapping on the tails.
	  (cons (apply f (smap car lis)) (apply map f (smap cdr lis)))))))

(define member (lambda (obj li) (if (empty? li) (list) #f (if (equal? obj (car li)) li (member obj (cdr li))))))

(define memq (lambda (obj li) (if (empty? li) #f (if (eq? obj (car li)) li (memq obj (cdr li))))))

(define memv (lambda (obj li) (if (empty? li) #f (if (eqv? obj (car li)) li (memv obj (cdr li))))))

(define remove (lambda (pred li)
	(if (empty? li) li
	  (if (pred (car li))
	    (remove pred (cdr li))
	    (cons (car li) (remove pred (cdr li)))))))

(define reverse (lambda (li) (if (empty? li) (list) (append (reverse (cdr li)) (list (car li))))))

(define some? (lambda (pred li) (if (empty? li) #f (if (pred (car li)) #t (some? pred (cdr li))))))

;Splits a list on any member satisfying pred.
;Returns a list of sublists that are contained between those members.
(define split (lambda (pred li)
	(if (empty? li) '()
	  (begin
	    ;split-help returns the first index that satisfies pred
	    ;it's pretty similar to list-index in SRFI-1 but returns (length li) instead of false
	    ;if nothing in the list satisfies pred
	    (define split-help (lambda (pred li k)
				 (if (empty? li) k
				   (if (pred (car li)) k (split-help pred (cdr li) (+ k 1))))))
	    ;findex is the index of the first list member satisfying pred
	    (define findex (split-help pred li 0))
	    ;take findex elements from list, append the resulting list to the result of calling split
	    ;on the rest of the list after the element at findex.
	    (cons (take li findex) (split pred (cdr (list-tail li findex))))))))

(define take (lambda (li k)
	;if the list is empty but we're not done taking, return an error
	(if (and (empty? li) (> k 0))
	  (error "take: Attempt to take more than length of list.")
	  ;otherwise if k is 0, return an empty list and stop recursing
	  (if (= k 0)
	    '()
	    (begin
	      (define lt (take (cdr li) (- k 1)))
	      ;if taking on the remaining tail returns an error, return that error.
	      (if (error? lt) lt (cons (car li) lt)))))))

(define zip (lambda lis (apply map list lis)))

;Because map runs in order in this implementation, map and for-each are equivalent.
(define for-each map)

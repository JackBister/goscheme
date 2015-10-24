;**x** A list to append to.
;**y** An element to append to the list. Can be another list.
;**z** The function is variadic so that any further arguments will be appended.
;Examples:
;`(append '(1 2 3) 4) => (1 2 3 4)`
;`(append '(1 2 3) '(4 5 6)) => (1 2 3 4 5 6)`
;`(append '(1 2 3) 4 '(5 6)) => (1 2 3 4 5 6)`
(define append (lambda (x y . z)
	(begin
	  (define sappend (lambda (li1 li2)
		(if (null? li1) li2 (if (null? li2) li1 (cons (car li1) (sappend (cdr li1) li2))))))
	  (if (null? z) (sappend x y)
	  (apply append (sappend x y) z)))))

;append two lists using a vector
;Should be faster than the sappend used in append right now, but buggy
;(define sappend (lambda (li1 li2)
;	(begin
;	  (define ret (make-vector (+ (length li1) (length li2))))
;	  (define tovec (lambda (vec li i)
;		(if (< i (vector-length vec))
;		(begin (vector-set! vec i (car li)) (tovec vec (cdr li) (+ i 1))))))
;	  (tovec ret li1 0)
;	  (tovec ret li2 (length li1))
;	  (vector->list ret))))

;**li** A list of lists to concatenate.
;Concatenates all lists in li into one list.
;Example: `(concatenate '((1 2 3) (4 5 6))) => (1 2 3 4 5 6)`
(define concatenate (lambda (li) (apply append li)))

;**pred** A predicate which is sought after in the list.
;**li** A list to filter.
;Returns all elements of li which satisfy pred.
;Example: `(filter number? '(1 2 a b c 3 d 4 e)) => (1 2 3 4)`
(define filter (lambda (pred li) (
	if (null? li) li (
		if (pred (car li)) (cons (car li) (filter pred (cdr li)))
			(filter pred (cdr li))))))

;**f** A function to map over the list.
;**li** The list to map over.
(define filter-map (lambda (f li)
	;TODO: uh...
	(filter (lambda (x) (not (not x))) (map f li))))

;**f** Function by which to fold the list.
;**z** The accumulator's starting value.
;**li** The list to fold.
;Folds left over the list. For an explanation that makes sense, see https://en.wikipedia.org/wiki/Fold_(higher-order_function)
;Examples:
;`(fold-left + 0 '(1 2 3)) => 6`
;`(fold-left cons '() '(1 2 3)) => (3 2 1)`
(define fold-left (lambda (f z li)
	(if (null? li) z (fold-left f (f (car li) z) (cdr li)))))


;**f** Function by which to fold the list.
;**z** The accumulator's starting value.
;**li** The list to fold.
;Folds right over the list. For an explanation that makes sense, see https://en.wikipedia.org/wiki/Fold_(higher-order_function)
;Examples:
;`(fold-right cons '() '(1 2 3)) => (1 2 3)`
(define fold-right (lambda (f z li)
	(if (null? li) z (f (car li) (fold-right f z (cdr li))))))

;**sep** The separator to insert between the lists.
;**lis** Any number of lists to join together with the separator.
;Variadic function which joins all its arguments with the separator inserted between them.
;Example: `(join 'a '(1 2 3) '(4 5 6) '(7 8 9)) => (1 2 3 a 4 5 6 a 7 8 9)`
(define join (lambda (e . lis)
	(fold-right (lambda (x y) (if (null? y) x (append x e y))) '() lis)))
	
;**li** A list to return the last element from.
;Returns the last element of li.
;Example: `(last '(1 2 3)) => 3`
(define last (lambda (li)
	(if (= (length li) 1) (car li) (last (cdr li)))))

;**x** A list to return the length of.
;Returns the length of x.
;Example: `(length '(1 2 3)) => 3`
(define length (lambda (x)
	(fold-right (lambda (_ y) (+ y 1)) 0 x)))

;**li** A list to get the tail from.
;**k** The index to start the tail at (0-based)
;Takes the remaining elements of li starting at position k.
;Examples:
;`(list-tail '(1 2 3 4) 0) => (1 2 3 4)`
;`(list-tail '(1 2 3 4) 1) => (2 3 4)`
;`(list-tail '(1 2 3 4) 4) => ()`
(define list-tail (lambda (li k) (if (zero? k) li (list-tail (cdr li) (- k 1)))))

;**li** List to get an element from.
;**k** Index of the element to get (0-based)
;Returns the *k*th element of li.
;Example: `(list-ref '(1 2 3) 1) => 2`
(define list-ref (lambda (li k) (car (list-tail li k))))

;**f** A function to apply to the lists. The arity of the function must be equal to the number of lists.
;**lis** (variadic) any number of lists to apply the function to.
;Returns a list of the same length as the shortest argument containing the result of applying f to each element of the lists.
;Examples:
;`(map number? '(1 a #t)) => (#t #f #f)`
;`(map + '(1 2 3) '(4 5 6 7) '(1 1)) => (6 8)`
(define map (lambda (f . lis)
;TODO: (map +) => freeze
	(begin
	  ;smap is a non-variadic map function, it is necessary for the implementation of the variadic version,
	  ;but it is not necessary outside this context.
	  (define smap (lambda (f li) (if (null? li) (list) (cons (f (car li)) (smap f (cdr li))))))
	  ;if any of the lists passed as arguments is empty, return the empty list.
	  (if (some? (lambda (x) (eqv? #t x)) (smap null? lis)) (list)
	  ;otherwise, apply f to the head of the lists, and append the result to the result of mapping on the tails.
	  (cons (apply f (smap car lis)) (apply map f (smap cdr lis)))))))

;**obj** Object to search for in li.
;**li** List to search for obj.
;Searches for an object equal to obj in list li. If it is found, the object and the remainder of the list after it is returned. Uses 'equal?' to check for equality. If the object is not found, returns #f.
;TODO: remove (list) causing it to return false all the time.
;Examples:
;`(member 2 '(1 2 3 4)) => (2 3 4)`
;`(member 5 '(1 2 3 4)) => #f`
(define member (lambda (obj li) (if (null? li) (list) #f (if (equal? obj (car li)) li (member obj (cdr li))))))

;**obj** Object to search for in li.
;**li** List to search for obj.
;See member. This function uses eq? instead of equal?.
;Example: `(memq '(1 2) '(3 4 (1 2) a b)) => ((1 2) a b)`
(define memq (lambda (obj li) (if (null? li) #f (if (eq? obj (car li)) li (memq obj (cdr li))))))

;**obj** Object to search for in li.
;**li** List to search for obj.
;See member. This function uses eqv? instead of equal?.
(define memv (lambda (obj li) (if (null? li) #f (if (eqv? obj (car li)) li (memv obj (cdr li))))))

;**x** A list to check.
;Returns #t if x is an empty list, otherwise returns #f.
;Examples:
;`(null? '()) => #t`
;`(null? 1) => #f`
(define null? (lambda (x) (eqv? '() x)))

;**pred** A predicate. Any element of li which satisfies it will be removed from the list.
;**li** The list from which the elements should be removed.
;Returns li with all elements which satisfy pred removed.
;Example: `(remove number? '(a 1 2 b c 3 4 d)) => (a b c d)`
(define remove (lambda (pred li)
	(if (null? li) li
	  (if (pred (car li))
	    (remove pred (cdr li))
	    (cons (car li) (remove pred (cdr li)))))))

;**li** A list which will be reversed.
;Returns the list with its elements in reverse order.
;Example: `(reverse '(1 2 3)) => (3 2 1)`
(define reverse (lambda (li) (if (null? li) (list) (append (reverse (cdr li)) (list (car li))))))

;**pred** A predicate against which to check the elements of li.
;**li** A list to check.
;Returns #t if any element of li satisfies pred, #f otherwise.
;Examples:
;`(some? number? '(a b c)) => #f`
;`(some? number? '(a b 1)) => #t`
(define some? (lambda (pred li) (if (null? li) #f (if (pred (car li)) #t (some? pred (cdr li))))))

;**pred** The predicate to split the list on.
;**li** The list to split.
;Splits a list on any member satisfying pred.
;Returns a list of sublists that are contained between those members.
;Example: `(split symbol? '(1 2 3 a 4 5 6 a 7 8 9)) => ((1 2 3) (4 5 6) (7 8 9))`
(define split (lambda (pred li)
	(if (null? li) '()
	  (begin
	    ;split-help returns the first index that satisfies pred
	    ;it's pretty similar to list-index in SRFI-1 but returns (length li) instead of false
	    ;if nothing in the list satisfies pred
	    (define split-help (lambda (pred li k)
				 (if (null? li) k
				   (if (pred (car li)) k (split-help pred (cdr li) (+ k 1))))))
	    ;findex is the index of the first list member satisfying pred
	    (define findex (split-help pred li 0))
	    ;take findex elements from list, append the resulting list to the result of calling split
	    ;on the rest of the list after the element at findex.
	    (cons (take li findex) (split pred (cdr (list-tail li findex))))))))

;**li** A list.
;**k** The number of elements to take from li.
;Returns the first k elements of li.
;Examples:
;`(take '(1 2 3) 0) => ()`
;`(take '(1 2 3) 2) => (1 2)`
;`(take '(1 2 3) 4) => Error`
(define take (lambda (li k)
	;if the list is empty but we're not done taking, return an error
	(if (and (null? li) (> k 0))
	  (error "take: Attempt to take more than length of list.")
	  ;otherwise if k is 0, return an empty list and stop recursing
	  (if (= k 0)
	    '()
	    (begin
	      (define lt (take (cdr li) (- k 1)))
	      ;if taking on the remaining tail returns an error, return that error.
	      (if (error? lt) lt (cons (car li) lt)))))))

;**lis** This is a variadic function that takes any number of lists.
;Returns the result of zipping the given lists.
;Examples:
;`(zip '(1 2 3) '(a b c)) => ((1 a) (2 b) (3 c))`
;`(zip '(1 a one) '(2 b two) '(3 c three)) => ((1 2 3) (a b c) (one two three))`
(define zip (lambda lis (apply map list lis)))

;Because map runs in order in this implementation, map and for-each are equivalent.
(define for-each map)

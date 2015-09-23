(define char=? (lambda (x y)
	(if (not (and (char? x) (char? y))) (error "char=?: Argument is not a char.")
	  (= (char->integer x) (char->integer y)))))

(define char<? (lambda (x y)
	(if (not (and (char? x) (char? y))) (error "char<?: Argument is not a char.")
	  (< (char->integer x) (char->integer y)))))

(define char>? (lambda (x y)
	(if (not (and (char? x) (char? y))) (error "char>?: Argument is not a char.")
	  (> (char->integer x) (char->integer y)))))

(define char<=? (lambda (x y)
	(if (not (and (char? x) (char? y))) (error "char<=?: Argument is not a char.")
	  (<= (char->integer x) (char->integer y)))))

(define char>=? (lambda (x y)
	(if (not (and (char? x) (char? y))) (error "char>=?: Argument is not a char.")
	  (>= (char->integer x) (char->integer y)))))

(define char-ci=? (lambda (x y)
	(if (not (and (char? x) (char? y))) (error "char-ci=?: Argument is not a char.")
	  (= (char->integer (char-downcase x)) (char->integer (char-downcase y))))))

(define char-ci<? (lambda (x y)
	(if (not (and (char? x) (char? y))) (error "char-ci<?: Argument is not a char.")
	  (< (char->integer (char-downcase x)) (char->integer (char-downcase y))))))

(define char-ci>? (lambda (x y)
	(if (not (and (char? x) (char? y))) (error "char-ci>?: Argument is not a char.")
	  (> (char->integer (char-downcase x)) (char->integer (char-downcase y))))))

(define char-ci<=? (lambda (x y)
	(if (not (and (char? x) (char? y))) (error "char-ci<=?: Argument is not a char.")
	  (<= (char->integer (char-downcase x)) (char->integer (char-downcase y))))))

(define char-ci>=? (lambda (x y)
	(if (not (and (char? x) (char? y))) (error "char-ci>=?: Argument is not a char.")
	  (>= (char->integer (char-downcase x)) (char->integer (char-downcase y))))))

(define string=? (lambda (x y)
	(eqv? x y)))

(define string<? (lambda (x y)
	(begin
	  (define xli (string->list x))
	  (define yli (string->list y))
	  (if (< (length xli) (length yli)) #t
	    (some? (lambda (z) (eqv? z #t)) (map char<? xli yli))))))

(define string>? (lambda (x y)
	(begin
	  (define xli (string->list x))
	  (define yli (string->list y))
	  (if (> (length xli) (length yli)) #t
	    (some? (lambda (z) (eqv? z #t)) (map char>? xli yli))))))

(define string<=? (lambda (x y)
	(begin
	  (define xli (string->list x))
	  (define yli (string->list y))
	  (if (< (length xli) (length yli)) #t
	    (if (> (length xli) (length yli)) #f
	    (not (some? (lambda (z) (eqv? z #f)) (map char<=? xli yli))))))))

(define string>=? (lambda (x y)
	(begin
	  (define xli (string->list x))
	  (define yli (string->list y))
	  (if (> (length xli) (length yli)) #t
	    (if (< (length xli) (length yli)) #f
	    (not (some? (lambda (z) (eqv? z #f)) (map char>=? xli yli))))))))

(define string-ci=? (lambda (x y)
	(eqv? (map char-downcase (string->list x)) (map char-downcase (string->list y)))))

(define string-ci<? (lambda (x y)
	(begin
	  (define xn (list->string (map char-downcase (string->list x))))
	  (define yn (list->string (map char-downcase (string->list y))))
	  (string<? xn yn))))

(define string-ci>? (lambda (x y)
	(begin
	  (define xn (list->string (map char-downcase (string->list x))))
	  (define yn (list->string (map char-downcase (string->list y))))
	  (string>? xn yn))))

(define string-ci<=? (lambda (x y)
	(begin
	  (define xn (list->string (map char-downcase (string->list x))))
	  (define yn (list->string (map char-downcase (string->list y))))
	  (string<=? xn yn))))

(define string-ci>=? (lambda (x y)
	(begin
	  (define xn (list->string (map char-downcase (string->list x))))
	  (define yn (list->string (map char-downcase (string->list y))))
	  (string>=? xn yn))))


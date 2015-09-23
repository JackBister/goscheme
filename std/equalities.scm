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


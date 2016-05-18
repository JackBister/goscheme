(define-syntax cond (syntax-rules (cond else)
	((cond (test expression) ...)
	 (if test expression (cond ...))
	 (cond (else expression))
	 expression)))

(define and (lambda x (if (not (car x)) (car x) (if (null? (cdr x)) (car x) (apply and (cdr x))))))

(define or (lambda x (if (car x) (car x) (if (null? (cdr x)) (car x) (apply or (cdr x))))))


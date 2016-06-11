(define-syntax cond (syntax-rules (cond else)
	((cond (test expression) expr ...)
	 (if test expression (cond expr ...))
	 (cond (else expression))
	 expression)))

(define-syntax let (syntax-rules (let)
	((let ((x v) ...) expr ...)
	 ((lambda (x ...) expr ...) v ...))))

(define and (lambda x (if (not (car x)) (car x) (if (null? (cdr x)) (car x) (apply and (cdr x))))))

(define or (lambda x (if (car x) (car x) (if (null? (cdr x)) (car x) (apply or (cdr x))))))


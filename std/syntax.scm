(define-syntax cond (syntax-rules (cond else)
	((cond (test expression) ...)
	 (if test expression (cond ...))
	 (cond (else expression))
	 expression)))


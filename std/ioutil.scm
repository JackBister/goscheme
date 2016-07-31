(define slurpfile (lambda (fn)
	(begin
	  (define sfiter (lambda (file)
		(begin
		  (if (not (char-ready? file))
		    '()
		    (cons (read-char file) (sfiter file))))))
	  (sfiter (open-input-file fn)))))


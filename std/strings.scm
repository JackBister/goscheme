;Convert a string to a list, apply a list function to it, and convert it back.
(define aslist (lambda (f s) (list->string (f (string->list s)))))

(define string (lambda chars
	(list->string chars)))

(define string-append (lambda strings
	(list->string (concatenate (map string->list strings)))))

(define string-length (lambda (s)
	(length (string->list s))))

(define string-ref (lambda (s k)
	(list-ref (string->list s) k)))

(define substring (lambda (s start end)
	(list->string (take (list-tail (string->list s) start) (- end start)))))

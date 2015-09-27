;**x** A string containing the file name of the file to generate a doc from.
;Returns a string containing documentation generated from the file.
(define docgen (lambda (x)
	(begin
	  (define slurpfile (lambda (x)
		(bytes->chars (read-bytes (file-size x) (open-input-file x)))))
	  (define islf (lambda (x) (eqv? x #\lf)))
	  (define isspace (lambda (x) (eqv? x #\space)))

	  (define flines (split islf (slurpfile x)))

	  (define iter (lambda (li accum)
		(if (null? li) '()
		(begin
		  (define splitline (split isspace (car li)))
		  (cond
		    ((eqv? #\; (caar li)) (iter (cdr li) (append accum (append (cdar li) #\space #\space #\lf))))
		    ((eqv? (string->list "(define") (car splitline)) (cons (append '(#\# #\# #\# #\space) (cadr splitline) #\lf accum #\lf) (iter (cdr li) '())))
		    (else (iter (cdr li) '())))))))
	  (list->string (concatenate (iter flines '()))))))

;**inname** The name of the file to generate a doc from.
;**outname** The name of the file to write the resulting documentation to.
;Uses docgen on inname and writes the result to the file with the name outname.
(define docgen-and-write (lambda (inname outname)
	(begin
	  (define outfile (open-output-file outname))
	  ;TODO: display instead of write
	  (write (docgen inname) outfile)
	  (flush outfile)
	  (close-output-port outfile))))


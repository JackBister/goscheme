(define assoc (lambda (obj alist) (if (equal? obj (car (car alist))) (car alist) (assv obj (cdr alist)))))

(define assq (lambda (obj alist) (if (eq? obj (car (car alist))) (car alist) (assv obj (cdr alist)))))

(define assv (lambda (obj alist) (if (eqv? obj (car (car alist))) (car alist) (assv obj (cdr alist)))))


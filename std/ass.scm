(define assoc (lambda (obj alist) (if (null? alist) #f (if (equal? obj (car (car alist))) (car alist) (assv obj (cdr alist))))))

(define assq (lambda (obj alist) (if (null? alist) #f (if (eq? obj (car (car alist))) (car alist) (assv obj (cdr alist))))))

(define assv (lambda (obj alist) (if (null? alist) #f (if (eqv? obj (car (car alist))) (car alist) (assv obj (cdr alist))))))


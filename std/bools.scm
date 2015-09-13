(define and (lambda x (if (not (car x)) (car x) (if (empty? (cdr x)) (car x) (apply and (cdr x))))))

(define or (lambda x (if (car x) (car x) (if (empty? (cdr x)) (car x) (apply or (cdr x))))))


### append
**x** A list to append to.  
**y** An element to append to the list. Can be another list.  
**z** The function is variadic so that any further arguments will be appended.  
Examples:  
`(append '(1 2 3) 4) => (1 2 3 4)`  
`(append '(1 2 3) '(4 5 6)) => (1 2 3 4 5 6)`  
`(append '(1 2 3) 4 '(5 6)) => (1 2 3 4 5 6)`  

### concatenate
**li** A list of lists to concatenate.  
Concatenates all lists in li into one list.  
Example: `(concatenate '((1 2 3) (4 5 6))) => (1 2 3 4 5 6)`  

### filter
**pred** A predicate which is sought after in the list.  
**li** A list to filter.  
Returns all elements of li which satisfy pred.  
Example: `(filter number? '(1 2 a b c 3 d 4 e)) => (1 2 3 4)`  

### filter-map
**f** A function to map over the list.  
**li** The list to map over.  

### fold-left
**f** Function by which to fold the list.  
**z** The accumulator's starting value.  
**li** The list to fold.  
Folds left over the list. For an explanation that makes sense, see https://en.wikipedia.org/wiki/Fold_(higher-order_function)  
Examples:  
`(fold-left + 0 '(1 2 3)) => 6`  
`(fold-left cons '() '(1 2 3)) => (3 2 1)`  

### fold-right
**f** Function by which to fold the list.  
**z** The accumulator's starting value.  
**li** The list to fold.  
Folds right over the list. For an explanation that makes sense, see https://en.wikipedia.org/wiki/Fold_(higher-order_function)  
Examples:  
`(fold-right cons '() '(1 2 3)) => (1 2 3)`  

### join
**sep** The separator to insert between the lists.  
**lis** Any number of lists to join together with the separator.  
Variadic function which joins all its arguments with the separator inserted between them.  
Example: `(join 'a '(1 2 3) '(4 5 6) '(7 8 9)) => (1 2 3 a 4 5 6 a 7 8 9)`  

### last
**li** A list to return the last element from.  
Returns the last element of li.  
Example: `(last '(1 2 3)) => 3`  

### length
**x** A list to return the length of.  
Returns the length of x.  
Example: `(length '(1 2 3)) => 3`  

### list-tail
**li** A list to get the tail from.  
**k** The index to start the tail at (0-based)  
Takes the remaining elements of li starting at position k.  
Examples:  
`(list-tail '(1 2 3 4) 0) => (1 2 3 4)`  
`(list-tail '(1 2 3 4) 1) => (2 3 4)`  
`(list-tail '(1 2 3 4) 4) => ()`  

### list-ref
**li** List to get an element from.  
**k** Index of the element to get (0-based)  
Returns the *k*th element of li.  
Example: `(list-ref '(1 2 3) 1) => 2`  

### map
**f** A function to apply to the lists. The arity of the function must be equal to the number of lists.  
**lis** (variadic) any number of lists to apply the function to.  
Returns a list of the same length as the shortest argument containing the result of applying f to each element of the lists.  
Examples:  
`(map number? '(1 a #t)) => (#t #f #f)`  
`(map + '(1 2 3) '(4 5 6 7) '(1 1)) => (6 8)`  

### member
**obj** Object to search for in li.  
**li** List to search for obj.  
Searches for an object equal to obj in list li. If it is found, the object and the remainder of the list after it is returned. Uses 'equal?' to check for equality. If the object is not found, returns #f.  
TODO: remove (list) causing it to return false all the time.  
Examples:  
`(member 2 '(1 2 3 4)) => (2 3 4)`  
`(member 5 '(1 2 3 4)) => #f`  

### memq
**obj** Object to search for in li.  
**li** List to search for obj.  
See member. This function uses eq? instead of equal?.  
Example: `(memq '(1 2) '(3 4 (1 2) a b)) => ((1 2) a b)`  

### memv
**obj** Object to search for in li.  
**li** List to search for obj.  
See member. This function uses eqv? instead of equal?.  

### null?
**x** A list to check.  
Returns #t if x is an empty list, otherwise returns #f.  
Examples:  
`(null? '()) => #t`  
`(null? 1) => #f`  

### remove
**pred** A predicate. Any element of li which satisfies it will be removed from the list.  
**li** The list from which the elements should be removed.  
Returns li with all elements which satisfy pred removed.  
Example: `(remove number? '(a 1 2 b c 3 4 d)) => (a b c d)`  

### reverse
**li** A list which will be reversed.  
Returns the list with its elements in reverse order.  
Example: `(reverse '(1 2 3)) => (3 2 1)`  

### some?
**pred** A predicate against which to check the elements of li.  
**li** A list to check.  
Returns #t if any element of li satisfies pred, #f otherwise.  
Examples:  
`(some? number? '(a b c)) => #f`  
`(some? number? '(a b 1)) => #t`  

### split
**pred** The predicate to split the list on.  
**li** The list to split.  
Splits a list on any member satisfying pred.  
Returns a list of sublists that are contained between those members.  
Example: `(split symbol? '(1 2 3 a 4 5 6 a 7 8 9)) => ((1 2 3) (4 5 6) (7 8 9))`  

### take
**li** A list.  
**k** The number of elements to take from li.  
Returns the first k elements of li.  
Examples:  
`(take '(1 2 3) 0) => ()`  
`(take '(1 2 3) 2) => (1 2)`  
`(take '(1 2 3) 4) => Error`  

### zip
**lis** This is a variadic function that takes any number of lists.  
Returns the result of zipping the given lists.  
Examples:  
`(zip '(1 2 3) '(a b c)) => ((1 a) (2 b) (3 c))`  
`(zip '(1 a one) '(2 b two) '(3 c three)) => ((1 2 3) (a b c) (one two three))`  

### for-each
Because map runs in order in this implementation, map and for-each are equivalent.  



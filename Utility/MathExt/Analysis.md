# Analysis

## Functions

### is_in_map

```pseudo
ALGORITHM is_in_map(integer n, map m) -> boolean
    FOR EACH key IN m DO:
        IF key is a factor of n THEN:
            RETURN true
    RETURN false
```

### PrimeFactorization

```pseudo
ALGORITHM PrimeFactorization(integer num) -> map
    IF num = 0 THEN:
        ERROR "Cannot factorize 0"

    IF num = 1 OR num = -1 THEN:
        RETURN (1: 1)

    IF num < 0 THEN:
        num <- num * -1

    LET factors BE an empty map of the factors of num
    current_factor <- 2

	WHILE num > 1 DO:
        WHILE is_in_map(current_factor, factors) DO:
            current_factor <- current_factor + 1

        factor_count <- 0

        WHILE num > 1 AND current_factor is a factor of num DO:
            factor_count <- factor_count + 1
            num <- num / current_factor

        IF factor_count != 0 THEN:
            factors <- factors + (current_factor: factor_count)

        current_factor <- current_factor + 1

    RETURN factors
```

Complexity (measured in number of operations):

$C(num) = 8 + u(5 + C_R(v+1) + v + 3k + C_Ix)$, where:
- $C_R$ is the cost of searching whether or not the *current_factor* is not prime;
- $C_I$ is the cost of inserting a new key-value pair into the map;
- $u$ is the number of primes from 2 to *num*;
- $v$ is the average distance between two consecutive primes from 2 to infinity;
- $k$ is the average number of the powers of each prime that make up *num*;
- $x$ is the number of prime factors of *num*.

If we define:
1. *num* as $\prod_{i=0}^{n}P_i^{\alpha_i}$, where $P_i$ is the $i$-th prime number of *num* and $\alpha_i$ is the power of $P_i$ in the prime factorization of *num*;
2. $m$ as the numbers of primes from 2 to *num*.

Then:
1. $u = m$;
2. $v = \frac{1}{m}\sum_{i=1}^{m}(F(i)-F(i-1))$, where $F(i)$ is a function that returns the $i$-th prime number (e.g. $F(0) = 2$, $F(1) = 3$, $F(2) = 5$, etc.);
3. $k = \frac{1}{n}\sum_{i=0}^{n}\alpha_i$;
4. $x = n$.
5. Since we are using a map, $C_I = C_In = O(log_n)$

Therefore: $C(num) = 8 + m(5 + O(\log_n)(\frac{1}{m}\sum_{i=1}^{m}(F(i)-F(i-1))+1) + \frac{1}{m}\sum_{i=1}^{m}(F(i)-F(i-1)) + 3\frac{1}{n}\sum_{i=0}^{n}\alpha_i + O(\log_n))$;

Then: $C(num) = 8 + 5m + O(\log_n)(\sum_{i=1}^{m}(F(i)-F(i-1))+1) + \sum_{i=1}^{m}(F(i)-F(i-1)) + 3\frac{m}{n}\sum_{i=0}^{n}\alpha_i + mO(\log_n)$;

Now, since we know that $F(0)=2$, we can rewrite $\sum_{i=1}^{m}(F(i)-F(i-1))$ as $\sum_{i=1}^{m}(F(i-1) + 2 - F(i-1))$, and so it can be rewritten as $\sum_{i=1}^{m}2$, which is equal to $2m$;

By substitution, we obtain: $C(num) = 8 + 5m + O(\log_n)(2m+1) + 2m + 3\frac{m}{n}\sum_{i=0}^{n}\alpha_i + mO(\log_n)$.

Now let's consider some base cases:
1. $m = 1$ and $n = 1$. This case represents any power of two since between 2 and *num* there is only one prime and it is composed of a single prime.
2. $m > 1$ and $n = 1$. This case represents $num = P^{\alpha}$ since between 2 and *num* there are more than one prime and it is composed of a single prime.
3. $m = n$. This case represents $num = \prod_{i=0}^{n}P_i$, where $P_i$ is the $i$-th prime number of *num*.

First case (when *num* is a power of two):\
$C(num) = 15 + 3\alpha_i$. Therefore: $C(num) = 3(5 + \log_2(num))$. (Or, in other words, $C(num) = O(\log_2(num))$)

Second case (when *num* is a power of a prime):\
$C(num) = 8 + 7m + 3m\alpha_i$. Therefore: $C(num) = 8 + m(7 + \log_P(num))$. (Or, in other words, $C(num) = O(\log(num))$)

Third case (when *num* is a product of primes):\
$C(num) = 8 + 7m + 3\sum_{i=0}^{m}\alpha_i + (3m + 1)O(\log_m)$.

Now, if all $\alpha$ are equal to 1, then $8 + 10m + (3m + 1)O(\log_m)$. (Or, in other words, $C(num) = O(m\log(m))$).

Since *n* can't be bigger than *m*, then the upper bound is $C(num) = O(m\log(m))$, while the lower bound is $C(num) = O(\log(num))$.
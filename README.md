go-naivebayes
=============

Yet Another Naive-Bayesian filter algorithm

SEE ALSO
========

[https://github.com/jbrukh/bayesian](https://github.com/jbrukh/bayesian), which I based my implementation on

Why did I fork? mostly stylistic problems, which seemed like required an API change. And API changes are never easy to press on others. Mainly:

* Remove use of panic() - panics are bad in golang
* Return errors
* Name things Get\*() et al
* Allow more parallism by working closer with channels

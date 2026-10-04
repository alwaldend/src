# Review-agent evaluation configuration

This suite stages the pinned upstream review-agent skill without modifying it.
Cases cover introduced defects, bounded review authority, applicable standards,
base-branch selection, and false-positive exclusion. Explicit review requests
select the upstream skill even though its metadata disables implicit invocation.

`eval_config_test` checks packaging and configuration offline. It does not
measure review accuracy or superiority. No online target is declared because
representative reviews require a controlled Git history, repository policy,
call sites, and tool evidence that these text scenarios do not supply.

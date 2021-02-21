## The Multi-Armed Bandit (MAB) problem

The [Multi-Armed Bandit](https://en.wikipedia.org/wiki/Multi-armed_bandit) problem is a great introduction to the exploration–exploitation dilemma and is most easily understood as a gambling metaphor.

In gambling circles, slot machines are also called *one-armed bandits*, due to the lever that used to power them. (By the way, we have a similar expression for slot machines in greek, *koulocheris*, meaning 'one-handed man').

So, imagine a bright-eyed mathematician walking in a casino hall (usually called the *agent*). 

Her goal is to maximize her profit for the night, before skipping town. Unfortunately, she starts off with no information on each machine's payoff probability (see *reward distribution*); i.e. how much money it gives out, how often it pays, and if it's profitable in the long run.

So, as she starts playing, a dilemma starts forming. 
- Is our heroine certain that the current machine is a profitable one? Is it *the most* profitable in the hall? Should she continue playing *exploiting* the situation?
- Or should she *explore* for a more profitable machine? If so, when is it the best time to leave?

If it still sounds like a dry math problem, such decisions are part of our everyday lives.

- Do I try that new, unknown chinese restaurant, or stick to the well-known quantity of the neighborhood pizzeria? 
- Do I go out for drinks with our co-workers, or try and meet new people from yoga class? 
- Do I watch Parks and Rec for the fourth time, or do I check out that new series on Netflix?

## A more formal definition

https://castlelab.princeton.edu/html/ORF544/Readings/OptimalLearning_Chapter6.pdf


## Let's pose some questions

Imagine a three-armed bandit problem; 
- one arm where you have pulled twice (one win and one loss)
- one arm where you have pulled seven times (four wins, three losses)
- one arm where you have pulled twenty times (one win, nineteen losses)

*Which one should I push next to maximize draw next to maximize my long-term cumulative reward?*

*Is my strategy better than just dumb luck, drawing a random lever each time?*

So let's try to answer these questions! 

## Basic Strategies

Here are some basic strategies that you can follow
[Win–stay, lose–switch](https://en.wikipedia.org/wiki/Win%E2%80%93stay,_lose%E2%80%93switch), 

[epsilon-greedy](https://en.wikipedia.org/wiki/Greedy_algorithm)

[Upper Confidence Bound]()

https://www.chrisstucchio.com/blog/2012/bandit_algorithms_vs_ab.html



## The Gittins Index

Fortunately, a solution to the problem exists!

A theorem, the Gittins index, first published by John C. Gittins, gives an optimal policy for maximizing the expected discounted reward.[7]

In Gittins' paper, he
```
He then moves on to the "Multi–armed bandit problem" where each pull on a "one armed bandit" lever is allocated a reward function for a successful pull, and a zero reward for an unsuccessful pull. The sequence of successes forms a Bernoulli process and has an unknown probability of success. There are multiple "bandits" and the distribution of successful pulls is calculated and different for each machine. Gittins states that the problem here is "to decide which arm to pull next at each stage so as to maximize the total expected reward from an infinite sequence of pulls."[1]
```

If the machines are independent from each other and only one machine at a time may evolve, the problem is called multi-armed bandit (one type of Stochastic scheduling problems) and the Gittins index policy is optimal.





## A simplified version

In this post, we're going to discuss a simplified version of the problem, the [*Bernoulli bandit*](https://en.wikipedia.org/wiki/Bernoulli_distribution) and use it to calculate values of the Gittins index.

The *reward distribution* of each bandit will follow the Bernoulli distribution; it is a *discrete* probability distribution. Each bandit will provide a reward (success) with probability *p* and no-reward (failure) with probability *q = 1-p*.


## Constraints

As we mentioned, these kind of problems are encountered in our everyday lives. Unfortunately, the real-world is *a little* more complicated.

First off, constraints; in the real-world the player will not have an infinite budget, and may go bust early on. She doesn't have infinite time to spend in that casino hall. She will also not have infinite focus; we're human beings, and even machines have finite computing power.

As such, all these constraints make up completely different problems in the real-world. Moving back to our restaurant metaphor, think of a time where you've just moved into a city and have a few years to explore all your dining options and form opinions, versus when you're moving out in a few weeks and want to try the restaurants which provided your best moments.

Moreover, the in practice, many times systems are related to each other, or they may themselves evolve, with the underlying reward distribution changing during play. 

Once again, in our restaurant metaphor, our favorite pizzeria may try to compete with a local Dominos by dropping prices and quality; o 



## More keywords

Markov Decision Process?

Stochastic processes? Reinforcement LEarning?

Restless bandit?

## Resources 
https://gdmarmerola.github.io/ts-for-bernoulli-bandit/

http://www.statslab.cam.ac.uk/~rrw1/oc/ocgittins.pdf

https://eprints.lancs.ac.uk/id/eprint/84589/1/2016edwardsphd.pdf

https://github.com/jedwards24/gittins

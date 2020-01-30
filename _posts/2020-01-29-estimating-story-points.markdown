---
layout: post
title:  Story Points
date:   2020-01-29
author: Paschalis Ts
tags:   [agile, scrum, personal]
mathjax: false
description: "Drink *that* kool-aid"
---

Whether you start working with an Agile squad, your company forces you to use Atlassian products, or you're looking for better ways to estimate software-related tasks, you'll probably hear about *Story Points*. Since this blog is all over the place, why not give my 2c?

We can all agree that estimation is *hard*; that's why Agile teams use Story Points as an abstract metric that represents the effort required to complete an objective (either you call it user-story, feature, task or whatever). The exact measure is not important, and cannot be used to compare performance between teams or be translated into man-hours, weeks or sprints.

The idea is that you use *Story Points* to score a task, using a simple formula
> SP
> \=
> Amount of Work +
> Complexity of Work +
> Risk or Uncertainty

In many cases, the Fibonacci sequence is used to provide these scores since it's a good way to grasp changes in magnitudes; smaller tasks can be easier to understand and handle, while larger tasks are usually more vague and should be broken down to more, smaller ones. That means that a task can have a score of 1, 2, 3, 5, 8, 13, 21, 34, 55, or ∞ (cannot be scored).

I've seen this work wonderfully in a team setting; following this mindset when estimating a task, has helped both me personally and the whole team to understand different factors that can speed up or delay a piece of work.


## An example

Let's see an example of an actual feature being scored! Imagine you're a developer on a team working for digital storefront/app.

The PO explains that she needs a brand new feature; Add a textbox in the main page.

The team scores this task with 2 Story Points.

The PO explains that she actually needs *five* new textboxes.
This task is not *five* times as hard; the **Amount of work** has increased, but the other factors are still the same. So the team might score this task with 3 SP.

The PO thinks it through; what she'd need is *five* textboxes, and *one date picker*.
The task suddenly adds some **Complexity**. Will the user be able to select a date, or select a specific time, down to the second? What timezone will the datepicker be at? Will it use each user's locale, or some browser setting? Due to the complexity rising, the team might rate this task with a 5 or an 8. Will we need a designer, to ensure the date picker works on mobile devices as well?

The PO explains that the Date Picker will actually be used so that the users enter their Date of Birth before accessing the storefront, so they can access games of different [PEGI ratings](https://en.wikipedia.org/wiki/Pan_European_Game_Information).
If the team has no clue what this entails, the score will increase due to **Risk/Uncertainty**. What happens when a user is under-age? Do we flat-out refuse access, display an error message, or redirect to another page? What if the user just refreshes, and re-enters another DoB; will there be a cooldown period? Is it tied to their profile? Do we have the right to ask for this data according to GDPR? If yes, are we allowed to store this data in a cache as well? How will this be tested? This task is now a *34*, which might take a whole week to actually understand, gather the requirements, develop, test and deliver.

This seemingly simple task, now requires back-and-forth communication between technical teams and business owners, requires more extensive testing and logic.


## Get your team started.

So how can you start doing it in a team setting? Well, I'm not really qualified, so I'd suggest you ask your local Agile coach, but I can describe how I've seen it play out.

Get everyone a deck of [planning poker](https://en.wikipedia.org/wiki/Planning_poker) cards.

Select a characteristic value; It should not be too large or too small, say you've selected an '8'.

You can assign a task that is often encountered as this characteristic value. Eg. writing a new API route that will write to the database and returns a JSON response, could be the team's baseline.

All tasks during the first three or four sprints should be scored compared to this characteristic value. After the team has completed tickets of various scores, they can revisit and re-establish a characteristic value for their cards. The team should have baseline 1, 2, 3, 5, 8, 13 tasks and so on.

Scoring your tasks should be a rapid activity, that doesn't derail the team and encourage hour-long meetings.

The PO/team lead explains the feature that is going to be scored, and answers any questions that the team may have. Each member makes up their mind, and selects a card which reflects their estimate. Everyone reveals their cards their cards simultaneously and if there's consensus, or near-consensus, the team moves on to the next task. If no, a brief debate takes place, with the larger and smaller estimates kicking the discussion off, and another round of voting takes place.

## But
This meeting should not be a deep-dive into the technical aspects of *how* the development will take place, or be used to to plan and estimate things far in the future. An hour should be enough to score a whole sprint's worth backlog.

Again, this metric is not comparable between teams; a team that finishes a 100-story-point sprint, is not more productive than a team whose sprint was composed of 60 story points. It's rather a tool for estimating the workload of a team, track progress, growth, and maturity. It will allow you to calculate and track *this team's velocity*, to take into account sick days, or holidays, and in the long run, allow you to score and handle larger-scale objectives with more confidence.

Finally, when a task seems to be too big eg. approaching half your team's velocity, or is on the upper half of your deck, it's a good indication that it should be broken down, or it will be awkward to handle in an Agile way. Don't be afraid to break tasks down, it's not an exact science, and broken-down tasks might be score a little higher.



## Further reading
I'm by no means well-versed in Agile, so take all this with a grain of salt; I've been lucky enough to work in great teams, with leadership who knew to use Agile as a tool to *get things done* in a sustainable way. If your coffee is still warm and you're procrastinating, you can read some more at :

- https://martinfowler.com/bliki/StoryPoint.html
- https://www.atlassian.com/agile/project-management/estimation


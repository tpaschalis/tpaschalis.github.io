---
layout: post
title:  Computational Physics MSc - Semester 1
date:   2018-03-03
author: Paschalis Ts
tags:   [academia, comphys, msc, wil]
mathjax: false
description: "Here's what I learned/did during my first semester at the Comphys MSc"
---



### Introduction

Last October I was accepted at the [Computational Physics](http://comphys.web.auth.gr/) MSc program, at the [Aristotle University of Thessaloniki](http://auth.gr/). My first semester just finished, and I thought it's a good opportunity to share what I've learned and done during this period.

First of all, [Computational Physics](https://en.wikipedia.org/wiki/Computational_physics) as a discipline uses numerical analysis and computers to solve problems. I think it's a field that bridges Applied Physics with Computer Science, using  *loads of Statistics* and pinches of Mathematics.

Just a few example research topics of our professors include:  
* Monte Caro simulations of elementary particle experiments
* Simulation of large-scale cosmological evolution
* Computational electromagnetism/fluid dynamics for industrial applications
* Time series and [Econophysics](https://en.wikipedia.org/wiki/Econophysics)
* Medical Physics

I'm lucky enough to be in a postgraduate program, with excellent researchers, using cutting-edge tech, and always searching for the newest advancements to share with us. It's difficult to summarize 14 weeks of hard work in a few words, but I'll try! Well, here's what I did in the first semester.

### Courses
I took part in all four courses available for the first semester    
* **Data Analysis**     
This was *by far* the most difficult and demanding course. We had to use `Matlab/Octave`, to well, analyze data, and write *a lot* of code. [Prof. Kougioumtzis](http://users.auth.gr/dkugiu/) did a great job teaching us Advanced Statistics Concepts, as well as the most efficient ways to apply them and had us work and get result on real-world data. As the professor is an experienced researcher, his expertise allowed him to offer crucial feedback on both the implementation as well as the theory behind it. I also developed some soft skills as we had to present our solutions and take part in code reviews.    
*I learned about (and worked on): Distributions and their parameters, the Likelihood function, Statistical tests, Null hypotheses, Confidence Intervals, Fitting, Chi-squared gof, Least Squares, Linear/non-linear and multivariate regression, Time-series*

* **Computational Mathematics**     
This was the second most difficult course. It was quite math-heavy, and the curriculum was super dense in information, but [Prof. Kosmidis](http://kelifos.physics.auth.gr/MEMBERS/kosmas.html) managed to provide examples and mandated a hands-on approach on how the methods and their extensions are used in real-world projects. I chose the Python/NumPy/PyPlot stack for my reports, but also experimented with C and Fortran on my own time. I was very interested because not only we learned the methods as well as their mathematical basis, but also common pitfalls of their implementations, and how to mitigate them (eg. Chebyshev nodes for Runge's phenomenon, or how to deal with Stiff ODEs or non-well defined systems).   
*I learned about (and implemented): Root Finding methods, Numerical Differentiation/Integration, (Ordinary/Partial) Differential Equation Solvers, the Runge-Kutta family, Systems of Equations, Polynomial Interpolation, Linear Regression*. 

* **Programming Tools**   
This was a four-part course, where we had two to four weeks per "subject", to get the basics, and write a report, solving problems with the tools we learned about. Emphasis was placed on the efficiency of the results.    
These "subjects" were *Parallel Programming using OpenMP, Fortran, Wolfram Mathematica, and Matlab Plotting*. Example of problems we had to solve for the project report were ODE solving, Numerical Differentiation and Integration, Numerical approximation of pi, Root finding using various methods.

* **Programming Languages**   
This was one of the course that seemed relatively easier as I had already worked with C and C++, but ultimately the course was of more importance than I thought. I "re-learned" some core concepts I *thought* I knew well, worked on a lower-level language after some time (Python et. al does a good job of hiding issues under the carpet), and got a better understanding of some GNU utils and libraries.


### Research Work
Along with the mandatory coursework, I also put my hands on some research projects. Most of the efforts were in exploratory steps, mainly to develop some proof-of-concepts results, to "test the waters" and start deciding on a path with my advisor.   

Paving the ground for my Master Thesis (and hopefully, a publication), we on setting up and running some scientific software (Particle Collision Event Generators and Data Analysis tools) on both local (AUTh) and remote (CERN) HPC clusters. We replicated the workflow of some papers, to start gathering some data, and also starting reading on novel applications of Neural Networks and Feature Selection on Particle Physics.

Finally, I was also assigned in some TA-esque tasks with helping students in their undergraduate theses, with their problems on Linux and programming.


*Semester 2 is just starting, let's do it!*



PS. I consider AUTh to be my [alma mater](https://en.wikipedia.org/wiki/Alma_mater). It seems to me so far away, that on 2011, at eighteen years of age I moved here alone, with little to my name. During this time, Thessaloniki has become my new home. It's where I worked hard, got my Physics degree, got my first paper published, got pushed into research, internships, and felt the immense support of my professors-slash-mentors for which I will ever be grateful.



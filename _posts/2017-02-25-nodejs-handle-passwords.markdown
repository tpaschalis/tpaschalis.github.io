---
layout: post
title:  Storing and Authenticating passwords for a Node.js application with Passport and Bcrypt
date:   2017-02-25
author: Paschalis Ts
tags:   [tutorial, code, nodejs, infosec]
mathjax: true
description: "A simple solution to store and authenticate passwords for a Node.js application"  
---




### What can I read in this post?
You will follow along the code required to store and authenticate passwords for a ``Node.js`` application the way it should be done in 2017, using ``Passport`` and ``Bcrypt``.
The final result will be on [this GitHub repo](https://github.com/), where you will be able to create a user, salt-round the password to generate and store a hash, as well as authenticate login attempts and manage login sessions.

### Yet another javascript post in the wild?
Well, to tell the truth, I don't think that me and Javascript are on great terms, and I can't promise you that we'll ever be. But when I recently got to mess around with ``Node.js``, it actually ticked a lot of boxes, and provided a great opportunity to learn new stuff, outside the Apache/nginx ecosystem, using a very different approach-paradigm.

### How to store passwords in 2017 ?
Well, technically the answer "just don't use passwords any more" could be correct (we'll talk more about that soon). But let's just roll with something that is familiar for the end-user. 

The main reasons behind selecting [Passport](https://github.com/jaredhanson/passport), are its Open-Source state, and the good available documentation. ``Passport`` goodies include various login "strategies" available as standalone modules, and the handling of persistent login sessions and authentication states.

As for [bcrypt](https://github.com/kelektiv/node.bcrypt.js) , I feel that [it enjoys](https://codahale.com/how-to-safely-store-a-password/) [community-wide](https://pthree.org/2016/06/28/lets-talk-password-hashing/) [support](https://rietta.com/blog/2016/02/05/bcrypt-not-sha-for-passwords/). It's one of the most accesible *slow* hashing algorithms, and many implementations allow to modify the 'strength' of the resulting hash. I'm not saying it's *the* state-of-the-art or the most future-proof solution, but as it stands, it's a good starting point.

([Salting](https://en.wikipedia.org/wiki/Salt_(cryptography)) assures you that identical passwords do not have the same hash that you will store.)

### Bring up the code!
As I mentioned, the full code is available on this [github repository](https://github.com). Feel free to tinker with the whole thing, or use it  as a starting point for one of your own modules.

To keep this short and sweet, the most important parts of the code are highlighted below.

The `Node.js` implementation of `bcrypt` features an asynchronous and a synchronous way to hash the password. In order to not leave your node server hanging while the password is being salted and hashed (which can range from some milliseconds, to as much as you wish), you should use the async way, as in this example.

When creating a user, you want to  
```javascript
var bcrypt = require('bcrypt');

const saltRounds = 10;			// Current bcrypt default | 1,024 iterations | 152.4 ms on an Intel Core i7-2700

const myPlaintextPassword = 'ex4mplep4ss';
 
bcrypt.hash(myPlaintextPassword, saltRounds, function(err, hash) 
{
  // Now store 'hash' in your password DB, for the specific user.
  // console.log(hash);
});
```
You need to create a `Passport` strategy to be used for login attempts.
```javascript
passport.use(new Strategy(
  function(username, password, cb) 
    {
	db.users.findByUsername(username, function(err, user) 
	{
	if (err) { return cb(err); }
	if (!user) { return cb(null, false); }
	//if (user.password != password) { return cb(null, false); }
 	//return cb(null, user);

  	// Locate the user you want to authenticate, and use this for your login form.
  	bcrypt.compare(password, user.password, function(err, res) 
   	{
    	if (err) return cb(err);
    	if (res === false) 
    	{
      		return cb(null, false);
    	} else 
    	{
      		return cb(null, user);
    	}
 	});

	});
    }));
```
You might choose to use application level middleware for some functionality eg. logging with `morgan` and session handling with `express`. The authentication state (initializing, restoring) is handled by `Passport`.
{% highlight javascript %}
app.use(require('morgan')('combined'));
app.use(require('express-session')({ secret: 'keyboard cat', resave: false, saveUninitialized: false }));  
//  The express-session secret is used to sign the session ID cookie. More information on https://github.com/expressjs/session 

app.use(passport.initialize());
app.use(passport.session());
{% endhighlight %}

Finally, when defining your application routes you need to assign your `Passport` strategy to some 'login' page/form, as well as take care of the logout functionality.  
To "protect" a page you `ensureLoggedIn()` for your specific `user`.
```javascript
app.post('/login', 
  passport.authenticate('local', { failureRedirect: '/login' }), 
  function(req, res) {
    res.redirect('/');
  }); 
  
app.get('/logout',
  function(req, res){
    req.logout();
    res.redirect('/');
  }); 

app.get('/profile',
  require('connect-ensure-login').ensureLoggedIn(),
  function(req, res){
    res.render('profile', { user: req.user }); 
  }); 
```

### Conclusion
That's about it! You can see that it's possible to set up a simple `Node.js` application with `express`, store and authenticate passwords in less than 60 lines of code.  

You can use [this post's](https://github.com/tpaschalis/blog) github page to comment by raising an issue. I'd love to hear some opinions, criticism or whatever you have in mind. I hope you have a nice day!

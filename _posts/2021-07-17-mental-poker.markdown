---
layout: post
title:  Mental Poker in Go
date:   2021-07-17
author: Paschalis Ts
tags:   [code]
mathjax: false
description: ""
---

**Standard Disclaimer** : Don't trust me or any random internet person with anything related to cryptography and production-grade systems. This is just a toy so that I could introduce my SO to some cryptography fundamentals with a 'fun' example that didn't involve too much math. The Go code is a hacky mess, just to showcase the idea, don't judge ^^

So without further ado, let's get to it!

## Intro 
I recently read the ['Mental Poker (1979)'](https://people.csail.mit.edu/rivest/pubs/SRA81.pdf) paper by Shamir, Rivest and Adleman (of RSA fame). The acknowledgment section also mentions Robert W Floyd (as in the Floyd-Warshall algorithm) as well as Michael Rabin (as in the Rabin-Karp search algorithm), so it's a star-spangled effort all in all!

At just 8 pages, it's easy for anyone to read and understand. The basic question it poses is *"Can two potentially dishonest players play a fair game of poker without using any cards (eg. via the mail, or via the telephone) and without using any mutually trusted intermediate"*?

Before jumping in the paper, I tried to reason for about twenty minutes, and decided that no, it wouldn't be possible. I was smiling at the paper's first part that agreed with my intuition, but then was intrigued when they provided such a concise and easy way to achieve it!

Actually, all it takes is a commutative encryption scheme!

## The commutative property

In mathematics, an operation is commutative if changing the order of the operands does not affect the result. Addition and multiplication are commutative operations, while subtraction and division are *noncommutative* operations.

Another two examples of commutative operations would be the symmetry of second derivatives `∂xy f = ∂yx f` for functions satisfying Young's theorem, or the dot product. On the other hand, the three-dimensional cross product is *anti-commutative* since `b × a = −(a × b)`. 

The authors introduce the concept of a commutative encryption scheme. That is, given a key *K*, we can agree on a pair of encryption and decryption functions *Ek* and *Dk*. For all messages *X*, and keys *K* and *J*, the order which we use to apply the encryption functions should not matter; that is Ek(Ej(X)) = Ej(Ek(X)).

We also suppose that this encryption key is *strong*, meaning that  
a) given a message *X* and its encrypted version *Ek(X)*, we can never infer the key *K* and  
b) given any messages *X* and *Y*, we cannot find any keys *J* and *K* such as Ej(X) = Ek(Y),   
since any of the above would undermine the integrity of the communication and our encryption scheme.

The authors propose a simple function that utilizes Euler's totient function and fullfills all the above criteria.


## The Protocol

If you haven't read the paper yet (it takes no more than 5 minutes), here's the protocol steps for our two classic suspects, Bob and Alice.

1) Bob and Alice agree on the encryption scheme. Each selects their key/passphrase, and creates their enryption/decryption function pairs *Eb/Db* and *Ea/Da*. The keys will remain secret until the end of the game, when they will be revealed to verify that no cheating has occured

2) It's Bob's turn to deal, so he creates a new deck and encrypts all of the cards using his encryption function. He then shuffles the deck.

3) Alice receives the deck, selects five cards at random, and sends them back to Bob. These five cards will be Bob's hand; he can use his decryption function to actually see what he has been dealt.

4) Now Alice encrypts the rest of the deck with her own key, and selects five other cards at random for her own hand. Each card in Alice's hand is now doubly-encrypted using Ea(Eb(X)), and she sends the hand over to Bob.

5) Bob decrypts Alice's hand with his own decryptor, and sends it back to her. Alice's hand is now only encrypted with her own key, and Bob has no knowledge of it.

6) Alice receives her hand, and can use her decryption function to see what she has been dealt.

7) The players can bet on their hand

8) Finally, for the showdown, each player reveals their hand, *and* their secret key/passphrase. This way, each player can now check that the other was actually dealt the hand that he claims to have.


## The Code

What feels so great about this process (at least for me) is that the physical-world analogy of using padlocks and sending over each card in a little box is *just as good* for understanding what's up, without having to mention anything about mathematics.

Since we could gain a good feel of how this works, let's translate it into some hacky Go code! For the commutative encryption scheme, we can use a *ChaCha20 stream* to simply XOR with the key stream. A downside is that we have to re-create the Cipher object each time, since it maintains state and multiple calls to the `XORKeyStream` breaks the commutativity. Please reach out if you have a better alternative in mind!

```go
package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/chacha20poly1305"
)

// Do not use anywhere near serious or production-grade systems.
// Commutative scheme idea from https://asecuritysite.com/encryption/go_comm

func main() {
	// Bob and Alice select their passphrace, and get their
	// encoding/decoding functions pairs; Eb/Db and Ea/Da.
	// These will remain secrets until the end of the game, when they will be revealed
	// to verify that no cheating has occured.
	passBob := "pass1"
	passAlice := "pass2"
	Eb, Db := getED(passBob)
	Ea, Da := getED(passAlice)

	// Bob creates a new Deck, and encrypts each card using his key. He then shuffles it.
	deck := buildNewDeck()

	encryptedDeck := make([][]byte, 0)
	for _, card := range deck {
		encryptedDeck = append(encryptedDeck, Eb(card.name()))
	}
	bobRng := rand.New(rand.NewSource(time.Now().Unix()))

	bobRng.Shuffle(
		len(encryptedDeck),
		func(i, j int) {
			encryptedDeck[i], encryptedDeck[j] = encryptedDeck[j], encryptedDeck[i]
		})

	// Alice receives the deck, and selects five cards at random, and sends them back to Bob.
	// Τhis will be Bob's hand, and he can laterdecrypt the values to see what he has been dealt.
	aliceRng := rand.New(rand.NewSource(time.Now().Unix()))

	aliceRng.Shuffle(
		len(encryptedDeck),
		func(i, j int) {
			encryptedDeck[i], encryptedDeck[j] = encryptedDeck[j], encryptedDeck[i]
		})

	bobHand := encryptedDeck[0:5]
	encryptedDeck = encryptedDeck[5:]
	for _, c := range bobHand {
		fmt.Println(string(Db(c)))
	}

	// Now Alice selects five other cards at random, encrypts them with her key.
	// Each card in Alice's hand is doubly-encrypted now, and she sends them over to Bob.
	aliceRng.Shuffle(
		len(encryptedDeck),
		func(i, j int) {
			encryptedDeck[i], encryptedDeck[j] = encryptedDeck[j], encryptedDeck[i]
		})
	aliceHand := encryptedDeck[0:5]
	encryptedDeck = encryptedDeck[5:]
	for i := range aliceHand {
		aliceHand[i] = Ea(aliceHand[i])
	}

	// Bob decrypts Alice's hand with his own decryptor.
	// The hand is now only encrypted with her key, and Bob has no knowledge of it.
	for i := range aliceHand {
		aliceHand[i] = Db(aliceHand[i])
	}

	// Each player now has a hand, only encrypted with his own key,
	// and thus can decode it to see what was dealt to him, and who has the better hand.
	realBobHand := make([][]byte, 0)
	realAliceHand := make([][]byte, 0)

	for _, cb := range bobHand {
		realBobHand = append(realBobHand, Db(cb))
	}
	for _, ca := range aliceHand {
		realAliceHand = append(realAliceHand, Da(ca))
	}

	// The two players reveal their secret keys; now each player
	// can check that the other was actually dealt the cards he claimed to have played.
	fmt.Println(realBobHand)
	fmt.Println(realAliceHand)

	// Bob reveals his passphrase, and Alice verifies that she can replicate the above claim by Bob
	revealedBobPassphrase := passBob
	_, revealedDb := getED(revealedBobPassphrase)
	for i := range bobHand {
		if !reflect.DeepEqual(revealedDb(bobHand[i]), realBobHand[i]) {
			panic("Bob may have cheated!")
		}
	}

	// Alice reveals her passphrase, and Bob can do the same
	revealedAlicePassphrase := passAlice
	_, revealedDa := getED(revealedAlicePassphrase)

	for i := range aliceHand {
		if !reflect.DeepEqual(revealedDa(aliceHand[i]), realAliceHand[i]) {
			panic("Alice may have cheated!")
		}
	}

}

func getED(passphrase string) (func([]byte) []byte, func([]byte) []byte) {

	e := func(src []byte) []byte {
		key := sha256.Sum256([]byte(passphrase))
		nonce := make([]byte, chacha20poly1305.NonceSizeX)

		ecipher, _ := chacha20.NewUnauthenticatedCipher(key[:32], nonce)

		res := make([]byte, len(src))
		ecipher.XORKeyStream(res, src)
		return res
	}

	d := func(src []byte) []byte {
		key := sha256.Sum256([]byte(passphrase))
		nonce := make([]byte, chacha20poly1305.NonceSizeX)

		dcipher, _ := chacha20.NewUnauthenticatedCipher(key[:32], nonce)

		res := make([]byte, len(src))
		dcipher.XORKeyStream(res, src)
		return res
	}

	return d, e
}

type card struct {
	Rank string
	Suit string
}

func buildNewDeck() []card {
	ranks := []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}
	suits := []string{"Clubs", "Diamonds", "Hearts", "Spades"}

	var deck []card
	for _, s := range suits {
		for _, r := range ranks {
			deck = append(deck, card{Rank: r, Suit: s})
		}
	}

	return deck
}

func (c card) name() []byte {
	return []byte(c.Rank + " of " + c.Suit)
}
```


## Outro

That's all for today! I don't claim to know anything about crypto, so don't hesitate to reach out for any corrections and issues. See you around!

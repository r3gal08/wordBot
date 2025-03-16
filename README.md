# wordBot

This is a new word app project I've been building for quickly getting definitions for words I want to learn.
The end goal here is to eventually have a backend containerized service for handling json requests and returning 
word definitions, storing a database of words a user wants to learn, "smart-quizzing" the user to ensure they are learning
the words (quizzing more often on new words, less-often on old ones that the user has shown they understand). I may add
additional functionality to allow a user to fill in additional information such as where they found the word (from a book
, news article, etc) which may assist in giving a more tailored response to what the user is genuinely interested in
(likely leveraging some kind of LLM). By doing so, my hypothesis is a user will learn a word more easily as it is already
tied to a pre-existing thought processes and will further solidify their understanding.

- For my frontend I am choosing to use the Flutter framework once again
- For my backend I am choosing to use Golang as it is a language I am wanting to learn :)

---

## Codebase Guide

./backend - backend application code built in go. Also includes docker-compose file and db schema.sql file
./lib - Core frontend codebase built with flutter and the dart programming language

The remander of the codebase is mostly auto-generated flutter files...



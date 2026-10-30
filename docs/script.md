# TypeScript's Type Wizardry: the script

FOSSCOMM 2026 · about 23 minutes · 2839 words

One section per slide, numbered like the counter on the slide, so go-to with that number jumps to it. *[Italics in brackets]* are stage directions, not lines. Go to the next slide at the end of each section. The times are where each slide starts and ends at about 130 words a minute.

---

## 01 / 28 · TypeScript's Type Wizardry · 00:00 → 00:30

Hi everyone, and thanks for coming. *[pause]*

This talk is called TypeScript's Type Wizardry, and where to use it. For the next twenty minutes or so, we'll look at the part of TypeScript that looks like programming, but runs inside the type checker, before your code ever runs. And at the end, we'll talk about where it's actually worth using.

## 02 / 28 · A little about me · 00:30 → 01:30

First, a little about me. I'm Fanis, and I'm a software engineer at Typesense. Typesense is an open source search engine, written in C++. I work on everything around it: all of its official API clients and integrations. So JavaScript and TypeScript, Python, Go, PHP, Ruby, Java, Dart, and a few more.

In my own time, I build tools. *[point at the cards]* Phonto plays video wallpapers, on Wayland and on macOS. Stochos lets you drive your mouse from the keyboard. And Kyma turns markdown into presentations, in the terminal. That's what you're looking at right now.

And then there's typesense-ts. It's a fully type-safe client for Typesense, and most of the wizardry you'll see today comes from it.

Enough about me. Let's start with something that shouldn't be possible.

## 03 / 28 · The compiler read my string · 01:30 → 02:45

Here's a search, with typesense-ts. *[point at the left card]* First, we describe a collection called users. It has three fields: a name, an age, which we're allowed to sort by, and an email, which is optional.

Then we search it. We look for "fanis" in the name, and we ask for the results sorted by email, descending. *[point at line 13]*

There's a problem. We never said we can sort on email. Normally, you find that out when the server rejects your request. At runtime. Maybe in production.

Here, you find out in your editor, before anything runs. *[point at the right card]* And look at the error message: "email is not a sortable field." Nobody wrote that check by hand for this query. That message is a type. The compiler read the string "email colon desc", split it up, and checked it against the schema.

By the end of this talk, you'll know how this works.

## 04 / 28 · Types are sets · 02:45 → 04:15

To get there, we need one idea. Everything else today builds on it: a type is a set of values.

*[point at line 1]* "a" or "b" is a set with two values in it. String is the set of every string. And a single string, like "a", is a type too: a set with exactly one value. That's called a literal type, and we'll use them a lot.

Unknown is the set of every value. Never is the set with nothing in it, the empty set.

*[point at the drawing]* Here's the same thing, drawn. AB sits inside string, string sits inside unknown, and everything else, numbers, booleans, objects, lives in unknown too.

Now, the keyword extends. In types, extends asks one question: does the left side fit inside the right? *[point at T1 to T3]* Does "a" fit in AB? Yes. Does AB fit in string? Yes. Does string fit in AB? No, most strings aren't "a" or "b".

And then there's any. *[point at line 5]* Any isn't a set at all. It's TypeScript switching checking off. It fits inside everything, and everything fits inside it. So if you ask it our question, *[point at T4]* you get both answers at once: boolean. That's why we reach for unknown instead. It holds every value too, but it makes you check what you have before you use it.

Every trick today is one of these questions about sets.

## 05 / 28 · The type system is a small language · 04:15 → 05:00

The second idea: the type system is a small programming language. And you already know it, because you know JavaScript.

On the left are things you do every day, with values. On the right are the same ideas, one level up, with types. Each row is a tool we'll use for the rest of the talk, so let's go one row at a time.

## 06 / 28 · A small language: if · 05:00 → 05:30

If. In JavaScript, you'd write: typeof x equals string, then one answer or the other. In types, you write: X extends string, question mark, then the two answers. And extends is our question from the last slide: does X fit inside string? Yes gives you the first branch, no gives you the second. That's called a conditional type.

## 07 / 28 · A small language: destructure · 05:30 → 05:45

Destructuring. In JavaScript, const head, dot dot dot tail, pulls values out of an array. In types, the keyword infer does the same thing. You describe a shape, here a tuple with something first and the rest after it, and you give the pieces names. Head is the first element of L.

## 08 / 28 · A small language: loop over keys · 05:45 → 06:15

Looping over keys. For key in obj walks the keys of an object. A mapped type walks the keys of a type: Key in keyof Obj. For every key that goes in, one property comes out. This one copies Obj as it is, but you can change what each property becomes.

## 09 / 28 · A small language: recursion · 06:15 → 06:45

Recursion. Types have no for loop and no while loop. So when you need a loop, a type calls itself. F asks: is X an array? If it is, unwrap one level and call F again. If it isn't, we're done. So F of number, array, array, array, is just number.

## 10 / 28 · A small language: match a string · 06:45 → 07:15

And matching strings. Backticks work in types too. This is a template literal type, and it matches a string against a pattern. Does S start with "gid colon slash slash"? If it does, infer captures whatever comes after it, and names it Id. Remember this one: the parser from the start is built out of exactly this.

## 11 / 28 · Where did my autocomplete go? · 07:15 → 07:45

Let's put the sets idea to work, with a puzzle. Here's a search function. The query is either the wildcard, a star, or any string you like. We'd love the editor to suggest the star. *[point at line 3]* But it suggests nothing. And when you hover over q, it says: string.

We asked for star or string, and we got string. Where did the wildcard go? Think back to the sets.

## 12 / 28 · Where did my autocomplete go? The answer · 07:45 → 09:30

The star is a string. It's already inside the set of all strings. So star or string is just string, and TypeScript simplifies the star away. Your autocomplete goes with it.

*[point at Q2]* The fix looks strange: string, and an empty object type. The empty object type means anything that isn't null or undefined. Every string fits in it, so this is still every string. The set didn't change. So why doesn't TypeScript simplify this one too?

Because TypeScript doesn't simplify by checking which values a type holds. It follows a short list of rules about how a type looks. One of them says: if a union has a string literal and plain string, drop the literal. String and an empty object isn't plain string. To TypeScript, it's a string combined with an object type. So the rule doesn't fire, the star survives, and so does your autocomplete.

That's also why string and unknown wouldn't work: TypeScript always throws unknown away when you combine it with something, so you're back to plain string. And this isn't a loophole the next version will close. The compiler has a special case that keeps string and an empty object exactly as you wrote it, because so many libraries rely on this trick.

*[point at the right card]* typesense-ts uses exactly this for q.

## 13 / 28 · Mapped types are loops · 09:30 → 10:15

Now, the loop over keys, used for real. This comes from a tool I'm building. It syncs a few kinds of resources into a search engine: products, collections and articles. Each one needs a schema.

*[point at line 4]* Here's the mapped type. For every R in Resource, write one property that holds a Schema. *[point at the comment]* So Schemas becomes an object with all three keys.

*[point at the right card]* Now we fill it in, and we forget articles. It doesn't compile, and the error tells you exactly which one is missing.

The best part comes next year. Someone adds a fourth resource to that union, and this breaks until they give it a schema. You can't forget. And the built-in Record type does exactly this. Record is a mapped type.

## 14 / 28 · Parsing strings with infer · 10:15 → 10:45

So let's build the hook. This is a sort parser, in miniature.

*[point at lines 1 to 4]* It's the match-a-string row: a pattern with a colon in the middle, and infer on both sides. Give it "age colon desc", *[point at line 6]* and it splits at the colon. The field is age, and the direction is desc. Types, pulled out of a string.

## 15 / 28 · Parsing strings: checking the pieces · 10:45 → 11:15

Now we check the pieces. *[point at line 10]* Same pattern as before: a field and a direction, either side of the colon. *[point at line 8]* Say age and score are the fields you can sort by. If the field fits inside Sortable, and the direction is asc or desc, we give the string back, unchanged. If not, we give back a different string: an error message, as a type. Field, "is not a sortable field."

## 16 / 28 · Parsing strings: the error · 11:15 → 11:45

Then we make it the type of the parameter. *[point at lines 16 to 18]* A valid string is assignable to itself, so it passes. An invalid one would have to be assignable to the error message, and it isn't. So the compiler complains, and our message is right there in the error, where you typed it.

That's the hook from the start, in about twenty lines.

## 17 / 28 · Parsing strings: the real one · 11:45 → 12:30

The real one, in typesense-ts, does it properly: a tokenizer, then a validator, like a real parser. *[point at line 3]* ReadString takes the first character off the string, with infer Head and infer Tail. If it's part of a name, it keeps it, and calls itself with the rest. That's recursion again.

*[point at the right card]* So "age colon desc" becomes three tokens: a name, a colon, and a direction. Then each token is checked against your schema, and that's where the error from the hook came from.

## 18 / 28 · Why do I need recursion? · 12:30 → 13:00

There are no loops in the type system, so anything that loops has to recurse. *[point at TupleOf]* TupleOf keeps adding one more T to a tuple, until the length matches. TupleOf three, number, is three numbers.

And here's the trick that matters for later. A tuple's length isn't just number. It's a literal, like 3. So you can count. *[point at AddOne]* Make a tuple of three, add one more element, and read the length: four. That's addition, done by the type checker.

## 19 / 28 · Recursion: one weight per field · 13:00 → 13:30

Where is that useful? Typesense lets you give a weight to each field you search, one weight per field. *[point at lines 6 and 7]* Here we search three fields, but give only one weight.

The server would reject this at runtime. typesense-ts counts instead. The length of query by is three, so the weights have to be a tuple of exactly three numbers. *[point at the right card]* And the editor tells you: the source has one element, but the target requires three.

## 20 / 28 · How far can it go? · 13:30 → 15:00

So: types can branch, destructure, loop and count. How far can that go?

In February 2025, Dimitri Mitropoulos, from Michigan TypeScript, set out to show why Doom can't run in TypeScript types. The project's README says it best: "The funny thing is... it can." *[pause]*

TypeScript types can run Doom.

Here's how the game got in. Doom is written in C. It was compiled to WebAssembly, and then a script rewrote that WebAssembly as TypeScript types: ninety-nine thousand lines of them. Then a computer made entirely of types runs it.

Some numbers, from the launch video. Three and a half trillion lines of types went through the type checker, a hundred and seventy-seven terabytes of them. Twelve days, around the clock, for the first frame. For comparison: a normal app's whole type check is a few hundred thousand instantiations. This did twenty million of them every second, for twelve days, and the type checker could use a hundred gigabytes of memory.

*[point at the screenshot]* The frame comes out as a type. You hover over it in your editor, and there's Doom.

Every instruction of the game ran in the type checker.

## 21 / 28 · How far can it go? A computer made of types · 15:00 → 16:15

How do you build a computer out of types? The real machine is enormous, so here's the simplest one from the same project, made to explain the idea.

*[point at the left card]* A program is a tuple of instructions. Add "hello". Add "world". Uppercase. Add an exclamation mark. Next to it is the result, which starts out empty.

*[point at the right card]* executeInstruction is the processor. *[point at lines 4 to 7]* It destructures the instructions into the next one and the rest. That's infer, from the table. *[point at lines 8 to 10]* It runs that one instruction, and calls itself with the new state. That's recursion. When no instructions are left, the result is the answer.

*[point at the last line on the left]* This one gives you HELLO WORLD, in capitals.

Doom's machine has a hundred and sixteen instructions, plus memory and a stack. But it's the same loop: take the next instruction, run it, call yourself.

## 22 / 28 · How far can it go? Logic with a lookup table · 16:15 → 17:45

How does a computer made of types do math? Here's one of the simplest pieces: AND. In this machine, every number is a string of zeros and ones. AND compares two of them digit by digit, and the answer is one only where both digits are one.

*[point at And]* The rule is just a table. Look up the digit from A, then the digit from B. One and one gives one. Everything else gives zero.

*[point at BitwiseAnd]* Then we walk both strings together. This is the match-a-string row from earlier: take the first digit off each, look them up in the table, and call yourself with the rest. That's recursion.

*[point at the right card]* One-one-zero-zero, and one-zero-one-zero. One and one: one. One and zero: zero. Zero and one: zero. Zero and zero: zero. One-zero-zero-zero.

OR and XOR are the same walk, with a different table. The processor that runs Doom is built out of pieces like this.

One last detail. The type checker has limits, and our TupleOf hit one at a thousand. While building this, Dimitri lifted those limits in a patched compiler. The run itself stopped every fifty to a hundred instructions, saved the whole state, and carried on from there.

## 23 / 28 · TypeScript 7: the compiler, rewritten in Go · 17:45 → 18:15

Meanwhile, the compiler itself got rewritten. TypeScript 7 came out in July. It's the compiler, ported to Go. The same type system, about ten times faster.

*[point at the first row]* These are full builds of big open source projects, from the announcement. VS Code's own codebase went from two minutes to about ten seconds.

The examples from today give the same errors on it. And 7.1, which is in beta right now, adds something new.

## 24 / 28 · TypeScript 7.1: content mappers · 18:15 → 19:15

It's called content mappers.

Until now, tsc checked TypeScript and JavaScript files, and nothing else. If you imported a Vue file, a Svelte file or a SQL file, its type was whatever you declared by hand. Usually something loose, and often wrong.

*[point at the left card]* A content mapper is a plugin for the compiler. It takes a file, turns it into TypeScript, and remembers where every piece came from. tsc checks that TypeScript, and when it finds an error, it reports it in the original file.

*[point at the right card]* You register one per file extension, in your tsconfig. And because a mapper is a program that runs on your machine, you have to allow it, with run external code.

Vue, Svelte, Astro, MDX and Ember already have mappers, prototyped or in progress. But my favourite one is for SQL.

## 25 / 28 · Reylo: SQL files, typed by your database · 19:15 → 20:15

It's called Reylo, and it's by Colin McDonnell, who you might know as the author of Zod.

*[point at the left card]* You write a plain SQL query, in a .sql file. Select id, name and email from users, where the id is a parameter.

*[point at the right card]* Then you import that file from TypeScript, and run it. Look at the type of user. Nobody wrote it. Reylo asked your development database to prepare the query, without running it, and the database answered with the type of every column.

Look at the details. Email can be null, because the table doesn't say "not null". And the whole result can be null, because we matched on the primary key, so you get one row or none.

No hand-written row types, and no code generation step.

## 26 / 28 · Reylo: a typo in the SQL is a type error · 20:15 → 21:15

And when the SQL is wrong? *[point at the left card]* Let's misspell name.

The database rejects the query while it's being typed, and *[point at the right card]* tsc reports it inside the SQL file, at the right line and column. Your editor shows it too, as you type, with no extra plugin.

It works the other way around too. Rename a column in the database, and every query that reads it breaks on the next type check.

Earlier, we parsed strings with types. Reylo doesn't parse SQL in types at all. It asks the one thing that actually knows the answer: your database. Sometimes the best type-level trick is to not do it in types.

## 27 / 28 · So where should you use it? · 21:15 → 22:15

So, where should you use all this?

*[point at the top left]* In app code, a little goes a long way. The autocomplete trick. Mapped records, so you never forget a case. And let tools do the heavy lifting: let Reylo type your SQL, so you don't have to.

*[point at the top right]* In a library, spend it at the boundary, where your users hand you schemas, strings and options. Derive the types from what they give you. Check the strings they write. Put the fix in the error message. And test your types like the rest of your API, with expectTypeOf and ts-expect-error.

*[point at the bottom]* Because it costs something. It costs compile time: Doom took twelve days, and yours shouldn't, so measure it. It costs error quality: a typo in a filter can come back as "Unknown token: exclamation mark". And someone has to read your types next year.

So here's my advice: use the least wizardry that makes the wrong call impossible.

## 28 / 28 · Thank you. Questions? · 22:15 → 22:45

That's everything. Thank you! *[pause for applause]*

On the left are some places to keep practicing, plus the Doom video and Reylo. Some of today's puzzles were inspired by TypeHero. My GitHub and my site are on the right.

I'd love to take your questions. The weirder the type, the better.

*[leave this slide up for the Q&A]*

---

## If asked

**Why not just use Zod?**
Zod checks data while your program runs. Everything today checks your code before it runs. They work well together: a Shopify app I'm building uses Zod for the API responses it gets back, and types for everything else.

**Isn't all this slow?**
It can be, so measure it. tsc has a generate trace flag, and the analyze-trace tool points you at the slow types. typesense-ts generates a trace on every type check, so we notice when it gets slower.

**When is it too much?**
When the error messages stop helping, or when nobody on your team can maintain it.

**Does Reylo parse SQL in types?**
No. Its content mapper asks your dev database to prepare the query, and reads the types back. So you need a dev database running when you type-check.

**Did Doom need a patched compiler?**
While it was being built, yes: Dimitri lifted the compiler's limits, like the thousand-level recursion limit. The video says the end result runs in vanilla TypeScript, by stopping every fifty to a hundred instructions and picking up from the saved state.

**Will my code work on TypeScript 7?**
It's the same type system, ported to Go, and it's been stable since July. 7.0 shipped without the programmatic API that some tools rely on, and 7.1 adds it.

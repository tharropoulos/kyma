---
title: "TypeScript's Type Wizardry"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **TypeScript's Type Wizardry**  ·  until 00:30

  - Today: the parts of TypeScript's type system that look like programming
  - And where they're worth using

  **NEXT →** Before we begin, a little about who I am
---

![TypeScript's Type Wizardry|138x35](wizardry_cover.png)

----
---
title: "A little about me"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **A little about me**  ·  until 01:30

  - Software engineer at Typesense, an open source search engine
  - I work on all of its official API clients and integrations
  - My own projects: Phonto, Stochos, and Kyma, which made this deck
  - And typesense-ts: most of today's wizardry comes from it

  **NEXT →** Let's start with something that shouldn't be possible
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*$whoami ·   02 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Who Am I?
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row gap=1]
[col]
[/col]
[col width=26 align=right valign=middle]
**Software Engineer at**
[/col]
[col width=30]
![Typesense|30x3](typesense_logo.png)
[/col]
[col]
[/col]
[/row]
[row gap=3]
[col width=82]
[grid gap=1]
[row]
[col align=center]
Some of my own projects
[/col]
[/row]
[row gap=2]
[col border=rounded border_color="#3A57CD" pad="0 1"]
[grid gap=0]
[row height=2]
[col]
  museslabs/**phonto**
[/col]
[/row]
[row height=4]
[col]
Video wallpaper program for wayland compositors and macos
[/col]
[/row]
[row]
[col]
 Rust         436   11
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="0 1"]
[grid gap=0]
[row height=2]
[col]
  museslabs/**stochos**
[/col]
[/row]
[row height=4]
[col]
Keyboard driven mouse control for Linux and macOS
[/col]
[/row]
[row]
[col]
 Rust         148   17
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=2]
[col border=rounded border_color="#3A57CD" pad="0 1"]
[grid gap=0]
[row height=2]
[col]
  **typesense-ts**
[/col]
[/row]
[row height=4]
[col]
A fully typesafe Typesense client for typeheads
[/col]
[/row]
[row]
[col]
 TypeScript   53    7
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="0 1"]
[grid gap=0]
[row height=2]
[col]
  museslabs/**kyma**
[/col]
[/row]
[row height=4]
[col]
Presentations from markdown in the terminal with fancy transition animations
[/col]
[/row]
[row]
[col]
 Go           98    2
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
[/col]
[col width=53]
[grid gap=1]
[row]
[col align=center]
What I work on day to day
[/col]
[/row]
[row]
[col border=rounded border_color="#3A57CD" pad="0 1"]
[grid gap=0]
[row height=2]
[col]
  typesense/**typesense**
[/col]
[/row]
[row height=4]
[col]
A Fast, typo tolerant, in-memory fuzzy Search Engine for building delightful search experiences
[/col]
[/row]
[row]
[col]
 C++          27k   1k
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col border=rounded border_color="#3A57CD" pad="0 1"]
[grid gap=0]
[row height=2]
[col]
  typesense/*\**
[/col]
[/row]
[row height=4]
[col]
I work on all of the official API clients and integrations
[/col]
[/row]
[row]
[col]
JS/TS, Python, Go, PHP, Ruby, Java, Dart...
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
  `fanis.iusevimbtw.com`            `tharropoulos`
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]

----
---
title: "The compiler read my string"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **The compiler read my string**  ·  until 02:45

  - A search on a users collection, sorted by email, descending
  - email is optional, and you can't sort on it
  - Normally the server tells you, at runtime
  - Here the editor tells you, before anything runs
  - That error message is a type: the compiler parsed the string
  - By the end of the talk, you'll know how this works

  **NEXT →** First, a mental model
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*INTRO / THE HOOK   ·   03 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# The compiler read my string
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row gap=3]
[col width=66 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a typesense-ts search*
[/col]
[/row]
[row]
[col align=left]
```ts{6,13} --numbered
const users = collection({
  name: "users",
  fields: [
    { name: "name", type: "string" },
    { name: "age", type: "int32", sort: true },
    { name: "email", type: "string", optional: true },
  ],
});

await users.search({
  q: "fanis",
  query_by: ["name"],
  sort_by: "email:desc",
});
```
[/col]
[/row]
[/grid]
[/col]
[col valign=middle]
[grid gap=1]
[row]
[col border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TSC / the error, wrapped to fit*
[/col]
[/row]
[row]
[col align=left]
```text{99-99}
error TS2345: Argument of type
  '{ q: string; query_by: string[]; sort_by: string; }'
is not assignable to parameter of type
  '"[1;38;2;8;0;28;48;2;225;116;0m[Error on sort_by][0m[38;2;245;237;184m: Invalid token sequence:[0m
    [38;2;245;237;184mInvalid identifier: email is not a sortable field.[0m[38;2;128;128;128m"'[0m
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
*No request sent, no test run. Just the editor.*
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**By the end of this talk, you'll know how this works.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]

----
---
title: "Types are sets"
transition: slideUp
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Types are sets**  ·  until 04:15

  - The one idea everything else builds on: a type is a set of values
  - "a" | "b" has two values, string has every string
  - A single string like "a" is a type too, with one value: a literal type
  - unknown has every value, never has none: the empty set
  - extends asks one question: does the left side fit inside the right?
  - "a" fits in AB, AB fits in string, but string doesn't fit in AB
  - any is the exception: it isn't a set, it's type checking switched off
  - It fits inside everything, and everything fits inside it. Ask it the extends question and you get both answers: boolean
  - unknown is the safe one: it holds every value too, but you have to check what it is before you use it
  - Every trick today is a question about sets

  **NEXT →** And the type system is a small language
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MENTAL MODEL / TYPES ARE SETS   ·   04 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Types are sets
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=114]
[grid gap=1]
[row gap=3]
[col width=63 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / types as sets of values*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
type AB = "a" | "b";  // two values
type S  = string;     // every string
type U  = unknown;    // every value
type N  = never;      // no values at all
type A  = any;        // not a set: checking off

// does the left side fit inside the right?
type T1 = "a" extends AB ? true : false;     // true
type T2 = AB extends string ? true : false;  // true
type T3 = string extends AB ? true : false;  // false
type T4 = any extends AB ? true : false;     // boolean
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TEXT / the same types, drawn*
[/col]
[/row]
[row]
[col align=left]
```text{1-99}
┌─ unknown ───────────────────────────┐
│                                     │
│  ┌─ string ──────────────────────┐  │
│  │                               │  │
│  │  ┌─ AB ──────┐                │  │
│  │  │  "a" "b"  │    "hello"     │  │
│  │  └───────────┘                │  │
│  │                               │  │
│  └───────────────────────────────┘  │
│                                     │
│    42      true      { id: "1" }    │
│                                     │
└─────────────────────────────────────┘
          never: the empty set
       any: no box, it opts out
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**`extends` asks one question: does the left side fit inside the right?**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]

----
---
title: "The type system is a small language"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **The type system is a small language**  ·  until 05:00

  - You already know every idea on this slide, from JavaScript
  - Each row is a tool we'll use for the rest of the talk
  - One row at a time

  **NEXT →** Starting with if
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MENTAL MODEL / A SMALL LANGUAGE   ·   05 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# The type system is a small language
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row gap=2]
[col width=18 border=hidden pad="1 0"]
[grid gap=1]
[row]
[col align=left]
*THE IDEA*
[/col]
[/row]
[row]
[col align=left]
```text{1-99}
if
destructure
loop over keys
recursion
match a string
```
[/col]
[/row]
[/grid]
[/col]
[col width=54 border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*JAVASCRIPT / values*
[/col]
[/row]
[row]
[col align=left]
```js{1-99}
typeof x === "string" ? 1 : 0
const [head, ...tail] = list
for (const key in obj) out[key] = obj[key]
const f = (x) => Array.isArray(x) ? f(x[0]) : x
s.startsWith("gid://") ? s.slice(6) : null
```
[/col]
[/row]
[/grid]
[/col]
[col width=62 border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / types*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
X extends string ? 1 : 0
L extends [infer Head, ...infer Tail] ? Head : never
{ [Key in keyof Obj]: Obj[Key] }
type F<X> = X extends (infer I)[] ? F<I> : X
S extends `gid://${infer Id}` ? Id : null
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Ideas you already use, one level up. The rest of the talk fills in the right column.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "A small language: if"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **A small language: if**  ·  until 05:30

  - typeof x === "string" in JavaScript, X extends string in types
  - extends is the question from the last slide: does X fit inside string?
  - Yes takes the first branch, no takes the second

  **NEXT →** Destructuring
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MENTAL MODEL / A SMALL LANGUAGE   ·   06 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# The type system is a small language
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row gap=2]
[col width=18 border=hidden pad="1 0"]
[grid gap=1]
[row]
[col align=left]
*THE IDEA*
[/col]
[/row]
[row]
[col align=left]
```text{1}
if
destructure
loop over keys
recursion
match a string
```
[/col]
[/row]
[/grid]
[/col]
[col width=54 border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*JAVASCRIPT / values*
[/col]
[/row]
[row]
[col align=left]
```js{1}
typeof x === "string" ? 1 : 0
const [head, ...tail] = list
for (const key in obj) out[key] = obj[key]
const f = (x) => Array.isArray(x) ? f(x[0]) : x
s.startsWith("gid://") ? s.slice(6) : null
```
[/col]
[/row]
[/grid]
[/col]
[col width=62 border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / types*
[/col]
[/row]
[row]
[col align=left]
```ts{1}
X extends string ? 1 : 0
L extends [infer Head, ...infer Tail] ? Head : never
{ [Key in keyof Obj]: Obj[Key] }
type F<X> = X extends (infer I)[] ? F<I> : X
S extends `gid://${infer Id}` ? Id : null
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**if becomes a conditional type that asks whether `X` fits inside `string` or not.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "A small language: destructure"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **A small language: destructure**  ·  until 05:45

  - const [head, ...tail] = list pulls values out of an array
  - infer does the same to a type: match the shape, name the pieces
  - Here Head is the first element of the tuple L

  **NEXT →** Loops over keys
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MENTAL MODEL / A SMALL LANGUAGE   ·   07 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# The type system is a small language
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row gap=2]
[col width=18 border=hidden pad="1 0"]
[grid gap=1]
[row]
[col align=left]
*THE IDEA*
[/col]
[/row]
[row]
[col align=left]
```text{2}
if
destructure
loop over keys
recursion
match a string
```
[/col]
[/row]
[/grid]
[/col]
[col width=54 border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*JAVASCRIPT / values*
[/col]
[/row]
[row]
[col align=left]
```js{2}
typeof x === "string" ? 1 : 0
const [head, ...tail] = list
for (const key in obj) out[key] = obj[key]
const f = (x) => Array.isArray(x) ? f(x[0]) : x
s.startsWith("gid://") ? s.slice(6) : null
```
[/col]
[/row]
[/grid]
[/col]
[col width=62 border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / types*
[/col]
[/row]
[row]
[col align=left]
```ts{2}
X extends string ? 1 : 0
L extends [infer Head, ...infer Tail] ? Head : never
{ [Key in keyof Obj]: Obj[Key] }
type F<X> = X extends (infer I)[] ? F<I> : X
S extends `gid://${infer Id}` ? Id : null
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Destructuring becomes `infer` to match the shape and name the pieces.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "A small language: loop over keys"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **A small language: loop over keys**  ·  until 06:15

  - for...in walks the keys of an object
  - A mapped type walks the keys of a type: Key in keyof Obj
  - One property out for every key in. This one copies Obj as it is

  **NEXT →** Recursion
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MENTAL MODEL / A SMALL LANGUAGE   ·   08 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# The type system is a small language
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row gap=2]
[col width=18 border=hidden pad="1 0"]
[grid gap=1]
[row]
[col align=left]
*THE IDEA*
[/col]
[/row]
[row]
[col align=left]
```text{3}
if
destructure
loop over keys
recursion
match a string
```
[/col]
[/row]
[/grid]
[/col]
[col width=54 border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*JAVASCRIPT / values*
[/col]
[/row]
[row]
[col align=left]
```js{3}
typeof x === "string" ? 1 : 0
const [head, ...tail] = list
for (const key in obj) out[key] = obj[key]
const f = (x) => Array.isArray(x) ? f(x[0]) : x
s.startsWith("gid://") ? s.slice(6) : null
```
[/col]
[/row]
[/grid]
[/col]
[col width=62 border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / types*
[/col]
[/row]
[row]
[col align=left]
```ts{3}
X extends string ? 1 : 0
L extends [infer Head, ...infer Tail] ? Head : never
{ [Key in keyof Obj]: Obj[Key] }
type F<X> = X extends (infer I)[] ? F<I> : X
S extends `gid://${infer Id}` ? Id : null
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**A loop over keys becomes a mapped type: one property out for every key in.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "A small language: recursion"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **A small language: recursion**  ·  until 06:45

  - Types have no for or while: a type that refers to itself is the loop
  - F unwraps one array per call, until X is not an array anymore
  - F<number[][][]> is number

  **NEXT →** Matching strings
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MENTAL MODEL / A SMALL LANGUAGE   ·   09 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# The type system is a small language
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row gap=2]
[col width=18 border=hidden pad="1 0"]
[grid gap=1]
[row]
[col align=left]
*THE IDEA*
[/col]
[/row]
[row]
[col align=left]
```text{4}
if
destructure
loop over keys
recursion
match a string
```
[/col]
[/row]
[/grid]
[/col]
[col width=54 border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*JAVASCRIPT / values*
[/col]
[/row]
[row]
[col align=left]
```js{4}
typeof x === "string" ? 1 : 0
const [head, ...tail] = list
for (const key in obj) out[key] = obj[key]
const f = (x) => Array.isArray(x) ? f(x[0]) : x
s.startsWith("gid://") ? s.slice(6) : null
```
[/col]
[/row]
[/grid]
[/col]
[col width=62 border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / types*
[/col]
[/row]
[row]
[col align=left]
```ts{4}
X extends string ? 1 : 0
L extends [infer Head, ...infer Tail] ? Head : never
{ [Key in keyof Obj]: Obj[Key] }
type F<X> = X extends (infer I)[] ? F<I> : X
S extends `gid://${infer Id}` ? Id : null
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Types have no loops. A type that calls itself is the loop.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "A small language: match a string"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **A small language: match a string**  ·  until 07:15

  - Template literal types match a string against a pattern
  - infer Id captures whatever follows gid://
  - "gid://shopify/Product/1" gives "shopify/Product/1", anything else gives null
  - The sort_by parser from the start is built out of exactly this

  **NEXT →** A puzzle: where did my autocomplete go?
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MENTAL MODEL / A SMALL LANGUAGE   ·   10 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# The type system is a small language
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row gap=2]
[col width=18 border=hidden pad="1 0"]
[grid gap=1]
[row]
[col align=left]
*THE IDEA*
[/col]
[/row]
[row]
[col align=left]
```text{5}
if
destructure
loop over keys
recursion
match a string
```
[/col]
[/row]
[/grid]
[/col]
[col width=54 border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*JAVASCRIPT / values*
[/col]
[/row]
[row]
[col align=left]
```js{5}
typeof x === "string" ? 1 : 0
const [head, ...tail] = list
for (const key in obj) out[key] = obj[key]
const f = (x) => Array.isArray(x) ? f(x[0]) : x
s.startsWith("gid://") ? s.slice(6) : null
```
[/col]
[/row]
[/grid]
[/col]
[col width=62 border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / types*
[/col]
[/row]
[row]
[col align=left]
```ts{5}
X extends string ? 1 : 0
L extends [infer Head, ...infer Tail] ? Head : never
{ [Key in keyof Obj]: Obj[Key] }
type F<X> = X extends (infer I)[] ? F<I> : X
S extends `gid://${infer Id}` ? Id : null
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Matching a string becomes a template literal type: match a pattern, `infer` the rest.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "Where did my autocomplete go?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Where did my autocomplete go?**  ·  until 07:45

  - A search function: q is the wildcard "*", or any query string
  - We'd like the editor to suggest "*"
  - It suggests nothing, and the hover says q is just string
  - Think back to the sets slide

  **NEXT →** The answer
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*PUZZLE / AUTOCOMPLETE   ·   11 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Where did my autocomplete go?
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=88]
[grid gap=1]
[row gap=0]
[col]
[/col]
[col width=60 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a search function*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
function search(q: "*" | string) {}

search("");
//      ^ autocomplete suggests nothing

// hover over q:
// (parameter) q: string
```
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[row]
[col align=center]
**We asked for `"*" | string` but got `string` back. Where did the wildcard go?**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "Where did my autocomplete go? The answer"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Where did my autocomplete go? The answer**  ·  until 09:30

  - "*" is already inside string, so TypeScript simplifies "*" | string down to string
  - string & {}: every string is also a {}, so it's still every string
  - But TypeScript won't simplify it away, so "*" survives, and so does autocomplete
  - typesense-ts uses exactly this for q

  **NEXT →** Never forgetting a case
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*PUZZLE / AUTOCOMPLETE   ·   12 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Where did my autocomplete go?
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=124]
[grid gap=1]
[row gap=3]
[col width=64 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / the fix*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
// "*" fits inside string, so the union collapses
type Q1 = "*" | string;          // string

// string & {} is still every string,
// but TypeScript won't merge "*" into it
type Q2 = "*" | (string & {});   // "*" | (string & {})
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / typesense-ts, src/search.ts, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{3}
interface QueryParams<
  Fields extends CollectionField[],
  Q extends "*" | (string & {}),
> {
  q: Q;
  query_by: QueryBy<Fields>;
}
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**`string & {}` holds every string, but TypeScript never merges it with `"*"` so the wildcard stays.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "Mapped types are loops"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Mapped types are loops**  ·  until 10:15

  - The loop over keys from the table: for every R in Resource, write one property
  - So Schemas becomes an object with all three keys, each holding a Schema
  - From a tool I'm building: it syncs a few kinds of resources into search, and each needs a schema
  - Forget articles, and the error names the missing one
  - Add a fourth resource next year, and this breaks until you give it a schema
  - Record<Resource, Schema> is the built-in shorthand for exactly this

  **NEXT →** How the hook works
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MECHANISMS / MAPPED TYPES   ·   13 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Mapped types are loops
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=124]
[grid gap=1]
[row gap=3]
[col width=62 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a mapped type, and what it becomes*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
type Resource = "products" | "collections" | "articles";

// for every R in Resource, write one property
type Schemas = { [R in Resource]: Schema };

// Schemas is now:
// {
//   products: Schema;
//   collections: Schema;
//   articles: Schema;
// }
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: forget one*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
const schemas: Schemas = {
  products: productsSchema,
  collections: collectionsSchema,
};

// Property 'articles' is missing in type
//   '{ products: Schema; collections: Schema; }'
//   but required in type 'Schemas'.
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**One property for every member of the union. Forget one, and it doesn't compile.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]

----
---
title: "Parsing strings with infer"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Parsing strings with infer**  ·  until 10:45

  - infer inside backticks: match the pattern, capture the pieces
  - "age:desc" splits at the colon: field "age", direction "desc"

  **NEXT →** Check it against a schema
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MECHANISMS / PARSING STRINGS   ·   14 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Parsing strings with infer
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=104]
[grid gap=1]
[row gap=0]
[col]
[/col]
[col width=88 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a sort_by parser, in miniature*
[/col]
[/row]
[row]
[col align=left]
```ts{1-6}
type ParseSort<S> =
  S extends `${infer Field}:${infer Direction}`
    ? { field: Field; direction: Direction }
    : never;

type P = ParseSort<"age:desc">;   // { field: "age"; direction: "desc" }

type Sortable = "age" | "score";
type CheckSort<S> =
  S extends `${infer Field}:${infer Direction}`
    ? Field extends Sortable
      ? Direction extends "asc" | "desc" ? S : `${Direction} is not asc or desc.`
      : `${Field} is not a sortable field.`
    : "Expected field:asc or field:desc";

declare function search<const S extends string>(
  params: { sort_by: CheckSort<S> },
): void;

search({ sort_by: "email:desc" });
// Type '"email:desc"' is not assignable
//   to type '"email is not a sortable field."'.
```
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[row]
[col align=center]
**`infer` inside a template literal splits a string wherever the pattern says.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "Parsing strings: checking the pieces"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Parsing strings: checking the pieces**  ·  until 11:15

  - Same pattern as ParseSort, then check each piece
  - The field against the sortable ones, the direction against asc and desc
  - If both fit, give the string back. If not, give back an error message, as a string type

  **NEXT →** Put it in a function
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MECHANISMS / PARSING STRINGS   ·   15 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Parsing strings with infer
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=104]
[grid gap=1]
[row gap=0]
[col]
[/col]
[col width=88 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a sort_by parser, in miniature*
[/col]
[/row]
[row]
[col align=left]
```ts{8-14}
type ParseSort<S> =
  S extends `${infer Field}:${infer Direction}`
    ? { field: Field; direction: Direction }
    : never;

type P = ParseSort<"age:desc">;   // { field: "age"; direction: "desc" }

type Sortable = "age" | "score";
type CheckSort<S> =
  S extends `${infer Field}:${infer Direction}`
    ? Field extends Sortable
      ? Direction extends "asc" | "desc" ? S : `${Direction} is not asc or desc.`
      : `${Field} is not a sortable field.`
    : "Expected field:asc or field:desc";

declare function search<const S extends string>(
  params: { sort_by: CheckSort<S> },
): void;

search({ sort_by: "email:desc" });
// Type '"email:desc"' is not assignable
//   to type '"email is not a sortable field."'.
```
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[row]
[col align=center]
**Check the pieces against the schema. The answer can be an error message.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "Parsing strings: the error"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Parsing strings: the error**  ·  until 11:45

  - Use it as the parameter's type
  - A valid string is assignable to itself. An invalid one isn't assignable to the message
  - So the message lands in the error, right where you typed it
  - That's the hook from the start, in about twenty lines

  **NEXT →** The real one
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MECHANISMS / PARSING STRINGS   ·   16 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Parsing strings with infer
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=104]
[grid gap=1]
[row gap=0]
[col]
[/col]
[col width=88 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a sort_by parser, in miniature*
[/col]
[/row]
[row]
[col align=left]
```ts{16-22}
type ParseSort<S> =
  S extends `${infer Field}:${infer Direction}`
    ? { field: Field; direction: Direction }
    : never;

type P = ParseSort<"age:desc">;   // { field: "age"; direction: "desc" }

type Sortable = "age" | "score";
type CheckSort<S> =
  S extends `${infer Field}:${infer Direction}`
    ? Field extends Sortable
      ? Direction extends "asc" | "desc" ? S : `${Direction} is not asc or desc.`
      : `${Field} is not a sortable field.`
    : "Expected field:asc or field:desc";

declare function search<const S extends string>(
  params: { sort_by: CheckSort<S> },
): void;

search({ sort_by: "email:desc" });
// Type '"email:desc"' is not assignable
//   to type '"email is not a sortable field."'.
```
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[row]
[col align=center]
**Make it the parameter's type, and the error shows up where you call it.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "Parsing strings: the real one"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Parsing strings: the real one**  ·  until 12:30

  - The real thing does it properly: a tokenizer, then a validator
  - ReadString peels one character off with infer Head and infer Tail
  - Then it calls itself with the rest: recursion
  - The tokens get checked against the schema, and that's where the hook's error came from

  **NEXT →** Why recursion?
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MECHANISMS / PARSING STRINGS   ·   17 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Parsing strings with infer
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=132]
[grid gap=1]
[row gap=3]
[col width=60 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / typesense-ts, src/lexer/types.ts, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{3,5}
// read an identifier, one character at a time
type ReadString<T, Acc extends string = ""> =
  T extends `${infer Head}${infer Tail}`
    ? Head extends Char
      ? ReadString<Tail, `${Acc}${Head}`>
      : [Acc, T]
    : [Acc, T];
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / typesense-ts, tests/sort.test.ts*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
expectTypeOf<SortTokenizer<"age:desc">>().toEqualTypeOf<
  [Ident<"age", FieldType>, Colon, "desc"]
>();
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**typesense-ts reads `sort_by` one character at a time, then checks every token against your schema.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "Why do I need recursion?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Why do I need recursion?**  ·  until 13:00

  - There's no for or while in the type system, so a type that loops has to recurse
  - TupleOf keeps adding one T until the length matches
  - A tuple's length is a literal type, like 3, not just number
  - So tuples are how you do arithmetic: add an element, read the length

  **NEXT →** Where that's useful
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MECHANISMS / RECURSION   ·   18 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Why do I need recursion?
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=96]
[grid gap=1]
[row gap=0]
[col]
[/col]
[col width=74 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / loops and counting*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
// no for, no while: a type that calls itself is the loop
type TupleOf<L extends number, T, Acc extends T[] = []> =
  Acc["length"] extends L ? Acc : TupleOf<L, T, [...Acc, T]>;

type Three = TupleOf<3, number>;   // [number, number, number]

// a tuple's length is a number literal, so you can count
type AddOne<N extends number> = [...TupleOf<N, 1>, 1]["length"];
type Four = AddOne<3>;   // 4
```
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[row]
[col align=center]
**Types have no loops. Recursion is the loop, and tuple lengths are how you count.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "Recursion: one weight per field"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Recursion: one weight per field**  ·  until 13:30

  - typesense-ts: one weight for each field you search
  - The length of query_by is a number literal, 3
  - query_by_weights is a tuple of exactly that many numbers
  - Typesense would reject it at runtime. Here the editor does first

  **NEXT →** How far can it go?
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*MECHANISMS / RECURSION   ·   19 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Why do I need recursion?
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=132]
[grid gap=1]
[row gap=3]
[col width=64 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / typesense-ts*
[/col]
[/row]
[row]
[col align=left]
```ts{6-7}
// query_by_weights?: TupleOfLength<L, number>
// where L is the length of query_by

await books.search({
  q: "dune",
  query_by: ["title", "authors", "categories"],
  query_by_weights: [3],
});
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TSC / the error*
[/col]
[/row]
[row]
[col align=left]
```text{1-99}
error TS2322: Type '[number]' is not assignable
  to type '[number, number, number]'.
  Source has 1 element(s) but target requires 3.
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**One weight per `query_by` field, checked by counting at the type level.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "How far can it go?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **How far can it go?**  ·  until 15:00

  - So types can loop, and tuples can count. How far can that go?
  - Dimitri Mitropoulos, of Michigan TypeScript, set out to show why Doom can't run in TypeScript types
  - The README: "The funny thing is.. It can." February 2025
  - Doom's C source, compiled to WebAssembly, then rewritten as types by a script
  - A computer made of types runs it, one instruction at a time
  - From the launch video: 3.5 trillion lines of types, 177 TB of them, 100 GB of memory
  - A whole app's type check is a few hundred thousand instantiations
  - This did 20 million a second, for 12 days, on a computer in a basement
  - The frame comes out as an object type, one line of pixels per key, that you hover over

  **NEXT →** How the computer works
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*LIMITS / DOOM   ·   20 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# TypeScript types can run DOOM
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=124]
[grid gap=1]
[row]
[col align=center]
![DOOM, as a type, in the editor|112x28](typescript-doom.jpg)
[/col]
[/row]
[row]
[col align=center]
**Every instruction of the game ran in the type checker.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]

----
---
title: "How far can it go? A computer made of types"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **How far can it go? A computer made of types**  ·  until 16:15

  - The real machine is huge, so here is the simplest version in the repo, made to explain it
  - A program is a tuple of instructions. The state is the instructions left, and a result
  - executeInstruction destructures the next instruction with infer
  - selectInstruction runs it, and executeInstruction calls itself with the new state
  - No instructions left: the result is the answer. Here, "HELLO WORLD!"
  - Doom's version has 116 instructions, memory and a stack, but the same loop

  **NEXT →** How it does logic
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*LIMITS / DOOM   ·   21 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# A computer made of types
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=124]
[grid gap=1]
[row gap=3]
[col width=44 border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a program*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
type Execution = executeInstruction<{
  instructions: [
    { kind: 'add', value: 'hello' },
    { kind: 'add', value: ' world' },
    { kind: 'uppercase' },
    { kind: 'add', value: '!' },
  ];
  result: '';
}>;
// "HELLO WORLD!"
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / the simplest machine in the repo, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{4-7,8-10}
type executeInstruction<
  state extends ProgramState
> =
  state['instructions'] extends [
    infer instruction extends Instruction,
    ...infer remainingInstructions extends Instruction[]
  ]
  ? executeInstruction<
      selectInstruction<state, remainingInstructions, instruction>
    >
  : state['result'];
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Take the next instruction, run it, call yourself. Destructure and recursion, from the table.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]

----
---
title: "How far can it go? Logic with a lookup table"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **How far can it go? Logic with a lookup table**  ·  until 17:45

  - Every number in the machine is a string of 0s and 1s
  - AND goes digit by digit: 1 only where both digits are 1
  - The rule is a table: look up A's digit, then B's
  - Walk both strings: first digit off each, look it up, recurse on the rest
  - 1100 AND 1010 is 1000. OR and XOR are the same walk, with a different table
  - Our TupleOf gave up at 1000 levels. Dimitri patched the compiler to lift those limits while developing
  - And paused every 50 to 100 instructions, printed the whole state, and carried on from there

  **NEXT →** Meanwhile, the compiler got faster
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*LIMITS / DOOM   ·   22 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Logic with a lookup table
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=124]
[grid gap=1]
[row gap=3]
[col width=62 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / ts-type-math/bitwise.ts, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
type Bit = '0' | '1';

type And = {
  '0': { '0': '0', '1': '0' },
  '1': { '0': '0', '1': '1' },
}

type BitwiseAnd<A extends string, B extends string> =
  A extends `${infer AHead extends Bit}${infer ATail}`
  ? B extends `${infer BHead extends Bit}${infer BTail}`
    ? `${And[AHead][BHead]}${BitwiseAnd<ATail, BTail>}`
    : ''
  : ''

type R = BitwiseAnd<'1100', '1010'>;   // '1000'
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TEXT / 1100 AND 1010*
[/col]
[/row]
[row]
[col align=left]
```text{1-99}
'1100'  AND  '1010'

1 and 1  →  1
1 and 0  →  0
0 and 1  →  0
0 and 0  →  0

→ '1000'
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Doom's processor is built from tables and recursion, like this one.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]

----
---
title: "TypeScript 7: the compiler, rewritten in Go"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **TypeScript 7: the compiler, rewritten in Go**  ·  until 18:15

  - Some news before we finish: TypeScript 7 came out in July
  - The compiler is rewritten in Go. Same type system, about 10x faster
  - Numbers from the 7.0 announcement, full builds of big open source projects
  - VS Code's own codebase: two minutes down to ten seconds
  - The miniature examples from today give the same errors on 7.1
  - And 7.1, in beta right now, adds something new

  **NEXT →** Content mappers
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*WHAT'S NEW / TYPESCRIPT 7   ·   23 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# TypeScript 7: the compiler, rewritten in Go
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=88]
[grid gap=1]
[row gap=0]
[col]
[/col]
[col width=66 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TEXT / full build, TypeScript 6 → 7*
[/col]
[/row]
[row]
[col align=left]
```text{1-99}
vscode        125.7s   →   10.6s     11.9x
sentry        139.8s   →   15.7s      8.9x
bluesky        24.3s   →    2.8s      8.7x
playwright     12.8s   →    1.47s     8.7x
tldraw         11.2s   →    1.46s     7.7x
```
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[row]
[col align=center]
**The same type system, about 10x faster.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]

----
---
title: "TypeScript 7.1: content mappers"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **TypeScript 7.1: content mappers**  ·  until 19:15

  - Until now tsc checked .ts and .js files, and nothing else
  - A .vue, .svelte or .sql import was whatever type you declared by hand
  - 7.1 adds content mappers: a plugin that turns a file into TypeScript
  - It also maps every piece back, so an error shows up in the original file
  - You register one per extension in tsconfig, and allow it with --runExternalCode
  - Vue, Svelte, Astro, MDX and Ember have mappers prototyped or in progress

  **NEXT →** The one I like most
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*WHAT'S NEW / TYPESCRIPT 7.1   ·   24 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# TypeScript 7.1: content mappers
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=124]
[grid gap=1]
[row gap=3]
[col width=60 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TEXT / what a content mapper does*
[/col]
[/row]
[row]
[col align=left]
```text{1-99}
before   tsc checks .ts and .js files, nothing else
         a .vue, .svelte or .sql import is a type
         you write by hand

7.1      a mapper turns the file into TypeScript
         tsc checks it, and reports errors
         in the original file
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*JSONC / tsconfig.json, trimmed*
[/col]
[/row]
[row]
[col align=left]
```jsonc{3}
{
  "contentMappers": [
    { "package": "reylo", "extensions": [".sql"] }
  ]
}
```
[/col]
[/row]
[row]
[col align=left]
```text{1-99}
$ tsc --noEmit --runExternalCode
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**A content mapper lets tsc check files that aren't TypeScript, and point errors back at them.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]

----
---
title: "Reylo: SQL files, typed by your database"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Reylo: SQL files, typed by your database**  ·  until 20:15

  - Reylo, by Colin McDonnell, who also wrote Zod. 1.0 is planned alongside 7.1 stable
  - You write plain SQL in a .sql file and import it
  - The mapper asks your dev database to prepare the query, without running it
  - The database answers with the type of every column and parameter
  - email can be null, because the schema doesn't say not null
  - And | null: a match on the primary key is one row or none
  - No hand-written row types, no codegen step

  **NEXT →** And when the SQL is wrong?
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*WHAT'S NEW / CONTENT MAPPERS   ·   25 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Reylo: SQL files, typed by your database
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=124]
[grid gap=1]
[row gap=3]
[col width=40 border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*SQL / get_user.sql*
[/col]
[/row]
[row]
[col align=left]
```sql{1-99}
select id, name, email
from users
where id = :id::int
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / main.ts, from Reylo's docs*
[/col]
[/row]
[row]
[col align=left]
```ts{1,3}
import getUser from "./get_user.sql";

const user = await getUser.run(db, { id: 1 });
//    ^? const user: {
//           id: number;
//           name: string;
//           email: string | null;
//       } | null
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Import a .sql file, and its types come from your database.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]

----
---
title: "Reylo: a typo in the SQL is a type error"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Reylo: a typo in the SQL is a type error**  ·  until 21:15

  - Misspell a column, and your database rejects the query while it's being typed
  - tsc reports that error inside the .sql file, at the line and column
  - Same in the editor, as you type. No editor plugin, the TypeScript extension shows it
  - Rename a column in the database, and every query that reads it breaks on the next check
  - Earlier we parsed strings with types. Reylo doesn't: it asks the one thing that knows the answer

  **NEXT →** So where should you use all this?
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*WHAT'S NEW / CONTENT MAPPERS   ·   26 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Reylo: a typo in the SQL is a type error
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=124]
[grid gap=1]
[row gap=3]
[col width=40 border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*SQL / with a typo*
[/col]
[/row]
[row]
[col align=left]
```sql{1}
select id, nmae, email
from users
where id = :id::int
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TSC / from Reylo's docs*
[/col]
[/row]
[row]
[col align=left]
```text{2,5}
$ tsc --noEmit --runExternalCode
get_user.sql:1:12 - error reylo6: column "nmae" does not exist (42703)

1 select id, nmae, email
             ~~~~

Found 1 error in get_user.sql:1
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**The database checks the query, and tsc reports it in the .sql file.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]

----
---
title: "So where should you use it?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **So where should you use it?**  ·  until 22:15

  - In app code, a little goes a long way: the autocomplete trick, mapped records
  - And let tools do the heavy lifting: Reylo types your SQL so you don't write those types
  - In a library, spend it at the boundary: where users pass schemas, strings and options
  - It costs compile time. Doom took 12 days. Yours shouldn't, so measure it
  - It costs error quality: a typo in a filter can surface as Unknown token: !
  - And someone has to maintain it. Test your types: expectTypeOf, and @ts-expect-error

  **NEXT →** Thank you
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*CLOSING / WHERE TO USE IT   ·   27 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# So where should you use it?
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=128]
[grid gap=1]
[row]
[col]
[grid gap=1]
[row gap=3]
[col border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*IN APP CODE / a little goes a long way*
[/col]
[/row]
[row]
[col align=left]
```text{1-99}
"*" | (string & {})   suggest values, take any string
mapped records        cover every member of a union
content mappers       let a tool type .sql and .vue
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*IN A LIBRARY / spend it at the boundary*
[/col]
[/row]
[row]
[col align=left]
```text{1-99}
derive types from the schema users give you
check the strings users write
put the fix in the error message
test the types like the rest of the API
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*WHAT IT COSTS*
[/col]
[/row]
[row]
[col align=left]
```text{1-99}
compile time     Doom took 12 days. Measure yours: tsc --generateTrace
error messages   a typo in a filter can come back as "Unknown token: !"
maintenance      someone has to read your types next year
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Use the least wizardry that makes the wrong call impossible.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]
----
---
title: "Thank you. Questions?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Thank you. Questions?**  ·  until 22:45

  - That's everything, thank you
  - QR on the left: the slides and every link, Doom video included
  - Some of today's puzzles were inspired by TypeHero
  - My GitHub and site on the right

  **IF ASKED**

  - Why not just zod? zod checks data at runtime, these types check code before it runs. They work together: the Shopify app uses zod for API responses
  - Isn't this slow? Measure it: --generateTrace, then analyze-trace. typesense-ts generates a trace on every type check
  - When is it too much? When the error messages stop helping, or nobody can maintain it
  - Does Reylo parse SQL in types? No. Its content mapper asks the dev database to prepare the query
  - Did Doom need a patched compiler? Dimitri lifted the limits while developing. The video says the result runs in vanilla TypeScript, by pausing every 50 to 100 instructions

  **END →** leave this slide up for the Q&A
---

[grid]
[row height=2]
[col]
[/col]
[/row]
[row height=2]
[col align=center]
*Q&A   ·   28 / 28*
[/col]
[/row]
[row height=1fr]
[col valign=middle]
[grid gap=1]
[row]
[col align=center]
# Thank you. Questions?
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=135]
[grid gap=1]
[row gap=4]
[col width=74 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*SLIDES & LINKS / scan me*
[/col]
[/row]
[row gap=3]
[col width=36 valign=middle]
![QR code for the slides and links|36x18](gist_qr.png)
[/col]
[col align=left valign=middle]
```text{1-99}
The slides, as a PDF
TypeHero
Type Challenges
TypeScript Handbook
TypeScript types can run DOOM
Michigan TypeScript
Reylo
ArkType
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*$ whoami*
[/col]
[/row]
[row gap=2]
[col width=18 valign=middle]
![Typebrain|18x9](typebrain.jpg)
[/col]
[col align=left]
```text{99-99}
[1;38;2;245;197;0mfanis[0m[38;2;245;237;184m@[0m[1;38;2;245;197;0mfosscomm[0m
[38;2;245;237;184m──────────────[0m
[1;38;2;245;197;0mName:[0m [38;2;245;237;184mFanis Tharropoulos[0m
[1;38;2;245;197;0mWork:[0m [38;2;245;237;184mTypesense[0m
[1;38;2;245;197;0mGitHub:[0m [38;2;245;237;184mgithub.com/tharropoulos[0m
[1;38;2;245;197;0mSite:[0m [38;2;245;237;184mfanis.iusevimbtw.com[0m
[1;38;2;245;197;0mLibrary:[0m [38;2;245;237;184mtypesense-ts[0m
[1;38;2;245;197;0mSlides:[0m [38;2;245;237;184mmade with Kyma[0m
[1;38;2;245;197;0mStatus:[0m [38;2;245;237;184mtaking questions[0m
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Ask me anything, the weirder the type the better.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
[/grid]
[/col]
[/row]
[/grid]

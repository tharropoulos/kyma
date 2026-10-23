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

[row gap=0]
[col]
[/col]
[col width=88]
[grid gap=1]
[row]
[col]
[grid gap=3]
[row]
[col align=center]
*FOSSCOMM 2026   ·   01 / 29*
[/col]
[/row]
[row]
[col align=center]
# TypeScript's Type Wizardry
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
## (and where to use it)
[/col]
[/row]
[row]
[col align=center]
![Typebrain|24x12](typebrain.jpg)
[/col]
[/row]
[row]
[col align=center pad="1 0"]
**Fanis Tharropoulos**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]

----
---
title: "A little about me"
transition: slideUp
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **A little about me**  ·  until 01:45

  - Software engineer at Typesense, an open source search engine
  - I work on all of its official API clients and integrations
  - My own projects: Phonto, Stochos, and Kyma, which made this deck
  - And typesense-ts: most of today's wizardry comes from it

  **NEXT →** Let's start with something that shouldn't be possible
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=3]
[row]
[col align=center]
*$whoami ·   02 / 29*
[/col]
[/row]
[row]
[col align=center]
# Who Am I?
[/col]
[/row]
[/grid]
[/col]
[/row]
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

----
---
title: "The compiler read my string"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **The compiler read my string**  ·  until 03:00

  - A search on a users collection, sorted by email, descending
  - email is optional, and you can't sort on it
  - Normally the server tells you, at runtime
  - Here the editor tells you, before anything runs
  - That error message is a type: the compiler parsed the string
  - By the end of the talk, you'll know how this works

  **NEXT →** First, a mental model
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*INTRO / THE HOOK   ·   03 / 29*
[/col]
[/row]
[row]
[col align=center]
# The compiler read my string
[/col]
[/row]
[/grid]
[/col]
[/row]
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

----
---
title: "Types are sets"
transition: slideUp
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Types are sets**  ·  until 04:30

  - The one idea everything else builds on: a type is a set of values
  - "a" | "b" has two values, string has every string
  - unknown has every value, never has none: the empty set
  - extends asks one question: does the left side fit inside the right?
  - "a" fits in AB, AB fits in string, but string doesn't fit in AB
  - Every trick today is a question about sets

  **NEXT →** And the type system is a small language
---

[row gap=0]
[col]
[/col]
[col width=114]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MENTAL MODEL / TYPES ARE SETS   ·   04 / 29*
[/col]
[/row]
[row]
[col align=center]
# Types are sets
[/col]
[/row]
[/grid]
[/col]
[/row]
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

// does the left side fit inside the right?
type T1 = "a" extends AB ? true : false;     // true
type T2 = AB extends string ? true : false;  // true
type T3 = string extends AB ? true : false;  // false
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

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MENTAL MODEL / A SMALL LANGUAGE   ·   05 / 29*
[/col]
[/row]
[row]
[col align=center]
# The type system is a small language
[/col]
[/row]
[/grid]
[/col]
[/row]
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
map
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
list.map((x) => [x])
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
U extends unknown ? [U] : never
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

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MENTAL MODEL / A SMALL LANGUAGE   ·   05 / 29*
[/col]
[/row]
[row]
[col align=center]
# The type system is a small language
[/col]
[/row]
[/grid]
[/col]
[/row]
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
map
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
list.map((x) => [x])
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
U extends unknown ? [U] : never
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
----
---
title: "A small language: destructure"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **A small language: destructure**  ·  until 06:00

  - const [head, ...tail] = list pulls values out of an array
  - infer does the same to a type: match the shape, name the pieces
  - Here Head is the first element of the tuple L

  **NEXT →** Loops over keys
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MENTAL MODEL / A SMALL LANGUAGE   ·   05 / 29*
[/col]
[/row]
[row]
[col align=center]
# The type system is a small language
[/col]
[/row]
[/grid]
[/col]
[/row]
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
map
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
list.map((x) => [x])
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
U extends unknown ? [U] : never
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
----
---
title: "A small language: loop over keys"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **A small language: loop over keys**  ·  until 06:30

  - for...in walks the keys of an object
  - A mapped type walks the keys of a type: Key in keyof Obj
  - One property out for every key in. This one copies Obj as it is

  **NEXT →** map
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MENTAL MODEL / A SMALL LANGUAGE   ·   05 / 29*
[/col]
[/row]
[row]
[col align=center]
# The type system is a small language
[/col]
[/row]
[/grid]
[/col]
[/row]
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
map
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
list.map((x) => [x])
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
U extends unknown ? [U] : never
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
----
---
title: "A small language: map"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **A small language: map**  ·  until 07:00

  - A conditional type on a union runs once for every member, like .map
  - "a" | "b" becomes ["a"] | ["b"], not ["a" | "b"]
  - This is the union that splits itself apart. More on it later

  **NEXT →** Recursion
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MENTAL MODEL / A SMALL LANGUAGE   ·   05 / 29*
[/col]
[/row]
[row]
[col align=center]
# The type system is a small language
[/col]
[/row]
[/grid]
[/col]
[/row]
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
map
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
list.map((x) => [x])
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
U extends unknown ? [U] : never
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
**map becomes distribution: a conditional type runs once per union member.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "A small language: recursion"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **A small language: recursion**  ·  until 07:30

  - Types have no for or while: a type that refers to itself is the loop
  - F unwraps one array per call, until X is not an array anymore
  - F<number[][][]> is number

  **NEXT →** Matching strings
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MENTAL MODEL / A SMALL LANGUAGE   ·   05 / 29*
[/col]
[/row]
[row]
[col align=center]
# The type system is a small language
[/col]
[/row]
[/grid]
[/col]
[/row]
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
map
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
list.map((x) => [x])
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
U extends unknown ? [U] : never
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
----
---
title: "A small language: match a string"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **A small language: match a string**  ·  until 08:00

  - Template literal types match a string against a pattern
  - infer Id captures whatever follows gid://
  - "gid://shopify/Product/1" gives "shopify/Product/1", anything else gives null
  - The sort_by parser from the start is built out of exactly this

  **NEXT →** A puzzle: where did my autocomplete go?
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MENTAL MODEL / A SMALL LANGUAGE   ·   05 / 29*
[/col]
[/row]
[row]
[col align=center]
# The type system is a small language
[/col]
[/row]
[/grid]
[/col]
[/row]
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
```text{6}
if
destructure
loop over keys
map
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
```js{6}
typeof x === "string" ? 1 : 0
const [head, ...tail] = list
for (const key in obj) out[key] = obj[key]
list.map((x) => [x])
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
```ts{6}
X extends string ? 1 : 0
L extends [infer Head, ...infer Tail] ? Head : never
{ [Key in keyof Obj]: Obj[Key] }
U extends unknown ? [U] : never
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
----
---
title: "Where did my autocomplete go?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Where did my autocomplete go?**  ·  until 08:30

  - A search function: q is the wildcard "*", or any query string
  - We'd like the editor to suggest "*"
  - It suggests nothing, and the hover says q is just string
  - Think back to the sets slide

  **NEXT →** The answer
---

[row gap=0]
[col]
[/col]
[col width=88]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*PUZZLE / AUTOCOMPLETE   ·   06 / 29*
[/col]
[/row]
[row]
[col align=center]
# Where did my autocomplete go?
[/col]
[/row]
[/grid]
[/col]
[/row]
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
----
---
title: "Where did my autocomplete go? The answer"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Where did my autocomplete go? The answer**  ·  until 09:15

  - "*" is already inside string, so TypeScript simplifies "*" | string down to string
  - string & {}: every string is also a {}, so it's still every string
  - But TypeScript won't simplify it away, so "*" survives, and so does autocomplete
  - typesense-ts uses exactly this for q
  - And later, if Q is "*", the result type has no highlights

  **NEXT →** Why did true become boolean?
---

[row gap=0]
[col]
[/col]
[col width=124]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*PUZZLE / AUTOCOMPLETE   ·   06 / 29*
[/col]
[/row]
[row]
[col align=center]
# Where did my autocomplete go?
[/col]
[/row]
[/grid]
[/col]
[/row]
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
----
---
title: "Why did true become boolean?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Why did true become boolean?**  ·  until 09:45

  - let can be reassigned, so "desc" widens to string
  - const can't, so it stays the literal "desc"
  - An object's properties can change, so true widens to boolean
  - Same when you pass an object to a generic function
  - A library that reads options at the type level loses information here

  **NEXT →** How to keep the literal
---

[row gap=0]
[col]
[/col]
[col width=88]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / LITERAL TYPES   ·   07 / 29*
[/col]
[/row]
[row]
[col align=center]
# Why did true become boolean?
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=64 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / widening*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
let sort = "desc";                 // string
const order = "desc";              // "desc"
const opts = { facet: true };      // { facet: boolean }

function field<T>(options: T): T {
  return options;
}
field({ facet: true });            // { facet: boolean }
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
**TypeScript widens a literal whenever it thinks you might change it.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Keeping the literal"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Keeping the literal**  ·  until 10:30

  - as const keeps literals, and makes everything readonly
  - A const type parameter does the same for whatever is passed in, no as const at the call site
  - From a tool I'm building: a schema for search collections
  - f.string({ facet: true }) keeps true, so the type knows the field is facetable
  - The catch: pass a variable typed boolean, and the type can't know if the field is optional
  - In that tool, such a field silently becomes required

  **NEXT →** keyof, typeof and indexed access
---

[row gap=0]
[col]
[/col]
[col width=132]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / LITERAL TYPES   ·   07 / 29*
[/col]
[/row]
[row]
[col align=center]
# Why did true become boolean?
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=50 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / keeping the literal*
[/col]
[/row]
[row]
[col align=left]
```ts{1,4}
const opts = { facet: true } as const;
// { readonly facet: true }

function field<const T>(options: T): T {
  return options;
}
field({ facet: true });
// { readonly facet: true }
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: field options a schema can read*
[/col]
[/row]
[row]
[col align=left]
```ts{3,6,9}
type FieldOptions = { facet?: boolean; optional?: boolean };

const string = <const O extends FieldOptions = {}>(options?: O) =>
  ({ kind: "string", options: options ?? ({} as O) });

string({ facet: true }).options.facet;   // true

declare const optional: boolean;
string({ optional }).options.optional;   // boolean
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**A `const` type parameter keeps the literal. Pass a `boolean` and the schema can't tell.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "keyof, typeof and indexed access"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **keyof, typeof and indexed access**  ·  until 11:15

  - keyof: the union of an object type's keys
  - Product["price"]: look up a property's type, like obj[key]
  - typeof, in a type position: the type of a value
  - From a tool I'm building: naming things in a store's search index
  - One as const array is the runtime list and the type, so they never drift
  - [number]: index with any number and you get every element, as a union
  - And the return type is "s1_products". How? Template literal types, in a few slides

  **NEXT →** Mapped types
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / KEYOF, TYPEOF, T[K]   ·   08 / 29*
[/col]
[/row]
[row]
[col align=center]
# Types from types, and from values
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=60 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / the basics*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
type Product = { id: string; price: number };

type Keys = keyof Product;        // "id" | "price"
type Price = Product["price"];    // number

const product = { id: "1", price: 9 };
type FromValue = typeof product;
// { id: string; price: number }
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: names in a search index*
[/col]
[/row]
[row]
[col align=left]
```ts{4,11-12}
const ALIASES =
  ["products", "collections", "articles", "pages"] as const;

type Alias = (typeof ALIASES)[number];
// "products" | "collections" | "articles" | "pages"

const resourceName = <Ns extends string, R extends Alias>(
  ns: Ns, r: R,
): `${Ns}_${R}` => `${ns}_${r}`;

resourceName("s1", "products");   // "s1_products"
resourceName("s1", "product");    // error
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**One array is both the runtime list and the type. `[number]` reads every element.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Mapped types are loops"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Mapped types are loops**  ·  until 11:45

  - A mapped type walks the keys of a type, like for...in
  - For each key K you write the property you want
  - Add a ? and every property is optional: that's Partial
  - Add readonly and you have Readonly. Record and Pick are mapped types too

  **NEXT →** Where this pays off
---

[row gap=0]
[col]
[/col]
[col width=88]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / MAPPED TYPES   ·   09 / 29*
[/col]
[/row]
[row]
[col align=center]
# Mapped types are loops
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=66 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / mapped types*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
type Product = { id: string; price: number };

// for every key K in T, write one property
type MyPartial<T>  = { [K in keyof T]?: T[K] };
type MyReadonly<T> = { readonly [K in keyof T]: T[K] };

type Draft = MyPartial<Product>;
// { id?: string; price?: number }
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
**Like `for (const key in obj)` one level up: one property out for every key in.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Mapped types: one schema per resource"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Mapped types: one schema per resource**  ·  until 12:30

  - From a tool I'm building: it syncs a few kinds of resources into search
  - Every resource needs a schema, and SCHEMAS holds one for each
  - Typed as a mapped type over the union, it has to cover every member
  - Add "articles" to the union, forget its schema: compile error
  - Record<SyncResource, ...> does the same. Record is a mapped type

  **NEXT →** Conditional types
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / MAPPED TYPES   ·   09 / 29*
[/col]
[/row]
[row]
[col align=center]
# Mapped types are loops
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=76 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: a schema per synced resource*
[/col]
[/row]
[row]
[col align=left]
```ts{1,4}
type SyncResource = "products" | "collections" | "articles";

const SCHEMAS: {
  readonly [R in SyncResource]: (name: string) => CollectionSchema;
} = {
  products: productsSchema,
  collections: collectionsSchema,
};
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TSC / the error, abridged*
[/col]
[/row]
[row]
[col align=left]
```text{1}
error TS2741: Property 'articles' is missing
  in type '{ products: ...; collections: ...; }'
  but required in type
  '{ readonly products: ...;
     readonly collections: ...;
     readonly articles: ...; }'
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Add a resource, forget its schema, and it doesn't compile.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Conditional types and infer"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Conditional types and infer**  ·  until 13:00

  - extends, from the sets slide: does the left side fit inside the right?
  - Yes or no picks the branch
  - infer only works inside extends: describe a shape, name the part you want
  - That's how the built-in ReturnType is written
  - The same trick pulls the element type out of an array

  **NEXT →** From a tool I'm building
---

[row gap=0]
[col]
[/col]
[col width=100]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / CONDITIONAL TYPES   ·   10 / 29*
[/col]
[/row]
[row]
[col align=center]
# Conditional types and infer
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=84 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / if, and destructuring*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
// if: does T fit inside string?
type IsString<T> = T extends string ? true : false;
type A = IsString<"a">;    // true
type B = IsString<42>;     // false

// infer: match a shape, name the piece you want
type MyReturnType<F> = F extends (...args: any[]) => infer R ? R : never;
type R = MyReturnType<() => Promise<number>>;   // Promise<number>

type ElementOf<T> = T extends (infer E)[] ? E : never;
type E = ElementOf<string[]>;   // string
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
**Use `extends` for if and `infer` for destructuring: match a shape, name the piece.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Conditional types: what a field holds"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Conditional types: what a field holds**  ·  until 13:45

  - From a tool I'm building: a schema for search collections
  - A field definition carries its value type and its options as type parameters
  - FieldValue pulls both out with infer
  - Then one more conditional: array: true gives T[], otherwise T
  - One piece of turning a schema into a document type. The rest needs two more tricks

  **NEXT →** Dropping keys with as
---

[row gap=0]
[col]
[/col]
[col width=120]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / CONDITIONAL TYPES   ·   10 / 29*
[/col]
[/row]
[row]
[col align=center]
# Conditional types and infer
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=58 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: what a field holds*
[/col]
[/row]
[row]
[col align=left]
```ts{2}
type FieldValue<F> =
  F extends FieldDef<infer TBase, infer TOpts>
    ? TOpts extends { array: true }
      ? TBase[]
      : TBase
    : never;
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / using it*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
const tags = f.string({ array: true });
const price = f.float();

type Tags = FieldValue<typeof tags>;     // string[]
type Price = FieldValue<typeof price>;   // number
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Two `infer` clauses pull out the value type and the options, which pick `T[]` or `T` here.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Dropping keys with as"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Dropping keys with as**  ·  until 14:15

  - as lets a mapped type rewrite each key. Map one to never and it's dropped
  - From a tool I'm building: per-field options, passed through to the search engine
  - The client library's field type ends in an index signature: any key goes, typos included
  - string extends K is only true when K is string itself, the catch-all
  - "infix" is one string, not all of them, so it stays
  - Drop the catch-all and a typo is an error again, with a suggestion

  **NEXT →** Picking keys out of a type
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / KEY REMAPPING   ·   11 / 29*
[/col]
[/row]
[row]
[col align=center]
# Dropping keys with as
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=62 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / as rewrites each key*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
// map a key to never and it's gone
type NoId<T> = {
  [K in keyof T as K extends "id" ? never : K]: T[K];
};

type P = NoId<{ id: string; price: number }>;
// { price: number }
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: engine options per field*
[/col]
[/row]
[row]
[col align=left]
```ts{5,9,13-14}
// from the search engine's client library
interface CollectionFieldSchema {
  name: string;
  infix?: boolean;
  [t: string]: unknown;   // swallows any typo
}

type KnownKeys<T> = {
  [K in keyof T as string extends K ? never : K]: T[K];
};
type Override = Omit<KnownKeys<CollectionFieldSchema>, "name">;

const o: Override = { infux: true };
// Did you mean to write 'infix'?
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**`string extends K` is only true for the catch-all key, so that's the only one dropped.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Map, then index"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Map, then index**  ·  until 14:45

  - The pattern you'll see most: map, then index
  - Map: each key becomes its own name if it passes the test, never if it doesn't
  - Index with [keyof F]: a union of all the values
  - never is the empty set, so it vanishes from a union
  - What's left is exactly the keys that passed

  **NEXT →** Building a document type out of it
---

[row gap=0]
[col]
[/col]
[col width=92]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / MAP, THEN INDEX   ·   12 / 29*
[/col]
[/row]
[row]
[col align=center]
# Map, then index: picking keys
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=72 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / picking the optional keys*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
type Fields = {
  id:    { type: "string" };
  brand: { type: "string"; optional: true };
  price: { type: "float" };
};

type OptionalKeys<F> = {
  [K in keyof F]: F[K] extends { optional: true } ? K : never;
}[keyof F];

// map:    { id: never; brand: "brand"; price: never }
// index:  never | "brand" | never
// result: "brand", because never is the empty set
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
**Map each key to itself or `never` and index with `keyof` to drop the rest.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Map, then index: a document type"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Map, then index: a document type**  ·  until 15:30

  - From the same tool: the document type, derived from the schema
  - Required fields: every key except the optional ones. Optional fields: those, with a ?
  - Each value comes from FieldValue, the infer slide
  - The result is an intersection, and the hover shows it as one
  - Prettify maps over it once, so the editor shows one flat object

  **NEXT →** The same pattern, for business rules
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / MAP, THEN INDEX   ·   12 / 29*
[/col]
[/row]
[row]
[col align=center]
# Map, then index: picking keys
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=76 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: a schema's document type*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
type InferDoc<F> = Prettify<
  { [K in Exclude<keyof F, OptionalKeys<F>>]: FieldValue<F[K]> } &
  { [K in OptionalKeys<F>]?: FieldValue<F[K]> }
>;

type Prettify<T> = { [K in keyof T]: T[K] } & {};
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*EDITOR / the hover*
[/col]
[/row]
[row]
[col align=left]
```text{5-7}
without Prettify
  { id: string; price: number; }
  & { brand?: string | undefined; }

with Prettify
  { id: string; price: number;
    brand?: string | undefined; }
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Required keys and optional keys, intersected. `Prettify` flattens it for the hover.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Map, then index: business rules"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Map, then index: business rules**  ·  until 16:00

  - Same pattern, different test
  - From the same tool: which field can be the collection's default sort
  - Numbers sort by default, strings only when asked, optional fields never
  - So the allowed keys are a union, computed from the schema
  - Pick title, a string without sort: true, and the error says what's allowed

  **NEXT →** A puzzle: why did my union split apart?
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / MAP, THEN INDEX   ·   12 / 29*
[/col]
[/row]
[row]
[col align=center]
# Map, then index: picking keys
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=66 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: the default sort, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
type DefaultSortKeys<F> = {
  [K in keyof F]: F[K] extends { optional: true } ? never
    : F[K] extends { type: "float" } ? K
    : F[K] extends { sort: true } ? K
    : never;
}[keyof F];
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / using it*
[/col]
[/row]
[row]
[col align=left]
```ts{7,9}
defineCollection({
  fields: {
    id: { type: "string" },
    title: { type: "string" },
    price: { type: "float" },
  },
  defaultSort: "title",
});
// Type '"title"' is not assignable to type '"price"'.
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**A business rule becomes a union of allowed keys, and the error names them.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Why did my union split apart?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Why did my union split apart?**  ·  until 16:30

  - T extends unknown is always true, so this should just be T[]
  - Pass string | number and you'd expect one array that holds both
  - You get two arrays, each holding only one of them
  - The union split itself apart

  **NEXT →** Why
---

[row gap=0]
[col]
[/col]
[col width=88]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*PUZZLE / DISTRIBUTION   ·   13 / 29*
[/col]
[/row]
[row]
[col align=center]
# Why did my union split apart?
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=60 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / one union in*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
type ToArray<T> = T extends unknown ? T[] : never;

type A = ToArray<string | number>;

// expected: (string | number)[]
// got:      string[] | number[]
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
**We passed in one union. It came back as two arrays.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Why did my union split apart? Distribution"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Why did my union split apart? Distribution**  ·  until 17:00

  - The rule: a conditional type on a bare type parameter distributes over a union
  - It runs once per member and unions the results: the map row from the table
  - Extract is built on it: keep the members that fit, never for the rest
  - From a tool I'm building: a batch of writes is retried only if writing it twice is safe
  - create isn't, it fails the second time. So passing it is a compile error

  **NEXT →** Turning distribution off
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*PUZZLE / DISTRIBUTION   ·   13 / 29*
[/col]
[/row]
[row]
[col align=center]
# Why did my union split apart?
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=56 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / what happened*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
// a conditional type on a bare type parameter
// runs once for every member of a union
ToArray<string | number>
  = ToArray<string> | ToArray<number>
  = string[] | number[]

// the built-in Extract relies on it
type Extract<T, U> = T extends U ? T : never;
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: only writes that are safe to retry*
[/col]
[/row]
[row]
[col align=left]
```ts{3,6-7}
type Action = "create" | "update" | "upsert" | "emplace";

type RetrySafe = Extract<Action, "upsert" | "update">;
// "update" | "upsert"

importDocuments(docs, { action: "create" });
// Type '"create"' is not assignable to type 'RetrySafe'.
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Distribution is `.map` over a union. `Extract` uses it to filter one.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Turning distribution off"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Turning distribution off**  ·  until 17:45

  - Distribution only happens when the checked type is a bare type parameter
  - Wrap both sides in a tuple and the union is compared as one
  - Brain teaser: never is an empty union, so distributing over it maps over nothing and gives never
  - From a tool I'm building: it lists duration units that have to match a library's list
  - Fitting in both directions means equal, and Assert fails to compile if not
  - A test that runs in the compiler, in production code, at zero runtime cost

  **NEXT →** Template literal types
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / NO DISTRIBUTION   ·   14 / 29*
[/col]
[/row]
[row]
[col align=center]
# Turning distribution off
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=68 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / wrap it in a tuple*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
type ToArray<T> = [T] extends [unknown] ? T[] : never;
type A = ToArray<string | number>;   // (string | number)[]

// a brain teaser: never is an empty union
type IsNever<T> = T extends never ? true : false;
type X = IsNever<never>;   // never: nothing to map over

type IsNever2<T> = [T] extends [never] ? true : false;
type Y = IsNever2<never>;  // true
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: a check that never runs*
[/col]
[/row]
[row]
[col align=left]
```ts{1,5-6}
type Assert<T extends true> = T;

// our duration units must match the library's
type _UnitsMatch = Assert<
  [Ours] extends [Theirs]
    ? [Theirs] extends [Ours] ? true : false
    : false
>;
// if they drift apart: Type 'false' does
// not satisfy the constraint 'true'.
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Only a bare type parameter distributes. Wrap both sides in `[ ]` and the union stays whole.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Template literal types"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Template literal types**  ·  until 18:30

  - Backticks in a type: a pattern a string has to match
  - ${number} is any number, so "1.5 seconds" fits and "200 milis" doesn't
  - From a tool I'm building: retry delays as readable strings, typos caught
  - Put unions inside and you get every combination
  - And the return type from earlier: the namespace, an underscore, the resource

  **NEXT →** Parsing strings
---

[row gap=0]
[col]
[/col]
[col width=132]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / TEMPLATE LITERAL TYPES   ·   15 / 29*
[/col]
[/row]
[row]
[col align=center]
# Template literal types
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=60 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: readable retry delays*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
type Unit = "millis" | "seconds" | "minutes";
type Duration = number | `${number} ${Unit}`;

const a: Duration = 200;            // ok
const b: Duration = "1.5 seconds";  // ok
const c: Duration = "200 milis";    // error

// unions multiply
type Method = `${"get" | "set"}${"Id" | "Name"}`;
// "getId" | "getName" | "setId" | "setName"
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / back to the names from earlier*
[/col]
[/row]
[row]
[col align=left]
```ts{3,5}
const resourceName = <Ns extends string, R extends Alias>(
  ns: Ns, r: R,
): `${Ns}_${R}` => `${ns}_${r}`;

resourceName("s1", "products");   // "s1_products"
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**A type can describe the shape of a string, and build new strings out of old ones.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Parsing strings with infer"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Parsing strings with infer**  ·  until 19:00

  - infer inside backticks: match the pattern, capture the pieces
  - "age:desc" splits at the colon: field "age", direction "desc"

  **NEXT →** Check it against a schema
---

[row gap=0]
[col]
[/col]
[col width=104]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / PARSING STRINGS   ·   16 / 29*
[/col]
[/row]
[row]
[col align=center]
# Parsing strings with infer
[/col]
[/row]
[/grid]
[/col]
[/row]
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
```ts{1-4}
type ParseSort<S> =
  S extends `${infer Field}:${infer Dir}` ? { field: Field; dir: Dir } : never;

type P = ParseSort<"age:desc">;   // { field: "age"; dir: "desc" }

type Sortable = "age" | "score";
type CheckSort<S> =
  S extends `${infer Field}:${"asc" | "desc"}`
    ? Field extends Sortable ? S : `${Field} is not a sortable field.`
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
----
---
title: "Parsing strings: checking the pieces"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Parsing strings: checking the pieces**  ·  until 19:30

  - Now check the field against the sortable ones
  - If it fits, give the string back. If not, give back an error message, as a string type

  **NEXT →** Put it in a function
---

[row gap=0]
[col]
[/col]
[col width=104]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / PARSING STRINGS   ·   16 / 29*
[/col]
[/row]
[row]
[col align=center]
# Parsing strings with infer
[/col]
[/row]
[/grid]
[/col]
[/row]
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
```ts{6-10}
type ParseSort<S> =
  S extends `${infer Field}:${infer Dir}` ? { field: Field; dir: Dir } : never;

type P = ParseSort<"age:desc">;   // { field: "age"; dir: "desc" }

type Sortable = "age" | "score";
type CheckSort<S> =
  S extends `${infer Field}:${"asc" | "desc"}`
    ? Field extends Sortable ? S : `${Field} is not a sortable field.`
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
----
---
title: "Parsing strings: the error"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Parsing strings: the error**  ·  until 20:00

  - Use it as the parameter's type
  - A valid string is assignable to itself. An invalid one isn't assignable to the message
  - So the message lands in the error, right where you typed it
  - That's the hook from the start, in about fifteen lines

  **NEXT →** The real one
---

[row gap=0]
[col]
[/col]
[col width=104]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / PARSING STRINGS   ·   16 / 29*
[/col]
[/row]
[row]
[col align=center]
# Parsing strings with infer
[/col]
[/row]
[/grid]
[/col]
[/row]
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
```ts{12-18}
type ParseSort<S> =
  S extends `${infer Field}:${infer Dir}` ? { field: Field; dir: Dir } : never;

type P = ParseSort<"age:desc">;   // { field: "age"; dir: "desc" }

type Sortable = "age" | "score";
type CheckSort<S> =
  S extends `${infer Field}:${"asc" | "desc"}`
    ? Field extends Sortable ? S : `${Field} is not a sortable field.`
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
----
---
title: "Parsing strings: the real one"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Parsing strings: the real one**  ·  until 20:30

  - The real thing does it properly: a tokenizer, then a validator
  - ReadString peels one character off with infer Head and infer Tail
  - Then it calls itself with the rest: recursion
  - The tokens get checked against the schema, and that's where the hook's error came from

  **NEXT →** Why recursion?
---

[row gap=0]
[col]
[/col]
[col width=132]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / PARSING STRINGS   ·   16 / 29*
[/col]
[/row]
[row]
[col align=center]
# Parsing strings with infer
[/col]
[/row]
[/grid]
[/col]
[/row]
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
----
---
title: "Why do I need recursion?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Why do I need recursion?**  ·  until 21:00

  - There's no for or while in the type system, so a type that loops has to recurse
  - TupleOf keeps adding one T until the length matches
  - A tuple's length is a literal type, like 3, not just number
  - So tuples are how you do arithmetic: add an element, read the length

  **NEXT →** Where that's useful
---

[row gap=0]
[col]
[/col]
[col width=96]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / RECURSION   ·   17 / 29*
[/col]
[/row]
[row]
[col align=center]
# Why do I need recursion?
[/col]
[/row]
[/grid]
[/col]
[/row]
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
----
---
title: "Recursion: one weight per field"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Recursion: one weight per field**  ·  until 21:45

  - typesense-ts: one weight for each field you search
  - The length of query_by is a number literal, 3
  - query_by_weights is a tuple of exactly that many numbers
  - Typesense would reject it at runtime. Here the editor does first

  **NEXT →** How deep can it go?
---

[row gap=0]
[col]
[/col]
[col width=132]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / RECURSION   ·   17 / 29*
[/col]
[/row]
[row]
[col align=center]
# Why do I need recursion?
[/col]
[/row]
[/grid]
[/col]
[/row]
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
----
---
title: "How deep can it go?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **How deep can it go?**  ·  until 22:30

  - The compiler protects itself: too deep and it gives up with TS2589
  - The same type written two ways. Build the result after the call returns, and 48 is already too deep
  - Make the recursive call the whole answer, carry the result in an accumulator: TypeScript treats that as tail recursion, up to 999
  - When type checking gets slow: --generateTrace, then analyze-trace points at the slow types
  - typesense-ts generates a trace on every type check

  **NEXT →** Part 3: designing library APIs
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*MECHANISMS / LIMITS   ·   18 / 29*
[/col]
[/row]
[row]
[col align=center]
# How deep can it go?
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=76 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / the same type, written two ways*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
// builds its result after the recursive call returns
type Naive<L extends number, T, Acc extends unknown[] = []> =
  Acc["length"] extends L ? [] : [T, ...Naive<L, T, [...Acc, 1]>];

type A = Naive<47, 0>;       // ok
type B = Naive<48, 0>;       // error TS2589

// the recursive call is the whole answer: tail recursion
type C = TupleOf<999, 0>;    // ok
type D = TupleOf<1000, 0>;   // error TS2589
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TEXT / when it gets slow*
[/col]
[/row]
[row]
[col align=left]
```text{1-2}
TS2589: Type instantiation is excessively
deep and possibly infinite.

- keep the recursive call in tail position,
  and carry the result in an accumulator
- find slow types:
  tsc --generateTrace trace/
  npx @typescript/analyze-trace trace/
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**As tail recursion it reaches 999. Written naively it fails at 48.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Where does TypeScript infer from?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Where does TypeScript infer from?**  ·  until 23:15

  - From a tool I'm building: it syncs rows from a database into search engines
  - TypeScript infers a type parameter from the arguments where it appears as is
  - C comes from collection, TRow from source
  - InferDoc<C> is computed from C. TypeScript never tries to infer C from it
  - So by the time you write map, row already has a type, wherever map sits in the object
  - And map has to return exactly the document type the schema describes

  **NEXT →** When the wrong argument decides
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / INFERENCE   ·   19 / 29*
[/col]
[/row]
[row]
[col align=center]
# Where does TypeScript infer from?
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=76 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: syncing a database into search*
[/col]
[/row]
[row]
[col align=left]
```ts{2-4}
interface SyncConfig<C, TRow> {
  collection: C;                       // C comes from here
  source: Source<TRow>;                // TRow comes from here
  map: (row: TRow) => InferDoc<C>;     // computed, never inferred
  targets: readonly Target<InferDoc<C>>[];
}

declare function createSync<const C extends CollectionDef, TRow>(
  config: SyncConfig<C, TRow>,
): Sync<C>;
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / using it*
[/col]
[/row]
[row]
[col align=left]
```ts{2-3,10}
createSync({
  collection: products,
  source: memorySource(rows),
  map: (row) => ({
    id: String(row.id),
    name: row.name,
  }),
  targets: [typesense],
});
// row: { id: number; name: string; price: number }
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Give each type parameter one obvious source. Compute everything else from it.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "When the wrong argument decides the type"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **When the wrong argument decides the type**  ·  until 23:45

  - From a tool I'm building: it stores sync state, and every document has a kind
  - Ids start with their kind, so a schema id looks like schema:products
  - put takes a kind, an id and data, and all three should agree
  - Here they don't: a build id, stored as a schema
  - It compiles. Every place K appears gets a vote, and K widens until both fit

  **NEXT →** Taking away a vote
---

[row gap=0]
[col]
[/col]
[col width=100]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / NOINFER   ·   20 / 29*
[/col]
[/row]
[row]
[col align=center]
# When the wrong argument decides the type
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=80 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: storing sync state*
[/col]
[/row]
[row]
[col align=left]
```ts{5,9-10}
type StateId<K extends StateKind> = `${K}:${string}`;

declare function put<K extends StateKind>(
  kind: K,
  id: StateId<K>,
  data: StateData<K>,
): void;

put("schema", stateIds.build("products", "1"), { schemaHash: "a" });
// compiles: K widens to "schema" | "build", so both fit
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
**Every argument gets a vote on `K` so it widens until both arguments fit.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "When the wrong argument decides: NoInfer"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **When the wrong argument decides: NoInfer**  ·  until 24:30

  - NoInfer, since TypeScript 5.4: use the type here, but don't infer from it
  - Now only kind decides K
  - The id and the data have to match it, or it's an error at the call site
  - I found this one the hard way, while building that tool

  **NEXT →** Arguments that depend on earlier arguments
---

[row gap=0]
[col]
[/col]
[col width=100]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / NOINFER   ·   20 / 29*
[/col]
[/row]
[row]
[col align=center]
# When the wrong argument decides the type
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=80 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: storing sync state*
[/col]
[/row]
[row]
[col align=left]
```ts{3-4,8-9}
declare function put<K extends StateKind>(
  kind: K,
  id: StateId<NoInfer<K>>,
  data: NoInfer<StateData<K>>,
): void;

put("schema", stateIds.build("products", "1"), { schemaHash: "a" });
// Argument of type '`build:${string}:${string}`' is not
//   assignable to parameter of type '`schema:${string}`'.
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
**`NoInfer` takes away a vote. Only `kind` decides `K` and the rest must agree.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "What you can pass depends on what you passed"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **What you can pass depends on what you passed**  ·  until 25:15

  - The question from the abstract: arguments that depend on what you passed earlier
  - collection() keeps the whole schema as a type, thanks to const type parameters
  - search() reads it: facet_by takes facetable fields, sort_by sortable ones
  - query_by_weights, from the recursion slide, follows the length of query_by
  - And q: "*" turns highlights off in the result type

  **NEXT →** It works across calls, too
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / DEPENDENT ARGUMENTS   ·   21 / 29*
[/col]
[/row]
[row]
[col align=center]
# What you can pass depends on what you passed
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=60 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / typesense-ts: a schema*
[/col]
[/row]
[row]
[col align=left]
```ts{5-6}
const books = collection({
  name: "books",
  fields: [
    { name: "title", type: "string" },
    { name: "year", type: "int32", sort: true },
    { name: "genre", type: "string", facet: true },
  ],
});
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / typesense-ts: what search accepts*
[/col]
[/row]
[row]
[col align=left]
```ts{4-5,8-9}
await books.search({
  q: "dune",
  query_by: ["title"],
  facet_by: ["genre"],      // only facetable fields
  sort_by: "year:desc",     // only sortable fields
});

facet_by: ["title"]
// Type '"title"' is not assignable to type '"genre"'.
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**The schema you passed first decides which fields every later option accepts.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "What you can pass: one result per search"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **What you can pass: one result per search**  ·  until 26:00

  - multisearch takes a tuple of searches, each against any registered collection
  - const keeps it a tuple, so every position keeps its own type
  - A mapped type over a tuple gives back a tuple, so result 0 matches search 0
  - books gives you titles, authors gives you names

  **NEXT →** Making invalid states impossible to write
---

[row gap=0]
[col]
[/col]
[col width=106]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / DEPENDENT ARGUMENTS   ·   21 / 29*
[/col]
[/row]
[row]
[col align=center]
# What you can pass depends on what you passed
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=0]
[col]
[/col]
[col width=90 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / typesense-ts, src/http/fetch/multisearch.ts, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{1,4}
async function multisearch<const Searches extends readonly object[]>(
  params: { searches: [...Searches] },
): Promise<{
  results: { [K in keyof Searches]: MultiSearchResultEntry<Searches[K]> };
}>;
```
[/col]
[/row]
[row]
[col align=left]
```ts{8-9}
const { results } = await multisearch({
  searches: [
    multisearchEntry({ collection: "books", q: "dune", query_by: ["title"] }),
    multisearchEntry({ collection: "authors", q: "frank", query_by: ["name"] }),
  ],
});

results[0].hits.map((hit) => hit.document.title);   // string[]
results[1].hits.map((hit) => hit.document.name);    // string[]
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
**A mapped type over a tuple stays a tuple, so every result keeps its own collection's type.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Making invalid states impossible to write"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Making invalid states impossible to write**  ·  until 26:45

  - Some combinations of options are wrong, and the server rejects them at runtime
  - ?: never: the key exists in the type, but no value fits, so you can only leave it out
  - So a plain string field can't take vector options
  - A conditional type does the reverse: once there's an object field, the flag is required
  - Typesense rejects both at runtime. typesense-ts rejects them as you type

  **NEXT →** A registry, filled in by your code
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / INVALID STATES   ·   22 / 29*
[/col]
[/row]
[row]
[col align=center]
# Making invalid states impossible to write
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=64 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / typesense-ts, src/collection/base.ts, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{4-6,10-13}
// a plain field can't carry vector options
type RegularField = BaseField & {
  type: FieldType;
  embed?: never;
  num_dim?: never;
  vec_dist?: never;
};

// an object field needs the nested fields flag
type EnforceNestedFields<Fields> =
  ObjectFields<Fields> extends never
    ? { enable_nested_fields?: boolean }
    : { enable_nested_fields: true };
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 2" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / typesense-ts: what you get*
[/col]
[/row]
[row]
[col align=left]
```ts{3,5,9-10}
collection({
  name: "users",
  fields: [{ name: "address", type: "object" }],
});
// Property 'enable_nested_fields' is missing
//   in type '{ name: "users"; fields: [...]; }'
//   but required in type '...'.

{ name: "title", type: "string", vec_dist: "cosine" }
// compile error: vec_dist is never here
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**`?: never` forbids a key outright. A conditional type can make one required.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "A registry, filled in by your code"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **A registry, filled in by your code**  ·  until 27:30

  - Maybe the biggest wait-it-can-do-that of the talk
  - Interfaces with the same name merge, even across modules
  - The library declares an empty Collections interface and uses it in its types
  - Your code adds your collections to it, so the library can look them up by name
  - A join filter is then parsed against the other collection's schema
  - And it checks the two collections are actually linked by a reference

  **NEXT →** Everything at once
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / DECLARATION MERGING   ·   23 / 29*
[/col]
[/row]
[row]
[col align=center]
# A registry, filled in by your code
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=80 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / typesense-ts: registering collections*
[/col]
[/row]
[row]
[col align=left]
```ts{3,6}
// in the library: an empty interface
interface Collections {}

// in your code: the same interface, merged
declare module "typesense-ts" {
  interface Collections {
    customers: typeof customers.schema;
    orders: typeof orders.schema;  // customer_id references customers.id
  }
}
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / joins, checked across collections*
[/col]
[/row]
[row]
[col align=left]
```ts{4,7-9}
await customers.search({
  q: "*",
  query_by: ["name"],
  filter_by: "$orders(total:>100)",
});

filter_by: "$authors(name:=Frank)"
// [Error on filter_by]: Collection
//   `authors` not referenced in `customers`
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**The library declares an empty interface. Your code fills it in, and the library's types see it.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Everything at once: names from a client"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Everything at once: names from a client**  ·  until 28:15

  - From a tool I'm building: every request to Typesense is wrapped, and named for logs and errors
  - Before: the name was any string, so it could say export while the code imports
  - After: the name is one of 78 operations, derived from the client's own types
  - It autocompletes, a typo is an error, and it has to match the method f actually calls
  - How does the promise know which method made it? A label that only exists in the types

  **NEXT →** The type behind it
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / THE FINALE   ·   24 / 29*
[/col]
[/row]
[row]
[col align=center]
# Everything at once: names from a client
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=76 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: naming every request*
[/col]
[/row]
[row]
[col align=left]
```ts{8-11}
// before
request<A>(
  operation: string,
  f: (client: Client) => Promise<A>,
): Call<A>;

// after
request<O extends Operation, A>(
  operation: O,
  f: (client: LabelledClient) => Labelled<Promise<A>, NoInfer<O>>,
): Call<A>;
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / using it*
[/col]
[/row]
[row]
[col align=left]
```ts{4-7,9-11}
request("collections.documents.import", (c) =>
  c.collections(name).documents().import(docs));

request("collections.documents.export", (c) =>
  c.collections(name).documents().import(docs));
// '"collections.documents.import"' is not
//   assignable to '"collections.documents.export"'

request("collections.documents.mport", ...);
// not assignable to '"debug.retrieve" | ...
//   70 more ... | "curationSets.items.delete"'
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Every request is named, the name autocompletes, and it has to match the call.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Everything at once: the type"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Everything at once: the type**  ·  until 28:45

  - The type behind it, about fifteen lines
  - Every mechanism from the talk is in here
  - Let's go through them one at a time

  **NEXT →** Distribution
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / THE FINALE   ·   24 / 29*
[/col]
[/row]
[row]
[col align=center]
# Everything at once: names from a client
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=80 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: the Operation type, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{1-99}
type Requests<T, Path extends string> = T extends unknown
  ? {
      [K in keyof T & string]: T[K] extends Request
        ? `${Path}${K}`
        : Overloads<T[K]> extends object
          ? Requests<Overloads<T[K]>, `${Path}${K}.`>
          : never;
    }[keyof T & string]
  : never;

type Overloads<F> =
  F extends { (...a: infer _A): infer R1; (...a: infer _B): infer R2 }
    ? R1 | R2
    : F extends (...args: never) => infer R ? R : F;

type Operation = Requests<Omit<Client, "configuration" | "apiCall">, "">;
// "collections.create" | "collections.documents.import" | ... 78 in all
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TEXT / what's in it*
[/col]
[/row]
[row]
[col align=left]
```text{1-99}
distribution       T extends unknown
mapped type        [K in keyof T & string]
conditional type   T[K] extends Request
template literal   `${Path}${K}`
recursion          Requests<..., `${Path}${K}.`>
map, then index    }[keyof T & string]
infer              infer R1, infer R2
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Every mechanism from today, in one type that names 78 requests.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Everything at once: the type, distribution"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Everything at once: the type, distribution**  ·  until 29:00

  - Overloads gives back a union: collections() and collections(name) return different things
  - T extends unknown is always true, it's only there to distribute over that union

  **NEXT →** Mapped type
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / THE FINALE   ·   24 / 29*
[/col]
[/row]
[row]
[col align=center]
# Everything at once: names from a client
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=80 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: the Operation type, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{1,9}
type Requests<T, Path extends string> = T extends unknown
  ? {
      [K in keyof T & string]: T[K] extends Request
        ? `${Path}${K}`
        : Overloads<T[K]> extends object
          ? Requests<Overloads<T[K]>, `${Path}${K}.`>
          : never;
    }[keyof T & string]
  : never;

type Overloads<F> =
  F extends { (...a: infer _A): infer R1; (...a: infer _B): infer R2 }
    ? R1 | R2
    : F extends (...args: never) => infer R ? R : F;

type Operation = Requests<Omit<Client, "configuration" | "apiCall">, "">;
// "collections.create" | "collections.documents.import" | ... 78 in all
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TEXT / what's in it*
[/col]
[/row]
[row]
[col align=left]
```text{1}
distribution       T extends unknown
mapped type        [K in keyof T & string]
conditional type   T[K] extends Request
template literal   `${Path}${K}`
recursion          Requests<..., `${Path}${K}.`>
map, then index    }[keyof T & string]
infer              infer R1, infer R2
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**`T extends unknown` distributes, so a union of resources is walked one member at a time.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Everything at once: the type, mapped type"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Everything at once: the type, mapped type**  ·  until 29:15

  - For every key of the client: debug, collections, aliases and so on

  **NEXT →** Conditional type
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / THE FINALE   ·   24 / 29*
[/col]
[/row]
[row]
[col align=center]
# Everything at once: names from a client
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=80 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: the Operation type, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{3}
type Requests<T, Path extends string> = T extends unknown
  ? {
      [K in keyof T & string]: T[K] extends Request
        ? `${Path}${K}`
        : Overloads<T[K]> extends object
          ? Requests<Overloads<T[K]>, `${Path}${K}.`>
          : never;
    }[keyof T & string]
  : never;

type Overloads<F> =
  F extends { (...a: infer _A): infer R1; (...a: infer _B): infer R2 }
    ? R1 | R2
    : F extends (...args: never) => infer R ? R : F;

type Operation = Requests<Omit<Client, "configuration" | "apiCall">, "">;
// "collections.create" | "collections.documents.import" | ... 78 in all
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TEXT / what's in it*
[/col]
[/row]
[row]
[col align=left]
```text{2}
distribution       T extends unknown
mapped type        [K in keyof T & string]
conditional type   T[K] extends Request
template literal   `${Path}${K}`
recursion          Requests<..., `${Path}${K}.`>
map, then index    }[keyof T & string]
infer              infer R1, infer R2
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**A mapped type visits every key of the client, like the loop from the table.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Everything at once: the type, conditional type"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Everything at once: the type, conditional type**  ·  until 29:30

  - A method that returns a promise is a request
  - Anything else that returns an object is a resource with more requests inside

  **NEXT →** Template literal
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / THE FINALE   ·   24 / 29*
[/col]
[/row]
[row]
[col align=center]
# Everything at once: names from a client
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=80 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: the Operation type, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{3,5}
type Requests<T, Path extends string> = T extends unknown
  ? {
      [K in keyof T & string]: T[K] extends Request
        ? `${Path}${K}`
        : Overloads<T[K]> extends object
          ? Requests<Overloads<T[K]>, `${Path}${K}.`>
          : never;
    }[keyof T & string]
  : never;

type Overloads<F> =
  F extends { (...a: infer _A): infer R1; (...a: infer _B): infer R2 }
    ? R1 | R2
    : F extends (...args: never) => infer R ? R : F;

type Operation = Requests<Omit<Client, "configuration" | "apiCall">, "">;
// "collections.create" | "collections.documents.import" | ... 78 in all
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TEXT / what's in it*
[/col]
[/row]
[row]
[col align=left]
```text{3}
distribution       T extends unknown
mapped type        [K in keyof T & string]
conditional type   T[K] extends Request
template literal   `${Path}${K}`
recursion          Requests<..., `${Path}${K}.`>
map, then index    }[keyof T & string]
infer              infer R1, infer R2
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**A conditional type asks of every key: is this a request, or a resource to walk into?**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Everything at once: the type, template literal"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Everything at once: the type, template literal**  ·  until 29:45

  - collections. plus documents. plus import gives collections.documents.import

  **NEXT →** Recursion
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / THE FINALE   ·   24 / 29*
[/col]
[/row]
[row]
[col align=center]
# Everything at once: names from a client
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=80 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: the Operation type, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{4}
type Requests<T, Path extends string> = T extends unknown
  ? {
      [K in keyof T & string]: T[K] extends Request
        ? `${Path}${K}`
        : Overloads<T[K]> extends object
          ? Requests<Overloads<T[K]>, `${Path}${K}.`>
          : never;
    }[keyof T & string]
  : never;

type Overloads<F> =
  F extends { (...a: infer _A): infer R1; (...a: infer _B): infer R2 }
    ? R1 | R2
    : F extends (...args: never) => infer R ? R : F;

type Operation = Requests<Omit<Client, "configuration" | "apiCall">, "">;
// "collections.create" | "collections.documents.import" | ... 78 in all
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TEXT / what's in it*
[/col]
[/row]
[row]
[col align=left]
```text{4}
distribution       T extends unknown
mapped type        [K in keyof T & string]
conditional type   T[K] extends Request
template literal   `${Path}${K}`
recursion          Requests<..., `${Path}${K}.`>
map, then index    }[keyof T & string]
infer              infer R1, infer R2
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**A template literal type joins the path so far and the key into a name.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Everything at once: the type, recursion"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Everything at once: the type, recursion**  ·  until 30:00

  - A resource is walked with the same type, its path one level longer

  **NEXT →** Map, then index
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / THE FINALE   ·   24 / 29*
[/col]
[/row]
[row]
[col align=center]
# Everything at once: names from a client
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=80 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: the Operation type, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{6}
type Requests<T, Path extends string> = T extends unknown
  ? {
      [K in keyof T & string]: T[K] extends Request
        ? `${Path}${K}`
        : Overloads<T[K]> extends object
          ? Requests<Overloads<T[K]>, `${Path}${K}.`>
          : never;
    }[keyof T & string]
  : never;

type Overloads<F> =
  F extends { (...a: infer _A): infer R1; (...a: infer _B): infer R2 }
    ? R1 | R2
    : F extends (...args: never) => infer R ? R : F;

type Operation = Requests<Omit<Client, "configuration" | "apiCall">, "">;
// "collections.create" | "collections.documents.import" | ... 78 in all
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TEXT / what's in it*
[/col]
[/row]
[row]
[col align=left]
```text{5}
distribution       T extends unknown
mapped type        [K in keyof T & string]
conditional type   T[K] extends Request
template literal   `${Path}${K}`
recursion          Requests<..., `${Path}${K}.`>
map, then index    }[keyof T & string]
infer              infer R1, infer R2
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Recursion walks one level deeper, with a longer path.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Everything at once: the type, map, then index"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Everything at once: the type, map, then index**  ·  until 30:15

  - The mapped object holds a name, or never, under every key
  - Indexing it with its keys turns it into a union, and the nevers drop out

  **NEXT →** infer
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / THE FINALE   ·   24 / 29*
[/col]
[/row]
[row]
[col align=center]
# Everything at once: names from a client
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=80 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: the Operation type, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{2,8}
type Requests<T, Path extends string> = T extends unknown
  ? {
      [K in keyof T & string]: T[K] extends Request
        ? `${Path}${K}`
        : Overloads<T[K]> extends object
          ? Requests<Overloads<T[K]>, `${Path}${K}.`>
          : never;
    }[keyof T & string]
  : never;

type Overloads<F> =
  F extends { (...a: infer _A): infer R1; (...a: infer _B): infer R2 }
    ? R1 | R2
    : F extends (...args: never) => infer R ? R : F;

type Operation = Requests<Omit<Client, "configuration" | "apiCall">, "">;
// "collections.create" | "collections.documents.import" | ... 78 in all
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TEXT / what's in it*
[/col]
[/row]
[row]
[col align=left]
```text{6}
distribution       T extends unknown
mapped type        [K in keyof T & string]
conditional type   T[K] extends Request
template literal   `${Path}${K}`
recursion          Requests<..., `${Path}${K}.`>
map, then index    }[keyof T & string]
infer              infer R1, infer R2
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Map, then index: the object of names collapses into one union of names.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Everything at once: the type, infer"
transition: morph
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Everything at once: the type, infer**  ·  until 30:30

  - Matching a type with two call signatures captures both return types
  - That's how collections() and collections(name) both get walked
  - 78 names, derived from the client. A new client version brings its new methods along

  **NEXT →** Part 4: testing your types
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*LIBRARY APIS / THE FINALE   ·   24 / 29*
[/col]
[/row]
[row]
[col align=center]
# Everything at once: names from a client
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=80 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: the Operation type, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{12-14}
type Requests<T, Path extends string> = T extends unknown
  ? {
      [K in keyof T & string]: T[K] extends Request
        ? `${Path}${K}`
        : Overloads<T[K]> extends object
          ? Requests<Overloads<T[K]>, `${Path}${K}.`>
          : never;
    }[keyof T & string]
  : never;

type Overloads<F> =
  F extends { (...a: infer _A): infer R1; (...a: infer _B): infer R2 }
    ? R1 | R2
    : F extends (...args: never) => infer R ? R : F;

type Operation = Requests<Omit<Client, "configuration" | "apiCall">, "">;
// "collections.create" | "collections.documents.import" | ... 78 in all
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TEXT / what's in it*
[/col]
[/row]
[row]
[col align=left]
```text{7}
distribution       T extends unknown
mapped type        [K in keyof T & string]
conditional type   T[K] extends Request
template literal   `${Path}${K}`
recursion          Requests<..., `${Path}${K}.`>
map, then index    }[keyof T & string]
infer              infer R1, infer R2
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**`infer` reads both overloads, like `collections()` and `collections(name)` here.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "It compiles. Is it correct?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **It compiles. Is it correct?**  ·  until 31:15

  - Three findings from a review of a tool I'm building. Each one compiles, or fails, the wrong way
  - map can return extra keys: an inferred return type gets no excess property check
  - A target for a different document type is accepted, because method parameters are checked both ways
  - Writing upsert as a property, upsert: (docs) => ..., makes that check strict
  - Too strict: code is one plain union on the base class, so checking it narrows nothing
  - Your types are part of your public API. They need tests like the rest of it

  **NEXT →** How to test them
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*TESTING TYPES / WHY   ·   25 / 29*
[/col]
[/row]
[row]
[col align=center]
# It compiles. Is it correct?
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=80 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: too loose*
[/col]
[/row]
[row]
[col align=left]
```ts{4,7,10,14}
createSync({
  collection: products,
  source: memorySource(rows),
  map: (row) => ({ id: String(row.id), name: row.name, color: "red" }),
  targets: [typesense],
});
// compiles: an inferred return gets no excess property check

interface Target<TDoc> {
  upsert(docs: readonly TDoc[]): Promise<void>;  // a method
}
declare const wide: Target<Doc & { extra: string }>;
const target: Target<Doc> = wide;
// compiles: method parameters are checked both ways
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / a tool I'm building: too strict*
[/col]
[/row]
[row]
[col align=left]
```ts{2,8-10}
class ZetoError extends Error {
  readonly code: "SOURCE_ERROR" | "TARGET_ERROR";
}
class SourceError extends ZetoError {
  readonly source: string;
}

if (e.code === "SOURCE_ERROR") e.source;
// Property 'source' does not exist
//   on type 'ZetoError'.
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Types can be wrong both ways. Too loose lets bugs in, and too strict blocks working code.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Testing types"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Testing types**  ·  until 32:00

  - expectTypeOf: the type has to be exactly this. Vitest and bun:test both ship one
  - @ts-expect-error is a negative test: if the line ever compiles, the unused directive is an error itself
  - Always write the reason after it
  - typesense-ts pairs the compile-time check with a runtime one: the server rejects the same call
  - When the two disagree, one of them has a bug
  - Other tools: tsd, and attest from ArkType, which can also benchmark types

  **NEXT →** Are your type tests even running?
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*TESTING TYPES / HOW   ·   26 / 29*
[/col]
[/row]
[row]
[col align=center]
# Testing types
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=72 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / tools I'm building: the type, and the wrong call*
[/col]
[/row]
[row]
[col align=left]
```ts{1,7-8}
expectTypeOf<InferDoc<typeof products>>().toMatchObjectType<{
  id: string;
  brand?: string;
  tags: string[];
}>();

// @ts-expect-error a build id is not a schema id
const wrong: StateId<"schema"> = stateIds.build("products", "1");
```
[/col]
[/row]
[/grid]
[/col]
[col border=rounded border_color="#3A57CD" pad="1 1" valign=middle]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / typesense-ts, tests/search.test.ts, trimmed*
[/col]
[/row]
[row]
[col align=left]
```ts{5-6,8-10}
await expect(async () => {
  await books.search({
    q: "q",
    query_by: ["foo", "quux", "baz.qux"],
    // @ts-expect-error 3 fields, 1 weight
    query_by_weights: [3],
  });
}).rejects.toThrow(
  "Number of weights in `query_by_weights` does not match",
);
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**Assert the types you promise. `@ts-expect-error` fails the build if the error ever goes away.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "Are your type tests even running?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Are your type tests even running?**  ·  until 32:45

  - expectTypeOf does nothing at runtime. Only a type checker gives it meaning
  - Vitest's typecheck mode only picks up .test-d.ts files by default
  - My projects get away with it because tsc checks the test files too. Make sure yours does
  - @ts-expect-error covers exactly the next line, and a formatter can move the error off it
  - Then the directive is unused and the build fails, which is the good outcome
  - The bad one: it lands on a line with a different error, and hides it
  - One more: check the .d.ts files you publish, not just your source

  **NEXT →** So where should you use all this?
---

[row gap=0]
[col]
[/col]
[col width=138]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*TESTING TYPES / GOTCHAS   ·   27 / 29*
[/col]
[/row]
[row]
[col align=center]
# Are your type tests even running?
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=3]
[col width=66 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*TYPESCRIPT / vitest.config.ts*
[/col]
[/row]
[row]
[col align=left]
```ts{4-7}
export default defineConfig({
  test: {
    typecheck: {
      enabled: true,
      // the default only picks up *.test-d.ts files:
      // ["**/*.{test,spec}-d.?(c|m)[jt]s?(x)"]
      include: ["tests/**/*.test.ts"],
    },
  },
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
*TYPESCRIPT / a formatter moved the error*
[/col]
[/row]
[row]
[col align=left]
```ts{1-3,6-8}
// @ts-expect-error a create is not idempotent
expect(
  importDocuments({ action: "create" }),
).toBe(true);

// Unused '@ts-expect-error' directive.
// Type '"create"' is not assignable to type
//   '"upsert" | "update"'.
```
[/col]
[/row]
[/grid]
[/col]
[/row]
[row]
[col align=center]
**`expectTypeOf` does nothing at runtime. Make sure a type checker actually reads your tests.**
[/col]
[/row]
[/grid]
[/col]
[col]
[/col]
[/row]
----
---
title: "So where should you use it?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **So where should you use it?**  ·  until 33:45

  - In app code, a little goes a long way: as const, mapped records, satisfies, NoInfer, Extract
  - satisfies, if you haven't met it: checks a value against a type without widening the value
  - In a library, spend it at the boundary: where users pass schemas, strings and options
  - It costs compile time, so measure it
  - It costs error quality: a typo in a filter can surface as Unknown token: !
  - And someone has to maintain it. Tests help

  **NEXT →** Thank you
---

[row gap=0]
[col]
[/col]
[col width=128]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*CLOSING / WHERE TO USE IT   ·   28 / 29*
[/col]
[/row]
[row]
[col align=center]
# So where should you use it?
[/col]
[/row]
[/grid]
[/col]
[/row]
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
as const + [number]   one list, runtime and type
mapped records        cover every member of a union
satisfies             check a value, keep its literals
NoInfer               let one argument decide
Extract               narrow a library's union
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
compile time     measure it: tsc --generateTrace, then analyze-trace
error messages   a typo in a filter can come back as "Unknown token: !"
maintenance      someone has to read Requests<T, Path> next year
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
----
---
title: "Thank you. Questions?"
style:
  border: hidden
  theme: /Users/fanis/.cache/phonto/kyma.json
notes: |
  **Thank you. Questions?**  ·  until 34:15

  - That's everything, thank you
  - Links on the left to keep practicing. Some of today's puzzles were inspired by TypeHero
  - My GitHub and site on the right

  **IF ASKED**

  - Why not just zod? zod checks data at runtime, these types check code before it runs. They work together: the Shopify app uses zod for API responses
  - Isn't this slow? Measure it: --generateTrace, then analyze-trace. typesense-ts generates a trace on every type check
  - When is it too much? When the error messages stop helping, or nobody can maintain it

  **END →** leave this slide up for the Q&A
---

[row gap=0]
[col]
[/col]
[col width=124]
[grid gap=1]
[row]
[col]
[grid gap=0]
[row]
[col align=center]
*Q&A   ·   29 / 29*
[/col]
[/row]
[row]
[col align=center]
# Thank you. Questions?
[/col]
[/row]
[/grid]
[/col]
[/row]
[row gap=4]
[col width=74 border=rounded border_color="#3A57CD" pad="1 2"]
[grid gap=1]
[row]
[col align=left]
*LEARN MORE / keep practicing*
[/col]
[/row]
[row]
[col align=left]
```text{1-99}
TypeHero             typehero.dev
Type Challenges      github.com/type-challenges/type-challenges
TypeScript Handbook  typescriptlang.org/docs/handbook/2/
                     types-from-types.html
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
[row]
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

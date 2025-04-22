# research

Half baked ideas. Random shit.

## math expression parsing

RPN is a mathematical notation in which operators follow their operands. Evaluation on this can happen in O(n) time. The notation does not need any parentheses for as long as each operator has a fixed number of operands.

In comparison testing of reverse Polish notation with algebraic notation, reverse Polish has been found to lead to faster calculations, for two reasons. The first reason is that reverse Polish calculators do not need expressions to be parenthesized, so fewer operations need to be entered to perform typical calculations. Additionally, users of reverse Polish calculators made fewer mistakes than for other types of calculators.

Edsger W. Dijkstra invented the shunting-yard algorithm to convert infix expressions to postfix expressions (reverse Polish notation), so named because its operation resembles that of a railroad shunting yard.

Probably super useful for this project: https://en.wikipedia.org/wiki/Shunting_yard_algorithm

just write the parser: https://tiarkrompf.github.io/notes/?/just-write-the-parser

also look at: https://github.com/josdejong/mathjs

and this: https://blog.mbedded.ninja/programming/algorithms-and-data-structures/how-to-parse-mathematical-expressions/

need to be able to deal with variables too.

look into: https://github.com/gnebehay/parser

look into: https://en.wikipedia.org/wiki/Abstract_syntax_tree



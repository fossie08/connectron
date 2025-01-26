# connectron
A-Level Comp sci project - RMB

## checklist
- [x] 0-10 players
- [x] 6x6 to 100x100 grid size
- [x] line length of 4-10 counters needed to win
- [x] counters always drop to the bottom
- [ ] best of 1, 3, 5, 7, etc

### special rules
- [x] counters in any of the corners count as 2 counters
- [x] if a counter is completely surrounded by 1 other player's counters, it is destroyed 
- [x] bomb counter 1 per person per game. destroys all adjacent counters
- [x] overflow rule, if a counter completely fills a column then it overflows into the adjacent columns adding 1 or 2 extra counters

### alliances
- [ ] alliance game, counters from each allied player count as one colour for the sake of winning lines and solitare rule. if winning line is made up of 1 colour, that person gets all the points, otherwise shared.
- [ ] alliances can be changed at the start of each round

### data
- [ ] save games into a text file.
Goal helps you achieve your goals by using strategy and tactics.

> Strategy without tactics is the slowest route to victory. Tactics without strategy is the noise before defeat. — Sun Tzu.

Have you ever seen someone, for example you :-), being or at least appearing busy but not really getting anywhere? If you don't know where you are going, you're going nowhere. You need a goal. If you have a goal but haven't thought about how to achieve it, you will get lost. You need a strategy. If you have strategy you need to get going. You need tactics.

```
$ go install cmd/goal.go

$ goal -example
description: Become a black belt martial artist in under five years.
strategy: Get a personal trainer and train consistently over the next five years.
tactics:
  - do: Find a personal trainer.
    done: true
  - do: Set annual, monthly and weekly goals.
    done: false # can be omitted
  - do: Have a health/diet plan focused on mind, body and spirit.
  - do: Develop a series of minor milestones (to stay motivated).
  - do: Research martial arts instructors in this area.
  - do: Find a ‘training buddy’.
  - do: Find an online community to share ideas and get tips.
  - do: Train on Monday, Tuesday, Thursday and Friday (2 hours per session).
  - do: Write a diet plan.
  - do: Buy training equipment for home use.
  - do: Meditate daily (10 – 30 minutes).
  - do: Develop a ‘rewards’ scheme for minor milestones achieved.

$ mkdir ~/goal
$ goal -example > ~/goal/karate.yaml

$ goal
Goal                           Done
----                           ----
/Users/jozef/goal/karate.yaml  08% (1/12)
```

Inspired by the book The Art of Cyber Security by Gary Hibberd.
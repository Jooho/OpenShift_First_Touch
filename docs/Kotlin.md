IntRange & random()
```
(1..6).random()
```

When
~~~
fun main() {
    val myFirstDice = Dice(6)
    val rollResult = myFirstDice.roll()
    val luckyNumber = 4

    when (rollResult) {
        luckyNumber -> println("You won!")
        1 -> println("So sorry! You rolled a 1. Try again!")
        2 -> println("Sadly, you rolled a 2. Try again!")
        3 -> println("Unfortunately, you rolled a 3. Try again!")
        4 -> println("No luck! You rolled a 4. Try again!")
        5 -> println("Don't cry! You rolled a 5. Try again!")
        6 -> println("Apologies! you rolled a 6. Try again!")
   }
}

class Dice(val numSides: Int) {
    fun roll(): Int {
        return (1..numSides).random()
    }
}
~~~

# Reference
## Not read
- [Vocabulary for Android Basics in Kotlin](https://developer.android.com/codelabs/basic-android-kotlin-training-vocab/#0)
- [Random number generation (Wikipedia)](https://en.wikipedia.org/wiki/Random_number_generation#Practical_applications_and_uses)
- [The Mind-Boggling Challenge of Designing 120-sided Dice](https://www.wired.com/2016/05/mathematical-challenge-of-designing-the-worlds-most-complex-120-sided-dice/)
- [Classes in Kotlin](https://play.kotlinlang.org/byExample/01_introduction/05_Classes)
- [Function declarations in Kotlin](https://kotlinlang.org/docs/reference/functions.html#function-declarations)
- [Returning a value from a function](https://kotlinlang.org/docs/reference/basic-syntax.html#defining-functions)
- [Kotlin style guide](https://developer.android.com/kotlin/style-guide)
- [Control Flow](https://kotlinlang.org/docs/reference/control-flow.html)
- [Accessibility](https://developer.android.com/guide/topics/ui/accessibility)

## Already read

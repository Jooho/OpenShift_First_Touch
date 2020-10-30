**IntRange & random()**
```
(1..6).random()
```

**When**

Ex1
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

Ex2
~~~

val a = "trout"
 
when(a.length){
     1 -> println("length is 1 :$a")
     in 2..50 -> println("length is between 2 to 50 :$a")
     else -> println("length is more than 50 :$a")
 }
~~~

**Elvis Operator**

When we have a nullable reference b, we can say "if b is not null, use it, otherwise use some non-null value":
~~~
val l: Int = if (b != null) b.length else -1
~~~
Along with the complete if-expression, this can be expressed with the Elvis operator, written ?::
~~~
val l = b?.length ?: -1
~~~

Remainder with null value
~~~
println("what date is your birthday?")
var birthday = readLine()?.toIntOrNull() ?:1

or

println("what date is your birthday?")
var birthday = readLine()?.toIntOrNull()
var index = if (birthday == null) 1 else birthday%7
~~~


**For**

Ex1
~~~
val swarm = arrayOf("test1","test2")
for ((index, element) in swarm.withIndex()) {
     println("Fish at $index is $element")
 }
~~~

Ex2
~~~
for (i in 'b'..'g') print(i)
~~~

Ex3
~~~
for (i in 5 downTo 1) print(i)
~~~

Ex4
~~~
for ( i in 3..6 step 2) print(i)
~~~

**$ operator**
~~~
println("Good ${if(Calendar.getInstance().get(Calendar.HOUR_OF_DAY) <=12) "Morning" else "Night"} Kotlin")

val isUnit = println("This is an expression")
println(isUnit)

val temperature = 10
val isHot = temperature > 50
println(isHot)

val message = "You are ${if (temperature > 50) "fried" else "safe"} fish"
println(message)
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

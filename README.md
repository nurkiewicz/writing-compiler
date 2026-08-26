# Write yourself a compiler article series

https://nurkiewicz.com/categories/writing-compiler/

## JVM compiler

`jvm-compiler` accepts the same `integer operator integer` language as the
custom-IR compiler and writes a Java 8 class file to standard output. The
generated class is `com.nurkiewicz.PL0` and uses JVM integer push and
arithmetic instructions before calling `System.out.println(int)`. Its class
metadata identifies the source file as `PL0.pl0`.

```bash
make jvm-compiler
mkdir -p com/nurkiewicz
echo "2 + 3" | ./jvm-compiler > com/nurkiewicz/PL0.class
java -cp . com.nurkiewicz.PL0
# 5
```

The class exposes Java's standard `main(String[] args)` entry point. It can be
embedded in other Java code without generated source:

```java
public class HostApplication {
    public static void main(String[] args) {
        com.nurkiewicz.PL0.main(new String[0]);
    }
}
```

With `PL0.class` in `com/nurkiewicz`, compile and run the host:

```bash
javac -cp . HostApplication.java
java -cp . HostApplication
```

# [String]

Count by String

```
## Sort by Count
cat  $File_Name | tr " " "\n" | sort | uniq -c | sort -nr | awk '{printf("%s %s \n",$2,$1)}END{print}'

## Sort by Name
cat  $File_Name | tr " " "\n" | sort | uniq -c | sort -2knr | awk '{printf("%s %s \n",$2,$1)}END{print}'

Ex)
echo "10 7 2 5 10 2 8 10 2 7 7"|tr " " "\n" | sort | uniq -c | sort -nr | awk '{printf("%s %s \n",$2,$1)}END{print}'
7 3 
2 3 
10 3 
8 1 
5 1 

echo "10 7 2 5 10 2 8 10 2 7 7"|tr " " "\n" | sort | uniq -c | sort -k2nr | awk '{printf("%s %s \n",$2,$1)}END{print}'
10 3 
8 1 
7 3 
5 1 
2 3
```

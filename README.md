# cpic

cpic (Char PICutre) is a tool to draw ASCII char pictures for data strcture.

## Motivation:

Sometimes,I want add some a picture in comments in order to understand a algrithom better,and also inspired by compiler techs like lex and parse.


## Example:

```
tree:
    ->a
        ->b
            ->d
        ->c

to

a   
/\  
b c 
/   
d  
graph:
a -> 1 b,2 c,3 d
b -> a,2 c

to

a-3.000---------------+
 -2.000--------+      |
 -1.000+       |      |
^      |       |      |
|      v       |      |
+------b-2.000+|      |
              ||      |
              vv      |
              c       |
                      |
                      v
                      d

```


# TODO:

*  more error handling.
*  UML char picture.

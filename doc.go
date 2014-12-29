//cpic is a tool to draw data structure in ASCI character picture.
//
//Example:
//
//tree is parsed with top to left,and indent denotes depth
//
//	tree:
//	    -> black
//		    ->red
//              ->red
//	                ->red
//	                ->red
//	                    ->red
//	                    ->red
//              ->red
//	                ->red
//	                    ->red
//              ->red
//	                ->red
//
//	to
//
//	TREE
//	black
//	|          \   \
//	red         red red
//	|           |   |
//	red         red red
//	|  \        |
//	red red  	red
//	|  \
//	red red
package cpic

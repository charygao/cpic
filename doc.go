//Cpic is a tool to draw data structure in ASCII character picture.
//
//Tree:
//
//  Tree is conveted from top to left,and indent denotes depth.
//  Every tree needs to begin with `tree:` in order to direct tree parsing and
//	one indent next line with a `->` followe by a identifier (begin with '_' or letters)
//	to denote a root node.
//
//Graph:
//
// Graph is vertecies followed by a '->' and more than one edge '[int] id',first is weight,and it's optional,default is nil,
// id is an identifier(begin with '_' or letters').
// Currently connect vertex itself is not allowed.
//
package cpic

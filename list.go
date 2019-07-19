// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

import "strconv"

type ListType int

type List struct {
	*BaseNode

	Bullet bool
	Start  int
	Delim  string
	Tight  bool

	Marker            string
	startIndentSpaces int
	indentSpaces      int
}

func newList(marker string, bullet bool, start int, delim string, startIndentSpaces, indentSpaces int, t *Tree) (ret Node) {
	baseNode := &BaseNode{typ: NodeList}
	ret = &List{
		baseNode,
		bullet,
		start,
		delim,
		false,
		marker,
		startIndentSpaces,
		indentSpaces,
	}
	t.context.CurNode = ret

	return
}

func (t *Tree) parseListMarker(line items) (remains items, marker, delim string, bullet bool, start, startIndentSpaces, w, n int) {
	spaces, tabs, tokens, firstNonWhitespace := t.nonWhitespace(line)
	var markers items
	markers = append(markers, firstNonWhitespace)
	line = line[len(tokens):]
	bullet = true
	start = 1
	if firstNonWhitespace.isNumInt() {
		bullet = false
		start, _ = strconv.Atoi(firstNonWhitespace.val)
		markers = append(markers, line[0])
		line = line[1:]
	}
	switch markers[len(markers)-1].typ {
	case itemAsterisk:
		delim = " "
	case itemHyphen:
		delim = " "
	case itemPlus:
		delim = " "
	case itemCloseParen:
		delim = ")"
	case itemDot:
		delim = "."
	}
	startIndentSpaces = spaces + tabs*4
	marker = markers.rawText()
	spaces, tabs, _, firstNonWhitespace = t.nonWhitespace(line)
	w = len(marker)
	n = spaces + tabs*4
	if 4 < n {
		n = 1
	} else if 1 > n {
		n = 1
	}
	if line[0].isTab() {
		line = t.indentOffset(line, 2)
	} else {
		line = line[1:]
	}

	remains = line

	return
}

func (t *Tree) parseList(line items) (ret Node) {
	indentSpaces := t.context.IndentSpaces

	line, marker, delim, bullet, start, startIndentSpaces, w, n := t.parseListMarker(line)
	ret = newList(marker, bullet, start, delim, startIndentSpaces, startIndentSpaces+w+n, t)
	t.context.IndentSpaces += startIndentSpaces + w + n

	if line.isBlankLine() {
		t.context.IndentSpaces = startIndentSpaces + w + 1
	}

	var node Node
	var blankLineIndices []int
	i := 0

	node = t.parseListItem(line)
	if nil == node {
		return
	}
	ret.AppendChild(ret, node)
	t.context.IndentSpaces = indentSpaces

	blankLines := t.skipBlankLines()
	if 0 < len(blankLines) {
		blankLineIndices = append(blankLineIndices, i)
	}

	//line = t.nextLine()
	//if line.isEOF() {
	//	break
	//}
	//
	//if t.isThematicBreak(line) {
	//	t.backupLine(line)
	//	break
	//}
	//
	//if t.blockquoteMarkerCount(line) < t.context.BlockquoteLevel {
	//	t.backupLine(line)
	//	break
	//}
	//
	//start++
	//
	//if isList, _ := t.isList(line); !isList {
	//	t.backupLine(line)
	//	node.(*ListItem).Tight = endWithBlankLine
	//	break
	//}
	//
	//nextLine, nextMarker, nextDelim, startIndentSpaces, indentSpaces := t.parseListItemMarker(line, ret)
	//if bullet {
	//	if marker != nextMarker {
	//		t.backupLine(line)
	//		break
	//	}
	//} else {
	//	if delim != nextDelim || strconv.Itoa(start) != nextMarker[:1] {
	//		t.backupLine(line)
	//		break
	//	}
	//}
	//
	//if nextLine.isBlankLine() && t.context.IndentSpaces > line.spaceCountLeft() {
	//	t.backupLine(line)
	//	break
	//}
	//
	//if startIndentSpaces < t.context.IndentSpaces && t.context.IndentSpaces >= indentSpaces {
	//	t.backupLine(line)
	//	break
	//}
	//
	//if 3 < startIndentSpaces {
	//	t.backupLine(line)
	//	break
	//}

	t.context.IndentSpaces = indentSpaces
	line = t.indentOffset(line, t.context.IndentSpaces)

	//tight := true
	//for child := ret.FirstChild(); nil != child; child = child.Next() {
	//	if !child.(*ListItem).Tight {
	//		tight = false
	//		break
	//	}
	//}
	//ret.(*List).Tight = tight
	//for child := ret.FirstChild(); nil != child; child = child.Next() {
	//	child.(*ListItem).Tight = tight
	//}

	return
}

func (t *Tree) isList(line items) (isList bool, marker string) {
	if 2 > len(line) { // at least marker and newline
		return
	}

	_, line = line.trimLeft()
	if 1 > len(line) {
		return
	}

	firstNonWhitespace := line[0]

	if itemAsterisk == firstNonWhitespace.typ {
		isList = line[1].isWhitespace()
		marker = "*"
		return
	} else if itemHyphen == firstNonWhitespace.typ {
		isList = line[1].isWhitespace()
		marker = "-"
		return
	} else if itemPlus == firstNonWhitespace.typ {
		isList = line[1].isWhitespace()
		marker = "+"
		return
	} else if firstNonWhitespace.isNumInt() && 9 >= len(firstNonWhitespace.val) {
		isList = line[2].isWhitespace()
		if itemDot == line[1].typ {
			marker = firstNonWhitespace.val + "."
		} else if itemCloseParen == line[1].typ {
			marker = firstNonWhitespace.val + ")"
		}
		return
	}

	return
}

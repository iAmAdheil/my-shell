package main

const count = 256

type Trienode struct {
	children [count]*Trienode
	terminal bool
}

func InitTrie(words []string) *Trienode {
	root := &Trienode{
		children: [count]*Trienode{},
		terminal: false,
	}

	for _, word := range words {
		root.Insert(word)
	}

	return root
}

func backtrack_helper(suggs *[]string, tmp *Trienode, sug string) {
	if tmp.terminal {
		*suggs = append(*suggs, sug)
	}

	for i := 0; i < len(tmp.children); i++ {
		char := string([]byte{byte(i)})
		if tmp.children[i] != nil {
			sug += char
			backtrack_helper(suggs, tmp.children[i], sug)
			sug = sug[:len(sug)-1]
		}
	}
}

func (root *Trienode) Complete(word string) []string {
	suggs := []string{}
	tmp := root
	sugg := ""
	for i := 0; i < len(word); i++ {
		idx := int(word[i])
		if tmp.children[idx] == nil {
			return suggs
		}
		tmp = tmp.children[idx]
		sugg += string(word[i])
	}

	if tmp.terminal {
		suggs = append(suggs, sugg)
		return suggs
	}

	backtrack_helper(&suggs, tmp, sugg)

	return suggs
}

func (root *Trienode) Insert(word string) {
	tmp := root
	for i := 0; i < len(word); i++ {
		idx := int(word[i])
		if tmp.children[idx] == nil {
			nNode := &Trienode{
				children: [count]*Trienode{},
				terminal: false,
			}
			tmp.children[idx] = nNode
		}
		tmp = tmp.children[idx]
	}

	tmp.terminal = true
}

func GetComPrefix(suggs []string) string {
	root := &Trienode{
		children: [count]*Trienode{},
		terminal: false,
	}
	for _, sug := range suggs {
		root.Insert(sug)
	}

	tmp := root
	pref := []rune{}
	for {
		c := 0
		li := 0
		for i := 0; i < count; i++ {
			if tmp.children[i] != nil {
				c++
				li = i
			}
		}
		if c == 0 || c > 1 || tmp.terminal {
			break
		}
		pref = append(pref, rune(li))
		tmp = tmp.children[li]
	}

	return string(pref)
}

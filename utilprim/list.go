package utilprim

import "iter"

// header da lista
type List[T any] struct {
	head *node[T]
	len int64	
}


// node da lista
type node[T any] struct {
	data T
	r *node[T]
	l *node[T]
}

// NewList constroi uma nova lista
func NewList[T any] () *List[T] {
	return &List[T]{}
}

// InsertEnd insere um elemento ao final da lista
func (list *List[T]) InsertEnd(v T) {

	defer func(){
		list.len++
	}()

	no := &node[T]{data: v}

	// lista vazia
	if list.head == nil {
		no.l = no
		no.r = no
		list.head = no
		return
	}

	// lista com pelo menos 1 nó
	tail := list.head.l // último nó
	no.l = tail
	no.r = list.head
	tail.r = no
	list.head.l = no
}

// IterPrim cria uma sequencia que interage com range
func (list *List[T]) IterPrim() iter.Seq[T] {
	return func(yield func(T) bool) {
		if list.head == nil {
			return
		}
		
		cur := list.head

		// verifica se a lista so tem 1 elemento
		if cur.r == nil {
			yield(cur.data)
			return
		}

		// cria o loop para yield
		for {
			if !yield(cur.data) {
				return
			}

			cur = cur.r

			if cur == list.head {
				return
			}
		}
	}
}


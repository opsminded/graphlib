package graphlib

type unhealthyIterator struct {
	vertices []*vertex
	index    int
}

// HasNext retorna true se ainda houver vértices não saudáveis a serem iterados.
func (it *unhealthyIterator) HasNext() bool {
	return it.index < len(it.vertices)
}

// Next retorna o próximo vértice não saudável.
func (it *unhealthyIterator) Next() *vertex {
	if !it.HasNext() {
		return nil
	}
	v := it.vertices[it.index]
	it.index++
	return v
}

// Iterate aplica uma função a todos os vértices não saudáveis.
func (it *unhealthyIterator) Iterate(fn func(v *vertex) error) error {
	for it.HasNext() {
		if err := fn(it.Next()); err != nil {
			return err
		}
	}
	return nil
}

// Reset reinicia o índice do iterador.
func (it *unhealthyIterator) Reset() {
	it.index = 0
}

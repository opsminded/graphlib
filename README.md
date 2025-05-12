# 📦 graphlib
graphlib é uma biblioteca Go para representar e manipular grafos direcionados acíclicos (DAGs) de forma simples, segura e concorrente.
Ela permite criar vértices e arestas, navegar pelas relações de dependência, explorar caminhos e extrair subgrafos. A biblioteca é dividida em dois níveis:

- core.go: lógica interna (não exportada)
- graph.go: API pública, segura para concorrência

## ✨ Features
- Criação de vértices e arestas com controle de ciclo
- Suporte a dependentes e dependências
- Caminhos entre vértices (Path)
- Subgrafos com vizinhos, dependentes, dependências
- Thread-safe via sync.RWMutex

## 📦 Instalação
```bash
go get github.com/opsminded/graphlib@latest
```

🧑‍💻 Exemplo de uso
```go
package main

import (
	"fmt"
	"github.com/opsminded/graphlib"
)

func main() {
	g := graphlib.NewGraph()

	g.NewVertex("A")
	g.NewVertex("B")
	g.NewVertex("C")
	g.NewVertex("D")

	g.NewEdge("A->B", "A", "B")
	g.NewEdge("B->C", "B", "C")
	g.NewEdge("A->D", "A", "D")

	// Neighbors
	sub := g.Neighbors("A")
	fmt.Println("Vizinhos de A:", sub.Vertices)

	// Caminho de A para C
	path := g.Path("A", "C")
	fmt.Println("Caminho A -> C:", path.Vertices)
}
```

## 🚧 Restrições
- Arestas bidirecionais entre dois nós não são permitidas
- Ciclos não são permitidos (DAG)
- O nome dos vértices (Label) deve ser único

## ✅ Testes
- A biblioteca inclui testes unitários para:
- Neighbors
- Dependents/Dependencies (diretos e recursivos)
- Path (com múltiplos caminhos)
- Subgrafos

Você pode executar os testes clonando o repositório:
```bash
git clone https://github.com/opsminded/graphlib.git
go test ./...
```

## 📄 Licença
[Licença MIT](https://github.com/opsminded/graphlib/blob/main/LICENSE)
```
MIT License

Copyright (c) 2025 opsminded

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
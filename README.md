# graphlib

Biblioteca Go para modelagem de grafos direcionados e acíclicos com controle de classes e validações.

## Conceitos
- **Vertex (vértice)**: Representa um nó no grafo, contendo um identificador, classe, label e status de saúde.
- **Edge (aresta)**: Representa uma dependência entre dois vértices.
- **Classe**: Uma categorização textual para vértices e arestas (ex: DATABASE, FIREWALL).
- **Label**: Um nome único e descritivo do recurso (ex: DB1, FW_MAIN).

```go
g := graphlib.NewGraph()

s1, _ := g.NewVertex("SERVER01", "SERVER")
s2, _ := g.NewVertex("SERVER02", "SERVER")

e1, _ := g.NewEdge("CONEXAO", "SERVER_CONN", s1, s2)

```

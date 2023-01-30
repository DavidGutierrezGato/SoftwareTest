package main

import (
	"testing"
)

func BenchmarkCargar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cargar("D:/pruebaSoftware/back/prueba.txt")
	}

}

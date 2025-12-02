package main

import (
	"fmt"
	"github.com/AndreSS-ntp/univ_algs/lab5/internal/app/balanced"
	"github.com/AndreSS-ntp/univ_algs/lab5/internal/app/optimal"
	"github.com/AndreSS-ntp/univ_algs/lab5/internal/domain"
	"github.com/AndreSS-ntp/univ_algs/lab5/internal/pkg"
	"github.com/AndreSS-ntp/univ_algs/lab5/internal/repository"
)

func main() {
	keys, p, q := repository.GetLabData()

	fmt.Println("Ключи:", keys)
	fmt.Println("p[i]:", p[1:]) // p[0] не используем
	fmt.Println("q[i]:", q)

	// 1. Полностью сбалансированное дерево
	balancedRoot := balanced.BuildBalanced(keys)
	fmt.Println("\nПолностью сбалансированное дерево:")
	domain.PrintTree(balancedRoot, 0)

	balancedCost := pkg.ComputeCost(balancedRoot, keys, p, q)
	fmt.Printf("Цена поиска по сбалансированному дереву: %d\n", balancedCost)

	// 2. Оптимальное дерево поиска
	optimalRoot, optimalDpCost := optimal.BuildOptimal(keys, p, q)
	fmt.Println("\nОптимальное дерево поиска:")
	domain.PrintTree(optimalRoot, 0)

	optimalCost := pkg.ComputeCost(optimalRoot, keys, p, q)
	fmt.Printf("Цена поиска по оптимальному дереву (через обход): %d\n", optimalCost)
	fmt.Printf("Цена поиска по оптимальному дереву (из DP C(1,N)): %d\n", optimalDpCost)
}

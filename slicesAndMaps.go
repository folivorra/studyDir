package dirForStudy

import (
	"errors"
	"sort"
)

/*
На вход приходят бесконечное число слайсов интов разной длины. На выходе вернуть один слайс из отсортированных значений.
*/

func sortSlices(arrays ...[]int) []int {
	res := make([]int, 0)
	for _, array := range arrays {
		res = append(res, array...)
	}
	sort.Ints(res)
	return res
}

/*
Заполнение мапы из нескольких неотсортированных слайсов только теми ключами, которых нет в мапе.
Затем перезаписать в новый слайс (и вернуть его) отсортированные значения из мапы.
*/

func sortFromMapToSlice(arrays ...[]int) []int {
	tempArr := make([]int, 0)
	res := make([]int, 0)
	m := make(map[int]struct{})
	for _, array := range arrays {
		tempArr = append(tempArr, array...)
	}
	for _, element := range tempArr {
		if _, exists := m[element]; !exists {
			m[element] = struct{}{}
		}
	}
	for key := range m {
		res = append(res, key)
	}
	sort.Ints(res)
	return res
}

/*
Взять 2 "первых" ключа из мапы и добавить новый составленный из них в
мапу, если такого ключа ещё нет в мапе. После чего вернуть мапу из функции.
*/

func addSumKey(m map[string]struct{}) map[string]struct{} {
	slice := make([]string, 0)

	for k := range m {
		slice = append(slice, k)
	}
	sort.Strings(slice)

	if len(slice) < 2 {
		return m
	}

	newKey := slice[1] + slice[0] // если наоборот то кейс при котором их комбинация уже существует - невозможен
	if _, exists := m[newKey]; !exists {
		m[newKey] = struct{}{}
	}

	return m
}

/*
Найти разницу двух слайсов и записать её в третий и вернуть его.
*/

func subtractSlices(arr1, arr2 []int) []int {
	m2 := make(map[int]struct{}, len(arr2))
	res := make([]int, 0)

	for _, v := range arr2 {
		m2[v] = struct{}{}
	}

	for _, v := range arr1 {
		if _, exists := m2[v]; !exists {
			res = append(res, v)
		}
	}

	return res
}

/*
Найти пересечение в двух слайсах и заполнить её в третий и вернуть его.
*/

func intersectionSlices(arr1, arr2 []int) []int {
	m2 := make(map[int]struct{}, len(arr2))
	res := make([]int, 0)

	for _, v := range arr2 {
		m2[v] = struct{}{}
	}

	for _, v := range arr1 {
		if _, exists := m2[v]; exists {
			res = append(res, v)
		}
	}

	return res
}

/*
Создание "зеркальной" мапы
*/

func mirrorMap(m map[string]int) map[int][]string {
	res := make(map[int][]string, len(m))

	for k, v := range m {
		res[v] = append(res[v], k)
	}

	return res
}

/*
Поиск максимума и минимума
*/

func minAndMax(arr []float64) (min, max float64) {
	min = arr[0]
	max = arr[0]

	for _, v := range arr {
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}

	return
}

/*
Фильтрация слайса
*/

func filterSlice(arr []int, filter func(int) bool) []int {
	res := make([]int, 0)

	for i := 0; i < len(arr); i++ {
		if filter(arr[i]) {
			res = append(res, arr[i])
		}
	}

	return res
}

/*
Разделение слайса (split)
*/

func splitMaps(arr []int, num int) [][]int {
	res := make([][]int, 0)
	iter := 0

	if len(arr) == 0 || num < 1 {
		return res
	}

	for i := 0; iter+num < len(arr); i++ {
		res = append(res, arr[iter:iter+num])
		iter += num
	}
	res = append(res, arr[iter:])

	return res
}

/*
Объединение двух мап
*/

func combineMaps(m1, m2 map[string]int) map[string]int {
	res := make(map[string]int)

	for k, v := range m1 {
		res[k] = v
	}

	for k, v := range m2 {
		if _, exists := res[k]; !exists {
			res[k] = v
		} else {
			if res[k] < v {
				res[k] = v
			}
		}
	}

	return res
}

/*
Группировка элементов по признаку
*/

type Item struct {
	Category string
	Value    int
}

func groupByStruct(s []Item) map[string][]int {
	res := make(map[string][]int)

	for _, item := range s {
		res[item.Category] = append(res[item.Category], item.Value)
	}

	return res
}

/*
Удаление дубликатов из слайса
*/

func deleteDuplicates(arr []string) []string {
	res := make([]string, 0)
	m := make(map[string]struct{})

	for _, v := range arr {
		if _, exists := m[v]; !exists {
			m[v] = struct{}{}
			res = append(res, v)
		}
	}

	return res
}

/*
На входе слайс с интами, второй слайс с теми интами которые нужно найти.
Использовать алгоритм бинарного посика и сформировать новый слайс с найденными значениями и вернуть его.
*/

func findBinary(input []int, targets []int) []int {
	res := make([]int, 0)

	if len(targets) == 0 || targets == nil {
		return res
	}

	for _, target := range targets {
		if _, err := binarySearch(input, target); err == nil {
			res = append(res, target)
		}
	}

	return res
}

func binarySearch(arr []int, target int) (int, error) {
	mid := len(arr) / 2
	if arr[mid] == target {
		return mid, nil
	} else if arr[mid] > target && len(arr) > 1 {
		return binarySearch(arr[:mid], target)
	} else if arr[mid] < target && len(arr) > 1 {
		return binarySearch(arr[mid+1:], target)
	} else {
		return -1, errors.New("not found")
	}
}

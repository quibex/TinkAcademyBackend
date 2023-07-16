package main

import . "fmt"

func main(){
	var n int
	Scan(&n)

	que := make([]int, 0)

	com := 0

	for i := 0; i < n; i++ {
		Scanf("%d",&com)

		if com == 1 {
			var val int
			Scan(&val)
			que = append(que, val);
		}
		if com == 2 {
			que2 := make([]int, 0, len(que)*2)
			for _, v := range que {
				que2 = append(que2, v)
				que2 = append(que2, v)
			}
			que = que2
		}
		if com == 3 {
			val := que[0]
			Println(val)
			que = que[1:] //Память будет использоваться не эффективно, т.к. начальные элемента будут занимать место, но не будут доступны для использования. Лучше сделать колцевой буфер.
		}
	}
}
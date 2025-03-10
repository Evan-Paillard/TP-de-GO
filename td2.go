package main
    import (
        "fmt"
        "math/rand"
        "sync"
    )


//Question 1:
func message(c chan string, mess string){

    c <- mess

}


//Question 2:
func inc(c chan int, i int){

    nb := i+1
    c <- nb

}

//Question 3:

type resultat struct{
    Id int
    Carre int
}


func maitre(c chan int, i int){

    c <- i

    
}

func ouvriere(id int, chin chan int, chout chan resultat){

    nombre := <- chin
    fmt.Println(nombre)

    carre := nombre*nombre
    chout <- resultat{Id: id+1, Carre: carre}



}

//Question 4:
type resu struct{
    Id int
    Mess string
}


func messagefunc(c chan resu, mess string, id int, listeRoutine *[]resu, wg *sync.WaitGroup){

    defer wg.Done()
    nb := resu{Id: id+1, Mess: mess}
    *listeRoutine = append(*listeRoutine, nb)
    c <- nb
    

}






func main(){


    //Question 1:
    canal := make(chan string)
    go message(canal, "Bonjour, Goroutine !")
    fmt.Println(<- canal)

    //Question 2:
    canal1 := make(chan int)
    canal2 := make(chan int)
    canal3 := make(chan int)
    canal4 := make(chan int)
    canal5 := make(chan int)
    canal6 := make(chan int)
    canal7 := make(chan int)
    canal8 := make(chan int)
    canal9 := make(chan int)
    canal10 := make(chan int)
    go inc(canal1, 0)
    go inc(canal2, <- canal1)
    go inc(canal3, <- canal2)
    go inc(canal4, <- canal3)
    go inc(canal5, <- canal4)
    go inc(canal6, <- canal5)
    go inc(canal7, <- canal6)
    go inc(canal8, <- canal7)
    go inc(canal9, <- canal8)
    go inc(canal10, <- canal9)
    fmt.Println(<- canal10)

    //Question 3:
    maitr := make(chan int)
    chin := make(chan int)
    chout := make(chan resultat)
    go maitre(maitr, rand.Intn(100))
    num := <- maitr

    for i := 0; i < 2; i++{
        go ouvriere(i, chin, chout)
    }


    for i:= 0; i < 2; i++{
        chin <- num

    }

    close(chin)


    for i:= 0; i < 2; i++{
        
        res := <- chout
        fmt.Println("ouvriere", res.Id,res.Carre)


    }
    
    close(chout)


    //Question 4:

    channel1 := make(chan resu)
    var listeRoutine []resu 

    message := "Bonjour"

    var wg sync.WaitGroup

    for i:= 0; i < 3; i++{
        wg.Add(1)
        go messagefunc(channel1, message, i, &listeRoutine, &wg)

    }


    wg.Wait() 



    fmt.Println(listeRoutine)
    // for i:= 0; i < len(listeRoutine)







    
    



}
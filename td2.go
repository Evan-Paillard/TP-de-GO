package main
    import (
        "fmt"
        "math/rand"
        "sync"
        "time"
    )


//Question 2:
func sendmessage(c chan string, mess string){

    c <- mess

}


//Question 3:
func inc(after chan int, before chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    nb := <-before
    after <- nb + 1
}

//Question 4:

type resultat struct{
    Id int
    Carre int
}



func maitre(randomValue int, ouvrieres []chan int, result chan int, wg *sync.WaitGroup) {

    defer wg.Done()


	for i, ch := range ouvrieres {
		fmt.Printf("Le maître envoie %d à l'ouvrière %d\n", randomValue, i+1)
		ch <- randomValue
	}


	for i := 0; i < len(ouvrieres); i++ {
		res := <-result
		fmt.Printf("Le maître a reçu : %d\n", res)
	}
}

func ouvrieresfunc(chin chan int, chout chan int,id int, wg *sync.WaitGroup){

    defer wg.Done()

    nombre := <- chin

    fmt.Printf("L'ouvrière %d à reçu %d de la part du maitre\n",id,nombre)

    carre := nombre*nombre

    chout <- carre



}

//Question 5:
type Message struct {
	Texte string
	Tour  int
	Index int 
}
func anneau(id, total int, ch chan Message) {
	for {
		msg, ok := <-ch
		if !ok {
			return
		}

		if msg.Index == id {
			fmt.Printf("Goroutine %d a reçu : %s (Tour %d)\n", id, msg.Texte, msg.Tour)

			if msg.Tour == 0 {
				close(ch)
				return
			}

	
			msg.Index = (id + 1) % total
			if msg.Index == 0 {
				msg.Tour-- 
			}
			ch <- msg
		} else {
			ch <- msg
		}
	}
}


func messagePart2(chout chan string, chin chan string, id int, wg *sync.WaitGroup, k int) {
	defer wg.Done()

	transmissions := 0

	for {
		message, ok := <-chin
		if !ok {
			return
		}

		fmt.Printf("La goroutine %d a reçu le message : %s\n", id, message)

		transmissions++

		if transmissions >= k {
			fmt.Printf("La goroutine %d a terminé ses %d transmissions et arrête.\n", id, k)
			return 
		}

		chout <- message
	}
}



//Question 6

func genererTableauAleatoire(taille, max int) []int {
	rand.Seed(time.Now().UnixNano())
	tableau := make([]int, taille)
	for i := range tableau {
		tableau[i] = rand.Intn(max)
	}
	return tableau
}

func mergeSortParallele(tableau []int, N int) []int {
	if len(tableau) <= 1 {
		return tableau
	}


    taillePartie := len(tableau) / N
	if taillePartie < 1 {
		taillePartie = 1
		N = len(tableau)
	}


    resultats := make(chan []int, N)


    var wg sync.WaitGroup
	wg.Add(N)


    for i := 0; i < N; i++ {
		debut := i * taillePartie
		fin := (i + 1) * taillePartie
		if i == N-1 {
			fin = len(tableau) 
		}

		partie := tableau[debut:fin]

		go func(partie []int) {
			defer wg.Done() 
			resultats <- triSimple(partie)
		}(partie)
	}

	go func() {
		wg.Wait()
		close(resultats)
	}()

	var tabsTriés [][]int
	for tab := range resultats {
		tabsTriés = append(tabsTriés, tab)
	}

	for len(tabsTriés) > 1 {
		tabsTriés[0] = fusionner(tabsTriés[0], tabsTriés[1])
		tabsTriés = append(tabsTriés[:1], tabsTriés[2:]...)
	}

	if len(tabsTriés) == 0 {
		return []int{}
	}

	return tabsTriés[0]
}

func triSimple(tableau []int) []int {
	if len(tableau) <= 1 {
		return tableau
	}

	milieu := len(tableau) / 2
	gauche := triSimple(tableau[:milieu])
	droite := triSimple(tableau[milieu:])

	return fusionner(gauche, droite)
}

func fusionner(gauche, droite []int) []int {
	resultat := make([]int, 0, len(gauche)+len(droite))
	i, j := 0, 0

	for i < len(gauche) && j < len(droite) {
		if gauche[i] <= droite[j] {
			resultat = append(resultat, gauche[i])
			i++
		} else {
			resultat = append(resultat, droite[j])
			j++
		}
	}

	resultat = append(resultat, gauche[i:]...)
	resultat = append(resultat, droite[j:]...)

	return resultat
}




func main() {
    // Question 2:
    canal := make(chan string)
    go sendmessage(canal, "Bonjour, Goroutine !")
    fmt.Println("Le message retourné à la fin de la go routine est : ",<-canal)

}

func main(){

        // Question 3:
        var wg sync.WaitGroup
        n := 10
    
        wg.Add(n)
        channels := make([]chan int, n)
    
        for i := 0; i < n; i++ {
            channels[i] = make(chan int,1)
        }
    
        go func() {
            defer wg.Done()
            channels[0] <- 0
        }()
    
        for i := 1; i < n; i++ {
            go inc(channels[i], channels[i-1], &wg)
        }    
    
        fmt.Println("La valeur à la dernière go routines est : ",<-channels[len(channels)-1])
        fmt.Println("\n")
        wg.Wait()

}




func main(){
    var wg sync.WaitGroup

    //Question 4:
    m := 5
    wg.Add(m)

    result := make(chan int)
    randomValue := rand.Intn(100)

    ouvrieres := make([]chan int, m)

    for i:=0; i<m; i++{

        ouvrieres[i] = make(chan int,1)

        go ouvrieresfunc(ouvrieres[i], result,i+1, &wg)

    }

    go maitre(randomValue, ouvrieres, result, &wg)


    wg.Wait()
}



func main(){

    //question 5.1
	P := 5    
	K := 3    
	ch := make(chan Message)

	for i := 0; i < P; i++ {
		go anneau(i, P, ch)
	}

	ch <- Message{
		Texte: "hello",
		Tour:  K,
		Index: 0,
	}

	fmt.Println("Message a fait", K, "tours. Programme terminé.")

    // Question 5.2:
    var wg sync.WaitGroup
	p2 := 5  
	k := 3  

	message := "Bonjour"

	messageRoutines := make([]chan string, p2)
	for i := 0; i < p2; i++ {
		messageRoutines[i] = make(chan string)
	}


	for i := 0; i < p2; i++ {
		next := (i + 1) % p2
		wg.Add(1)
		go messagePart2(messageRoutines[next], messageRoutines[i], i, &wg, k)
	}

	fmt.Println("Début de la transmission du message")

	messageRoutines[0] <- message

	wg.Wait()

	fmt.Println("Fin de la transmission")
}


func main(){

    tableau := genererTableauAleatoire(20, 100)
	fmt.Println("Tableau non trié:", tableau)

	N := 4

	resultat := mergeSortParallele(tableau, N)

	fmt.Println("Tableau trié:", resultat)



}



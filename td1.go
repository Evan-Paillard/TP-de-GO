// package main
// import (
//     "fmt"
//     "math/rand"
//     "math"
// )

// //Exercice 1: 

// //Question 1:
// func estBissextile(date int) bool {

//     if date%4==0{

//         if date%100 == 0 {
            
//             if date%400 ==0{
//                 return true
//             } else {
//                 return false
//             }
//         } else {
//             return true 
//         }
    
//     }
//     return false
// }


// //Question 2:

// func estPremier(entier int) bool {

//     if entier==1 {

//         return false
//     }


//     for i := 2; i<entier; i++ {
//         if entier%i==0 {
//             return false
//         }
//     }

//     return true

// }

// //Question 3:

// func premiersNombresPremiers(entier int) []int{

//     var tab []int

//     for i:=1; i<=entier; i++ {

//         if estPremier(i)==true {
//             tab = append(tab,i)
//         }
//     }

//     return tab

// }

// //Question 4:
// func genererTableauAleatoire(entier int) []int{

//     var tab []int

//     for i:=0; i<entier; i++ {

//         tab = append(tab,rand.Intn(100))

//     }

//     return tab

// }

// //Question 5:
// func triBulles(T []int) []int{

//     var tab = T
//     var annexe int

//     for i:=len(tab)-1; i >= 0; i--{
//         for j:=0; j<i ; j++ {
//             if tab[j+1] < tab[j] {
//                 annexe=tab[j]
//                 tab[j] = tab[j+1]
//                 tab[j+1] = annexe
//             }
//         }
//     }
//     return tab
// }


// //Question 6:

// func triSelection(T []int) []int{

//     var tab = T
//     var annexe int

//     for i:=0; i<len(tab)-1; i++{
//         var min int = i
//         for j:=i+1; j<len(tab); j++{
//             if tab[j] < tab[min] {
//                 min = j
//             }
//         }
//         if min !=i{
//             annexe= tab[i]
//             tab[i] = tab[min]
//             tab[min] = annexe
//         }
//     }
//     return tab
// }

// //Question 7:
// func rechercheDichotomique(T []int, x int) bool{

//     if len(T)==0{
//         return false
//     }
//     var m int = len(T)/2
//     if T[m] == x {
//         return true
//     } else if T[m]>x {
//         return rechercheDichotomique(T[:m], x) 
//     } else {
//         return rechercheDichotomique(T[m+1:], x)
//     }
// }



// func organiserParTaille(tab []string) map[int][]string{

//     dico := make(map[int][]string)

//     for _, valeur := range tab {

//         longueur := len(valeur)
//         if _, ok := dico[longueur]; ok {
//             dico[longueur] = append(dico[longueur],valeur)
//         } else {
//             dico[longueur] = []string{valeur}
//         }
//     }

//     return dico

// }


// //Exercice 2:
// //Question 1:
// func initGrille(n int,m int) [][]int {

//     var tab [][]int

//     for i := 0; i<n; i++{
//         var row []int
//         for j := 0; j<m; j++{
//             row = append(row,rand.Intn(2))
//         }
//         tab = append(tab, row)
//     }


//     return tab

// }

// //Question 2:
// func compterVoisin(grille [][]int, i int, j int) int {
//     var nb int = 0
//     var longueur int = len(grille[i]) - 1
//     var hauteur int = len(grille) - 1

//     // Vérifier le voisin du haut
//     if i > 0 && grille[i-1][j] == 1 {
//         nb=nb+1
//     }

//     // Vérifier le voisin du bas
//     if i < hauteur && grille[i+1][j] == 1 {
//         nb=nb+1
//     }

//     // Vérifier le voisin de gauche
//     if j > 0 && grille[i][j-1] == 1 {
//         nb=nb+1
//     }

//     // Vérifier le voisin de droite
//     if j < longueur && grille[i][j+1] == 1 {
//         nb=nb+1
//     }

//     // Vérifier le voisin en haut à gauche (diagonale)
//     if i > 0 && j > 0 && grille[i-1][j-1] == 1 {
//         nb=nb+1
//     }

//     // Vérifier le voisin en haut à droite (diagonale)
//     if i > 0 && j < longueur && grille[i-1][j+1] == 1 {
//         nb=nb+1
//     }

//     // Vérifier le voisin en bas à gauche (diagonale)
//     if i < hauteur && j > 0 && grille[i+1][j-1] == 1 {
//         nb=nb+1
//     }

//     // Vérifier le voisin en bas à droite (diagonale)
//     if i < hauteur && j < longueur && grille[i+1][j+1] == 1 {
//         nb=nb+1
//     }

//     return nb
// }


// //Question 3:
// // func update(jeu_vie [][]int) [][]int{

// //     var grille [][]int = jeu_vie

// //     for i:=0; i<len(grille);i++{
// //         for j:=0; j<len(grille[i]); j++{
// //             if grille[i][j]==1{
// //                 if compterVoisin(grille,i,j)==2 || compterVoisin(grille,i,j)==3{
// //                     grille[i][j]=1
// //                 }else{
// //                     grille[i][j]=0
// //                 }
// //             }else{
// //                 if compterVoisin(grille,i,j)==3{
// //                     grille[i][j]=1
// //                 }
// //             }
// //         }
// //     }

// //     return grille

// // }

// func update(jeu_vie [][]int) [][]int {
//     // Créer une nouvelle grille pour stocker les valeurs mises à jour
//     nouvelleGrille := make([][]int, len(jeu_vie))
//     for i := range jeu_vie {
//         nouvelleGrille[i] = make([]int, len(jeu_vie[i]))
//     }

//     for i := 0; i < len(jeu_vie); i++ {
//         for j := 0; j < len(jeu_vie[i]); j++ {
//             voisins := compterVoisin(jeu_vie, i, j)
//             if jeu_vie[i][j] == 1 {
//                 // Une cellule vivante reste vivante si elle a 2 ou 3 voisins
//                 if voisins == 2 || voisins == 3 {
//                     nouvelleGrille[i][j] = 1
//                 } else {
//                     nouvelleGrille[i][j] = 0
//                 }
//             } else {
//                 // Une cellule morte devient vivante si elle a exactement 3 voisins
//                 if voisins == 3 {
//                     nouvelleGrille[i][j] = 1
//                 } else {
//                     nouvelleGrille[i][j] = 0
//                 }
//             }
//         }
//     }

//     return nouvelleGrille
// }


// //Question 4:
// func afficherGrille(grille [][]int){


// }





// func main() {


//     //Fonction exercice 1:
//     fmt.Println(estBissextile(1900))
//     fmt.Println(estPremier(1))
//     fmt.Println(estPremier(10))
//     fmt.Println(premiersNombresPremiers(11))
//     fmt.Println(genererTableauAleatoire(20))
//     fmt.Println(triBulles([]int{1, 3, 2, 0}))
//     fmt.Println(triSelection([]int{1, 3, 2, 0}))
//     fmt.Println(rechercheDichotomique([]int{0, 1, 2, 3},0))
//     fmt.Println(organiserParTaille([]string{"abc", "a", "b", "bc"}))

//     //Fonction exercice 2:
//     var d [][]int = initGrille(6,6)
//     fmt.Println(d)
//     fmt.Println(compterVoisin(d, 3,3))
//     fmt.Println(update(d))


// }




// //Exercice 3

// //Question 1:





// type vect2i struct {

//     x int
//     y int

// }

// func (v vect2i) init(a int, b int){

//     v.x = a
//     v.y = b

// }

// func (v1 vect2i) addition(v2 vect2i) vect2i{

//     var v3 vect2i

//     v3.x = v1.x + v2.x
//     v3.y = v1.y + v2.y

//     return v3
    
// }

// func (v1 vect2i) soustraction(v2 vect2i) vect2i{

//     var v3 vect2i

//     v3.x = v2.x - v1.x
//     v3.y = v2.y - v1.y

//     return v3


// }

// func (v1 vect2i) multiplication(v2 vect2i) vect2i{

//     var v3 vect2i

//     v3.x = v2.x * v1.x
//     v3.y = v2.y * v1.y

//     return v3

// }


// func (v vect2i) norme() float64{
//     return math.Sqrt(float64(v.x*v.x + v.y*v.y))
// }

// func (v vect2i) normalization() (vect2i, error) {

//     var norme float64 = v.norme()

//     if norme == 0{

//         return vect2i{0,0}, fmt.Errorf("impossible de normaliser un vecteur nul")

//     }

//     return vect2i{int(float64(v.x) / norme), int(float64(v.y) / norme)}, nil

// }


// func (v1 vect2i) scalaire(v2 vect2i) int {
//     return v1.x*v2.x + v1.y*v2.y
// }

// func (v1 vect2i) vectoriel(v2 vect2i) int {

//     return v1.x*v2.y - v1.y*v2.x

// }


// func main(){

//     var v1 vect2i
//     var v2 vect2i

//     v1.init(2,4)
//     v2.init(3,8)

//     fmt.Println(v1.addition(v2))

//     fmt.Println(v1.soustraction(v2))

//     fmt.Println(v1.multiplication(v2))

//     fmt.Println(v1.norme())

//     fmt.Println(v2.norme())

//     fmt.Println(v1.normalization())

// 	normalized, err := v1.normalization()

// 	if err != nil {

// 		fmt.Println("Erreur:", err)

// 	} else {

// 		fmt.Println("Normalisation de v1:", normalized)
// 	}

//     fmt.Println(v1.scalaire(v2))
//     fmt.Println(v1.vectoriel(v2))

// }

// //Exercice 4:

// //Question 1:

// type Node struct {

//     data int
//     Next *Node

// }

// type LinkedList struct{

//     head *Node

// }

// //Question 4.2.1:

// func (list *LinkedList) append(value int){

//     newNode := &Node{Data: value}


//     if (list.head == nil){

//         liste.head = newNode
//     }

//     current := list.head

// 	for current.next != nil {

// 		current = current.next

// 	}

// 	current.next = newNode
// }


// func (list *LinkedList) Print(){

//     current := list.head

//     for current != nil {

//         Fmt.Println(current.data)

//         current = current.next

//     }

//     Fmt.Println("nil")

// }

// func (list *LinkedList) Delete(data int) {

// 	if list.head == nil {
// 		fmt.Println("La liste est vide.")
// 	}

// 	// Si l'élément à supprimer est la tête de liste

// 	if list.head.data == data {

// 		list.head = list.head.next
// 	}

// 	// Parcours de la liste pour trouver l'élément

// 	current := list.head

// 	for current != nil && current.next != nil {

// 		if current.next.data == data {

// 			current.next = current.next.next

// 		}

// 		current = current.next
// 	}
// }


// func main() {
// 	// Création d'une liste chaînée
// 	list := LinkedList{}

// 	// Ajout d'éléments à la liste
// 	fmt.Println("Ajout des éléments 10, 20, 30 à la liste :")
// 	list.append(10)
// 	list.append(20)
// 	list.append(30)

// 	// Affichage de la liste
// 	fmt.Println("Liste après ajout :")
// 	list.Print()

// 	// Suppression d'un élément (par exemple, 20)
// 	fmt.Println("Suppression de l'élément 20 :")
// 	list.Delete(20)
// 	list.Print()

// 	// Suppression d'un élément qui n'existe pas (par exemple, 50)
// 	fmt.Println("Suppression de l'élément 50 (inexistant) :")
// 	list.Delete(50)
// 	list.Print()
// }



package main
import (
    "fmt"
    "math"
    "math/rand"
    "time"
    "os"
    "os/exec"
)

// Exercice 1: 

// Question 1:
func estBissextile(date int) bool {

    if date%4==0{

        if date%100 == 0 {
            
            if date%400 ==0{
                return true
            } else {
                return false
            }
        } else {
            return true 
        }
    
    }
    return false
}


// Question 2:

func estPremier(entier int) bool {

    if entier==1 {

        return false
    }


    for i := 2; i<entier; i++ {
        if entier%i==0 {
            return false
        }
    }

    return true

}

// Question 3:

func premiersNombresPremiers(entier int) []int{

    var tab []int

    for i:=1; i<=entier; i++ {

        if estPremier(i)==true {
            tab = append(tab,i)
        }
    }

    return tab

}

// Question 4:
func genererTableauAleatoire(entier int) []int{

    var tab []int

    for i:=0; i<entier; i++ {

        tab = append(tab,rand.Intn(100))

    }

    return tab

}

// Question 5:
func triBulles(T []int) []int{

    var tab = T
    var annexe int

    for i:=len(tab)-1; i >= 0; i--{
        for j:=0; j<i ; j++ {
            if tab[j+1] < tab[j] {
                annexe=tab[j]
                tab[j] = tab[j+1]
                tab[j+1] = annexe
            }
        }
    }
    return tab
}


// Question 6:

func triSelection(T []int) []int{

    var tab = T
    var annexe int

    for i:=0; i<len(tab)-1; i++{
        var min int = i
        for j:=i+1; j<len(tab); j++{
            if tab[j] < tab[min] {
                min = j
            }
        }
        if min !=i{
            annexe= tab[i]
            tab[i] = tab[min]
            tab[min] = annexe
        }
    }
    return tab
}

// Question 7:
func rechercheDichotomique(T []int, x int) bool{

    if len(T)==0{
        return false
    }
    var m int = len(T)/2
    if T[m] == x {
        return true
    } else if T[m]>x {
        return rechercheDichotomique(T[:m], x) 
    } else {
        return rechercheDichotomique(T[m+1:], x)
    }
}



func organiserParTaille(tab []string) map[int][]string{

    dico := make(map[int][]string)

    for _, valeur := range tab {

        longueur := len(valeur)
        if _, ok := dico[longueur]; ok {
            dico[longueur] = append(dico[longueur],valeur)
        } else {
            dico[longueur] = []string{valeur}
        }
    }

    return dico

}


// func main() {


//     // Fonction exercice 1:
//     fmt.Println(estBissextile(1900))
//     fmt.Println(estPremier(1))
//     fmt.Println(estPremier(10))
//     fmt.Println(premiersNombresPremiers(11))
//     fmt.Println(genererTableauAleatoire(20))
//     fmt.Println(triBulles([]int{1, 3, 2, 0}))
//     fmt.Println(triSelection([]int{1, 3, 2, 0}))
//     fmt.Println(rechercheDichotomique([]int{0, 1, 2, 3},0))
//     fmt.Println(organiserParTaille([]string{"abc", "a", "b", "bc"}))

//     // Fonction exercice 2:
//     var d [][]int = initGrille(6,6)
//     fmt.Println(d)
//     fmt.Println(compterVoisin(d, 3,3))
//     fmt.Println(update(d))


// }


// Exercice 2:
// Question 1:
func initGrille(n int,m int) [][]int {

    var tab [][]int

    for i := 0; i<n; i++{
        var row []int
        for j := 0; j<m; j++{
            row = append(row,rand.Intn(2))
        }
        tab = append(tab, row)
    }


    return tab

}

// Question 2:
func compterVoisin(grille [][]int, i int, j int) int {
    var nb int = 0
    var longueur int = len(grille[i]) - 1
    var hauteur int = len(grille) - 1

    // Vérifier le voisin du haut
    if i > 0 && grille[i-1][j] == 1 {
        nb=nb+1
    }

    // Vérifier le voisin du bas
    if i < hauteur && grille[i+1][j] == 1 {
        nb=nb+1
    }

    // Vérifier le voisin de gauche
    if j > 0 && grille[i][j-1] == 1 {
        nb=nb+1
    }

    // Vérifier le voisin de droite
    if j < longueur && grille[i][j+1] == 1 {
        nb=nb+1
    }

    // Vérifier le voisin en haut à gauche (diagonale)
    if i > 0 && j > 0 && grille[i-1][j-1] == 1 {
        nb=nb+1
    }

    // Vérifier le voisin en haut à droite (diagonale)
    if i > 0 && j < longueur && grille[i-1][j+1] == 1 {
        nb=nb+1
    }

    // Vérifier le voisin en bas à gauche (diagonale)
    if i < hauteur && j > 0 && grille[i+1][j-1] == 1 {
        nb=nb+1
    }

    // Vérifier le voisin en bas à droite (diagonale)
    if i < hauteur && j < longueur && grille[i+1][j+1] == 1 {
        nb=nb+1
    }

    return nb
}



func update(jeu_vie [][]int) [][]int {
    // Créer une nouvelle grille pour stocker les valeurs mises à jour
    nouvelleGrille := make([][]int, len(jeu_vie))
    for i := range jeu_vie {
        nouvelleGrille[i] = make([]int, len(jeu_vie[i]))
    }

    for i := 0; i < len(jeu_vie); i++ {
        for j := 0; j < len(jeu_vie[i]); j++ {
            voisins := compterVoisin(jeu_vie, i, j)
            if jeu_vie[i][j] == 1 {
                // Une cellule vivante reste vivante si elle a 2 ou 3 voisins
                if voisins == 2 || voisins == 3 {
                    nouvelleGrille[i][j] = 1
                } else {
                    nouvelleGrille[i][j] = 0
                }
            } else {
                // Une cellule morte devient vivante si elle a exactement 3 voisins
                if voisins == 3 {
                    nouvelleGrille[i][j] = 1
                } else {
                    nouvelleGrille[i][j] = 0
                }
            }
        }
    }

    return nouvelleGrille
}


// Fonction pour afficher la grille
func afficherGrille(grille [][]int) {
    for _, row := range grille {
        for _, cell := range row {
            if cell == 1 {
                fmt.Print("  ") // Cellule vivante = espace (1)
            } else {
                fmt.Print("\u2588 ") // Cellule morte = bloc noir (0)
            }
        }
        fmt.Println()
    }
}

func jouerJeuDeLaVie(n int, m int, iterations int, delai time.Duration) {
    grille := initGrille(n, m)

    for iter := 0; iter < iterations; iter++ {
        // Effacer l'affichage du terminal
        c := exec.Command("clear") // Pour Linux/Mac
        c.Stdout = os.Stdout
        c.Run()

        fmt.Println("Iteration", iter+1)
        afficherGrille(grille)

        time.Sleep(delai) // Pause entre les mises à jour

        grille = update(grille) // Met à jour la grille
    }
}


func main() {
    var d [][]int = initGrille(6, 6)
    fmt.Println("Grille initiale:")
    afficherGrille(d)

    fmt.Println("\nNombre de voisins pour (3,3) :", compterVoisin(d, 3, 3))

    fmt.Println("\nGrille après mise à jour:")
    d = update(d)
    afficherGrille(d)

    jouerJeuDeLaVie(6, 6, 10, 500*time.Millisecond)
}













// Exercice 3

// Question 1:





type vect2i struct {

    x int
    y int

}

func (v *vect2i) init(a int, b int){

    v.x = a
    v.y = b

}

func (v1 vect2i) addition(v2 vect2i) vect2i{

    var v3 vect2i

    v3.x = v1.x + v2.x
    v3.y = v1.y + v2.y

    return v3
    
}

func (v1 vect2i) soustraction(v2 vect2i) vect2i{

    var v3 vect2i

    v3.x = v1.x - v2.x
    v3.y = v1.y - v2.y

    return v3


}

func (v1 vect2i) multiplication(v2 vect2i) vect2i{

    var v3 vect2i

    v3.x = v2.x * v1.x
    v3.y = v2.y * v1.y

    return v3

}


func (v vect2i) norme() float64{
    return math.Sqrt(float64(v.x*v.x + v.y*v.y))
}

func (v vect2i) normalization() (vect2i, error) {

    norm := v.norme()

    if norm == 0 {

        return vect2i{0, 0}, fmt.Errorf("impossible de normaliser un vecteur nul")
    }

    return vect2i{int(math.Round(float64(v.x) / norm * 100)), int(math.Round(float64(v.y) / norm * 100))}, nil
}



func (v1 vect2i) scalaire(v2 vect2i) int {
    return v1.x*v2.x + v1.y*v2.y
}

func (v1 vect2i) vectoriel(v2 vect2i) int {

    return v1.x*v2.y - v1.y*v2.x

}


// func main(){

//     var v1 vect2i
//     var v2 vect2i

//     v1.init(2,4)
//     v2.init(3,8)

//     fmt.Println(v1.addition(v2))

//     fmt.Println(v1.soustraction(v2))

//     fmt.Println(v1.multiplication(v2))

//     fmt.Println(v1.norme())

//     fmt.Println(v2.norme())

//     fmt.Println(v1.normalization())

// 	normalized, err := v1.normalization()

// 	if err != nil {

// 		fmt.Println("Erreur:", err)

// 	} else {

// 		fmt.Println("Normalisation de v1:", normalized)
// 	}

//     fmt.Println(v1.scalaire(v2))
//     fmt.Println(v1.vectoriel(v2))

// }

// Exercice 4:

// Question 1:


// Définir la structure Node
type Node struct {
    data int
    Next *Node
}

// Définir la structure LinkedList
type LinkedList struct {
    head *Node
}

// Méthode pour ajouter un élément à la fin de la liste
func (list *LinkedList) append(value int) {
    newNode := &Node{data: value}

    if list.head == nil {
        list.head = newNode
        return
    }

    current := list.head
    for current.Next != nil {
        current = current.Next
    }
    current.Next = newNode
}

// Méthode pour afficher la liste
func (list *LinkedList) Print() {
    current := list.head
    for current != nil {
        fmt.Println(current.data)
        current = current.Next
    }
    fmt.Println("nil")
}

// Méthode pour supprimer un élément de la liste
func (list *LinkedList) Delete(data int) {
    if list.head == nil {
        fmt.Println("La liste est vide.")
        return
    }

    // Si l'élément à supprimer est la tête de liste
    if list.head.data == data {
        list.head = list.head.Next
        return
    }

    // Parcours de la liste pour trouver l'élément à supprimer
    current := list.head
    for current.Next != nil {
        if current.Next.data == data {
            current.Next = current.Next.Next
            return
        }
        current = current.Next
    }

    // Si l'élément n'est pas trouvé
    fmt.Println("Élément non trouvé dans la liste.")
}


func (list *LinkedList) InsertAtPosition(data int, position int) {
    newNode := &Node{data: data}

    if position == 0 { // Insérer en tête
        newNode.Next = list.head
        list.head = newNode
        return
    }

    if list.head == nil { // Cas où la liste est vide
        fmt.Println("La liste est vide, ajout en tant que premier élément.")
        list.head = newNode
        return
    }

    current := list.head
    for i := 0; i < position-1; i++ {
        if current.Next == nil { // Si la position dépasse la liste, on ajoute à la fin
            fmt.Println("Position hors limites, ajout à la fin de la liste.")
            current.Next = newNode
            return
        }
        current = current.Next
    }

    // Insère le nouvel élément à la bonne position
    newNode.Next = current.Next
    current.Next = newNode
}





// func main() {
//     // Création d'une liste chaînée
//     list := LinkedList{}

//     // Ajout d'éléments à la liste
//     list.append(10)
//     list.append(20)
//     list.append(30)

//     // Affichage de la liste
//     fmt.Println("Liste après ajout :")
//     list.Print()

//     // Suppression d'un élément (par exemple, 20)
//     fmt.Println("Suppression de l'élément 20 :")
//     list.Delete(10)
//     list.Print()

//     // Suppression d'un élément qui n'existe pas (par exemple, 50)
//     fmt.Println("Suppression de l'élément 50 (inexistant) :")
//     list.Delete(50)
//     list.Print()

//     // Insertion à une position hors limites
//     fmt.Println("\nInsertion de 100 à la position 10 :")
//     list.InsertAtPosition(100, 10)
//     list.Print()
// }





//4.3 : 

// 1: Une liste chaînée simple possède des nœuds avec un seul pointeur vers le suivant, tandis qu'une liste chaînée double a des nœuds avec deux pointeurs : un vers le précédent et un vers le suivant.
// 2: Si on essaie de supprimer un élément qui n'existe pas, l'algorithme va parcourir toute la liste sans le trouver, ce qui entraîne une perte de temps et peut potentiellement provoquer des erreurs.
//3: L'ajout en début de liste chaînée se fait en **O(1)**, car il suffit de modifier le pointeur de tête. En fin ou à une position donnée, il faut parcourir la liste, donc c'est en **O(n)**. La suppression nécessite d'abord de retrouver l'élément, ce qui prend **O(n)** en moyenne. Enfin, la recherche d’un élément spécifique est aussi en **O(n)**, sauf si la liste est optimisée d’une certaine manière.
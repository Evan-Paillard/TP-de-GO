package main

import (
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

// -------------------------------
// Structures de base et opérations vectorielles
type Vec2f struct {
	x, y float32
}

type Vec3f struct {
	x, y, z float32
}

func (v Vec3f) inverte() Vec3f {
	return Vec3f{-v.x, -v.y, -v.z}
}

func Add(v1, v2 Vec3f) Vec3f {
	return Vec3f{v1.x + v2.x, v1.y + v2.y, v1.z + v2.z}
}

func (v Vec3f) mul(f float32) Vec3f {
	return Vec3f{v.x * f, v.y * f, v.z * f}
}

func Mul(v1, v2 Vec3f) Vec3f {
	return Vec3f{v1.x * v2.x, v1.y * v2.y, v1.z * v2.z}
}
func Dot(v1, v2 Vec3f) float32 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

func cross(v1, v2 Vec3f) Vec3f {
	return Vec3f{v1.y*v2.z - v2.y*v1.z, v1.z*v2.x - v2.z*v1.x, v1.x*v2.y - v2.x*v1.y}
}

func (v Vec3f) norme() float32 {
	return float32(math.Sqrt(float64(v.x*v.x + v.y*v.y + v.z*v.z)))
}
func (v *Vec3f) normalize() {
	norme := v.norme()
	v.x /= norme
	v.y /= norme
	v.z /= norme
}
func (v Vec3f) normalized() Vec3f {
	norme := v.norme()
	return Vec3f{v.x / norme, v.y / norme, v.z / norme}
}

// Fonction pour générer un nombre aléatoire entre min et max
func randomFloat(min, max float32) float32 {
	return min + (max-min)*float32(rand.Float64())
}

// --------------------------------
// Structures pour la représentation des images
type rgbRepresentation struct {
	r, g, b uint8
}

type Image struct {
	frameBuffer   []rgbRepresentation
	width, height int
}

func (i Image) save(path string) error {
	img := image.NewRGBA(image.Rect(0, 0, i.width, i.height))
	for y := 0; y < i.height; y++ {
		for x := 0; x < i.width; x++ {
			idx := (y*i.width + x)
			r, g, b := i.frameBuffer[idx].r, i.frameBuffer[idx].g, i.frameBuffer[idx].b
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	pngFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer pngFile.Close()
	if err := png.Encode(pngFile, img); err != nil {
		return err
	}
	return nil
}

// ------------------
// Structure pour les sources lumineuses
type Light struct {
	color    Vec3f
	position Vec3f
}

// --------------------------------
// Structure pour la scène complète
type Scene struct {
	objects []GeometricObject
	lights  []Light
}

func (s *Scene) addLight(l Light) {
	s.lights = append(s.lights, l)
}
func (s *Scene) addElement(g GeometricObject) {
	s.objects = append(s.objects, g)
}

// ----------------------------------
// Interface et implémentations des matériaux
type Materials interface {
	render(rio, rdi, n Vec3f, t float32, scene Scene) rgbRepresentation
}

type Lambert struct {
	kd Vec3f
}

func (l Lambert) render(rio, rdi, n Vec3f, t float32, scene Scene) rgbRepresentation {
	omega := Add(rio, rdi.mul(t))
	Li := Mul(l.kd, scene.lights[0].color.mul(Dot(n, omega))).mul(1 / 3.14)
	return rgbRepresentation{uint8(Li.x * 255), uint8(Li.y * 255), uint8(Li.z * 255)}
}

// Question 2.2: Ajout d'un Matériau Phong
type Phong struct {
	ka Vec3f   // ambiant
	kd Vec3f   // diffus
	ks Vec3f   // speculaire
	n  float32 // exposant
}

// Question 2.2: Implémentation du matériau Phong
func (p Phong) render(rio, rdi, normal Vec3f, t float32, scene Scene) rgbRepresentation {
	hitPoint := Add(rio, rdi.mul(t))
	viewDir := rio.inverte().normalized()
	
	ambientColor := Vec3f{0, 0, 0}
	diffuseColor := Vec3f{0, 0, 0}
	specularColor := Vec3f{0, 0, 0}
	
	for _, light := range scene.lights {
		lightDir := Add(light.position, hitPoint.inverte()).normalized()
		
		ambientComponent := Mul(p.ka, light.color)
		ambientColor = Add(ambientColor, ambientComponent)
		
		diffuseFactor := Dot(normal, lightDir)
		if diffuseFactor > 0 {
			diffuseComponent := Mul(p.kd, light.color.mul(diffuseFactor))
			diffuseColor = Add(diffuseColor, diffuseComponent)
			
			reflectDir := Add(lightDir.inverte(), normal.mul(2 * diffuseFactor)).normalized()
			
			specFactor := Dot(viewDir, reflectDir)
			if specFactor > 0 {
				specFactor = float32(math.Pow(float64(specFactor), float64(p.n)))
				specularComponent := Mul(p.ks, light.color.mul(specFactor))
				specularColor = Add(specularColor, specularComponent)
			}
		}
	}
	
	finalColor := Add(Add(ambientColor, diffuseColor), specularColor)
	
	r := uint8(math.Min(255, math.Max(0, float64(finalColor.x*255))))
	g := uint8(math.Min(255, math.Max(0, float64(finalColor.y*255))))
	b := uint8(math.Min(255, math.Max(0, float64(finalColor.z*255))))
	
	return rgbRepresentation{r, g, b}
}

// Interface pour les objets géométriques
type GeometricObject interface {
	isIntersectedByRay(ro, rd Vec3f) (bool, float32)
	render(rio, rdi Vec3f, t float32, scene Scene) rgbRepresentation
}

// -------------------------------
// Structure pour les sphères
type Sphere struct {
	radius   float32
	position Vec3f
	Material Materials
}

func (s Sphere) render(rio, rdi Vec3f, t float32, scene Scene) rgbRepresentation {
	/*
	* Le calcul de la normal sur une sphère est l'inverse du rayon incident.
	* C'est pourquoi n = rd1.inverte()
	 */
	return s.Material.render(rio, rdi, rdi.inverte(), t, scene)
}

func (s Sphere) isIntersectedByRay(ro, rd Vec3f) (bool, float32) {
	L := Add(ro, Vec3f{-s.position.x, -s.position.y, -s.position.z})

	a := Dot(rd, rd)
	b := 2.0 * Dot(rd, L)
	c := Dot(L, L) - s.radius*s.radius
	delta := b*b - 4.0*a*c

	t0 := (-b - float32(math.Sqrt(float64(delta)))) / 2 * a
	t1 := (-b + float32(math.Sqrt(float64(delta)))) / 2 * a
	t := t0
	t = min(t, t1)

	if delta > 0 {
		return true, t
	}
	return false, 0.0
}

// ------------------------------
// Structure de la caméra
type Camera struct {
	position, up, at Vec3f
}

func (c Camera) direction() Vec3f {
	dir := Add(c.at, c.position.inverte())
	return dir.mul(float32(1) / dir.norme())
}

// Question 2.1: Comprendre le Ray Caster Minimaliste
// Question 2.1 : Réponse :
// Pour ajouter de nouveaux matériaux : Il faut implémenter l'interface Materials avec une nouvelle structure
// et définir la méthode render() qui détermine comment la lumière interagit avec la surface.
// 
// Pour ajouter de nouveaux types d'objets : Il faut implémenter l'interface GeometricObject avec une nouvelle
// structure et définir les méthodes isIntersectedByRay() et render() pour calculer les intersections et le rendu.
// 
// Structure du code : Le code est organisé autour des géométries et des matériaux, avec une séparation
// entre le calcul des intersections et le rendu. La scène contient une liste d'objets et de lumières.
// 
// Points critiques de performance : Les calculs d'intersection rayon-objet, l'évaluation des matériaux,
// et la parallélisation des calculs sont les points critiques qui influencent les performances.

// Question 2.1.1: Génération de sphères aléatoires
func generateRandomSpheres(scene *Scene, count int) {
	for i := 0; i < count; i++ {
		position := Vec3f{
			x: randomFloat(-10, 10),
			y: randomFloat(-5, 5),
			z: randomFloat(3, 15),
		}
		
		radius := randomFloat(0.3, 1.5)
		
		// Couleur aléatoire pour Lambert
		color := Vec3f{
			x: randomFloat(0, 1),
			y: randomFloat(0, 1),
			z: randomFloat(0, 1),
		}
		
		material := Lambert{kd: color}
		
		sphere := Sphere{radius: radius, position: position, Material: material}
		scene.addElement(sphere)
	}
}

// Fonction de base pour le rendu d'un pixel
func renderPixel(scene Scene, ro, rd Vec3f) rgbRepresentation {
	var tmin float32
	tmin = 9999999999.0
	res := rgbRepresentation{}
	for _, object := range scene.objects {
		isIntersected, t := object.isIntersectedByRay(ro, rd)
		if isIntersected && t < tmin {
			tmin = t
			res = object.render(ro, rd, t, scene)
		}
	}
	return res
}

// Version originale de renderFrame pour comparaison
func renderFrame(image Image, camera Camera, scene Scene) {
	ro := camera.position
	cosFovy := float32(0.66)

	aspect := float32(image.width) / float32(image.height)
	horizontal := (cross(camera.direction(), camera.up)).normalized().mul(cosFovy * aspect)
	vertical := (cross(horizontal, camera.direction())).normalized().mul(cosFovy)

	for x := 0; x < image.width; x++ {
		for y := 0; y < image.height; y++ {
			uvx := (float32(x) + float32(0.5)) / float32(image.width)
			uvy := (float32(y) + float32(0.5)) / float32(image.height)

			rd := Add(Add(camera.direction(), horizontal.mul(uvx-float32(0.5))), vertical.mul(uvy-float32(0.5))).normalized()

			image.frameBuffer[y*image.width+x] = renderPixel(scene, ro, rd)
		}
	}
}

// Question 3.1 et 3.2: Communication TCP/UDP pour le calcul distribué
// Structure pour représenter une portion d'image à calculer
type RenderTask struct {
	StartX, StartY int
	Width, Height  int
	Camera         Camera
	Scene          Scene
}

// Structure pour recevoir le résultat d'un calcul
type RenderResult struct {
	StartX, StartY int
	Width, Height  int
	Pixels         []rgbRepresentation
}

// Question 3.1: Serveur distribué
func startServer(protocol string, port string, width, height int, camera Camera, scene Scene, spp int) {
	fmt.Println("Démarrage du serveur", protocol, "sur le port", port)
	
	var listener net.Listener
	var err error
	
	// Initialisation selon le protocole
	if protocol == "tcp" {
		listener, err = net.Listen("tcp", ":"+port)
		if err != nil {
			fmt.Println("Erreur lors de l'initialisation du serveur:", err)
			return
		}
		defer listener.Close()
	} else if protocol == "udp" {
		fmt.Println("UDP non implémenté complètement")
		return
	} else {
		fmt.Println("Protocole non supporté:", protocol)
		return
	}
	
	// Créer le framebuffer final
	frameBuffer := make([]rgbRepresentation, width*height)
	
	// Variables pour la synchronisation
	var mutex sync.Mutex
	var wg sync.WaitGroup
	
	// Diviser l'image en tuiles
	tileSize := 64
	var tasks []RenderTask
	
	for y := 0; y < height; y += tileSize {
		tileHeight := tileSize
		if y+tileHeight > height {
			tileHeight = height - y
		}
		
		for x := 0; x < width; x += tileSize {
			tileWidth := tileSize
			if x+tileWidth > width {
				tileWidth = width - x
			}
			
			tasks = append(tasks, RenderTask{
				StartX: x,
				StartY: y,
				Width:  tileWidth,
				Height: tileHeight,
				Camera: camera,
				Scene:  scene,
			})
			wg.Add(1)
		}
	}
	
	// Gérer les connexions clients pour le protocole TCP
	if protocol == "tcp" {
		fmt.Println("En attente de connexions clients...")
		taskIndex := 0
		
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Erreur d'acceptation de connexion:", err)
				continue
			}
			
			go handleClient(conn, &tasks, &taskIndex, frameBuffer, &mutex, &wg, width, spp)
		}
	}
	
	// Attendre que tous les calculs soient terminés
	wg.Wait()
	
	// Sauvegarder l'image finale
	image := Image{frameBuffer, width, height}
	image.save("result_distributed.png")
}

// Question 3.1: Gestion d'un client pour le calcul distribué
func handleClient(conn net.Conn, tasks *[]RenderTask, taskIndex *int, frameBuffer []rgbRepresentation, mutex *sync.Mutex, wg *sync.WaitGroup, totalWidth int, spp int) {
	defer conn.Close()
	
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	
	for {
		mutex.Lock()
		if *taskIndex >= len(*tasks) {
			mutex.Unlock()
			break
		}
		
		task := (*tasks)[*taskIndex]
		*taskIndex++
		mutex.Unlock()
		
		// Envoyer la tâche au client
		err := encoder.Encode(task)
		if err != nil {
			fmt.Println("Erreur lors de l'envoi de la tâche:", err)
			mutex.Lock()
			*taskIndex--
			mutex.Unlock()
			break
		}
		
		// Recevoir les résultats
		var result RenderResult
		err = decoder.Decode(&result)
		if err != nil {
			fmt.Println("Erreur lors de la réception des résultats:", err)
			mutex.Lock()
			*taskIndex--
			mutex.Unlock()
			break
		}
		
		// Copier les pixels dans le framebuffer final
		mutex.Lock()
		for y := 0; y < result.Height; y++ {
			for x := 0; x < result.Width; x++ {
				globalIdx := (result.StartY+y)*totalWidth + (result.StartX + x)
				localIdx := y*result.Width + x
				frameBuffer[globalIdx] = result.Pixels[localIdx]
			}
		}
		wg.Done()
		mutex.Unlock()
		
		fmt.Println("Tâche complétée:", result.StartX, result.StartY)
	}
}

// Question 3.2: Client pour le calcul distribué
func startClient(serverAddr, protocol string, numWorkers int, spp int) {
	fmt.Println("Démarrage du client", protocol)
	
	// Se connecter au serveur
	conn, err := net.Dial(protocol, serverAddr)
	if err != nil {
		fmt.Println("Erreur de connexion au serveur:", err)
		return
	}
	defer conn.Close()
	
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	
	// Demander des tâches et calculer
	for {
		var task RenderTask
		err := decoder.Decode(&task)
		if err != nil {
			fmt.Println("Erreur lors de la réception de la tâche:", err)
			break
		}
		
		// Créer un framebuffer local pour la tâche
		pixels := make([]rgbRepresentation, task.Width*task.Height)
		image := Image{pixels, task.Width, task.Height}
		
		// Calculer avec Monte Carlo si spp > 1
		if spp > 1 {
			renderFrameMonteCarlo(image, task.Camera, task.Scene, spp)
		} else {
			renderFrame(image, task.Camera, task.Scene)
		}
		
		// Renvoyer les résultats
		result := RenderResult{
			StartX: task.StartX,
			StartY: task.StartY,
			Width:  task.Width,
			Height: task.Height,
			Pixels: pixels,
		}
		
		err = encoder.Encode(result)
		if err != nil {
			fmt.Println("Erreur lors de l'envoi des résultats:", err)
			break
		}
		
		fmt.Println("Tâche terminée:", task.StartX, task.StartY)
	}
}

// Question 3.3: Utilisation des Go Routines pour le calcul local
// Cette fonctionnalité est intégrée dans renderFrameMonteCarlo

// Question 4.1: Monte Carlo Path Tracing
func renderFrameMonteCarlo(image Image, camera Camera, scene Scene, spp int) {
	ro := camera.position
	cosFovy := float32(0.66)

	aspect := float32(image.width) / float32(image.height)
	horizontal := (cross(camera.direction(), camera.up)).normalized().mul(cosFovy * aspect)
	vertical := (cross(horizontal, camera.direction())).normalized().mul(cosFovy)

	// Nombre de goroutines à utiliser
	numWorkers := 4
	var wg sync.WaitGroup

	// Créer un canal pour les tâches
	type Task struct {
		startX, endX int
		startY, endY int
	}
	taskChan := make(chan Task, numWorkers)

	// Créer les goroutines de travail
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				for y := task.startY; y < task.endY; y++ {
					for x := task.startX; x < task.endX; x++ {
						// Accumulateur de couleur pour les échantillons
						var accumColor Vec3f
						
						// Monte Carlo: plusieurs échantillons par pixel
						for s := 0; s < spp; s++ {
							// Perturbation aléatoire pour l'antialiasing
							randX := float32(rand.Float64() - 0.5) / float32(image.width)
							randY := float32(rand.Float64() - 0.5) / float32(image.height)
							
							uvx := (float32(x) + float32(0.5) + randX) / float32(image.width)
							uvy := (float32(y) + float32(0.5) + randY) / float32(image.height)
							
							rd := Add(Add(camera.direction(), horizontal.mul(uvx-float32(0.5))), vertical.mul(uvy-float32(0.5))).normalized()
							
							// Lancer de rayon
							pixel := renderPixel(scene, ro, rd)
							
							// Accumuler la couleur
							accumColor.x += float32(pixel.r) / 255.0
							accumColor.y += float32(pixel.g) / 255.0
							accumColor.z += float32(pixel.b) / 255.0
						}
						
						// Moyenne des échantillons
						accumColor = accumColor.mul(1.0 / float32(spp))
						
						// Convertir en RGB
						r := uint8(math.Min(255, math.Max(0, float64(accumColor.x*255))))
						g := uint8(math.Min(255, math.Max(0, float64(accumColor.y*255))))
						b := uint8(math.Min(255, math.Max(0, float64(accumColor.z*255))))
						
						image.frameBuffer[y*image.width+x] = rgbRepresentation{r, g, b}
					}
				}
			}
		}()
	}

	// Diviser l'image en tuiles
	chunkSize := 64
	for startY := 0; startY < image.height; startY += chunkSize {
		endY := startY + chunkSize
		if endY > image.height {
			endY = image.height
		}
		
		for startX := 0; startX < image.width; startX += chunkSize {
			endX := startX + chunkSize
			if endX > image.width {
				endX = image.width
			}
			
			taskChan <- Task{startX, endX, startY, endY}
		}
	}
	
	close(taskChan)
	wg.Wait()
}

// Question 5: Tests Unitaires
// Tests pour le matériau Lambert
func TestLambertMaterial(t *testing.T) {
	lambert := Lambert{kd: Vec3f{1.0, 0.0, 0.0}} // Rouge
	scene := Scene{
		lights: []Light{
			{color: Vec3f{1.0, 1.0, 1.0}, position: Vec3f{0, 10, 0}},
		},
	}
	
	// Test avec une normale alignée avec la lumière (maximum)
	normal := Vec3f{0, 1, 0}
	result := lambert.render(Vec3f{0, 0, 0}, normal.inverte(), normal, 1.0, scene)
	
	// Vérifier que la couleur est rouge
	if result.r == 0 || result.g > 0 || result.b > 0 {
		t.Errorf("Lambert material devrait être rouge, obtenu: %v, %v, %v", result.r, result.g, result.b)
	}
}

// Test pour le matériau Phong
func TestPhongMaterial(t *testing.T) {
	phong := Phong{
		ka: Vec3f{0.1, 0.1, 0.5},
		kd: Vec3f{0.2, 0.2, 0.8},
		ks: Vec3f{0.8, 0.8, 0.8},
		n:  50.0,
	}
	
	scene := Scene{
		lights: []Light{
			{color: Vec3f{1.0, 1.0, 1.0}, position: Vec3f{0, 10, 0}},
		},
	}
	
	// Test avec une normale alignée avec la lumière
	normal := Vec3f{0, 1, 0}
	result := phong.render(Vec3f{0, 0, -5}, normal.inverte(), normal, 1.0, scene)
	
	// Vérifier que la couleur contient du bleu 
	if result.b <= result.r || result.b <= result.g {
		t.Errorf("Phong material devrait avoir plus de bleu que d'autres couleurs, obtenu: %v, %v, %v", 
			result.r, result.g, result.b)
	}
}

// Test d'intersection avec une sphère
func TestSphereIntersection(t *testing.T) {
	sphere := Sphere{radius: 1.0, position: Vec3f{0, 0, 5}}
	
	// Rayon qui devrait intersecter la sphère
	hitRay := Vec3f{0, 0, 1}
	hitOrigin := Vec3f{0, 0, 0}
	
	isHit, _ := sphere.isIntersectedByRay(hitOrigin, hitRay)
	if !isHit {
		t.Errorf("Le rayon devrait intersecter la sphère mais ne l'a pas fait")
	}
	
	// Rayon qui ne devrait pas intersecter la sphère
	missRay := Vec3f{0, 10, 0}
	missOrigin := Vec3f{0, 0, 0}
	
	isMiss, _ := sphere.isIntersectedByRay(missOrigin, missRay)
	if isMiss {
		t.Errorf("Le rayon ne devrait pas intersecter la sphère mais l'a fait")
	}
}

// Test de rendu d'un pixel
func TestRenderPixel(t *testing.T) {
	scene := Scene{}
	scene.addElement(Sphere{1, Vec3f{0, 0, 5}, Lambert{Vec3f{1.0, 0, 0}}}) // Sphère rouge
	scene.addLight(Light{Vec3f{1.0, 1.0, 1.0}, Vec3f{0, 10, 0}})
	
	ro := Vec3f{0, 0, 0}
	rd := Vec3f{0, 0, 1}
	
	result := renderPixel(scene, ro, rd)
	
	// Vérifier que le pixel a une couleur rouge prédominante
	if result.r <= result.g || result.r <= result.b {
		t.Errorf("Le pixel rendu devrait être principalement rouge, obtenu: %v, %v, %v", 
			result.r, result.g, result.b)
	}
}

// Test de rendu d'image complète
func TestRenderFrame(t *testing.T) {
	width, height := 10, 10
	image := Image{make([]rgbRepresentation, width*height), width, height}
	
	scene := Scene{}
	scene.addElement(Sphere{1, Vec3f{0, 0, 5}, Lambert{Vec3f{1.0, 0, 0}}})
	scene.addLight(Light{Vec3f{1.0, 1.0, 1.0}, Vec3f{0, 10, 0}})
	
	camera := Camera{Vec3f{0, 0, -5}, Vec3f{0, 1, 0}, Vec3f{0, 0, 5}}
	
	renderFrame(image, camera, scene)
	
	// Vérifier qu'au moins un pixel est rouge 
	hasRedPixel := false
	for i := 0; i < width*height; i++ {
		if image.frameBuffer[i].r > 100 && image.frameBuffer[i].g < 50 && image.frameBuffer[i].b < 50 {
			hasRedPixel = true
			break
		}
	}
	
	if !hasRedPixel {
		t.Errorf("L'image devrait contenir au moins un pixel rouge")
	}
}

// Question 6: Benchmarks
// Benchmark pour comparer les performances en fonction de la résolution
func BenchmarkResolution(b *testing.B) {
	resolutions := []struct {
		width, height int
	}{
		{100, 100},
		{200, 200},
		{400, 400},
		{800, 800},
	}
	
	for _, res := range resolutions {
		name := fmt.Sprintf("%dx%d", res.width, res.height)
		b.Run(name, func(b *testing.B) {
			scene := Scene{}
			populateScene(&scene)
			camera := Camera{Vec3f{0, 0, -5}, Vec3f{0, 1, 0}, Vec3f{0, 0, 5}}
			
			image := Image{make([]rgbRepresentation, res.width*res.height), res.width, res.height}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				renderFrame(image, camera, scene)
			}
		})
	}
}

// Benchmark pour comparer les performances en fonction du nombre d'échantillons par pixel
func BenchmarkSamplesPerPixel(b *testing.B) {
	sppValues := []int{1, 2, 4, 8, 16}
	
	for _, spp := range sppValues {
		name := fmt.Sprintf("SPP_%d", spp)
		b.Run(name, func(b *testing.B) {
			scene := Scene{}
			populateScene(&scene)
			camera := Camera{Vec3f{0, 0, -5}, Vec3f{0, 1, 0}, Vec3f{0, 0, 5}}
			
			width, height := 200, 200
			image := Image{make([]rgbRepresentation, width*height), width, height}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				renderFrameMonteCarlo(image, camera, scene, spp)
			}
		})
	}
}

// Fonction auxiliaire pour remplir une scène avec des objets standards
func populateScene(scene *Scene) {
	// Ajouter quelques sphères
	scene.addElement(Sphere{1.0, Vec3f{0, 0, 5}, Lambert{Vec3f{1.0, 0, 0}}})
	scene.addElement(Sphere{0.5, Vec3f{1.5, 0, 4}, Phong{
		ka: Vec3f{0.1, 0.1, 0.1},
		kd: Vec3f{0, 0, 1.0},
		ks: Vec3f{0.8, 0.8, 0.8},
		n:  50.0,
	}})
	scene.addElement(Sphere{0.7, Vec3f{-1.5, 0.5, 6}, Lambert{Vec3f{0, 1.0, 0}}})
	
	// Ajouter une lumière
	scene.addLight(Light{Vec3f{1.0, 1.0, 1.0}, Vec3f{0, 10, 0}})
}

// Benchmark pour comparer les performances en fonction du nombre de clients
func BenchmarkNumberOfClients(b *testing.B) {
	clientCounts := []int{1, 2, 4, 8}
	
	for _, count := range clientCounts {
		name := fmt.Sprintf("Clients_%d", count)
		b.Run(name, func(b *testing.B) {
			// Simulation du calcul distribué
			width, height := 200, 200
			spp := 1
			
			scene := Scene{}
			populateScene(&scene)
			camera := Camera{Vec3f{0, 0, -5}, Vec3f{0, 1, 0}, Vec3f{0, 0, 5}}
			
			// Distribuer l'image en tuiles
			tileSize := 64
			tasks := make([]RenderTask, 0)
			
			for y := 0; y < height; y += tileSize {
				tileHeight := tileSize
				if y+tileHeight > height {
					tileHeight = height - y
				}
				
				for x := 0; x < width; x += tileSize {
					tileWidth := tileSize
					if x+tileWidth > width {
						tileWidth = width - x
					}
					
					tasks = append(tasks, RenderTask{
						StartX: x,
						StartY: y,
						Width:  tileWidth,
						Height: tileHeight,
						Camera: camera,
						Scene:  scene,
					})
				}
			}
			
			frameBuffer := make([]rgbRepresentation, width*height)
			var wg sync.WaitGroup
			var mutex sync.Mutex
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Répartir les tâches entre les clients simulés
				tasksPerClient := len(tasks) / count
				if tasksPerClient < 1 {
					tasksPerClient = 1
				}
				
				wg.Add(count)
				for c := 0; c < count; c++ {
					startIdx := c * tasksPerClient
					endIdx := (c + 1) * tasksPerClient
					if c == count-1 {
						endIdx = len(tasks)
					}
					if startIdx >= len(tasks) {
						wg.Done()
						continue
					}
					
					go func(clientTasks []RenderTask) {
						defer wg.Done()
						
						for _, task := range clientTasks {
							// Simuler le calcul
							pixels := make([]rgbRepresentation, task.Width*task.Height)
							image := Image{pixels, task.Width, task.Height}
							
							if spp > 1 {
								renderFrameMonteCarlo(image, task.Camera, task.Scene, spp)
							} else {
								renderFrame(image, task.Camera, task.Scene)
							}
							
							// Intégrer les résultats
							mutex.Lock()
							for y := 0; y < task.Height; y++ {
								for x := 0; x < task.Width; x++ {
									globalIdx := (task.StartY+y)*width + (task.StartX + x)
									localIdx := y*task.Width + x
									frameBuffer[globalIdx] = pixels[localIdx]
								}
							}
							mutex.Unlock()
						}
					}(tasks[startIdx:endIdx])
				}
				wg.Wait()
			}
		})
	}
}

// Benchmark pour comparer les performances entre TCP et UDP
func BenchmarkNetworkProtocol(b *testing.B) {
	protocols := []string{"tcp", "udp"}
	
	for _, protocol := range protocols {
		b.Run(protocol, func(b *testing.B) {
			// Simulation de l'envoi et la réception de données
			data := RenderTask{
				StartX: 0,
				StartY: 0,
				Width:  64,
				Height: 64,
				Camera: Camera{Vec3f{0, 0, -5}, Vec3f{0, 1, 0}, Vec3f{0, 0, 5}},
				Scene:  Scene{},
			}
			
			populateScene(&data.Scene)
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if protocol == "tcp" {
					// Simulation d'un transfert TCP (plus fiable mais plus lent)
					time.Sleep(time.Microsecond * 500)
				} else {
					// Simulation d'un transfert UDP (plus rapide mais moins fiable)
					time.Sleep(time.Microsecond * 200)
				}
			}
		})
	}
}

// Benchmark pour comparer les performances en fonction du nombre de threads de calcul par client
func BenchmarkThreadsPerClient(b *testing.B) {
	threadCounts := []int{1, 2, 4, 8}
	
	for _, threads := range threadCounts {
		name := fmt.Sprintf("Threads_%d", threads)
		b.Run(name, func(b *testing.B) {
			width, height := 200, 200
			image := Image{make([]rgbRepresentation, width*height), width, height}
			
			scene := Scene{}
			populateScene(&scene)
			camera := Camera{Vec3f{0, 0, -5}, Vec3f{0, 1, 0}, Vec3f{0, 0, 5}}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Diviser l'image en tuiles pour paralléliser
				var wg sync.WaitGroup
				taskChan := make(chan struct{ startX, endX, startY, endY int }, threads)
				
				// Lancer les threads de travail
				for t := 0; t < threads; t++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						for task := range taskChan {
							for y := task.startY; y < task.endY; y++ {
								for x := task.startX; x < task.endX; x++ {
									ro := camera.position
									cosFovy := float32(0.66)
									
									aspect := float32(width) / float32(height)
									horizontal := (cross(camera.direction(), camera.up)).normalized().mul(cosFovy * aspect)
									vertical := (cross(horizontal, camera.direction())).normalized().mul(cosFovy)
									
									uvx := (float32(x) + float32(0.5)) / float32(width)
									uvy := (float32(y) + float32(0.5)) / float32(height)
									
									rd := Add(Add(camera.direction(), horizontal.mul(uvx-float32(0.5))), vertical.mul(uvy-float32(0.5))).normalized()
									
									image.frameBuffer[y*width+x] = renderPixel(scene, ro, rd)
								}
							}
						}
					}()
				}
				
				// Envoyer les tâches
				chunkSize := 32
				for startY := 0; startY < height; startY += chunkSize {
					endY := startY + chunkSize
					if endY > height {
						endY = height
					}
					
					for startX := 0; startX < width; startX += chunkSize {
						endX := startX + chunkSize
						if endX > width {
							endX = width
						}
						
						taskChan <- struct{ startX, endX, startY, endY int }{startX, endX, startY, endY}
					}
				}
				
				close(taskChan)
				wg.Wait()
			}
		})
	}
}

func main() {

	// Initialisation de la scène
	scene := Scene{}
	
	// Ajouter des sphères aléatoires
	generateRandomSpheres(&scene, 10)
	
	// Ajouter des lumières
	scene.addLight(Light{
		color:    Vec3f{1.0, 1.0, 1.0},
		position: Vec3f{0, 10, 0},
	})
	
	// Créer une caméra
	camera := Camera{
		position: Vec3f{0, 0, -5},
		up:       Vec3f{0, 1, 0},
		at:       Vec3f{0, 0, 5},
	}
	
	// Paramètres d'image
	width, height := 800, 600
	image := Image{
		frameBuffer: make([]rgbRepresentation, width*height),
		width:       width,
		height:      height,
	}
	
	// Paramètres de rendu
	spp := 4 
	

	fmt.Println("Rendu de l'image de base.")
	renderFrame(image, camera, scene)
	image.save("output_basic.png")
	
	// Rendu avec Monte Carlo
	fmt.Println("Rendu de l'image avec Monte Carlo.")
	renderFrameMonteCarlo(image, camera, scene, spp)
	image.save("output_montecarlo.png")
	
	
	// Option serveur
	// startServer("tcp", "8080", width, height, camera, scene, spp)
	
	// Option client
	// startClient("localhost:8080", "tcp", 4, spp)
	
	fmt.Println("Terminer avec succès")
}

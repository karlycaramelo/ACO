package main  

import (  
    "fmt"
    "io/ioutil"
    "strings"
    "strconv"
    //"math"
    "math/rand"
)

//Estructura que represneta un vertice
//Tiene aparte del indice del vertice un valor de la feromona inicial
//Y un valor que sera el valor de la feromona
type Vertex struct {
    index int
    pheromone_init float64
    pheromone float64
}

//Estructura para representar las aristas
//Las aristas solo tienen 2 entero que son los indices 
//De los vertices que conecta
//Tambien tiene un peso weight
type Edge struct {
    v1_index int
    v2_index int
    weight float64
}

//Estructura que representa una grafica
//Tiene un slice de veritces vertexes que es un aputnador
//Para que todas las graficas que definamos para las hormigas hagan referencia al mismo slice
//Tiene un slice de aristas edges que van a representar las aristas de la grafica original
//Tiene un slice de aristas full que van a repesentar la grafica completa
type Graph struct {
    vertexes *[]Vertex
    edges *[]Edge
    full *[]Edge
}

//Estructura que representa una hormiga
//La hormiga tiene un indice index
//Un apuntador a una grafica grap
//Un apuntador a un generador aleatorio random
//Un slice de enteros para ir guardandao la solucion 
//El parametro q0
//Un parametro para la evaporacion de la hormona (evaporation_rate) se usa para la rregal de actualizacion global
//Un parametro para el ajusta de la hormona (pheromone_adjust) se usa para la regla de actualizacion local
//Tambien va a tener el parametro beta que sirve para controlar el peso que tiene la funcion NR
type Ant struct {
    index int
    graph *Graph
    random *rand.Rand
    solution *[]int
    q0 float64
    evaporation_rate float64
    pheromone_adjust float64
    //beta float64
}


//Funicon que hace la actualizacion de manera local a la fermonona del vertice con indice indexVertex
func (a Ant) ActualizaFeromonaLocalmente(indexVertex int){
    pheromone0 := (*(*a.graph).vertexes)[indexVertex].pheromone_init
    pheromonei := (*(*a.graph).vertexes)[indexVertex].pheromone
    (*(*a.graph).vertexes)[indexVertex].pheromone = ((1-a.pheromone_adjust)*pheromonei)+(a.pheromone_adjust*pheromone0)
}


//Funcion que simula un pao para la hormiga
func (a Ant) Paso() {
    //Primero calculamos el valor aleatorio q
    q := a.random.Float64()
    //Calculamos el valor RN poara todos los vertices 
    rnvalues := a.GetVertexesRNValue() 
    //sumrnValues nos rive para guardar la suma de los valores rn de todos los vertices
    sumrnValues := 0.0
    //maxrnvalue nos srive para guardar el valor maximo de los valores  rn
    maxrnValue := 0.0
    //indexmasrnValue nos sirve para guardar el indice del vertice con el valor rn maximo
    indexmaxrnValue := -1 
    count := 0;
    //Dentro de este ciclo calculamos sumrnValues, maxrnValue y indexmaxrnValue
    for count < len(rnvalues) {
        //fmt.Printf("IndexV: %d RN: %f\n", count, rnvalues[count])
        if maxrnValue <= rnvalues[count] {
            maxrnValue = rnvalues[count]
            indexmaxrnValue = count
        }
        sumrnValues = sumrnValues + rnvalues[count]
        count = count+1
    }    


    //Si el q es menor al valor dado q0 entonces vamos a agregar a la solucion el vertice con indice indexmaxrnValue
    //Actualzamos la feromona de manera local para el mismo vertice, y por ulitmo para ese vertice aremos  a todas
    //Las aristas que inciden en el con valor 0
    if q < a.q0{
        //fmt.Printf("Feromona antes: %f\n", (*(*a.graph).vertexes)[indexmaxrnValue].pheromone)
        a.AgregaASolucion(indexmaxrnValue)
        a.ActualizaFeromonaLocalmente(indexmaxrnValue)
        (*a.graph).FullEdgesTo0ForVertex(indexmaxrnValue)
        //fmt.Printf("Se agrego %d\n", indexmaxrnValue)
        //fmt.Printf("Feromona Despues: %f\n", (*(*a.graph).vertexes)[indexmaxrnValue].pheromone)
        //a.PrintSolution()

    //Si q no es menor al valor dado q0
    }else{
        //Mientras que no da un paso la hormiga haremos lo sigueinte
        dioPaso := false
        for dioPaso != true{
            //Calculamos un indice de los vertices de manera aleatoria
            randomIndex := rand.Intn(len(rnvalues))
            //Calcualmos la probabilidad (como en el pdf) que la hormiga de un paso a este vertice
            indexProb := (rnvalues[randomIndex]/sumrnValues)
            //Calculamos otro numero aletaorioa entre 0 y 1
            randonNumber := a.random.Float64()
            //fmt.Printf("Random Index: %d indexProb: %f RandomNumber: %f\n", randomIndex, indexProb, randonNumber)
            //Si el numeor aleatorio calculado es menor o igual a la probabilidad de que la hormiga diera el paso a ese vertice
            //Entonces damos el paso a ese vertice
            //Vamos a agregar a la solucion el vertice con indice randomIndex
            //Actualzamos la feromona de manera local para el mismo vertice, y por ulitmo para ese vertice aremos  a todas
            //Las aristas que inciden en el con valor 0
            if randonNumber <= indexProb{
                a.AgregaASolucion(randomIndex)
                a.ActualizaFeromonaLocalmente(randomIndex)
                (*a.graph).FullEdgesTo0ForVertex(randomIndex)
                //fmt.Printf("Se agrego %d\n", randomIndex)
                //Si damos el paso entonces hacemos este fvalor true de  tal manera que esta fucnion se detendra
                dioPaso = true
            }
        }
    }
    //fmt.Printf("%f\n",  q)
    //fmt.Printf("%t q=%f < %f=q0\n", q < a.q0, q, a.q0)
}

//Meotod que nos dice si la hormiga puede dar un paso es decir si alguan de sus
//Aristas todavia tiene peso 1
func (a Ant) PuedeDarUnPaso() bool {
    peso := (*a.graph).FullWeight()
    if peso != 0{
        return true
    }else{
        return false
    }
}

//Calcula el valor RN para cada uno de los vertices para la hormiga y los regresa en un 
//Slice de float donde el indice i del slice corresponde al vertice i
func (a Ant) GetVertexesRNValue() []float64 {
    numVertex := len((*(*a.graph).vertexes))
    rnvalues := make([]float64, numVertex)
    //fmt.Printf("Num vertexes %d\n", numVertex)
    count := 0;
    for count < numVertex {
        indexVertex := (*(*a.graph).vertexes)[count].index
        pheromoneVertex := (*(*a.graph).vertexes)[count].pheromone
        weightVertex := (*a.graph).FullWeightOfVertex(indexVertex)
        rnvalues[indexVertex] = pheromoneVertex*weightVertex
        //fmt.Printf("Index: %d Pheromone: %f Weight: %f RN: %f\n", indexVertex,pheromoneVertex,weightVertex, rnvalues[indexVertex])
        count = count+1
    }
    return rnvalues
}

//Funcion que inicializa el slice solucion en 0
func (a Ant) BorraSolucion(){
    (*a.solution) = make([]int, 0)
}

//Funcion que agrega el indice de un vertice al slice de solucion
func (a Ant) AgregaASolucion(v int){
    (*a.solution) = append((*a.solution),v)
}


//Funcion dada una grafica imprime el contenido del slice de los idices de los vertices con "la solucion"
func (a Ant) PrintSolution(){
    count := 0;
    fmt.Printf("Solucion: ")
    for count < len((*a.solution)){
        fmt.Printf("%d ", (*a.solution)[count])
        count = count +1
    }
    fmt.Printf("\n")
}

//Funcion que dada una grafica nos da el peso total de las aristas en "FULL"
func (g Graph) FullWeight() float64 {
    count := 0;
    weight := 0.0
    for count < len(*g.full) {
        weight = weight + (*g.full)[count].weight 
        count = count + 1
    }
    return weight
}

//Funcion que dada un vertice v nos regresa la suma de los pesos de las aristas en "FULL" que inciden en este
//Vertice
func (g Graph) FullWeightOfVertex(v int) float64 {
    count := 0;
    weight := 0.0
    for count < len(*g.full) {
        v1temp := (*g.full)[count].v1_index
        v2temp := (*g.full)[count].v2_index
        if (v == v1temp || v == v2temp){
            weight = weight + (*g.full)[count].weight
        }
         
        count = count + 1
    }
    return weight
}

//Funcion que convierte el peso de todas aristas en "FULL" que inciden en el vertice v
//a 0
func (g Graph) FullEdgesTo0ForVertex(v int)  {
    count := 0;
    for count < len(*g.full) {
        v1temp := (*g.full)[count].v1_index
        v2temp := (*g.full)[count].v2_index
        if (v == v1temp || v == v2temp){
           //fmt.Printf("%d ---- %d\n", v1temp, v2temp)
           (*g.full)[count].weight = 0.0
        }
         
        count = count + 1
    }
}


//Funcion que dada una grafica y dos vertices v1 y v2
//En la lista de vertices "FULL" asigna el valor del peso a value
func (g Graph) SetEdge(v1 int, v2 int, value float64){
    count := 0;
    for count < len(*g.full) {
        v1temp := (*g.full)[count].v1_index
        v2temp := (*g.full)[count].v2_index
        if (v1 == v1temp && v2 == v2temp){
            //fmt.Printf("%d ---- %d\n", v1temp, v2temp)
            (*g.full)[count].weight = value
        }

        if (v2 == v1temp && v1 == v2temp){
            //fmt.Printf("%d ---- %d\n", v2temp, v1temp)
            (*g.full)[count].weight = value
        }
        count = count + 1
    }
}


//Funcion que nos dice si una par de vertices son elemento de alguna de las aristas originales
//Es decir si existe una arista en las "aristas orginales" que conecte al vertice v1 con el v2
func (g Graph) ExisteEnEdges(v1 int, v2 int) bool {
    count := 0;
    for count < len(*g.edges) {
        v1temp := (*g.edges)[count].v1_index
        v2temp := (*g.edges)[count].v2_index
        if (v1 == v1temp && v2 == v2temp){
            //fmt.Printf("%d ---- %d\n", v1temp, v2temp)
            return true
        }

        if (v2 == v1temp && v1 == v2temp){
            //fmt.Printf("%d ---- %d\n", v2temp, v1temp)
            return true
        }
        count = count + 1
    }
    return false
}


//Funcion que dada una grafica tomas sus vertices y las aristas orginales y a la grafica le
//Asigna los pesos de tal forma que las aristas en las "aristas originales" tiene peso 1 
//y las demas peso 0 
func (g Graph) initFull(){
    numVertex := len(*g.vertexes)
    count := 0
    i := 0
    for i < numVertex{
        j := i+1
        for j < numVertex{
            if g.ExisteEnEdges(i,j){
                (*g.full)[count] = Edge{i,j,1}
            }else{
                (*g.full)[count] = Edge{i,j,0}
            }
            count = count +1
            j = j +1
        }
        i = i +1
    }
}


//Funcion que imprime la represetancion en cadena de la grafica completa de  una grafica
func (g Graph) PrintFull() {
    count := 0;
    for count < len(*g.full) {
        (*g.full)[count].Print()
        count = count + 1
    }
}

//Funcion que imprime la representacion en cadena de una grafica (su vertices y las aristas "originales")
func (g Graph) Print() {
    count := 0;
    for count < len(*g.vertexes) {
        (*g.vertexes)[count].Print()
        count = count + 1
    }
    count = 0;
    for count < len(*g.edges) {
        (*g.edges)[count].Print()
        count = count + 1
    }
}

//Funcion que regrfesa el indice de un vertice
func (v Vertex) GetIndex() int {
    return v.index
}

//Funcion que regrfesa el valor de la feromona de un vertice
func (v Vertex) GetPheromone() float64 {
    return v.pheromone
}

//Funcion que imprime la representacion en cadena de un vertice
func (v Vertex) Print() {
    fmt.Printf("Vertice: %d\n    Feromona Ini: %f\n    Feromona Act: %f\n", v.index, v.pheromone_init, v.pheromone)
}

//Funcion que regresa el peso de una arista
func (e Edge) GetWeight() float64 {
    return e.weight
}

//Funcion que regesa el vertice1 de una arista
func (e Edge) GetVertex1() int {
    return e.v1_index
}

//Funcion que regesa el vertice2 de una arista
func (e Edge) GetVertex2() int {
    return e.v2_index
}

//Funcincion que imprime la representacion en cadena de una arista
func (e Edge) Print() {
    fmt.Printf("%d --(%f)-- %d\n", e.v1_index, e.weight, e.v2_index)
}


//Funcion principal que corre ACO
func main() {
    //Leemos el archivo
    data, err := ioutil.ReadFile("g.txt")
    if err != nil {
        fmt.Println("File reading error", err)
        return
    }
    rows := strings.Split(string(data), "\n")
    //Inicializamos los vertices
    numOfVertexes, _ := strconv.Atoi(rows[0])
    vertexesG := make([]Vertex, numOfVertexes)
    count := 0;
    for count < numOfVertexes {
        vertexesG[count] = Vertex{count, 0.1, 0.2}
        count = count + 1
    }

    //Inicializamos las aristas con los valores del archivo
    numOfEdges := len(rows)-2
    edgesG := make([]Edge, numOfEdges)
    count = 0;
    for count < numOfEdges {
        vertIndx := strings.Split(rows[count+1], ",")  
        vertex1, _ := strconv.Atoi(vertIndx[0])
        vertex2, _ := strconv.Atoi(vertIndx[1])
        edgesG[count] = Edge{vertex1,vertex2,1}
        count = count + 1
    }

    //De lo anterior tenemos el conjuinto de vertices vertexesG y el conjunto de aristas edgesG
    //El apuntador a estos dos conjuntos se les pasara  a todas las hormigas para que los compartan
    //Y todas tengan los mismos conjuntos


    //VARIABLE QUE NOS DIRA CUANTAS VECES SE HA ENCONTRADO UNA SOLUCION CON EL MISMO TAMAÑO DE MANERA CONSECUTIVA 
    numSinCambios := 0
    //Variable que nos dice el numero de veces que si se encuentra la solcucion con el mismo tamaño se detendra el algoritmo
    numIteracionSeguidasParaParar := 200
    //Defefine el numero de hormigas que tendremos
    numberOfAnts := 20
    //Calculamos la cantidad de aristas que debe de tener la greafica comompleta dada la cantidad de vertices que tenemos
    numOfEdgesFull := (numOfVertexes*(numOfVertexes-1))/2

    //VARIABLES QUE DEFINE LOS PARAMETROS DEL ALGORITMO
    q := 0.5
    evaporation_rate := 0.2
    pheromone_adjust := 0.12
    //beta := 1
    //VARIABLES QUE DEFINE LOS PARAMETROS DEL ALGORITMO

    //Semilla para la funcion de numeros aleatorios
    randomSeed := 1
    //La funcion para generar los numeros aletarios que usaremos
    randomFun := rand.New(rand.NewSource(int64(randomSeed)))
    min_sol := make([]int, numOfVertexes)
    
    //Inicializamos un slice de hormigas con la cantida de hormigas definida
    ants := make([]Ant, numberOfAnts)

    antCount := 0
    //Para cada hormiga le vamos a crear su slice de edgesFull, una grafica con el apuntaos a los vertifces
    //el aputnador a las aristas y sus aristas que rempresentaran a la grafica completa
    //Depues vamos a crear ala hormiga con los parametros
    //La grafica que le creamos, un slice de enteros para sus soluciones, la funcion aleatorio y los parametros 
    for antCount < numberOfAnts {
        edgesFull := make([]Edge, numOfEdgesFull)
        graphG := Graph{&vertexesG,&edgesG,&edgesFull}
        slice := make([]int, 0)
        ants[antCount] = Ant{antCount,&graphG,randomFun, &slice,q,evaporation_rate,pheromone_adjust}
        antCount = antCount +1
    }
    
    ///////////////////////////////////
    //////////ALGORITMO ACO////////////
    ///////////////////////////////////
    //Mientras no se tengan numIteracionSeguidasParaParar sin cambios en la solucion
    //Ejecutaremos el siguiente ciclo

////////////////////////CICLO////////////////////////////
    for numSinCambios <= numIteracionSeguidasParaParar{

        //fmt.Printf("Sin cambios: %d\n", numSinCambios)

        //Inicializamos a cada una de las hormigas estos es 
        //Inicializar la grafica full
        //Inicialar el slice de soluciones
        //no necesitamos poner a la hormiga en un vertice arbitratio en un principio por que desde todos los vertices tenemos conexion a todos por ser grafica
        //a todos por ser grafica completa
        antCount = 0
        for antCount < numberOfAnts {
            (*ants[antCount].graph).initFull()
            ants[antCount].BorraSolucion()
            antCount = antCount +1
        }    
    
        //Mientras alguno de las hormigas pueda dar un paso
        sePuedeDarPaso := true
        for sePuedeDarPaso != false {

            otroPaso := false
            antCount = 0
            //Verificamos si alguna de las hormigas todavia puede dar un paso
            for antCount < numberOfAnts {
                if ants[antCount].PuedeDarUnPaso(){
                    otroPaso = true
                }
                antCount = antCount +1
            }
            
            //Si alguna hormiga todavia puede dar un paso
            sePuedeDarPaso = otroPaso
            if sePuedeDarPaso{
                antCount = 0
                for antCount < numberOfAnts {
                    //Verificamos si la hormiga con index antCount puede dar un paso y si es asi
                    //Esta damos el paso con esta hormiga
                    if ants[antCount].PuedeDarUnPaso(){

////////////////////////PASO//////////////////////////// (VER EL LA FUNCION)
                        ants[antCount].Paso()        
                    }
                    antCount = antCount +1
                }
            }
        }

///////////TODAS LAS HORMIGAS COMPLETAN SU TRAYECTO/////
        //Una vez que ya se dieron todos los pasos posibles debemo encontrar la mejor solucion guardarla y actualziar la feromona
        antCount = 0
        //El tamaño de la solucion minima
        minSolLen := -1
        //El indice de la hormiga que tiene la solucion minima
        minSolIndex := -1
        //Buscamos la solucion minima entre todas las hormigas
        for antCount < numberOfAnts {
            solLen := len((*ants[antCount].solution))
            if minSolLen > solLen || minSolLen < 0{
                minSolLen = solLen
                minSolIndex = antCount
            }
            antCount = antCount +1
        }
        ants[minSolIndex].PrintSolution()

        //Verificamos que la solucioon mejor encontrada en este ciclo sea mejor que minima solucion actual
        if len(min_sol) >= len((*ants[minSolIndex].solution)){
           //Si la solucion tiene el mismo tamaño entonces sumaos 1 al contador de ciclos sin con el mismo 
           //Tamaño de solucion
           if len(min_sol) == len((*ants[minSolIndex].solution)){
               numSinCambios = numSinCambios+1
           //Si la solucion es mas pequeña regreamos el contador a 0 
           }else{
               numSinCambios = 0
           }

            //Borramos la mejor solucion anterior
            min_sol = make([]int, len((*ants[minSolIndex].solution)))
            countSolIndex := 0
            //Copiamos la nueva solucion minima
            for countSolIndex < len(min_sol){
                min_sol[countSolIndex] = (*ants[minSolIndex].solution)[countSolIndex]
                countSolIndex = countSolIndex +1
            }
        }

        countSolIndex := 0
        //Imprimimos la mejor solucion hasta el momento
        fmt.Printf("MejorSolucion: ")
        for countSolIndex < len(min_sol){
            fmt.Printf("%d ", min_sol[countSolIndex])
            countSolIndex = countSolIndex +1
        }
        fmt.Printf("\n")

////////////////////////ACTUALIZACION GLOBAL////////////////////////////
        //Por ultimo vamos a hacer la actualizacion de la feromona de manera GLOBAL
        countVertexIndex := 0
        //Para cada uno de los vertices calculamos el nuevo valor de la hormona considerando la evaporacion
        for countVertexIndex < len(vertexesG){
            vertexIndex := vertexesG[countVertexIndex].index
            vertexPheromone := vertexesG[countVertexIndex].pheromone
            newPheromoneValue := (1-evaporation_rate)*vertexPheromone
            addedPheromone := 0.0
            countSolIndex := 0
            //Si el vertice es parte de la solucion minima actual entonces tambien calculamos la feromona extra que se le sumara
            for countSolIndex < len(min_sol){
                if vertexIndex == min_sol[countSolIndex]{
                    addedPheromone = evaporation_rate*(1.0/float64((len(min_sol))))
                    //fmt.Printf("AddedPhero %f\n",addedPheromone)
                }
                countSolIndex = countSolIndex +1
            } 
            //Actualizamos el valor de la feromona
            vertexesG[countVertexIndex].pheromone = newPheromoneValue + addedPheromone
            countVertexIndex = countVertexIndex +1
        }

    }

}

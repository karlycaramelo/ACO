package main  

import (  
    "fmt"
    "io/ioutil"
    "strings"
    "strconv"
    //"math"
    "math/rand"
)

type Vertex struct {
    index int
    pheromone_init float64
    pheromone float64
}

type Edge struct {
    v1_index int
    v2_index int
    weight float64
}

type Graph struct {
    vertexes *[]Vertex
    edges *[]Edge
    full *[]Edge
}

type Ant struct {
    index int
    graph *Graph
    random *rand.Rand
    solution *[]int
    q0 float64
    evaporation_rate float64
    pheromone_adjust float64
}

func (a Ant) ActualizaFeromona(indexVertex int){
    pheromone0 := (*(*a.graph).vertexes)[indexVertex].pheromone_init
    pheromonei := (*(*a.graph).vertexes)[indexVertex].pheromone
    (*(*a.graph).vertexes)[indexVertex].pheromone = ((1-a.pheromone_adjust)*pheromonei)+(a.pheromone_adjust*pheromone0)
}

func (a Ant) Paso() {
    q := a.random.Float64()
    rnvalues := a.GetVertexesRNValue() 
    sumrnValues := 0.0
    maxrnValue := 0.0
    indexmaxrnValue := -1 
    count := 0;
    for count < len(rnvalues) {
        //fmt.Printf("IndexV: %d RN: %f\n", count, rnvalues[count])
        if maxrnValue <= rnvalues[count] {
            maxrnValue = rnvalues[count]
            indexmaxrnValue = count
        }
        sumrnValues = sumrnValues + rnvalues[count]
        count = count+1
    }
    count = 0;
    //fmt.Printf("q<q0 Maximo Indice: %d Maximo Value: %f\n", indexmaxrnValue, maxrnValue)
    for count < len(rnvalues) {
        //fmt.Printf("IndexV: %d q>=q0: %f\n", count, (rnvalues[count]/sumrnValues))
        count = count+1
    }
    

    if q < a.q0{
        //fmt.Printf("Feromona antes: %f\n", (*(*a.graph).vertexes)[indexmaxrnValue].pheromone)
        a.AgregaASolucion(indexmaxrnValue)
        a.ActualizaFeromona(indexmaxrnValue)
        (*a.graph).FullEdgesTo0ForVertex(indexmaxrnValue)
        //fmt.Printf("Se agrego %d\n", indexmaxrnValue)
        //fmt.Printf("Feromona Despues: %f\n", (*(*a.graph).vertexes)[indexmaxrnValue].pheromone)
        //a.PrintSolution()
    }else{
        dioPaso := false
        for dioPaso != true{
            randomIndex := rand.Intn(len(rnvalues))
            indexProb := (rnvalues[randomIndex]/sumrnValues)
            randonNumber := a.random.Float64()
            //fmt.Printf("Random Index: %d indexProb: %f RandomNumber: %f\n", randomIndex, indexProb, randonNumber)
            if randonNumber <= indexProb{
                a.AgregaASolucion(randomIndex)
                a.ActualizaFeromona(randomIndex)
                (*a.graph).FullEdgesTo0ForVertex(randomIndex)
                //fmt.Printf("Se agrego %d\n", randomIndex)
                dioPaso = true
            }


        }
    }
    //fmt.Printf("%f\n",  q)
    //fmt.Printf("%t q=%f < %f=q0\n", q < a.q0, q, a.q0)
}

func (a Ant) PuedeDarUnPaso() bool {
    peso := (*a.graph).FullWeight()
    if peso != 0{
        return true
    }else{
        return false
    }
}

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

func (a Ant) BorraSolucion(){
    (*a.solution) = make([]int, 0)
}

func (a Ant) AgregaASolucion(v int){
    (*a.solution) = append((*a.solution),v)
}

func (a Ant) PrintSolution(){
    count := 0;
    fmt.Printf("Solucion: ")
    for count < len((*a.solution)){
        fmt.Printf("%d ", (*a.solution)[count])
        count = count +1
    }
    fmt.Printf("\n")
}

func (g Graph) FullWeight() float64 {
    count := 0;
    weight := 0.0
    for count < len(*g.full) {
        weight = weight + (*g.full)[count].weight 
        count = count + 1
    }
    return weight
}

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

func (g Graph) PrintFull() {
    count := 0;
    for count < len(*g.full) {
        (*g.full)[count].Print()
        count = count + 1
    }
}

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

func (v Vertex) GetIndex() int {
    return v.index
}

func (v Vertex) GetPheromone() float64 {
    return v.pheromone
}

func (v Vertex) Print() {
    fmt.Printf("Vertice: %d\n    Feromona Ini: %f\n    Feromona Act: %f\n", v.index, v.pheromone_init, v.pheromone)
}

func (e Edge) GetWeight() float64 {
    return e.weight
}

func (e Edge) GetVertex1() int {
    return e.v1_index
}

func (e Edge) GetVertex2() int {
    return e.v2_index
}

func (e Edge) Print() {
    fmt.Printf("%d --(%f)-- %d\n", e.v1_index, e.weight, e.v2_index)
}

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
/*
    count = 0;
    for count < numOfVertexes {
        vertexesG[count].Print()
        count = count + 1
    }
*/
    //Inicializamos las aristas
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

/*  
    count = 1;
    for count < numOfEdges {
        edgesG[count-1].Print()
        count = count + 1
    }
*/

    numSinCambios := 0
    q := 0.5
    evaporation_rate := 0.2
    pheromone_adjust := 0.12
    numOfEdgesFull := (numOfVertexes*(numOfVertexes-1))/2
    numberOfAnts := 2
    randomSeed := rand.New(rand.NewSource(int64(1)))
    min_sol := make([]int, numOfVertexes)
    
    ants := make([]Ant, numberOfAnts)
    antCount := 0
    for antCount < numberOfAnts {
        edgesFull := make([]Edge, numOfEdgesFull)
        graphG := Graph{&vertexesG,&edgesG,&edgesFull}
        slice := make([]int, 0)
        ants[antCount] = Ant{antCount,&graphG,randomSeed, &slice,q,evaporation_rate,pheromone_adjust}
        antCount = antCount +1
    }
    




    for numSinCambios <= 200{
    fmt.Printf("Sin cambios: %d\n", numSinCambios)
    //Iniciamos la graficaFull de las hormigas
    antCount = 0
    for antCount < numberOfAnts {
        (*ants[antCount].graph).initFull()
        ants[antCount].BorraSolucion()
        antCount = antCount +1
    }    
    
    sePuedeDarPaso := true
    for sePuedeDarPaso != false {

        otroPaso := false
        antCount = 0
        for antCount < numberOfAnts {
            if ants[antCount].PuedeDarUnPaso(){
                otroPaso = true
            }
            antCount = antCount +1
        }
        sePuedeDarPaso = otroPaso
        if sePuedeDarPaso{
            antCount = 0
            for antCount < numberOfAnts {
                if ants[antCount].PuedeDarUnPaso(){
                    ants[antCount].Paso()        
                }
                antCount = antCount +1
            }
        }
    }

    antCount = 0
    minSolLen := -1
    minSolIndex := -1
    for antCount < numberOfAnts {
        solLen := len((*ants[antCount].solution))
        if minSolLen > solLen || minSolLen < 0{
            minSolLen = solLen
            minSolIndex = antCount
        }
        antCount = antCount +1
    }
    ants[minSolIndex].PrintSolution()
    if len(min_sol) >= len((*ants[minSolIndex].solution)){
    fmt.Printf("Entra en igual o mayor: ")
       if len(min_sol) == len((*ants[minSolIndex].solution)){
           numSinCambios = numSinCambios+1
       }else{
           numSinCambios = 0
       }

        min_sol = make([]int, len((*ants[minSolIndex].solution)))
        countSolIndex := 0
        for countSolIndex < len(min_sol){
            min_sol[countSolIndex] = (*ants[minSolIndex].solution)[countSolIndex]
            countSolIndex = countSolIndex +1
        }
    }

    countSolIndex := 0
    fmt.Printf("MejorSolucion: ")
    for countSolIndex < len(min_sol){
        fmt.Printf("%d ", min_sol[countSolIndex])
        countSolIndex = countSolIndex +1
    }
    fmt.Printf("\n")

    //Actualiza feromona GLOBAL
    countVertexIndex := 0
    for countVertexIndex < len(vertexesG){
        //vertexesG[countVertexIndex].Print()
        countVertexIndex = countVertexIndex +1
    }

    countVertexIndex = 0
    for countVertexIndex < len(vertexesG){
        vertexIndex := vertexesG[countVertexIndex].index
        vertexPheromone := vertexesG[countVertexIndex].pheromone
        newPheromoneValue := (1-evaporation_rate)*vertexPheromone
        addedPheromone := 0.0
        countSolIndex := 0
        for countSolIndex < len(min_sol){
            if vertexIndex == min_sol[countSolIndex]{
                addedPheromone = evaporation_rate*(1.0/float64((len(min_sol))))
                //fmt.Printf("AddedPhero %f\n",addedPheromone)
            }
            countSolIndex = countSolIndex +1
        } 
        vertexesG[countVertexIndex].pheromone = newPheromoneValue + addedPheromone
        countVertexIndex = countVertexIndex +1
    }

    countVertexIndex = 0
    for countVertexIndex < len(vertexesG){
        //vertexesG[countVertexIndex].Print()
        countVertexIndex = countVertexIndex +1
    }
    }

    //fmt.Printf("Minimo %d\n",minSolIndex)
    //fmt.Printf("Ant %d\n",0)
    //ants[0].PrintSolution()
    //fmt.Printf("Ant %d\n",1)
    //ants[1].PrintSolution()


//(*ants[0].graph).Print()
//(*ants[1].graph).Print()

/*
    (*ants[0].graph).Print()
    (*ants[1].graph).Print()
    (*(*ants[0].graph).vertexes)[2].pheromone = 3.3
    (*(*ants[1].graph).vertexes)[4].pheromone = 4.3
    (*ants[0].graph).Print()
    (*ants[1].graph).Print()
    (*ants[0].graph).initFull()
    (*ants[1].graph).initFull()
    ants[0].PrintSolution()
    ants[1].PrintSolution()
    ants[0].AgregaASolucion(1)
    ants[0].AgregaASolucion(2)
    ants[0].AgregaASolucion(3)
    ants[1].AgregaASolucion(2)
    ants[0].PrintSolution()
    ants[1].PrintSolution()
    ants[0].BorraSolucion()
    ants[0].PrintSolution()
    ants[1].PrintSolution()

    //fmt.Printf("Existe %d %d %t \n", 7,5,graphG1.ExisteEnEdges(7,5))

    //graphG1.initFull()
    //graphG1.PrintFull()
    //fmt.Printf("Weight %f \n", graphG1.FullWeight())

    //graphG.SetEdge(4, 6, 2.2)
    //graphG.PrintFull()

    fmt.Printf("Weight %f \n", graphG.FullWeight())
    fmt.Printf("Weight of %d %f \n", 0,graphG.FullWeightOfVertex(0))
    fmt.Printf("Weight of %d %f \n", 1,graphG.FullWeightOfVertex(1))
    fmt.Printf("Weight of %d %f \n", 2,graphG.FullWeightOfVertex(2))
    fmt.Printf("Weight of %d %f \n", 3,graphG.FullWeightOfVertex(3))
    fmt.Printf("Weight of %d %f \n", 4,graphG.FullWeightOfVertex(4))
    fmt.Printf("Weight of %d %f \n", 5,graphG.FullWeightOfVertex(5))
    fmt.Printf("Weight of %d %f \n", 6,graphG.FullWeightOfVertex(6))
    fmt.Printf("Weight of %d %f \n", 7,graphG.FullWeightOfVertex(7))
*/

}

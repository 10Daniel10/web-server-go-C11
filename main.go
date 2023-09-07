package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Product representa la estructura de un producto
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

var (
	products  []Product
	idCounter int
	mutex     sync.Mutex
)

func main() {

	// Agregar productos de prueba
	products = []Product{
		{
			ID:          1,
			Name:        "Cheese - St. Andre",
			Quantity:    60,
			CodeValue:   "S73191A",
			IsPublished: true,
			Expiration:  "12/04/2022",
			Price:       50.15,
		},
		{
			ID:          2,
			Name:        "Apples",
			Quantity:    100,
			CodeValue:   "A12345",
			IsPublished: true,
			Expiration:  "25/12/2022",
			Price:       1.99,
		},
	}

	router := gin.Default()

	// Ruta para agregar un producto (POST)
	router.POST("/products", addProduct)

	// Ruta para obtener un producto por ID (GET)
	router.GET("/products/:id", getProductByID)

	router.Run(":8080")
}

func getProductByID(c *gin.Context) {
	// Obtener el ID del producto de los parámetros de la URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de producto inválido"})
		return
	}

	// Buscar el producto en la lista
	for _, p := range products {
		if p.ID == id {
			c.JSON(http.StatusOK, p)
			return
		}
	}

	// Si no se encuentra el producto, responder con un error
	c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
}

func addProduct(c *gin.Context) {
	var newProduct Product

	// Decodificar el cuerpo JSON del request en una estructura Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar el producto"})
		return
	}

	// Validar los campos del producto
	if newProduct.Name == "" || newProduct.Quantity <= 0 || newProduct.CodeValue == "" || newProduct.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos obligatorios incompletos"})
		return
	}

	if !isValidDate(newProduct.Expiration) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fecha de vencimiento inválida"})
		return
	}

	// Verificar si el campo is_published está vacío y establecerlo en false si es el caso
	if !newProduct.IsPublished {
		newProduct.IsPublished = false
	}

	// Verificar la unicidad del campo code_value
	for _, p := range products {
		if p.CodeValue == newProduct.CodeValue {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El código ya existe"})
			return
		}
	}

	// Generar un nuevo ID
	mutex.Lock()
	idCounter++
	newProduct.ID = idCounter
	mutex.Unlock()

	// Agregar el nuevo producto a la lista
	products = append(products, newProduct)

	// Responder con el nuevo producto creado
	c.JSON(http.StatusCreated, newProduct)
}

func isValidDate(dateStr string) bool {
	_, err := time.Parse("02/01/2006", dateStr)
	return err == nil
}

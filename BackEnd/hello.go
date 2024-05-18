package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

type Customers struct {
	Customers []Customer `json:"customers"`
}
type Customer struct {
	CustomerID      int    `json:"customerId"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	IsPrimeCustomer bool   `json:"isPrimeCustomer"`
}

func readJsonFile() (Customers, error) {
	jsonFile, err := os.Open("sample.json")

	if err != nil {
		return Customers{}, errors.New("file opening error")
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var customers Customers

	json.Unmarshal(byteValue, &customers)

	return customers, nil
}

func writeJsonFile(bytes []byte) error {
	jsonFile, err := os.OpenFile("sample.json", os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	if err := jsonFile.Truncate(0); err != nil {
		return err
	}

	if _, err := jsonFile.Seek(0, 0); err != nil { // Move the file cursor to the beginning
		return err
	}

	_, err = jsonFile.Write(bytes)

	if err != nil {
		return err
	}

	return nil
}

func getAllCustomers(ctx *gin.Context) {
	allCusts, err := readJsonFile()

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, allCusts.Customers)
}

func addOrUpdateCustomer(ctx *gin.Context) {
	var newOrExistingCustomer Customer
	var isExistingCustomer bool

	if err := ctx.BindJSON(&newOrExistingCustomer); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, "Invalid Json Given")
		return
	}

	fmt.Printf("the customerId we got was", newOrExistingCustomer.FirstName)

	//add or update the customer
	//Firstly read all customers
	allCusts, err := readJsonFile()

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	for i := 0; i < len(allCusts.Customers); i++ {
		if allCusts.Customers[i].CustomerID == newOrExistingCustomer.CustomerID {
			allCusts.Customers[i].FirstName = newOrExistingCustomer.FirstName
			allCusts.Customers[i].LastName = newOrExistingCustomer.LastName
			allCusts.Customers[i].IsPrimeCustomer = newOrExistingCustomer.IsPrimeCustomer
			isExistingCustomer = true
			break
		}
	}

	if !isExistingCustomer {
		allCusts.Customers = append(allCusts.Customers, newOrExistingCustomer)
	}

	returnBytes, err := json.Marshal(allCusts)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	err = writeJsonFile(returnBytes)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, "File Saved Successfully")
}

func deleteCustomer(ctx *gin.Context) {
	var customerToBeDeleted Customer
	indexOfDelete := -1

	if err := ctx.BindJSON(&customerToBeDeleted); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, "Invalid Json Given")
		return
	}

	allCusts, err := readJsonFile()

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	for i := 0; i < len(allCusts.Customers); i++ {
		if allCusts.Customers[i].CustomerID == customerToBeDeleted.CustomerID {
			indexOfDelete = i
			break
		}
	}

	log.Println("found value of i", indexOfDelete)

	if indexOfDelete != -1 {
		allCusts.Customers = append(allCusts.Customers[:indexOfDelete], allCusts.Customers[indexOfDelete+1:]...)
		log.Println("new size of custoemrs", len(allCusts.Customers))
	} else {
		ctx.IndentedJSON(http.StatusBadRequest, "Customer Doesn't Exist")
		return
	}

	returnBytes, err := json.Marshal(allCusts)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	err = writeJsonFile(returnBytes)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, "Customer Delete successfully")
}

func main() {
	//gin setup
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/allCustomers", getAllCustomers)
	router.POST("/addCustomer", addOrUpdateCustomer)
	router.POST("/deleteCustomer", deleteCustomer)
	router.Run()
}

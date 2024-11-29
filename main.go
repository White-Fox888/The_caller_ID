package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/add", handleAddNumber)
	http.HandleFunc("/", handleGetNumber)
	http.ListenAndServe(":8080", nil)
}

func handleAddNumber(w http.ResponseWriter, r *http.Request) {
	AddNumberGet := r.Header.Get("addNumber")
	TokenGet := r.Header.Get("getToken")
	token := checkToken(TokenGet)

	file, err := os.OpenFile("C:/Users/DFC/Documents/GolandProjects/awesomeProject/numberBlocked.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(w, "Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	if token {
		_, err = fmt.Fprintln(file, AddNumberGet)
		if err != nil {
			fmt.Fprintln(w, "Ошибка при записи в файл:", err)
		}
		fmt.Fprintln(w, "Номер успешно добавлен в файл!")
	} else {
		fmt.Fprintln(w, "Токен некорректен")
	}
}

func checkToken(Token string) bool {
	file, err := os.Open("C:/Users/DFC/Documents/GolandProjects/awesomeProject/token.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		savedtoken := scanner.Text()
		if Token == savedtoken {
			return true
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Ошибка при попытке сканирования:", err)
		}
	}
	return false
}

func handleGetNumber(w http.ResponseWriter, r *http.Request) {
	phoneNumberGet := r.Header.Get("getNumber")
	formattedNumber := formatPhoneNumber(phoneNumberGet)
	numbers := sliceOfBlocked()
	if blackList(formattedNumber, numbers) {
		fmt.Fprintf(w, "Number %s is blocked", formattedNumber)
	} else {
		fmt.Fprintf(w, "Number %s is not blocked", formattedNumber)
	}
}

func sliceOfBlocked() []string {
	file, err := os.Open("C:/Users/DFC/Documents/GolandProjects/awesomeProject/numberBlocked.txt")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	numbers := []string{}
	for scanner.Scan() {
		numbers = append(numbers, scanner.Text())
	}
	return numbers
}

func formatPhoneNumber(phoneNumber string) string {
	replacer := strings.NewReplacer("(", "", ")", "", "-", "")
	cleanNumber := replacer.Replace(phoneNumber)
	formattedNumber := "+7" + cleanNumber[len(cleanNumber)-10:]
	return formattedNumber
}

func blackList(formattedNumber string, numberBlocked []string) bool {
	for _, num := range numberBlocked {
		if num == formattedNumber {
			return true
		}
	}
	return false
}

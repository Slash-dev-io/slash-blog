package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const pageTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Header}}</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        .container {
            max-width: 900px;
            margin: 50px auto;
            background: white;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            font-size: 28px;
            font-weight: bold;
        }
        .content {
            padding: 30px;
            line-height: 1.8;
            color: #333;
        }
        .navigation {
            background: #f8f9fa;
            padding: 20px 30px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .nav-link {
            color: #667eea;
            text-decoration: none;
            padding: 8px 16px;
            border-radius: 5px;
            transition: all 0.3s;
        }
        .nav-link:hover {
            background: #667eea;
            color: white;
        }
        .page-info {
            color: #999;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            {{.Header}}
        </div>
        <div class="content">
            {{.Body}}
        </div>
        <div class="navigation">
            <div>
                <a href="/" class="nav-link">🏠 Главная</a>
                <a href="/about" class="nav-link">📋 О нас</a>
                <a href="/contact" class="nav-link">📧 Контакты</a>
            </div>
            <div class="page-info">
                Страница #{{.Page}}
            </div>
        </div>
    </div>
</body>
</html>
`

type Page struct {
	Page   int    `yaml:"page"`
	Header string `yaml:"header"`
	Body   string `yaml:"body"`
}

func loadPage(filename string) (*Page, error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	var page Page

	// ИСПРАВЛЕНО: передаём &page, а не err
	err = yaml.Unmarshal(data, &page)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	return &page, nil
}

func pageHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/favicon.ico" {
		http.ServeFile(w, r, "favicon.ico")
		return
	}

	path := strings.Trim(r.URL.Path, "/")
	if path == "" {
		path = "home"
	}

	filename := filepath.Join("page", path+".yaml")

	page, err := loadPage(filename)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Страница не найдена", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			log.Printf("Ошибка: %v", err)
		}
		return
	}
	tmpl, err := template.New("page").Parse(pageTemplate)
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, page)
	if err != nil {
		http.Error(w, "Ошибка отображения", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Обработчик для всех страниц
	http.HandleFunc("/", pageHandler)

	// Также можно добавить обработчик для статических файлов (CSS, JS, изображения)
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("🚀 Сервер запущен на http://localhost:8080")
	fmt.Println("Доступные страницы:")
	fmt.Println("  http://localhost:8080/       - Главная")
	fmt.Println("  http://localhost:8080/about  - О нас")
	fmt.Println("  http://localhost:8080/contact - Контакты")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

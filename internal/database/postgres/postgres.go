package postgres

import (
	"WST_lab1_server_new1/config"
	"WST_lab1_server_new1/internal/logging"
	"WST_lab1_server_new1/internal/models"

	
	"strconv"
	"strings"

	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*

 */

type Storage struct {
	DB               *gorm.DB
	PersonRepository *PersonRepository
}

type PersonRepository struct {
	DB *gorm.DB
}

/*
Инициализация
*/
func Init() (*Storage, error) {
	logging.InitializeLogger()
	var err error
	//Уровень логирования из файла конфигурации
	var logLevel logger.LogLevel
	switch config.GeneralServerSetting.LogLevel {
	case "fatal":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info", "debug":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}
	//Строка подключения к базе данных
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.DatabaseSetting.Host,
		config.DatabaseSetting.User,
		config.DatabaseSetting.Password,
		config.DatabaseSetting.Name,
		config.DatabaseSetting.Port,
		config.DatabaseSetting.SSLMode)
	//Подключаемся к базе данных
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	//Выводим при удачном подключении
	logging.Logger.Info("Database connection established successfully.")
	//Миграция базы данных
	db := conn
	err = db.AutoMigrate(&models.Person{})
	if err != nil {
		log.Fatalf("error creating table: %v", err)
		return nil, fmt.Errorf("error creating table: %v", err)
	}
	logging.Logger.Info("Migration completed successfully.")
	//Удаляем таблицу
	db.Exec("DELETE FROM people")
	//Заполняем таблицу из фаила конфигурации
	result := db.Create(&config.GeneralServerSetting.DataSet)
	if result.Error != nil {
		log.Fatalf("error creating table: %v", result.Error)
	}
	//Выводим при удачном заполнении таблицы
	logging.Logger.Info("Database updated successfully.")

	/*
		//Debug: Запрос к базе и вывод всех данных
	*/
	var results []models.Person
	if err := db.Find(&results).Error; err != nil {
		log.Fatalf("query failed: %v", err)
	}
	for _, record := range results {
		fmt.Println(record)

	}
	fmt.Println("database content in quantity:", len(results), "\n id max:", results[len(results)-1].ID, "id min:", results[0].ID)
	/*
		----
	*/
	//Возвращаем указатель

	personRepo := &PersonRepository{DB: db}
	return &Storage{
		DB:               db,
		PersonRepository: personRepo,
	}, nil

}

/*
//
Метод поиска в базе данных по запросу
//
*/
func (pr *PersonRepository) SearchPerson(searchString string) ([]models.Person, error) {
	var persons []models.Person
	query := pr.DB.Model(&models.Person{})
	//Удаляем пробелы из строки поиска
	searchString = strings.TrimSpace(searchString)
	// Проверяем строка является числом, если число ищем по возрасту
	if age, err := strconv.Atoi(searchString); err == nil {
		query = query.Where("age = ?", age)
	} else {
		//Если строка не может быть конвертирована в число ищем по строковым полям
		query = query.Where("name LIKE ? OR surname LIKE ? OR email LIKE ? OR telephone LIKE ?",
			"%"+searchString+"%", "%"+searchString+"%", "%"+searchString+"%", "%"+searchString+"%")
	}
	//Выполняем запрос и сохраняем результат в структуру
	if err := query.Find(&persons).Error; err != nil {
		return nil, err
	}
	//Возвращаем результат
	return persons, nil
}

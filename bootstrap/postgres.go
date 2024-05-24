package bootstrap

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

// Тип базы данных
type Storage struct {
	DB *sql.DB
}


// Все SQL запросы и функции работы с базой данных храняться в файле graph/resolver.go



// Функция инициализации базы данных и подключение к базе данных
func InitPostgresDatabase(cfg *Config) *Storage  {
	const op = "postgres.InitPostgresDatabase"

	dbHost := cfg.Postgres.PostgresHost
	dbPort := cfg.Postgres.PostgresPort
	dbUser := cfg.Postgres.PostgresUser
	dbPasswd := cfg.Postgres.PostgresPassword
	dbName := cfg.Postgres.DatabaseName

	postgresUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",dbHost, dbPort, dbUser, dbPasswd, dbName)
	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}

	// Создание таблицы поользователя, у которго есть зашифрованный пароль, имя и уникальный ID 
	createUserTable, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		username VARCHAR(20) NOT NULL UNIQUE,
		password CHAR(60) NOT NULL UNIQUE
	);`)
	if err != nil {	log.Fatalf("%s: %v", op, err) }
	_, err = createUserTable.Exec()
	if err != nil {	log.Fatalf("%s: %v", op, err) }

	// Создание таблицы постов, у которых есть текст, уникальный ID, а также ID пользователя, написавшего пост
	createPostTable, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS post (
		id UUID PRIMARY KEY,
		text TEXT NOT NULL,
		author_id UUID NOT NULL,
		FOREIGN KEY (author_id) REFERENCES users(id));
	`)	
	if err != nil {	log.Fatalf("%s: %v", op, err) }
	_, err = createPostTable.Exec()
	if err != nil {	log.Fatalf("%s: %v", op, err) }

	// Создание таблицы комментариев, у которых есть сам текст комменатария, ID пользователя, оставившего комментарий, а также ID того, на что был создан комментарий. Поскольку привязка к посту или комментарию происходит с помощью item_id, который не является FOREIGN KEY ни для какой таблицы, это позволяет нам оставлять коментарии любой глубины и уменьшит нагрузку на базу данных 
	createCommentTable, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS comment (
		id UUID PRIMARY KEY,
		comment VARCHAR(2000),
		author_id UUID NOT NULL,
		item_id UUID NOT NULL,
		FOREIGN KEY (author_id) REFERENCES users(id)
	);`)
	if err != nil {	log.Fatalf("%s: %v", op, err) }
	_, err = createCommentTable.Exec()
	if err != nil {	log.Fatalf("%s: %v", op, err) }

	return &Storage{DB: db}
}

// Функция закрытия соединения с базой данных
func CloseDB(db *Storage) error {
	return db.DB.Close()
}

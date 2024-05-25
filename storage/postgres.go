package bootstrap

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
	"github.com/VadimRight/GraphQLOzon/bootstrap"
)

// Тип базы данных
type PostgresStorage struct {
	DB *sql.DB
}


// Все SQL запросы и функции работы с базой данных храняться в файле graph/resolver.go, а также вспомогательные запросы для обеспечения функционала схемы и резольвера храняться в сервисе пользователей в internal/service/user.go

// Функция возвращающая объект PostgresStorage
func NewPostgresStorage(db *sql.DB) *PostgresStorage {
    return &PostgresStorage{DB: db}
}

// Функция инициализации базы данных и подключение к базе данных
func InitPostgresDatabase(cfg *bootstrap.Config) *PostgresStorage  {
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
		commentable BOOLEAN NOT NULL,
		FOREIGN KEY (author_id) REFERENCES users(id));
	`)	
	if err != nil {	log.Fatalf("%s: %v", op, err) }
	_, err = createPostTable.Exec()
	if err != nil {	log.Fatalf("%s: %v", op, err) }

	// Создание таблицы комментариев, у которых есть сам текст комменатария, ID пользователя, оставившего комментарий, а также есть ID поста, под которым комментарий был написан и это поле всегда заполяется даже если комментарий оставлен не прямо к посту, а также есть ID коментария - это поле заполяется только тогда, когда комментарий оставлен к другому коментарию. Такая конструкция сущности комментария позволяет нам создавать иерархическую структуру данных.
	createCommentTable, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS comment (
		id UUID PRIMARY KEY,
		comment VARCHAR(2000),
		author_id UUID NOT NULL,
		post_id UUID NOT NULL,
		parent_comment_id UUID,
    		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (author_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES post(id),
		FOREIGN KEY (parent_comment_id) REFERENCES comment(id)
	);`)
	if err != nil {	log.Fatalf("%s: %v", op, err) }
	_, err = createCommentTable.Exec()
	if err != nil {	log.Fatalf("%s: %v", op, err) }

	return &PostgresStorage{DB: db}
}

// Функция закрытия соединения с базой данных
func CloseDB(db *PostgresStorage) error {
	return db.DB.Close()
}

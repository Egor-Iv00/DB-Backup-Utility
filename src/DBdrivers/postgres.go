package DBdrivers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	FilePath string
}

func CreateConnectionString(conf Config) string {
	sPort := fmt.Sprintf("%d", conf.Port)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.User, conf.Password, conf.Host, sPort, conf.DBName)
}

func Recreate(ctx context.Context, conf Config) error {
	connString := CreateConnectionString(conf)
	pool, poolErr := pgxpool.New(ctx, connString)

	if poolErr != nil {
		return fmt.Errorf("failed to connect to system db: %w", poolErr)
	}
	defer pool.Close()

	dropSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", conf.DBName)
	if _, poolErr = pool.Exec(ctx, dropSQL); poolErr != nil {
		return fmt.Errorf("drop database failed: %w", poolErr)
	}

	dbIdentifier := pgx.Identifier{conf.DBName}.Sanitize()
	createSQL := fmt.Sprintf("CREATE DATABASE %s;", dbIdentifier)
	if _, poolErr = pool.Exec(ctx, createSQL); poolErr != nil {
		return fmt.Errorf("create database failed: %w", poolErr)
	}

	return nil
}

func ConnectToPostgres(conf Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connString := CreateConnectionString(conf)
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return fmt.Errorf("Failed to connect to postgres: %w", err)
	}

	closeCtx, closeCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer closeCancel()

	conn.Close(closeCtx)
	fmt.Println("Success!")

	return nil
}

func BackupPostgres(ctx context.Context, conf Config) error {
	if _, err := exec.LookPath("pg_dump"); err != nil {
		return fmt.Errorf("pg_dump not found: %w", err)
	}
	args := []string{
		"-h", conf.Host,
		"-p", fmt.Sprintf("%d", conf.Port),
		"-U", conf.User,
		"-d", conf.DBName,
		"-F", "c",
		"-f", conf.FilePath,
		"--no-owner",
		"--no-privileges",
	}
	//терминальная команда, только в коде и с использованием контекста
	cmd := exec.CommandContext(ctx, "pg_dump", args...)
	//magic
	//тема для того чтобы пароль в логах не отслеживался
	env := os.Environ()
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", conf.Password))
	//добавляем в переменные окружения пароль через данную команду
	cmd.Env = env // и подменяем старые переменные на новые

	fmt.Println("Starting backup")
	startTime := time.Now()
	err := cmd.Run()

	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("backup cancelled: %w", ctx.Err())
		}
		return fmt.Errorf("pg_dump failed: %w", err)
	}
	duration := time.Since(startTime)
	fmt.Printf("Backup completed successfully in %v\nSaved to: %s\n", duration, conf.FilePath)
	return nil
}

func RestorePostgres(ctx context.Context, conf Config) error {
	if _, err := exec.LookPath("pg_restore"); err != nil {
		return fmt.Errorf("pg_restore not found: %w", err)
	}
	if _, err := os.Stat(conf.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("File not found: %w", err)
	}

	args := []string{
		"-h", conf.Host,
		"-p", fmt.Sprintf("%d", conf.Port),
		"-U", conf.User,
		"-d", conf.DBName,
		"-F", "c",
		conf.FilePath,
	}
	cmd := exec.CommandContext(ctx, "pg_restore", args...)
	env := os.Environ()
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", conf.Password))
	cmd.Env = env

	fmt.Printf("WARNING!\n DB %s will be delete if it already exists!\n", conf.DBName)
	time.Sleep(3 * time.Second)
	fmt.Println("Starting restore")
	startTime := time.Now()
	if RecreateErr := Recreate(ctx, conf); RecreateErr != nil {
		return RecreateErr
	}
	err := cmd.Run()

	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("backup cancelled: %w", ctx.Err())
		}
		return fmt.Errorf("pg_dump failed: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("Backup completed successfully in %v\nSaved to: %s\n", duration, conf.FilePath)
	return nil
}

package MediaControllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SimpleOG/Social_Network/internal/service"
	"github.com/SimpleOG/Social_Network/pkg/util/httpResponse"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type MediaControllersInteface interface {
	UploadFile(ctx *gin.Context)
}

type MediaControllers struct {
}

var StorageUrl = "http://localhost:9090/"

func NewMediaControllers(service service.Service) MediaControllersInteface {
	return &MediaControllers{}
}

type PhotoMetadata struct {
	ID         string    `bson:"_id"`     // Уникальный идентификатор файла
	UserID     int32     `bson:"user_id"` // ID пользователя
	Name       string    `bson:"name"`    // Имя файла
	Size       int64     `bson:"size"`    // Размер файла
	Type       string    `bson:"type"`    // MIME-тип файла
	Tags       []string  `bson:"tags"`    // Кастомные теги
	UploadedAt time.Time `bson:"uploaded_at"`
}

type UploadPhotoParams struct {
	Tags []string `json:"tags"`
}

func (m *MediaControllers) UploadFile(ctx *gin.Context) {
	formFile, header, err := ctx.Request.FormFile("img")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpResponse.ErrorResponse(err))
		return
	}
	user, _ := ctx.Get("id")
	userID, ok := user.(int32)
	if !ok {
		ctx.JSON(http.StatusUnprocessableEntity, httpResponse.ErrorResponse(errors.New("ошибка парсинга id пользователя")))
		return
	}
	var PhotoParams UploadPhotoParams
	data := ctx.Request.FormValue("metadata")
	if err = json.Unmarshal([]byte(data), &PhotoParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	var metadata = &PhotoMetadata{
		ID:     "",
		UserID: userID,
		Name:   header.Filename,
		Size:   header.Size,
		Type:   "",
		Tags:   PhotoParams.Tags,
	}
	if err = sendFileIntoStorage(metadata, formFile); err != nil {
		log.Println("Ошибка в filestorage")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, metadata)

}

func sendFileIntoStorage(metadata *PhotoMetadata, file multipart.File) error {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	filePart, err := writer.CreateFormFile("img", metadata.Name)
	if err != nil {
		return fmt.Errorf("ошибка при создании части файла: %v", err)
	}
	_, err = io.Copy(filePart, file)
	if err != nil {
		return fmt.Errorf("ошибка при копировании файла: %v", err)
	}
	metaJson, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	err = writer.WriteField("metadata", string(metaJson))
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("ошибка при закрытии writer: %v", err)
	}
	fmt.Printf("Размер multipart-запроса: %d байт\n", requestBody.Len())
	req, err := http.NewRequest("POST", StorageUrl+"upload_image", &requestBody)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Ошибка при чтении тела ответа: %v\n", err)
			return errors.New(fmt.Sprintf("Ошибка при чтении тела ответа: %v\n", err))
		}
		log.Printf("Ошибка от сервера (код %d): %s\n", resp.StatusCode, string(body))
		return errors.New(fmt.Sprintf("Ошибка от сервера (код %d): %s\n", resp.StatusCode, string(body)))
	}
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &metadata)
	if err != nil {
		return err
	}
	return nil
}

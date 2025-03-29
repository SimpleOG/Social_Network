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
	"mime"
	"mime/multipart"
	"net/http"
	"time"
)

type MediaControllersInteface interface {
	UploadImage(ctx *gin.Context)
	DownloadImage(ctx *gin.Context)
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

func (m *MediaControllers) UploadImage(ctx *gin.Context) {
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
	if err = sendFileIntoUrl(http.MethodPost, StorageUrl+"upload_image", metadata, formFile); err != nil {
		log.Println("Ошибка в filestorage")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, metadata)

}
func CreateMultipartForm(metadata *PhotoMetadata, file io.Reader) (bytes.Buffer, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	filePart, err := writer.CreateFormFile("img", metadata.Name)
	if err != nil {
		return body, writer.FormDataContentType(), fmt.Errorf("ошибка при создании части файла: %v", err)
	}

	_, err = io.Copy(filePart, file)
	if err != nil {
		return body, writer.FormDataContentType(), fmt.Errorf("ошибка при копировании файла: %v", err)
	}

	metaJson, err := json.Marshal(metadata)
	if err != nil {
		return body, writer.FormDataContentType(), err
	}

	err = writer.WriteField("metadata", string(metaJson))
	if err != nil {
		return body, writer.FormDataContentType(), err
	}

	err = writer.Close()
	if err != nil {
		return body, writer.FormDataContentType(), fmt.Errorf("ошибка при закрытии writer: %v", err)
	}

	fmt.Printf("Размер multipart-запроса: %d байт\n", body.Len())
	return body, writer.FormDataContentType(), nil
}
func sendFileIntoUrl(httpMethod, url string, metadata *PhotoMetadata, file io.Reader) error {
	requestBody, header, err := CreateMultipartForm(metadata, file)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(httpMethod, url, &requestBody)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	req.Header.Set("Content-Type", header)

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

func (m *MediaControllers) DownloadImage(ctx *gin.Context) {
	//user, _ := ctx.Get("id")
	//userID, ok := user.(int32)
	//if !ok {
	//	ctx.JSON(http.StatusUnprocessableEntity, httpResponse.ErrorResponse(errors.New("ошибка парсинга id пользователя")))
	//	return
	//}
	userID := int32(1)
	meta, file, err := DownloadImageFromStorage(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpResponse.ErrorResponse(err))
		return
	}
	body, header, err := CreateMultipartForm(meta, file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpResponse.ErrorResponse(err))
	}
	ctx.Header("Content-Type", header)
	ctx.Data(http.StatusOK, header, body.Bytes())
}

type Tags struct {
	Tags []string `json:"tags"`
}

func DownloadImageFromStorage(id int32) (*PhotoMetadata, io.Reader, error) {
	var tags Tags
	tags.Tags = []string{"avatar"}
	body, err := json.Marshal(tags)
	resp, err := http.Post(StorageUrl+fmt.Sprintf("get_image/?id=%d", id), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}
	//Парсинг ответа
	if resp.StatusCode > 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Ошибка при чтении тела ответа: %v\n", err)
			return nil, nil, errors.New(fmt.Sprintf("Ошибка при чтении тела ответа: %v\n", err))
		}
		log.Printf("Ошибка от сервера (код %d): %s\n", resp.StatusCode, string(body))
		return nil, nil, errors.New(fmt.Sprintf("Ошибка от сервера (код %d): %s\n", resp.StatusCode, string(body)))
	}
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		log.Fatalf("Error parsing Content-Type: %v", err)
	}
	boundary, ok := params["boundary"]
	if !ok {
		log.Fatalf("No boundary in Content-Type")
	}
	reader := multipart.NewReader(resp.Body, boundary)

	// Переменные для хранения данных
	var metadata PhotoMetadata
	var imgData bytes.Buffer
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break // Закончили чтение
		}
		if err != nil {
			log.Fatalf("Error reading part: %v", err)
		}

		// Обрабатываем каждую часть
		switch part.FormName() {
		case "metadata":
			// Читаем метаданные (JSON)
			if err := json.NewDecoder(part).Decode(&metadata); err != nil {
				log.Fatalf("Error decoding metadata: %v", err)
			}

		case "img":
			// Читаем бинарные данные изображения
			if _, err := io.Copy(&imgData, part); err != nil {
				log.Fatalf("Error reading image data: %v", err)
			}

		default:
			log.Printf("Unknown part: %s\n", part.FormName())
		}
	}
	return &metadata, &imgData, nil

}

package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-bot/repository"
)

type Client interface {
	Updates(offset, limit int) ([]tgbotapi.Update, error)

	SendMessage(chatID int64,
		message string,
		replyMarkup *tgbotapi.ReplyKeyboardMarkup) error

	SendMessageWithFile(chatID int64,
		fileInfo repository.File,
		caption string,
		replyMarkup *tgbotapi.ReplyKeyboardMarkup) error

	SendMessageWithFiles(chatID int64,
		filesInfo []repository.File, caption string) error
}

//
//type client struct {
//	host     string
//	basePath string
//	client   http.Client
//}
//
////func New(host, token string) Client {
////	return &client{
////		host:     host,
////		basePath: newBasePath(token),
////		client:   http.Client{},
////	}
////}
//
//func newBasePath(token string) string {
//	return "bot" + token
//}
//
//const (
//	contentTypeJson     = "application/json"
//	getUpdateMethod     = "getUpdates"
//	sendMessageMethod   = "sendMessage"
//	sendPhotoMethod     = "sendPhoto"
//	setMyCommandsMethod = "setMyCommands"
//)
//
//func (c *client) Updates(offset, limit int) ([]Update, error) {
//	q := url.Values{}
//	q.Add("offset", strconv.Itoa(offset))
//	q.Add("limit", strconv.Itoa(limit))
//
//	data, err := c.doGetRequest(getUpdateMethod, q)
//	if err != nil {
//		return nil, err
//	}
//
//	var res UpdatesResponse
//	if err := json.Unmarshal(data, &res); err != nil {
//		return nil, err
//	}
//
//	return res.Result, nil
//}
//
//func (c *client) doGetRequest(method string, query url.Values) (data []byte, err error) {
//	const errMsg = "can't do request with GET method"
//
//	defer func() { err = e.WrapIfErr(errMsg, err) }()
//
//	u := url.URL{
//		Scheme: "https",
//		Host:   c.host,
//		Path:   path.Join(c.basePath, method),
//	}
//
//	u.RawQuery = query.Encode()
//	resp, err := c.client.Get(u.String())
//	if err != nil {
//		return nil, err
//	}
//
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	return body, nil
//}
//
//func (c *client) SendMessage(chatID int64, message string, buttons *ReplyMarkupData) error {
//	data := SendMessageRequest{
//		ChatID:      chatID,
//		Text:        message,
//		ReplyMarkup: buttons,
//	}
//	u := url.URL{
//		Scheme: "https",
//		Host:   c.host,
//		Path:   path.Join(c.basePath, sendMessageMethod),
//	}
//
//	return c.doPostRequest(u, data)
//}
//
//func (c *client) SendPhotoMessage(chatID int64, photo []byte, caption string, buttons *ReplyMarkupData) error {
//	requestBody := &bytes.Buffer{}
//	writer := multipart.NewWriter(requestBody)
//	data := SendPhotoRequest{
//		ChatID:      chatID,
//		Photo:       photo,
//		Caption:     caption,
//		ReplyMarkup: buttons,
//	}
//	u := url.URL{
//		Scheme: "https",
//		Host:   c.host,
//		Path:   path.Join(c.basePath, sendPhotoMethod),
//	}
//
//	err := c.writeMultipartStruct(writer, data)
//	req, err := http.NewRequest(http.MethodPost, u.String(), requestBody)
//	if err != nil {
//		return err
//	}
//
//	req.Header.Set("Content-Type", writer.FormDataContentType())
//
//	resp, err := c.client.Do(req)
//	if err != nil {
//		return err
//	}
//	log.Println("http.Post done!")
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return c.failedRequest(resp)
//	}
//
//	return nil
//}
//
//func (c *client) writeMultipartStruct(writer *multipart.Writer, data interface{}) error {
//	value := reflect.ValueOf(data)
//	if value.Kind() != reflect.Struct {
//		return fmt.Errorf("data is not a struct")
//	}
//
//	typ := value.Type()
//	for i := 0; i < typ.NumField(); i++ {
//		field := typ.Field(i)
//		fieldValue := value.Field(i)
//
//		if fieldValue.Kind() == reflect.Ptr {
//			log.Println(field.Tag.Get("json"))
//			if fieldValue.IsNil() {
//				continue
//			}
//			fieldValue = fieldValue.Elem()
//		}
//
//		key := field.Tag.Get("json")
//		if key == "" {
//			key = field.Name
//		}
//
//		switch fieldValue.Kind() {
//		case reflect.String:
//			log.Println(key, "string")
//			writer.WriteField(key, fieldValue.String())
//
//		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//			log.Println(key, "int")
//			writer.WriteField(key, strconv.FormatInt(fieldValue.Int(), 10))
//
//		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//			log.Println(key, "uint")
//			writer.WriteField(key, strconv.FormatUint(fieldValue.Uint(), 10))
//
//		case reflect.Bool:
//			log.Println(key, "bool")
//			writer.WriteField(key, strconv.FormatBool(fieldValue.Bool()))
//
//		case reflect.Slice:
//			log.Println(key, "slice")
//			if reflect.TypeOf(fieldValue.Interface()).Elem().Kind() == reflect.Uint8 {
//				log.Println(key, "[]byte"+
//					"")
//				// This is a byte slice (for a file, for example)
//				filename := uuid.New().String()
//				part, err := writer.CreateFormFile(key, filename)
//				if err != nil {
//					return err
//				}
//
//				if _, err := io.Copy(part, bytes.NewReader(fieldValue.Bytes())); err != nil {
//					return err
//				}
//			} else {
//				// This is a slice of some other type
//				for j := 0; j < fieldValue.Len(); j++ {
//					err := c.writeMultipartStruct(writer, fieldValue.Index(j).Interface())
//					if err != nil {
//						return err
//					}
//				}
//			}
//
//		case reflect.Struct:
//			log.Println(key, "struct")
//			err := c.writeMultipartStruct(writer, fieldValue.Interface())
//			if err != nil {
//				return err
//			}
//
//		case reflect.Ptr:
//			log.Println(key, "ptr")
//			if !fieldValue.IsNil() {
//				err := c.writeMultipartStruct(writer, fieldValue.Interface())
//				if err != nil {
//					return err
//				}
//			}
//
//		default:
//			return fmt.Errorf("unsupported field type: %v", fieldValue.Kind())
//		}
//	}
//
//	return nil
//}
//
//func (c *client) doPostRequest(u url.URL, data interface{}) (err error) {
//	log.Println("do Post Request URL:", u.String())
//	errMsg := "can't do request with POST method"
//	defer func() { err = e.WrapIfErr(errMsg, err) }()
//	body, err := json.Marshal(data)
//	if err != nil {
//		return err
//	}
//
//	resp, err := c.client.Post(u.String(), contentTypeJson, bytes.NewBuffer(body))
//	if err != nil {
//		return err
//	}
//
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return c.failedRequest(resp)
//	}
//
//	return nil
//}
//
//func (c *client) failedRequest(resp *http.Response) error {
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return err
//	}
//	errResp := ErrResponse{}
//	err = json.Unmarshal(body, &errResp)
//	if err != nil {
//		return err
//	}
//
//	return fmt.Errorf("%s", errResp.Description)
//}

package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/GodCratos/google_sheets/configs"
)

func RetailGetOrdersByPages(page int) ([]interface{}, error) {
	fmt.Println("----------------------------------------------------------------------------")
	url := fmt.Sprintf(configs.RetailGetOrders(), page)
	log.Println("[SERVICES:RETAIL] Start sending request : ", url)
	clientGet := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("[SERVICES:RETAIL] Error while creating a new request. Error description : ", err)
		return nil, errors.New(fmt.Sprintf("[SERVICES:RETAIL] Ошибка при создании нового запроса. Error description : %s", err.Error()))
	}
	resp, err := clientGet.Do(req)
	if err != nil {
		log.Println("[SERVICES:RETAIL] Error while sending request. Error description : ", err)
		return nil, errors.New(fmt.Sprintf("[SERVICES:RETAIL] Ошибка при отправке запроса в Retail. Error description : %s", err.Error()))
	}
	defer resp.Body.Close()
	respByte, err := ioutil.ReadAll(resp.Body)
	log.Println("[SERVICES:RETAIL] Response from Retail")
	retailStruct, err := RetailParserJSON(respByte)
	if err != nil {
		return nil, err
	}
	if len(retailStruct["orders"].([]interface{})) == 0 {
		log.Println("[SERVICES:RETAIL] Orders array is empty")
		return nil, nil
	}
	return retailStruct["orders"].([]interface{}), nil
}

func RetailGetNameStatusOrder(status string) string {
	url := "https://love-piano.retailcrm.ru/api/v5/reference/statuses?apiKey=wPOqarGJ6EvfxzjmlN40C24NOqqp0YUr"
	clientGet := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := clientGet.Do(req)
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	retailStruct, _ := RetailParserJSON(respByte)
	statuses := retailStruct["statuses"].(map[string]interface{})
	return statuses[status].(map[string]interface{})["name"].(string)
}

func RetailGetNameOrderMethod(method string) string {
	url := "https://love-piano.retailcrm.ru/api/v5/reference/order-methods?apiKey=wPOqarGJ6EvfxzjmlN40C24NOqqp0YUr"
	clientGet := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := clientGet.Do(req)
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	retailStruct, _ := RetailParserJSON(respByte)
	methods := retailStruct["orderMethods"].(map[string]interface{})
	return methods[method].(map[string]interface{})["name"].(string)
}

func RetailGetNameManager(managerID float64) string {
	url := fmt.Sprintf("https://love-piano.retailcrm.ru/api/v5/users/%v?apiKey=wPOqarGJ6EvfxzjmlN40C24NOqqp0YUr", managerID)
	clientGet := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := clientGet.Do(req)
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	retailStruct, _ := RetailParserJSON(respByte)
	user := retailStruct["user"].(map[string]interface{})
	if _, ok := user["lastName"]; ok {
		return fmt.Sprintf("%s %s", user["firstName"].(string), user["lastName"].(string))
	}
	return user["firstName"].(string)
}

func RetailGetNameShop(siteName string) string {
	url := "https://love-piano.retailcrm.ru/api/v5/reference/sites?apiKey=wPOqarGJ6EvfxzjmlN40C24NOqqp0YUr"
	clientGet := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := clientGet.Do(req)
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	retailStruct, _ := RetailParserJSON(respByte)
	sites := retailStruct["sites"].(map[string]interface{})
	return sites[siteName].(map[string]interface{})["name"].(string)
}

func RetailParserJSON(value []byte) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal(value, &jsonMap)
	if err != nil {
		log.Println("[SERVICES:RETAIL] Error while parsing JSON. Error description : ", err)
		return nil, errors.New(fmt.Sprintf("[SERVICES:RETAIL] Ошибка при разборе JSON. Error description : %s", err.Error()))
	}
	return jsonMap, nil
}

func RetailStructGenerationForGoogleSheets(orderStruct map[string]interface{}) []interface{} {

	var structSheets []interface{}
	if value, ok := orderStruct["site"]; ok {
		structSheets = append(structSheets, RetailGetNameShop(value.(string)))
	} else {
		structSheets = append(structSheets, "")
	}

	if value, ok := orderStruct["createdAt"]; ok {
		structSheets = append(structSheets, value.(string))
	} else {
		structSheets = append(structSheets, "")
	}

	if value, ok := orderStruct["number"]; ok {
		structSheets = append(structSheets, value.(string))
	} else {
		structSheets = append(structSheets, "")
	}

	if value, ok := orderStruct["status"]; ok {
		structSheets = append(structSheets, RetailGetNameStatusOrder(value.(string)))
	} else {
		structSheets = append(structSheets, "")
	}

	if value, ok := orderStruct["managerId"]; ok {
		structSheets = append(structSheets, RetailGetNameManager(value.(float64)))
	} else {
		structSheets = append(structSheets, "")
	}

	if value, ok := orderStruct["totalSumm"]; ok {
		structSheets = append(structSheets, value.(float64))
	} else {
		structSheets = append(structSheets, "")
	}

	if value, ok := orderStruct["delivery"].(map[string]interface{})["date"]; ok {
		structSheets = append(structSheets, value.(string))
	} else {
		structSheets = append(structSheets, "")
	}

	if value, ok := orderStruct["orderMethod"]; ok {
		structSheets = append(structSheets, RetailGetNameOrderMethod(value.(string)))
	} else {
		structSheets = append(structSheets, "")
	}

	structItem := ""
	article := ""
	maxPrice := 0.00
	maxItemName := ""
	for _, item := range orderStruct["items"].([]interface{}) {
		if articleCRM, ok := item.(map[string]interface{})["offer"].(map[string]interface{})["article"]; ok {
			article = articleCRM.(string)
		}
		structItem += fmt.Sprintf("%s/%s,%v шт.;", item.(map[string]interface{})["offer"].(map[string]interface{})["name"].(string), article, item.(map[string]interface{})["quantity"].(float64))
		if maxPrice < item.(map[string]interface{})["initialPrice"].(float64) {
			maxPrice = item.(map[string]interface{})["initialPrice"].(float64)
			maxItemName = item.(map[string]interface{})["offer"].(map[string]interface{})["name"].(string)
		}
	}
	structSheets = append(structSheets, structItem)
	if maxPrice != 0.00 {
		structSheets = append(structSheets, fmt.Sprintf("%s; %v RUB", maxItemName, maxPrice))
	} else {
		structSheets = append(structSheets, 0)
	}

	if value, ok := orderStruct["delivery"].(map[string]interface{})["address"]; ok {
		if _, ok := value.(map[string]interface{})["region"]; ok {
			structSheets = append(structSheets, value.(map[string]interface{})["region"].(string))
		} else {
			structSheets = append(structSheets, "")
		}
	} else {
		structSheets = append(structSheets, "")
	}

	if value, ok := orderStruct["delivery"].(map[string]interface{})["address"]; ok {
		if _, ok := value.(map[string]interface{})["city"]; ok {
			structSheets = append(structSheets, value.(map[string]interface{})["city"].(string))
		} else {
			structSheets = append(structSheets, "")
		}
	} else {
		structSheets = append(structSheets, "")
	}

	if value, ok := orderStruct["customFields"].(map[string]interface{})["cancellation_reason"]; ok {
		structSheets = append(structSheets, value.(string))
	} else {
		structSheets = append(structSheets, "")
	}

	return structSheets
}
